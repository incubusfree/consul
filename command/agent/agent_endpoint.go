package agent

import (
	"fmt"
	"github.com/hashicorp/consul/consul/structs"
	"net/http"
	"strings"
)

func (s *HTTPServer) AgentServices(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	services := s.agent.state.Services()
	return services, nil
}

func (s *HTTPServer) AgentChecks(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	checks := s.agent.state.Checks()
	return checks, nil
}

func (s *HTTPServer) AgentMembers(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	// Check if the WAN is being queried
	wan := false
	if other := req.URL.Query().Get("wan"); other != "" {
		wan = true
	}
	if wan {
		return s.agent.WANMembers(), nil
	} else {
		return s.agent.LANMembers(), nil
	}
}

func (s *HTTPServer) AgentJoin(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	// Check if the WAN is being queried
	wan := false
	if other := req.URL.Query().Get("wan"); other != "" {
		wan = true
	}

	// Get the address
	addr := strings.TrimPrefix(req.URL.Path, "/v1/agent/join/")
	if wan {
		_, err := s.agent.JoinWAN([]string{addr})
		return err, nil
	} else {
		_, err := s.agent.JoinLAN([]string{addr})
		return err, nil
	}
}

func (s *HTTPServer) AgentForceLeave(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	addr := strings.TrimPrefix(req.URL.Path, "/v1/agent/force-leave/")
	return s.agent.ForceLeave(addr), nil
}

func (s *HTTPServer) AgentRegisterCheck(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	var args CheckDefinition
	if err := decodeBody(req, &args); err != nil {
		resp.WriteHeader(400)
		resp.Write([]byte(fmt.Sprintf("Request decode failed: %v", err)))
		return nil, nil
	}

	// Verify the check has a name
	if args.Name == "" {
		resp.WriteHeader(400)
		resp.Write([]byte("Missing check name"))
		return nil, nil
	}

	// Construct the health check
	health := structs.HealthCheck{
		Node:    s.agent.config.NodeName,
		CheckID: args.ID,
		Name:    args.Name,
		Status:  structs.HealthUnknown,
		Notes:   args.Notes,
	}

	// Fixup the ID if not given
	if health.CheckID == "" && health.Name != "" {
		health.CheckID = health.Name
	}

	// Verify the check type
	chkType := &args.CheckType
	if !chkType.Valid() {
		resp.WriteHeader(400)
		resp.Write([]byte("Must provide TTL or Script and Interval!"))
		return nil, nil
	}

	// Add the check
	return s.agent.AddCheck(&health, chkType), nil
}

func (s *HTTPServer) AgentDeregisterCheck(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	checkID := strings.TrimPrefix(req.URL.Path, "/v1/agent/check/deregister/")
	return s.agent.RemoveCheck(checkID), nil
}

func (s *HTTPServer) AgentCheckPass(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	checkID := strings.TrimPrefix(req.URL.Path, "/v1/agent/check/pass/")
	note := req.URL.Query().Get("note")
	return s.agent.UpdateCheck(checkID, structs.HealthPassing, note), nil
}

func (s *HTTPServer) AgentCheckWarn(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	checkID := strings.TrimPrefix(req.URL.Path, "/v1/agent/check/warn/")
	note := req.URL.Query().Get("note")
	return s.agent.UpdateCheck(checkID, structs.HealthWarning, note), nil
}

func (s *HTTPServer) AgentCheckFail(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	checkID := strings.TrimPrefix(req.URL.Path, "/v1/agent/check/fail/")
	note := req.URL.Query().Get("note")
	return s.agent.UpdateCheck(checkID, structs.HealthCritical, note), nil
}

func (s *HTTPServer) AgentRegisterService(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	var args ServiceDefinition
	if err := decodeBody(req, &args); err != nil {
		resp.WriteHeader(400)
		resp.Write([]byte(fmt.Sprintf("Request decode failed: %v", err)))
		return nil, nil
	}

	// Verify the service has a name
	if args.Name == "" {
		resp.WriteHeader(400)
		resp.Write([]byte("Missing service name"))
		return nil, nil
	}

	// Construct the health check
	ns := structs.NodeService{
		ID:      args.ID,
		Service: args.Name,
		Tag:     args.Tag,
		Port:    args.Port,
	}

	// Fixup the ID if not given
	if ns.ID == "" && ns.Service != "" {
		ns.ID = ns.Service
	}

	// Verify the check type
	chkType := args.Check
	if chkType != nil && !chkType.Valid() {
		resp.WriteHeader(400)
		resp.Write([]byte("Must provide TTL or Script and Interval!"))
		return nil, nil
	}

	// Add the check
	return s.agent.AddService(&ns, chkType), nil
}

func (s *HTTPServer) AgentDeregisterService(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	serviceID := strings.TrimPrefix(req.URL.Path, "/v1/agent/service/deregister/")
	return s.agent.RemoveService(serviceID), nil
}
