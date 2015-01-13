package api

import (
	"encoding/json"
	"fmt"
	"path"
	"strings"
	"sync"
	"time"
)

const (
	// DefaultSemaphoreSessionName is the Session Name we assign if none is provided
	DefaultSemaphoreSessionName = "Consul API Semaphore"

	// DefaultSemaphoreSessionTTL is the default session TTL if no Session is provided
	// when creating a new Semaphore. This is used because we do not have another
	// other check to depend upon.
	DefaultSemaphoreSessionTTL = "15s"

	// DefaultSemaphoreWaitTime is how long we block for at a time to check if semaphore
	// acquisition is possible. This affects the minimum time it takes to cancel
	// a Semaphore acquisition.
	DefaultSemaphoreWaitTime = 15 * time.Second

	// DefaultSemaphoreRetryTime is how long we wait after a failed lock acquisition
	// before attempting to do the lock again. This is so that once a lock-delay
	// is in affect, we do not hot loop retrying the acquisition.
	DefaultSemaphoreRetryTime = 5 * time.Second

	// DefaultSemaphoreKey is the key used within the prefix to
	// use for coordination between all the contenders.
	DefaultSemaphoreKey = "_lock"
)

var (
	// ErrSemaphoreHeld is returned if we attempt to double lock
	ErrSemaphoreHeld = fmt.Errorf("Semaphore already held")

	// ErrSemaphoreNotHeld is returned if we attempt to unlock a lock
	// that we do not hold.
	ErrSemaphoreNotHeld = fmt.Errorf("Semaphore not held")
)

// Semaphore is used to implement a distributed semaphore
// using the Consul KV primitives.
type Semaphore struct {
	c    *Client
	opts *SemaphoreOptions

	isHeld       bool
	sessionRenew chan struct{}
	lockSession  string
	l            sync.Mutex
}

// SemaphoreOptions is used to parameterize the Semaphore
type SemaphoreOptions struct {
	Prefix      string // Must be set and have write permissions
	Limit       int    // Must be set, and be positive
	Value       []byte // Optional, value to associate with the contender entry
	Session     string // OPtional, created if not specified
	SessionName string // Optional, defaults to DefaultLockSessionName
	SessionTTL  string // Optional, defaults to DefaultLockSessionTTL
}

// semaphoreLock is written under the DefaultSemaphoreKey and
// is used to coordinate between all the contenders.
type semaphoreLock struct {
	// Limit is the integer limit of holders. This is used to
	// verify that all the holders agree on the value.
	Limit int

	// Holders is a list of all the semaphore holders.
	// It maps the session ID to true. It is used as a set effectively.
	Holders map[string]bool
}

// SemaphorePrefix is used to created a Semaphore which will operate
// at the given KV prefix and uses the given limit for the semaphore.
// The prefix must have write privileges, and the limit must be agreed
// upon by all contenders.
func (c *Client) SemaphorePrefix(prefix string, limit int) (*Semaphore, error) {
	opts := &SemaphoreOptions{
		Prefix: prefix,
		Limit:  limit,
	}
	return c.SemaphoreOpts(opts)
}

// SemaphoreOpts is used to create a Semaphore with the given options.
// The prefix must have write privileges, and the limit must be agreed
// upon by all contenders. If a Session is not provided, one will be created.
func (c *Client) SemaphoreOpts(opts *SemaphoreOptions) (*Semaphore, error) {
	if opts.Prefix == "" {
		return nil, fmt.Errorf("missing prefix")
	}
	if opts.Limit <= 0 {
		return nil, fmt.Errorf("semaphore limit must be positive")
	}
	if opts.SessionName == "" {
		opts.SessionName = DefaultSemaphoreSessionName
	}
	if opts.SessionTTL == "" {
		opts.SessionTTL = DefaultSemaphoreSessionTTL
	} else {
		if _, err := time.ParseDuration(opts.SessionTTL); err != nil {
			return nil, fmt.Errorf("invalid SessionTTL: %v", err)
		}
	}
	s := &Semaphore{
		c:    c,
		opts: opts,
	}
	return s, nil
}

// Acquire attempts to reserve a slot in the semaphore, blocking until
// success, interrupted via the stopCh or an error is encounted.
// Providing a non-nil stopCh can be used to abort the attempt.
// On success, a channel is returned that represents our slot.
// This channel could be closed at any time due to session invalidation,
// communication errors, operator intervention, etc. It is NOT safe to
// assume that the slot is held until Release() unless the Session is specifically
// created without any associated health checks. By default Consul sessions
// prefer liveness over safety and an application must be able to handle
// the session being lost.
func (s *Semaphore) Acquire(stopCh chan struct{}) (chan struct{}, error) {
	// Hold the lock as we try to acquire
	s.l.Lock()
	defer s.l.Unlock()

	// Check if we already hold the semaphore
	if s.isHeld {
		return nil, ErrSemaphoreHeld
	}

	// Check if we need to create a session first
	s.lockSession = s.opts.Session
	if s.lockSession == "" {
		if sess, err := s.createSession(); err != nil {
			return nil, fmt.Errorf("failed to create session: %v", err)
		} else {
			s.sessionRenew = make(chan struct{})
			s.lockSession = sess
			go s.renewSession(sess, s.sessionRenew)

			// If we fail to acquire the lock, cleanup the session
			defer func() {
				if !s.isHeld {
					close(s.sessionRenew)
					s.sessionRenew = nil
				}
			}()
		}
	}

	// Create the contender entry
	kv := s.c.KV()
	made, _, err := kv.Acquire(s.contenderEntry(s.lockSession), nil)
	if err != nil || !made {
		return nil, fmt.Errorf("failed to make contender entry: %v", err)
	}

	// Setup the query options
	qOpts := &QueryOptions{
		WaitTime: DefaultSemaphoreWaitTime,
	}

WAIT:
	// Check if we should quit
	select {
	case <-stopCh:
		return nil, nil
	default:
	}

	// Read the prefix
	pairs, meta, err := kv.List(s.opts.Prefix, qOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to read prefix: %v", err)
	}

	// Decode the lock
	lockPair := s.findLock(pairs)
	lock, err := s.decodeLock(lockPair)
	if err != nil {
		return nil, err
	}

	// Verify we agree with the limit
	if lock.Limit != s.opts.Limit {
		return nil, fmt.Errorf("semaphore limit conflict (lock: %d, local: %d)",
			lock.Limit, s.opts.Limit)
	}

	// Prune the dead holders
	s.pruneDeadHolders(lock, pairs)

	// Check if the lock is held
	if len(lock.Holders) >= lock.Limit {
		qOpts.WaitIndex = meta.LastIndex
		goto WAIT
	}

	// Create a new lock with us as a holder
	lock.Holders[s.lockSession] = true
	newLock, err := s.encodeLock(lock, lockPair.ModifyIndex)
	if err != nil {
		return nil, err
	}

	// Attempt the acquisition
	didSet, _, err := kv.CAS(newLock, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to update lock: %v", err)
	}
	if !didSet {
		// Update failed, could have been a race with another contender,
		// retry the operation
		goto WAIT
	}

	// Watch to ensure we maintain ownership of the slot
	lockCh := make(chan struct{})
	go s.monitorLock(s.lockSession, lockCh)

	// Set that we own the lock
	s.isHeld = true

	// Acquired! All done
	return lockCh, nil
}

// Release is used to voluntarily give up our semaphore slot. It is
// an error to call this if the semaphore has not been acquired.
func (s *Semaphore) Release() error {
	// Hold the lock as we try to release
	s.l.Lock()
	defer s.l.Unlock()

	// Ensure the lock is actually held
	if !s.isHeld {
		return ErrSemaphoreNotHeld
	}

	// Set that we no longer own the lock
	s.isHeld = false

	// Stop the session renew
	if s.sessionRenew != nil {
		defer func() {
			close(s.sessionRenew)
			s.sessionRenew = nil
		}()
	}

	// Get and clear the lock session
	lockSession := s.lockSession
	s.lockSession = ""

	// Remove ourselves as a lock holder
	kv := s.c.KV()
	key := path.Join(s.opts.Prefix, DefaultSemaphoreKey)
READ:
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		return err
	}
	if pair == nil {
		pair = &KVPair{}
	}
	lock, err := s.decodeLock(pair)
	if err != nil {
		return err
	}

	// Create a new lock without us as a holder
	if _, ok := lock.Holders[lockSession]; ok {
		delete(lock.Holders, lockSession)
		newLock, err := s.encodeLock(lock, pair.ModifyIndex)
		if err != nil {
			return err
		}

		// Swap the locks
		didSet, _, err := kv.CAS(newLock, nil)
		if err != nil {
			return fmt.Errorf("failed to update lock: %v", err)
		}
		if !didSet {
			goto READ
		}
	}

	// Destroy the contender entry
	contenderKey := path.Join(s.opts.Prefix, lockSession)
	if _, err := kv.Delete(contenderKey, nil); err != nil {
		return err
	}
	return nil
}

// createSession is used to create a new managed session
func (s *Semaphore) createSession() (string, error) {
	session := s.c.Session()
	se := &SessionEntry{
		Name:     s.opts.SessionName,
		TTL:      s.opts.SessionTTL,
		Behavior: SessionBehaviorDelete,
	}
	id, _, err := session.Create(se, nil)
	if err != nil {
		return "", err
	}
	return id, nil
}

// renewSession is a long running routine that maintians a session
// by doing a periodic Session renewal.
func (s *Semaphore) renewSession(id string, doneCh chan struct{}) {
	session := s.c.Session()
	ttl, _ := time.ParseDuration(s.opts.SessionTTL)
	for {
		select {
		case <-time.After(ttl / 2):
			entry, _, err := session.Renew(id, nil)
			if err != nil || entry == nil {
				return
			}

			// Handle the server updating the TTL
			ttl, _ = time.ParseDuration(entry.TTL)

		case <-doneCh:
			// Attempt a session destroy
			session.Destroy(id, nil)
			return
		}
	}
}

// contenderEntry returns a formatted KVPair for the contender
func (s *Semaphore) contenderEntry(session string) *KVPair {
	return &KVPair{
		Key:     path.Join(s.opts.Prefix, session),
		Value:   s.opts.Value,
		Session: session,
	}
}

// findLock is used to find the KV Pair which is used for coordination
func (s *Semaphore) findLock(pairs KVPairs) *KVPair {
	key := path.Join(s.opts.Prefix, DefaultSemaphoreKey)
	for _, pair := range pairs {
		if pair.Key == key {
			return pair
		}
	}
	return &KVPair{}
}

// decodeLock is used to decode a semaphoreLock from an
// entry in Consul
func (s *Semaphore) decodeLock(pair *KVPair) (*semaphoreLock, error) {
	// Handle if there is no lock
	if pair == nil || pair.Value == nil {
		return &semaphoreLock{Limit: s.opts.Limit}, nil
	}

	l := &semaphoreLock{}
	if err := json.Unmarshal(pair.Value, l); err != nil {
		return nil, fmt.Errorf("lock decoding failed: %v", err)
	}
	return l, nil
}

// encodeLock is used to encode a semaphoreLock into a KVPair
// that can be PUT
func (s *Semaphore) encodeLock(l *semaphoreLock, oldIndex uint64) (*KVPair, error) {
	enc, err := json.Marshal(l)
	if err != nil {
		return nil, fmt.Errorf("lock encoding failed: %v", err)
	}
	pair := &KVPair{
		Key:         path.Join(s.opts.Prefix, DefaultSemaphoreKey),
		Value:       enc,
		ModifyIndex: oldIndex,
	}
	return pair, nil
}

// pruneDeadHolders is used to remove all the dead lock holders
func (s *Semaphore) pruneDeadHolders(lock *semaphoreLock, pairs KVPairs) {
	// Gather all the live holders
	alive := make(map[string]struct{}, len(pairs))
	for _, pair := range pairs {
		session := strings.TrimPrefix(pair.Key, s.opts.Prefix)
		alive[session] = struct{}{}
	}

	// Remove any holders that are dead
	for holder := range lock.Holders {
		if _, ok := alive[holder]; !ok {
			delete(lock.Holders, holder)
		}
	}
}

// monitorLock is a long running routine to monitor a semaphore ownership
// It closes the stopCh if we lose our slot.
func (s *Semaphore) monitorLock(session string, stopCh chan struct{}) {
	defer close(stopCh)
	kv := s.c.KV()
	opts := &QueryOptions{RequireConsistent: true}
	key := path.Join(s.opts.Prefix, DefaultSemaphoreKey)
WAIT:
	pair, meta, err := kv.Get(key, opts)
	if err != nil {
		return
	}
	lock, err := s.decodeLock(pair)
	if err != nil {
		return
	}
	if _, ok := lock.Holders[session]; ok {
		opts.WaitIndex = meta.LastIndex
		goto WAIT
	}
}
