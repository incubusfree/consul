package api

import (
	"encoding/json"
	"time"
)

type ServiceRouterConfigEntry struct {
	Kind      string
	Name      string
	Namespace string `json:",omitempty"`

	Routes []ServiceRoute `json:",omitempty"`

	CreateIndex uint64
	ModifyIndex uint64
}

func (e *ServiceRouterConfigEntry) GetKind() string        { return e.Kind }
func (e *ServiceRouterConfigEntry) GetName() string        { return e.Name }
func (e *ServiceRouterConfigEntry) GetCreateIndex() uint64 { return e.CreateIndex }
func (e *ServiceRouterConfigEntry) GetModifyIndex() uint64 { return e.ModifyIndex }

type ServiceRoute struct {
	Match       *ServiceRouteMatch       `json:",omitempty"`
	Destination *ServiceRouteDestination `json:",omitempty"`
}

type ServiceRouteMatch struct {
	HTTP *ServiceRouteHTTPMatch `json:",omitempty"`
}

type ServiceRouteHTTPMatch struct {
	PathExact  string `json:",omitempty" alias:"path_exact"`
	PathPrefix string `json:",omitempty" alias:"path_prefix"`
	PathRegex  string `json:",omitempty" alias:"path_regex"`

	Header     []ServiceRouteHTTPMatchHeader     `json:",omitempty"`
	QueryParam []ServiceRouteHTTPMatchQueryParam `json:",omitempty" alias:"query_param"`
	Methods    []string                          `json:",omitempty"`
}

type ServiceRouteHTTPMatchHeader struct {
	Name    string
	Present bool   `json:",omitempty"`
	Exact   string `json:",omitempty"`
	Prefix  string `json:",omitempty"`
	Suffix  string `json:",omitempty"`
	Regex   string `json:",omitempty"`
	Invert  bool   `json:",omitempty"`
}

type ServiceRouteHTTPMatchQueryParam struct {
	Name    string
	Present bool   `json:",omitempty"`
	Exact   string `json:",omitempty"`
	Regex   string `json:",omitempty"`
}

type ServiceRouteDestination struct {
	Service               string        `json:",omitempty"`
	ServiceSubset         string        `json:",omitempty" alias:"service_subset"`
	Namespace             string        `json:",omitempty"`
	PrefixRewrite         string        `json:",omitempty" alias:"prefix_rewrite"`
	RequestTimeout        time.Duration `json:",omitempty" alias:"request_timeout"`
	NumRetries            uint32        `json:",omitempty" alias:"num_retries"`
	RetryOnConnectFailure bool          `json:",omitempty" alias:"retry_on_connect_failure"`
	RetryOnStatusCodes    []uint32      `json:",omitempty" alias:"retry_on_status_codes"`
}

func (e *ServiceRouteDestination) MarshalJSON() ([]byte, error) {
	type Alias ServiceRouteDestination
	exported := &struct {
		RequestTimeout string `json:",omitempty"`
		*Alias
	}{
		RequestTimeout: e.RequestTimeout.String(),
		Alias:          (*Alias)(e),
	}
	if e.RequestTimeout == 0 {
		exported.RequestTimeout = ""
	}

	return json.Marshal(exported)
}

func (e *ServiceRouteDestination) UnmarshalJSON(data []byte) error {
	type Alias ServiceRouteDestination
	aux := &struct {
		RequestTimeout string
		*Alias
	}{
		Alias: (*Alias)(e),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	var err error
	if aux.RequestTimeout != "" {
		if e.RequestTimeout, err = time.ParseDuration(aux.RequestTimeout); err != nil {
			return err
		}
	}
	return nil
}

type ServiceSplitterConfigEntry struct {
	Kind      string
	Name      string
	Namespace string `json:",omitempty"`

	Splits []ServiceSplit `json:",omitempty"`

	CreateIndex uint64
	ModifyIndex uint64
}

func (e *ServiceSplitterConfigEntry) GetKind() string        { return e.Kind }
func (e *ServiceSplitterConfigEntry) GetName() string        { return e.Name }
func (e *ServiceSplitterConfigEntry) GetCreateIndex() uint64 { return e.CreateIndex }
func (e *ServiceSplitterConfigEntry) GetModifyIndex() uint64 { return e.ModifyIndex }

type ServiceSplit struct {
	Weight        float32
	Service       string `json:",omitempty"`
	ServiceSubset string `json:",omitempty" alias:"service_subset"`
	Namespace     string `json:",omitempty"`
}

type ServiceResolverConfigEntry struct {
	Kind      string
	Name      string
	Namespace string `json:",omitempty"`

	DefaultSubset  string                             `json:",omitempty" alias:"default_subset"`
	Subsets        map[string]ServiceResolverSubset   `json:",omitempty"`
	Redirect       *ServiceResolverRedirect           `json:",omitempty"`
	Failover       map[string]ServiceResolverFailover `json:",omitempty"`
	ConnectTimeout time.Duration                      `json:",omitempty" alias:"connect_timeout"`

	// LoadBalancer determines the load balancing policy and configuration for services
	// issuing requests to this upstream service.
	LoadBalancer *LoadBalancer `json:",omitempty" alias:"load_balancer"`

	CreateIndex uint64
	ModifyIndex uint64
}

func (e *ServiceResolverConfigEntry) MarshalJSON() ([]byte, error) {
	type Alias ServiceResolverConfigEntry
	exported := &struct {
		ConnectTimeout string `json:",omitempty"`
		*Alias
	}{
		ConnectTimeout: e.ConnectTimeout.String(),
		Alias:          (*Alias)(e),
	}
	if e.ConnectTimeout == 0 {
		exported.ConnectTimeout = ""
	}

	return json.Marshal(exported)
}

func (e *ServiceResolverConfigEntry) UnmarshalJSON(data []byte) error {
	type Alias ServiceResolverConfigEntry
	aux := &struct {
		ConnectTimeout string
		*Alias
	}{
		Alias: (*Alias)(e),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	var err error
	if aux.ConnectTimeout != "" {
		if e.ConnectTimeout, err = time.ParseDuration(aux.ConnectTimeout); err != nil {
			return err
		}
	}
	return nil
}

func (e *ServiceResolverConfigEntry) GetKind() string        { return e.Kind }
func (e *ServiceResolverConfigEntry) GetName() string        { return e.Name }
func (e *ServiceResolverConfigEntry) GetCreateIndex() uint64 { return e.CreateIndex }
func (e *ServiceResolverConfigEntry) GetModifyIndex() uint64 { return e.ModifyIndex }

type ServiceResolverSubset struct {
	Filter      string `json:",omitempty"`
	OnlyPassing bool   `json:",omitempty" alias:"only_passing"`
}

type ServiceResolverRedirect struct {
	Service       string `json:",omitempty"`
	ServiceSubset string `json:",omitempty" alias:"service_subset"`
	Namespace     string `json:",omitempty"`
	Datacenter    string `json:",omitempty"`
}

type ServiceResolverFailover struct {
	Service       string   `json:",omitempty"`
	ServiceSubset string   `json:",omitempty" alias:"service_subset"`
	Namespace     string   `json:",omitempty"`
	Datacenters   []string `json:",omitempty"`
}

// LoadBalancer determines the load balancing policy and configuration for services
// issuing requests to this upstream service.
type LoadBalancer struct {
	// EnvoyConfig contains Envoy-specific load balancing configuration for this upstream
	EnvoyConfig *EnvoyLBConfig `json:",omitempty" alias:"envoy_config"`

	// OpaqueConfig contains load balancing configuration opaque to Consul for 3rd party proxies
	OpaqueConfig string `json:",omitempty" alias:"opaque_config"`
}

type EnvoyLBConfig struct {
	// Policy is the load balancing policy used to select a host
	Policy string `json:",omitempty"`

	// RingHashConfig contains configuration for the "ring_hash" policy type
	RingHashConfig *RingHashConfig `json:",omitempty" alias:"ring_hash_config"`

	// LeastRequestConfig contains configuration for the "least_request" policy type
	LeastRequestConfig *LeastRequestConfig `json:",omitempty" alias:"least_request_config"`

	// HashPolicies is a list of hash policies to use for hashing load balancing algorithms.
	// Hash policies are evaluated individually and combined such that identical lists
	// result in the same hash.
	// If no hash policies are present, or none are successfully evaluated,
	// then a random backend host will be selected.
	HashPolicies []HashPolicy `json:",omitempty" alias:"hash_policies"`
}

// RingHashConfig contains configuration for the "ring_hash" policy type
type RingHashConfig struct {
	// MinimumRingSize determines the minimum number of entries in the hash ring
	MinimumRingSize uint64 `json:",omitempty" alias:"minimum_ring_size"`

	// MaximumRingSize determines the maximum number of entries in the hash ring
	MaximumRingSize uint64 `json:",omitempty" alias:"maximum_ring_size"`
}

// LeastRequestConfig contains configuration for the "least_request" policy type
type LeastRequestConfig struct {
	// ChoiceCount determines the number of random healthy hosts from which to select the one with the least requests.
	ChoiceCount uint32 `json:",omitempty" alias:"choice_count"`
}

// HashPolicy defines which attributes will be hashed by hash-based LB algorithms
type HashPolicy struct {
	// Field is the attribute type to hash on.
	// Must be one of "header","cookie", or "query_parameter".
	// Cannot be specified along with SourceIP.
	Field string `json:",omitempty"`

	// FieldValue is the value to hash.
	// ie. header name, cookie name, URL query parameter name
	// Cannot be specified along with SourceIP.
	FieldValue string `json:",omitempty" alias:"field_value"`

	// CookieConfig contains configuration for the "cookie" hash policy type.
	CookieConfig *CookieConfig `json:",omitempty" alias:"cookie_config"`

	// SourceIP determines whether the hash should be of the source IP rather than of a field and field value.
	// Cannot be specified along with Field or FieldValue.
	SourceIP bool `json:",omitempty" alias:"source_ip"`

	// Terminal will short circuit the computation of the hash when multiple hash policies are present.
	// If a hash is computed when a Terminal policy is evaluated,
	// then that hash will be used and subsequent hash policies will be ignored.
	Terminal bool `json:",omitempty"`
}

// CookieConfig contains configuration for the "cookie" hash policy type.
// This is specified to have Envoy generate a cookie for a client on its first request.
type CookieConfig struct {
	// TTL for generated cookies
	TTL time.Duration `json:",omitempty"`

	// The path to set for the cookie
	Path string `json:",omitempty"`
}
