// Package cache provides caching features for data from a Consul server.
//
// While this is similar in some ways to the "agent/ae" package, a key
// difference is that with anti-entropy, the agent is the authoritative
// source so it resolves differences the server may have. With caching (this
// package), the server is the authoritative source and we do our best to
// balance performance and correctness, depending on the type of data being
// requested.
//
// The types of data that can be cached is configurable via the Type interface.
// This allows specialized behavior for certain types of data. Each type of
// Consul data (CA roots, leaf certs, intentions, KV, catalog, etc.) will
// have to be manually implemented. This usually is not much work, see
// the "agent/cache-types" package.
package cache

import (
	"container/heap"
	"context"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/consul/lib"
	"golang.org/x/time/rate"
)

//go:generate mockery -all -inpkg

// Constants related to refresh backoff. We probably don't ever need to
// make these configurable knobs since they primarily exist to lower load.
const (
	CacheRefreshBackoffMin = 3               // 3 attempts before backing off
	CacheRefreshMaxWait    = 1 * time.Minute // maximum backoff wait time
)

// Cache is a agent-local cache of Consul data. Create a Cache using the
// New function. A zero-value Cache is not ready for usage and will result
// in a panic.
//
// The types of data to be cached must be registered via RegisterType. Then,
// calls to Get specify the type and a Request implementation. The
// implementation of Request is usually done directly on the standard RPC
// struct in agent/structs.  This API makes cache usage a mostly drop-in
// replacement for non-cached RPC calls.
//
// The cache is partitioned by ACL and datacenter. This allows the cache
// to be safe for multi-DC queries and for queries where the data is modified
// due to ACLs all without the cache having to have any clever logic, at
// the slight expense of a less perfect cache.
//
// The Cache exposes various metrics via go-metrics. Please view the source
// searching for "metrics." to see the various metrics exposed. These can be
// used to explore the performance of the cache.
type Cache struct {
	// types stores the list of data types that the cache knows how to service.
	// These can be dynamically registered with RegisterType.
	typesLock sync.RWMutex
	types     map[string]typeEntry

	// entries contains the actual cache data. Access to entries and
	// entriesExpiryHeap must be protected by entriesLock.
	//
	// entriesExpiryHeap is a heap of *cacheEntry values ordered by
	// expiry, with the soonest to expire being first in the list (index 0).
	//
	// NOTE(mitchellh): The entry map key is currently a string in the format
	// of "<DC>/<ACL token>/<Request key>" in order to properly partition
	// requests to different datacenters and ACL tokens. This format has some
	// big drawbacks: we can't evict by datacenter, ACL token, etc. For an
	// initial implementation this works and the tests are agnostic to the
	// internal storage format so changing this should be possible safely.
	entriesLock       sync.RWMutex
	entries           map[string]cacheEntry
	entriesExpiryHeap *expiryHeap

	// stopped is used as an atomic flag to signal that the Cache has been
	// discarded so background fetches and expiry processing should stop.
	stopped uint32
	// stopCh is closed when Close is called
	stopCh chan struct{}
	// options includes a per Cache Rate limiter specification to avoid performing too many queries
	options          Options
	rateLimitContext context.Context
	rateLimitCancel  context.CancelFunc
}

// typeEntry is a single type that is registered with a Cache.
type typeEntry struct {
	// Name that was used to register the Type
	Name string
	Type Type
	Opts *RegisterOptions
}

// ResultMeta is returned from Get calls along with the value and can be used
// to expose information about the cache status for debugging or testing.
type ResultMeta struct {
	// Hit indicates whether or not the request was a cache hit
	Hit bool

	// Age identifies how "stale" the result is. It's semantics differ based on
	// whether or not the cache type performs background refresh or not as defined
	// in https://www.consul.io/api/index.html#agent-caching.
	//
	// For background refresh types, Age is 0 unless the background blocking query
	// is currently in a failed state and so not keeping up with the server's
	// values. If it is non-zero it represents the time since the first failure to
	// connect during background refresh, and is reset after a background request
	// does manage to reconnect and either return successfully, or block for at
	// least the yamux keepalive timeout of 30 seconds (which indicates the
	// connection is OK but blocked as expected).
	//
	// For simple cache types, Age is the time since the result being returned was
	// fetched from the servers.
	Age time.Duration

	// Index is the internal ModifyIndex for the cache entry. Not all types
	// support blocking and all that do will likely have this in their result type
	// already but this allows generic code to reason about whether cache values
	// have changed.
	Index uint64
}

// Options are options for the Cache.
type Options struct {
	// EntryFetchMaxBurst max burst size of RateLimit for a single cache entry
	EntryFetchMaxBurst int
	// EntryFetchRate represents the max calls/sec for a single cache entry
	EntryFetchRate rate.Limit
}

// New creates a new cache with the given RPC client and reasonable defaults.
// Further settings can be tweaked on the returned value.
func New(options Options) *Cache {
	// Initialize the heap. The buffer of 1 is really important because
	// its possible for the expiry loop to trigger the heap to update
	// itself and it'd block forever otherwise.
	h := &expiryHeap{NotifyCh: make(chan struct{}, 1)}
	heap.Init(h)
	ctx, cancel := context.WithCancel(context.Background())
	c := &Cache{
		types:             make(map[string]typeEntry),
		entries:           make(map[string]cacheEntry),
		entriesExpiryHeap: h,
		stopCh:            make(chan struct{}),
		options:           options,
		rateLimitContext:  ctx,
		rateLimitCancel:   cancel,
	}

	// Start the expiry watcher
	go c.runExpiryLoop()

	return c
}

// RegisterOptions are options that can be associated with a type being
// registered for the cache. This changes the behavior of the cache for
// this type.
type RegisterOptions struct {
	// LastGetTTL is the time that the values returned by this type remain
	// in the cache after the last get operation. If a value isn't accessed
	// within this duration, the value is purged from the cache and
	// background refreshing will cease.
	LastGetTTL time.Duration

	// Refresh configures whether the data is actively refreshed or if
	// the data is only refreshed on an explicit Get. The default (false)
	// is to only request data on explicit Get.
	Refresh bool

	// SupportsBlocking should be set to true if the type supports blocking queries.
	// Types that do not support blocking queries will not be able to use
	// background refresh nor will the cache attempt blocking fetches if the
	// client requests them with MinIndex.
	SupportsBlocking bool

	// RefreshTimer is the time to sleep between attempts to refresh data.
	// If this is zero, then data is refreshed immediately when a fetch
	// is returned.
	//
	// Using different values for RefreshTimer and QueryTimeout, various
	// "refresh" mechanisms can be implemented:
	//
	//   * With a high timer duration and a low timeout, a timer-based
	//     refresh can be set that minimizes load on the Consul servers.
	//
	//   * With a low timer and high timeout duration, a blocking-query-based
	//     refresh can be set so that changes in server data are recognized
	//     within the cache very quickly.
	//
	RefreshTimer time.Duration

	// QueryTimeout is the default value for the maximum query time for a fetch
	// operation. It is set as FetchOptions.Timeout so that cache.Type
	// implementations can use it as the MaxQueryTime.
	QueryTimeout time.Duration
}

// RegisterType registers a cacheable type.
//
// This makes the type available for Get but does not automatically perform
// any prefetching. In order to populate the cache, Get must be called.
func (c *Cache) RegisterType(n string, typ Type) {
	opts := typ.RegisterOptions()
	if opts.LastGetTTL == 0 {
		opts.LastGetTTL = 72 * time.Hour // reasonable default is days
	}

	c.typesLock.Lock()
	defer c.typesLock.Unlock()
	c.types[n] = typeEntry{Name: n, Type: typ, Opts: &opts}
}

// Get loads the data for the given type and request. If data satisfying the
// minimum index is present in the cache, it is returned immediately. Otherwise,
// this will block until the data is available or the request timeout is
// reached.
//
// Multiple Get calls for the same Request (matching CacheKey value) will
// block on a single network request.
//
// The timeout specified by the Request will be the timeout on the cache
// Get, and does not correspond to the timeout of any background data
// fetching. If the timeout is reached before data satisfying the minimum
// index is retrieved, the last known value (maybe nil) is returned. No
// error is returned on timeout. This matches the behavior of Consul blocking
// queries.
func (c *Cache) Get(ctx context.Context, t string, r Request) (interface{}, ResultMeta, error) {
	c.typesLock.RLock()
	tEntry, ok := c.types[t]
	c.typesLock.RUnlock()
	if !ok {
		// Shouldn't happen given that we successfully fetched this at least
		// once. But be robust against panics.
		return nil, ResultMeta{}, fmt.Errorf("unknown type in cache: %s", t)
	}
	return c.getWithIndex(ctx, newGetOptions(tEntry, r))
}

// getOptions contains the arguments for a Get request. It is used in place of
// Request so that internal functions can modify Info without having to extract
// it from the Request each time.
type getOptions struct {
	// Fetch is a closure over tEntry.Type.Fetch which provides the original
	// Request from the caller.
	Fetch     func(opts FetchOptions) (FetchResult, error)
	Info      RequestInfo
	TypeEntry typeEntry
}

func newGetOptions(tEntry typeEntry, r Request) getOptions {
	return getOptions{
		Fetch: func(opts FetchOptions) (FetchResult, error) {
			return tEntry.Type.Fetch(opts, r)
		},
		Info:      r.CacheInfo(),
		TypeEntry: tEntry,
	}
}

// getEntryLocked retrieves a cache entry and checks if it is ready to be
// returned given the other parameters. It reads from entries and the caller
// has to issue a read lock if necessary.
func (c *Cache) getEntryLocked(
	tEntry typeEntry,
	key string,
	info RequestInfo,
) (entryExists bool, entryValid bool, entry cacheEntry) {
	entry, ok := c.entries[key]
	if !entry.Valid {
		return ok, false, entry
	}

	// Check index is not specified or lower than value, or the type doesn't
	// support blocking.
	if tEntry.Opts.SupportsBlocking && info.MinIndex > 0 && info.MinIndex >= entry.Index {
		// MinIndex was given and matches or is higher than current value so we
		// ignore the cache and fallthrough to blocking on a new value below.
		return true, false, entry
	}

	// Check MaxAge is not exceeded if this is not a background refreshing type
	// and MaxAge was specified.
	if !tEntry.Opts.Refresh && info.MaxAge > 0 && entryExceedsMaxAge(info.MaxAge, entry) {
		return true, false, entry
	}

	// Check if re-validate is requested. If so the first time round the
	// loop is not a hit but subsequent ones should be treated normally.
	if !tEntry.Opts.Refresh && info.MustRevalidate {
		return true, false, entry
	}

	return true, true, entry
}

func entryExceedsMaxAge(maxAge time.Duration, entry cacheEntry) bool {
	return !entry.FetchedAt.IsZero() && maxAge < time.Since(entry.FetchedAt)
}

// getWithIndex implements the main Get functionality but allows internal
// callers (Watch) to manipulate the blocking index separately from the actual
// request object.
func (c *Cache) getWithIndex(ctx context.Context, r getOptions) (interface{}, ResultMeta, error) {
	if r.Info.Key == "" {
		metrics.IncrCounter([]string{"consul", "cache", "bypass"}, 1)

		// If no key is specified, then we do not cache this request.
		// Pass directly through to the backend.
		result, err := r.Fetch(FetchOptions{MinIndex: r.Info.MinIndex})
		return result.Value, ResultMeta{}, err
	}

	key := makeEntryKey(r.TypeEntry.Name, r.Info.Datacenter, r.Info.Token, r.Info.Key)

	// First time through
	first := true

	// timeoutCh for watching our timeout
	var timeoutCh <-chan time.Time

RETRY_GET:
	// Get the current value
	c.entriesLock.RLock()
	_, entryValid, entry := c.getEntryLocked(r.TypeEntry, key, r.Info)
	c.entriesLock.RUnlock()

	if entryValid {
		meta := ResultMeta{Index: entry.Index}
		if first {
			metrics.IncrCounter([]string{"consul", "cache", r.TypeEntry.Name, "hit"}, 1)
			meta.Hit = true
		}

		// If refresh is enabled, calculate age based on whether the background
		// routine is still connected.
		if r.TypeEntry.Opts.Refresh {
			meta.Age = time.Duration(0)
			if !entry.RefreshLostContact.IsZero() {
				meta.Age = time.Since(entry.RefreshLostContact)
			}
		} else {
			// For non-background refresh types, the age is just how long since we
			// fetched it last.
			if !entry.FetchedAt.IsZero() {
				meta.Age = time.Since(entry.FetchedAt)
			}
		}

		// Touch the expiration and fix the heap.
		c.entriesLock.Lock()
		entry.Expiry.Update(r.TypeEntry.Opts.LastGetTTL)
		c.entriesExpiryHeap.Fix(entry.Expiry)
		c.entriesLock.Unlock()

		// We purposely do not return an error here since the cache only works with
		// fetching values that either have a value or have an error, but not both.
		// The Error may be non-nil in the entry in the case that an error has
		// occurred _since_ the last good value, but we still want to return the
		// good value to clients that are not requesting a specific version. The
		// effect of this is that blocking clients will all see an error immediately
		// without waiting a whole timeout to see it, but clients that just look up
		// cache with an older index than the last valid result will still see the
		// result and not the error here. I.e. the error is not "cached" without a
		// new fetch attempt occurring, but the last good value can still be fetched
		// from cache.
		return entry.Value, meta, nil
	}

	// If this isn't our first time through and our last value has an error, then
	// we return the error. This has the behavior that we don't sit in a retry
	// loop getting the same error for the entire duration of the timeout.
	// Instead, we make one effort to fetch a new value, and if there was an
	// error, we return. Note that the invariant is that if both entry.Value AND
	// entry.Error are non-nil, the error _must_ be more recent than the Value. In
	// other words valid fetches should reset the error. See
	// https://github.com/hashicorp/consul/issues/4480.
	if !first && entry.Error != nil {
		return entry.Value, ResultMeta{Index: entry.Index}, entry.Error
	}

	if first {
		// We increment two different counters for cache misses depending on
		// whether we're missing because we didn't have the data at all,
		// or if we're missing because we're blocking on a set index.
		missKey := "miss_block"
		if r.Info.MinIndex == 0 {
			missKey = "miss_new"
		}
		metrics.IncrCounter([]string{"consul", "cache", r.TypeEntry.Name, missKey}, 1)
	}

	// Set our timeout channel if we must
	if r.Info.Timeout > 0 && timeoutCh == nil {
		timeoutCh = time.After(r.Info.Timeout)
	}

	// At this point, we know we either don't have a value at all or the
	// value we have is too old. We need to wait for new data.
	waiterCh := c.fetch(key, r, true, 0, false)

	// No longer our first time through
	first = false

	select {
	case <-ctx.Done():
		return nil, ResultMeta{}, ctx.Err()
	case <-waiterCh:
		// Our fetch returned, retry the get from the cache.
		r.Info.MustRevalidate = false
		goto RETRY_GET

	case <-timeoutCh:
		// Timeout on the cache read, just return whatever we have.
		return entry.Value, ResultMeta{Index: entry.Index}, nil
	}
}

func makeEntryKey(t, dc, token, key string) string {
	return fmt.Sprintf("%s/%s/%s/%s", t, dc, token, key)
}

// fetch triggers a new background fetch for the given Request. If a
// background fetch is already running for a matching Request, the waiter
// channel for that request is returned. The effect of this is that there
// is only ever one blocking query for any matching requests.
//
// If allowNew is true then the fetch should create the cache entry
// if it doesn't exist. If this is false, then fetch will do nothing
// if the entry doesn't exist. This latter case is to support refreshing.
func (c *Cache) fetch(key string, r getOptions, allowNew bool, attempt uint, ignoreExisting bool) <-chan struct{} {
	// We acquire a write lock because we may have to set Fetching to true.
	c.entriesLock.Lock()
	defer c.entriesLock.Unlock()
	ok, entryValid, entry := c.getEntryLocked(r.TypeEntry, key, r.Info)

	// This handles the case where a fetch succeeded after checking for its existence in
	// getWithIndex. This ensures that we don't miss updates.
	if ok && entryValid && !ignoreExisting {
		ch := make(chan struct{})
		close(ch)
		return ch
	}

	// If we aren't allowing new values and we don't have an existing value,
	// return immediately. We return an immediately-closed channel so nothing
	// blocks.
	if !ok && !allowNew {
		ch := make(chan struct{})
		close(ch)
		return ch
	}

	// If we already have an entry and it is actively fetching, then return
	// the currently active waiter.
	if ok && entry.Fetching {
		return entry.Waiter
	}

	// If we don't have an entry, then create it. The entry must be marked
	// as invalid so that it isn't returned as a valid value for a zero index.
	if !ok {
		entry = cacheEntry{
			Valid:  false,
			Waiter: make(chan struct{}),
			FetchRateLimiter: rate.NewLimiter(
				c.options.EntryFetchRate,
				c.options.EntryFetchMaxBurst,
			),
		}
	}

	// Set that we're fetching to true, which makes it so that future
	// identical calls to fetch will return the same waiter rather than
	// perform multiple fetches.
	entry.Fetching = true
	c.entries[key] = entry
	metrics.SetGauge([]string{"consul", "cache", "entries_count"}, float32(len(c.entries)))

	tEntry := r.TypeEntry
	// The actual Fetch must be performed in a goroutine.
	go func() {
		// If we have background refresh and currently are in "disconnected" state,
		// waiting for a response might mean we mark our results as stale for up to
		// 10 minutes (max blocking timeout) after connection is restored. To reduce
		// that window, we assume that if the fetch takes more than 31 seconds then
		// they are correctly blocking. We choose 31 seconds because yamux
		// keepalives are every 30 seconds so the RPC should fail if the packets are
		// being blackholed for more than 30 seconds.
		var connectedTimer *time.Timer
		if tEntry.Opts.Refresh && entry.Index > 0 && tEntry.Opts.QueryTimeout > 31*time.Second {
			connectedTimer = time.AfterFunc(31*time.Second, func() {
				c.entriesLock.Lock()
				defer c.entriesLock.Unlock()
				entry, ok := c.entries[key]
				if !ok || entry.RefreshLostContact.IsZero() {
					return
				}
				entry.RefreshLostContact = time.Time{}
				c.entries[key] = entry
			})
		}

		fOpts := FetchOptions{}
		if tEntry.Opts.SupportsBlocking {
			fOpts.MinIndex = entry.Index
			fOpts.Timeout = tEntry.Opts.QueryTimeout

			if fOpts.Timeout == 0 {
				fOpts.Timeout = 10 * time.Minute
			}
		}
		if entry.Valid {
			fOpts.LastResult = &FetchResult{
				Value: entry.Value,
				State: entry.State,
				Index: entry.Index,
			}
		}
		if err := entry.FetchRateLimiter.Wait(c.rateLimitContext); err != nil {
			if connectedTimer != nil {
				connectedTimer.Stop()
			}
			entry.Error = fmt.Errorf("rateLimitContext canceled: %s", err.Error())
			return
		}
		// Start building the new entry by blocking on the fetch.
		result, err := r.Fetch(fOpts)
		if connectedTimer != nil {
			connectedTimer.Stop()
		}

		// Copy the existing entry to start.
		newEntry := entry
		newEntry.Fetching = false

		// Importantly, always reset the Error. Having both Error and a Value that
		// are non-nil is allowed in the cache entry but it indicates that the Error
		// is _newer_ than the last good value. So if the err is nil then we need to
		// reset to replace any _older_ errors and avoid them bubbling up. If the
		// error is non-nil then we need to set it anyway and used to do it in the
		// code below. See https://github.com/hashicorp/consul/issues/4480.
		newEntry.Error = err

		if result.Value != nil {
			// A new value was given, so we create a brand new entry.
			if !result.NotModified {
				newEntry.Value = result.Value
			}
			newEntry.State = result.State
			newEntry.Index = result.Index
			newEntry.FetchedAt = time.Now()
			if newEntry.Index < 1 {
				// Less than one is invalid unless there was an error and in this case
				// there wasn't since a value was returned. If a badly behaved RPC
				// returns 0 when it has no data, we might get into a busy loop here. We
				// set this to minimum of 1 which is safe because no valid user data can
				// ever be written at raft index 1 due to the bootstrap process for
				// raft. This insure that any subsequent background refresh request will
				// always block, but allows the initial request to return immediately
				// even if there is no data.
				newEntry.Index = 1
			}

			// This is a valid entry with a result
			newEntry.Valid = true
		} else if result.State != nil && err == nil {
			// Also set state if it's non-nil but Value is nil. This is important in the
			// case we are returning nil due to a timeout or a transient error like rate
			// limiting that we want to mask from the user - there is no result yet but
			// we want to manage retrying internally before we return an error to user.
			// The retrying state is in State so we need to still update that in the
			// entry even if we don't have an actual result yet (e.g. hit a rate limit
			// on first request for a leaf certificate).
			newEntry.State = result.State
		}

		// Error handling
		if err == nil {
			labels := []metrics.Label{{Name: "result_not_modified", Value: strconv.FormatBool(result.NotModified)}}
			metrics.IncrCounterWithLabels([]string{"consul", "cache", "fetch_success"}, 1, labels)
			metrics.IncrCounterWithLabels([]string{"consul", "cache", tEntry.Name, "fetch_success"}, 1, labels)

			if result.Index > 0 {
				// Reset the attempts counter so we don't have any backoff
				attempt = 0
			} else {
				// Result having a zero index is an implicit error case. There was no
				// actual error but it implies the RPC found in index (nothing written
				// yet for that type) but didn't take care to return safe "1" index. We
				// don't want to actually treat it like an error by setting
				// newEntry.Error to something non-nil, but we should guard against 100%
				// CPU burn hot loops caused by that case which will never block but
				// also won't backoff either. So we treat it as a failed attempt so that
				// at least the failure backoff will save our CPU while still
				// periodically refreshing so normal service can resume when the servers
				// actually have something to return from the RPC. If we get in this
				// state it can be considered a bug in the RPC implementation (to ever
				// return a zero index) however since it can happen this is a safety net
				// for the future.
				attempt++
			}

			// If we have refresh active, this successful response means cache is now
			// "connected" and should not be stale. Reset the lost contact timer.
			if tEntry.Opts.Refresh {
				newEntry.RefreshLostContact = time.Time{}
			}
		} else {
			metrics.IncrCounter([]string{"consul", "cache", "fetch_error"}, 1)
			metrics.IncrCounter([]string{"consul", "cache", tEntry.Name, "fetch_error"}, 1)

			// Increment attempt counter
			attempt++

			// If we are refreshing and just failed, updated the lost contact time as
			// our cache will be stale until we get successfully reconnected. We only
			// set this on the first failure (if it's zero) so we can track how long
			// it's been since we had a valid connection/up-to-date view of the state.
			if tEntry.Opts.Refresh && newEntry.RefreshLostContact.IsZero() {
				newEntry.RefreshLostContact = time.Now()
			}
		}

		// Create a new waiter that will be used for the next fetch.
		newEntry.Waiter = make(chan struct{})

		// Set our entry
		c.entriesLock.Lock()

		// If this is a new entry (not in the heap yet), then setup the
		// initial expiry information and insert. If we're already in
		// the heap we do nothing since we're reusing the same entry.
		if newEntry.Expiry == nil || newEntry.Expiry.HeapIndex == -1 {
			newEntry.Expiry = &cacheEntryExpiry{Key: key}
			newEntry.Expiry.Update(tEntry.Opts.LastGetTTL)
			heap.Push(c.entriesExpiryHeap, newEntry.Expiry)
		}

		c.entries[key] = newEntry
		c.entriesLock.Unlock()

		// Trigger the old waiter
		close(entry.Waiter)

		// If refresh is enabled, run the refresh in due time. The refresh
		// below might block, but saves us from spawning another goroutine.
		if tEntry.Opts.Refresh {
			// Check if cache was stopped
			if atomic.LoadUint32(&c.stopped) == 1 {
				return
			}

			// If we're over the attempt minimum, start an exponential backoff.
			if wait := backOffWait(attempt); wait > 0 {
				time.Sleep(wait)
			}

			// If we have a timer, wait for it
			if tEntry.Opts.RefreshTimer > 0 {
				time.Sleep(tEntry.Opts.RefreshTimer)
			}

			// Trigger. The "allowNew" field is false because in the time we were
			// waiting to refresh we may have expired and got evicted. If that
			// happened, we don't want to create a new entry.
			r.Info.MustRevalidate = false
			r.Info.MinIndex = 0
			c.fetch(key, r, false, attempt, true)
		}
	}()

	return entry.Waiter
}

func backOffWait(failures uint) time.Duration {
	if failures > CacheRefreshBackoffMin {
		shift := failures - CacheRefreshBackoffMin
		waitTime := CacheRefreshMaxWait
		if shift < 31 {
			waitTime = (1 << shift) * time.Second
		}
		if waitTime > CacheRefreshMaxWait {
			waitTime = CacheRefreshMaxWait
		}
		return waitTime + lib.RandomStagger(waitTime)
	}
	return 0
}

// runExpiryLoop is a blocking function that watches the expiration
// heap and invalidates entries that have expired.
func (c *Cache) runExpiryLoop() {
	var expiryTimer *time.Timer
	for {
		// If we have a previous timer, stop it.
		if expiryTimer != nil {
			expiryTimer.Stop()
		}

		// Get the entry expiring soonest
		var entry *cacheEntryExpiry
		var expiryCh <-chan time.Time
		c.entriesLock.RLock()
		if len(c.entriesExpiryHeap.Entries) > 0 {
			entry = c.entriesExpiryHeap.Entries[0]
			expiryTimer = time.NewTimer(time.Until(entry.Expires))
			expiryCh = expiryTimer.C
		}
		c.entriesLock.RUnlock()

		select {
		case <-c.stopCh:
			return
		case <-c.entriesExpiryHeap.NotifyCh:
			// Entries changed, so the heap may have changed. Restart loop.

		case <-expiryCh:
			c.entriesLock.Lock()

			// Entry expired! Remove it.
			delete(c.entries, entry.Key)
			heap.Remove(c.entriesExpiryHeap, entry.HeapIndex)

			// This is subtle but important: if we race and simultaneously
			// evict and fetch a new value, then we set this to -1 to
			// have it treated as a new value so that the TTL is extended.
			entry.HeapIndex = -1

			// Set some metrics
			metrics.IncrCounter([]string{"consul", "cache", "evict_expired"}, 1)
			metrics.SetGauge([]string{"consul", "cache", "entries_count"}, float32(len(c.entries)))

			c.entriesLock.Unlock()
		}
	}
}

// Close stops any background work and frees all resources for the cache.
// Current Fetch requests are allowed to continue to completion and callers may
// still access the current cache values so coordination isn't needed with
// callers, however no background activity will continue. It's intended to close
// the cache at agent shutdown so no further requests should be made, however
// concurrent or in-flight ones won't break.
func (c *Cache) Close() error {
	wasStopped := atomic.SwapUint32(&c.stopped, 1)
	if wasStopped == 0 {
		// First time only, close stop chan
		close(c.stopCh)
		c.rateLimitCancel()
	}
	return nil
}

// Prepopulate puts something in the cache manually. This is useful when the
// correct initial value is know and the cache shouldn't refetch the same thing
// on startup. It is used to set the ConnectRootCA and AgentLeafCert when
// AutoEncrypt.TLS is turned on. The cache itself cannot fetch that the first
// time because it requires a special RPCType. Subsequent runs are fine though.
func (c *Cache) Prepopulate(t string, res FetchResult, dc, token, k string) error {
	key := makeEntryKey(t, dc, token, k)
	newEntry := cacheEntry{
		Valid:     true,
		Value:     res.Value,
		State:     res.State,
		Index:     res.Index,
		FetchedAt: time.Now(),
		Waiter:    make(chan struct{}),
		Expiry:    &cacheEntryExpiry{Key: key},
		FetchRateLimiter: rate.NewLimiter(
			c.options.EntryFetchRate,
			c.options.EntryFetchMaxBurst,
		),
	}
	c.entriesLock.Lock()
	c.entries[key] = newEntry
	c.entriesLock.Unlock()
	return nil
}
