package agent

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/hashicorp/consul/consul/structs"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/hashicorp/raft"
)

// OperatorRaftConfiguration is used to inspect the current Raft configuration.
// This supports the stale query mode in case the cluster doesn't have a leader.
func (s *HTTPServer) OperatorRaftConfiguration(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	if req.Method != "GET" {
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return nil, nil
	}

	var args structs.DCSpecificRequest
	if done := s.parse(resp, req, &args.Datacenter, &args.QueryOptions); done {
		return nil, nil
	}

	var reply structs.RaftConfigurationResponse
	if err := s.agent.RPC("Operator.RaftGetConfiguration", &args, &reply); err != nil {
		return nil, err
	}

	return reply, nil
}

// OperatorRaftPeer supports actions on Raft peers. Currently we only support
// removing peers by address.
func (s *HTTPServer) OperatorRaftPeer(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	if req.Method != "DELETE" {
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return nil, nil
	}

	var args structs.RaftPeerByAddressRequest
	s.parseDC(req, &args.Datacenter)
	s.parseToken(req, &args.Token)

	params := req.URL.Query()
	if _, ok := params["address"]; ok {
		args.Address = raft.ServerAddress(params.Get("address"))
	} else {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("Must specify ?address with IP:port of peer to remove"))
		return nil, nil
	}

	var reply struct{}
	if err := s.agent.RPC("Operator.RaftRemovePeerByAddress", &args, &reply); err != nil {
		return nil, err
	}
	return nil, nil
}

type keyringArgs struct {
	Key         string
	Token       string
	RelayFactor uint8
}

// OperatorKeyringEndpoint handles keyring operations (install, list, use, remove)
func (s *HTTPServer) OperatorKeyringEndpoint(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	var args keyringArgs
	if req.Method == "POST" || req.Method == "PUT" || req.Method == "DELETE" {
		if err := decodeBody(req, &args, nil); err != nil {
			resp.WriteHeader(400)
			resp.Write([]byte(fmt.Sprintf("Request decode failed: %v", err)))
			return nil, nil
		}
	}
	s.parseToken(req, &args.Token)

	// Parse relay factor
	if relayFactor := req.URL.Query().Get("relay-factor"); relayFactor != "" {
		n, err := strconv.Atoi(relayFactor)
		if err != nil {
			resp.WriteHeader(400)
			resp.Write([]byte(fmt.Sprintf("Error parsing relay factor: %v", err)))
			return nil, nil
		}

		args.RelayFactor, err = ParseRelayFactor(n)
		if err != nil {
			resp.WriteHeader(400)
			resp.Write([]byte(fmt.Sprintf("Invalid relay factor: %v", err)))
			return nil, nil
		}
	}

	// Switch on the method
	switch req.Method {
	case "GET":
		return s.KeyringList(resp, req, &args)
	case "POST":
		return s.KeyringInstall(resp, req, &args)
	case "PUT":
		return s.KeyringUse(resp, req, &args)
	case "DELETE":
		return s.KeyringRemove(resp, req, &args)
	default:
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return nil, nil
	}
}

// KeyringInstall is used to install a new gossip encryption key into the cluster
func (s *HTTPServer) KeyringInstall(resp http.ResponseWriter, req *http.Request, args *keyringArgs) (interface{}, error) {
	responses, err := s.agent.InstallKey(args.Key, args.Token, args.RelayFactor)
	if err != nil {
		return nil, err
	}

	return nil, keyringErrorsOrNil(responses.Responses)
}

// KeyringList is used to list the keys installed in the cluster
func (s *HTTPServer) KeyringList(resp http.ResponseWriter, req *http.Request, args *keyringArgs) (interface{}, error) {
	responses, err := s.agent.ListKeys(args.Token, args.RelayFactor)
	if err != nil {
		return nil, err
	}

	return responses.Responses, keyringErrorsOrNil(responses.Responses)
}

// KeyringRemove is used to list the keys installed in the cluster
func (s *HTTPServer) KeyringRemove(resp http.ResponseWriter, req *http.Request, args *keyringArgs) (interface{}, error) {
	responses, err := s.agent.RemoveKey(args.Key, args.Token, args.RelayFactor)
	if err != nil {
		return nil, err
	}

	return nil, keyringErrorsOrNil(responses.Responses)
}

// KeyringUse is used to change the primary gossip encryption key
func (s *HTTPServer) KeyringUse(resp http.ResponseWriter, req *http.Request, args *keyringArgs) (interface{}, error) {
	responses, err := s.agent.UseKey(args.Key, args.Token, args.RelayFactor)
	if err != nil {
		return nil, err
	}

	return nil, keyringErrorsOrNil(responses.Responses)
}

func keyringErrorsOrNil(responses []*structs.KeyringResponse) error {
	var errs error
	for _, response := range responses {
		if response.Error != "" {
			pool := response.Datacenter + " (LAN)"
			if response.WAN {
				pool = "WAN"
			}
			errs = multierror.Append(errs, fmt.Errorf("%s error: %s", pool, response.Error))
			for key, message := range response.Messages {
				errs = multierror.Append(errs, fmt.Errorf("%s: %s", key, message))
			}
		}
	}
	return errs
}

// OperatorAutopilotConfiguration is used to inspect the current Autopilot configuration.
// This supports the stale query mode in case the cluster doesn't have a leader.
func (s *HTTPServer) OperatorAutopilotConfiguration(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	// Switch on the method
	switch req.Method {
	case "GET":
		var args structs.DCSpecificRequest
		if done := s.parse(resp, req, &args.Datacenter, &args.QueryOptions); done {
			return nil, nil
		}

		var reply structs.AutopilotConfig
		if err := s.agent.RPC("Operator.AutopilotGetConfiguration", &args, &reply); err != nil {
			return nil, err
		}

		return reply, nil
	case "PUT":
		var args structs.AutopilotSetConfigRequest
		s.parseDC(req, &args.Datacenter)
		s.parseToken(req, &args.Token)

		// Check for cas value
		params := req.URL.Query()
		if _, ok := params["cas"]; ok {
			casVal, err := strconv.ParseUint(params.Get("cas"), 10, 64)
			if err != nil {
				resp.WriteHeader(400)
				resp.Write([]byte(fmt.Sprintf("Error parsing cas value: %v", err)))
				return nil, nil
			}
			args.Config.ModifyIndex = casVal
			args.CAS = true
		}

		if err := decodeBody(req, &args.Config, nil); err != nil {
			resp.WriteHeader(400)
			resp.Write([]byte(fmt.Sprintf("Error parsing autopilot config: %v", err)))
			return nil, nil
		}

		var reply bool
		if err := s.agent.RPC("Operator.AutopilotSetConfiguration", &args, &reply); err != nil {
			return nil, err
		}

		// Only use the out value if this was a CAS
		if !args.CAS {
			return true, nil
		} else {
			return reply, nil
		}
	default:
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return nil, nil
	}
}

// OperatorServerHealth is used to get the health of the servers in the local DC
func (s *HTTPServer) OperatorServerHealth(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	if req.Method != "GET" {
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return nil, nil
	}

	var args structs.DCSpecificRequest
	if done := s.parse(resp, req, &args.Datacenter, &args.QueryOptions); done {
		return nil, nil
	}

	var reply structs.OperatorHealthReply
	if err := s.agent.RPC("Operator.ServerHealth", &args, &reply); err != nil {
		return nil, err
	}

	// Reply with status 429 if something is unhealthy
	if !reply.Healthy {
		resp.WriteHeader(http.StatusTooManyRequests)
	}

	return reply, nil
}
