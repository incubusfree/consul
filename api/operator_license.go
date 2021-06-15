package api

import (
	"fmt"
	"io/ioutil"
	"time"
)

type License struct {
	// The unique identifier of the license
	LicenseID string `json:"license_id"`

	// The customer ID associated with the license
	CustomerID string `json:"customer_id"`

	// If set, an identifier that should be used to lock the license to a
	// particular site, cluster, etc.
	InstallationID string `json:"installation_id"`

	// The time at which the license was issued
	IssueTime time.Time `json:"issue_time"`

	// The time at which the license starts being valid
	StartTime time.Time `json:"start_time"`

	// The time after which the license expires
	ExpirationTime time.Time `json:"expiration_time"`

	// The time at which the license ceases to function and can
	// no longer be used in any capacity
	TerminationTime time.Time `json:"termination_time"`

	// The product the license is valid for
	Product string `json:"product"`

	// License Specific Flags
	Flags map[string]interface{} `json:"flags"`

	// Modules is a list of the licensed enterprise modules
	Modules []string `json:"modules"`

	// List of features enabled by the license
	Features []string `json:"features"`
}

type LicenseReply struct {
	Valid    bool
	License  *License
	Warnings []string
}

func (op *Operator) LicenseGet(q *QueryOptions) (*LicenseReply, error) {
	var reply LicenseReply
	if _, err := op.c.query("/v1/operator/license", &reply, q); err != nil {
		return nil, err
	} else {
		return &reply, nil
	}
}

func (op *Operator) LicenseGetSigned(q *QueryOptions) (string, error) {
	r := op.c.newRequest("GET", "/v1/operator/license")
	r.params.Set("signed", "1")
	r.setQueryOptions(q)
	_, resp, err := requireOK(op.c.doRequest(r))
	if err != nil {
		return "", err
	}
	defer closeResponseBody(resp)

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// LicenseReset will reset the license to the builtin one if it is still valid.
// If the builtin license is invalid, the current license stays active.
//
// DEPRECATED: Consul 1.10 removes the corresponding HTTP endpoint as licenses
// are now set via agent configuration instead of through the API
func (*Operator) LicenseReset(_ *WriteOptions) (*LicenseReply, error) {
	return nil, fmt.Errorf("Consul 1.10 no longer supports API driven license management.")
}

// LicensePut will configure the Consul Enterprise license for the target datacenter
//
// DEPRECATED: Consul 1.10 removes the corresponding HTTP endpoint as licenses
// are now set via agent configuration instead of through the API
func (*Operator) LicensePut(_ string, _ *WriteOptions) (*LicenseReply, error) {
	return nil, fmt.Errorf("Consul 1.10 no longer supports API driven license management.")
}
