package consul

import (
	"fmt"
	"github.com/hashicorp/consul/testutil"
	"github.com/hashicorp/consul/consul/structs"
	"github.com/hashicorp/serf/serf"
	"os"
	"testing"
	"time"
)

func TestLeader_RegisterMember(t *testing.T) {
	dir1, s1 := testServer(t)
	defer os.RemoveAll(dir1)
	defer s1.Shutdown()

	dir2, c1 := testClient(t)
	defer os.RemoveAll(dir2)
	defer c1.Shutdown()

	client := rpcClient(t, s1)
	testutil.WaitForLeader(t, client.Call)

	// Try to join
	addr := fmt.Sprintf("127.0.0.1:%d",
		s1.config.SerfLANConfig.MemberlistConfig.BindPort)
	if _, err := c1.JoinLAN([]string{addr}); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Wait for registration
	time.Sleep(10 * time.Millisecond)

	// Client should be registered
	state := s1.fsm.State()
	_, found, _ := state.GetNode(c1.config.NodeName)
	if !found {
		t.Fatalf("client not registered")
	}

	// Should have a check
	_, checks := state.NodeChecks(c1.config.NodeName)
	if len(checks) != 1 {
		t.Fatalf("client missing check")
	}
	if checks[0].CheckID != SerfCheckID {
		t.Fatalf("bad check: %v", checks[0])
	}
	if checks[0].Name != SerfCheckName {
		t.Fatalf("bad check: %v", checks[0])
	}
	if checks[0].Status != structs.HealthPassing {
		t.Fatalf("bad check: %v", checks[0])
	}

	// Server should be registered
	_, found, _ = state.GetNode(s1.config.NodeName)
	if !found {
		t.Fatalf("server not registered")
	}

	// Service should be registered
	_, services := state.NodeServices(s1.config.NodeName)
	if _, ok := services.Services["consul"]; !ok {
		t.Fatalf("consul service not registered: %v", services)
	}
}

func TestLeader_FailedMember(t *testing.T) {
	dir1, s1 := testServer(t)
	defer os.RemoveAll(dir1)
	defer s1.Shutdown()

	dir2, c1 := testClient(t)
	defer os.RemoveAll(dir2)
	defer c1.Shutdown()

	// Wait until we have a leader
	time.Sleep(100 * time.Millisecond)

	// Try to join
	addr := fmt.Sprintf("127.0.0.1:%d",
		s1.config.SerfLANConfig.MemberlistConfig.BindPort)
	if _, err := c1.JoinLAN([]string{addr}); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Fail the member
	c1.Shutdown()

	// Wait for failure detection
	time.Sleep(500 * time.Millisecond)

	// Should be registered
	state := s1.fsm.State()
	_, found, _ := state.GetNode(c1.config.NodeName)
	if !found {
		t.Fatalf("client not registered")
	}

	// Should have a check
	_, checks := state.NodeChecks(c1.config.NodeName)
	if len(checks) != 1 {
		t.Fatalf("client missing check")
	}
	if checks[0].CheckID != SerfCheckID {
		t.Fatalf("bad check: %v", checks[0])
	}
	if checks[0].Name != SerfCheckName {
		t.Fatalf("bad check: %v", checks[0])
	}
	if checks[0].Status != structs.HealthCritical {
		t.Fatalf("bad check: %v", checks[0])
	}
}

func TestLeader_LeftMember(t *testing.T) {
	dir1, s1 := testServer(t)
	defer os.RemoveAll(dir1)
	defer s1.Shutdown()

	dir2, c1 := testClient(t)
	defer os.RemoveAll(dir2)
	defer c1.Shutdown()

	// Wait until we have a leader
	time.Sleep(100 * time.Millisecond)

	// Try to join
	addr := fmt.Sprintf("127.0.0.1:%d",
		s1.config.SerfLANConfig.MemberlistConfig.BindPort)
	if _, err := c1.JoinLAN([]string{addr}); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Wait for registration
	time.Sleep(10 * time.Millisecond)

	// Should be registered
	state := s1.fsm.State()
	_, found, _ := state.GetNode(c1.config.NodeName)
	if !found {
		t.Fatalf("client not registered")
	}

	// Node should leave
	c1.Leave()
	c1.Shutdown()

	// Wait for failure detection
	time.Sleep(500 * time.Millisecond)

	// Should be deregistered
	_, found, _ = state.GetNode(c1.config.NodeName)
	if found {
		t.Fatalf("client registered")
	}
}

func TestLeader_ReapMember(t *testing.T) {
	dir1, s1 := testServer(t)
	defer os.RemoveAll(dir1)
	defer s1.Shutdown()

	dir2, c1 := testClient(t)
	defer os.RemoveAll(dir2)
	defer c1.Shutdown()

	// Wait until we have a leader
	time.Sleep(100 * time.Millisecond)

	// Try to join
	addr := fmt.Sprintf("127.0.0.1:%d",
		s1.config.SerfLANConfig.MemberlistConfig.BindPort)
	if _, err := c1.JoinLAN([]string{addr}); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Wait for registration
	time.Sleep(10 * time.Millisecond)

	// Should be registered
	state := s1.fsm.State()
	_, found, _ := state.GetNode(c1.config.NodeName)
	if !found {
		t.Fatalf("client not registered")
	}

	// Simulate a node reaping
	mems := s1.LANMembers()
	var c1mem serf.Member
	for _, m := range mems {
		if m.Name == c1.config.NodeName {
			c1mem = m
			c1mem.Status = StatusReap
			break
		}
	}
	s1.reconcileCh <- c1mem

	// Wait to reconcile
	time.Sleep(10 * time.Millisecond)

	// Should be deregistered
	_, found, _ = state.GetNode(c1.config.NodeName)
	if found {
		t.Fatalf("client registered")
	}
}

func TestLeader_Reconcile_ReapMember(t *testing.T) {
	dir1, s1 := testServer(t)
	defer os.RemoveAll(dir1)
	defer s1.Shutdown()

	// Wait until we have a leader
	time.Sleep(100 * time.Millisecond)

	// Register a non-existing member
	dead := structs.RegisterRequest{
		Datacenter: s1.config.Datacenter,
		Node:       "no-longer-around",
		Address:    "127.1.1.1",
		Check: &structs.HealthCheck{
			Node:    "no-longer-around",
			CheckID: SerfCheckID,
			Name:    SerfCheckName,
			Status:  structs.HealthCritical,
		},
	}
	var out struct{}
	if err := s1.RPC("Catalog.Register", &dead, &out); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Force a reconciliation
	if err := s1.reconcile(); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Node should be gone
	state := s1.fsm.State()
	_, found, _ := state.GetNode("no-longer-around")
	if found {
		t.Fatalf("client registered")
	}
}

func TestLeader_Reconcile(t *testing.T) {
	dir1, s1 := testServer(t)
	defer os.RemoveAll(dir1)
	defer s1.Shutdown()

	dir2, c1 := testClient(t)
	defer os.RemoveAll(dir2)
	defer c1.Shutdown()

	// Join before we have a leader, this should cause a reconcile!
	addr := fmt.Sprintf("127.0.0.1:%d",
		s1.config.SerfLANConfig.MemberlistConfig.BindPort)
	if _, err := c1.JoinLAN([]string{addr}); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Should not be registered
	state := s1.fsm.State()
	_, found, _ := state.GetNode(c1.config.NodeName)
	if found {
		t.Fatalf("client registered")
	}

	// Wait for leader
	time.Sleep(100 * time.Millisecond)

	// Should be registered
	_, found, _ = state.GetNode(c1.config.NodeName)
	if !found {
		t.Fatalf("client not registered")
	}
}

func TestLeader_LeftServer(t *testing.T) {
	dir1, s1 := testServer(t)
	defer os.RemoveAll(dir1)
	defer s1.Shutdown()

	dir2, s2 := testServerDCBootstrap(t, "dc1", false)
	defer os.RemoveAll(dir2)
	defer s2.Shutdown()

	dir3, s3 := testServerDCBootstrap(t, "dc1", false)
	defer os.RemoveAll(dir3)
	defer s3.Shutdown()
	servers := []*Server{s1, s2, s3}

	// Try to join
	addr := fmt.Sprintf("127.0.0.1:%d",
		s1.config.SerfLANConfig.MemberlistConfig.BindPort)
	if _, err := s2.JoinLAN([]string{addr}); err != nil {
		t.Fatalf("err: %v", err)
	}
	if _, err := s3.JoinLAN([]string{addr}); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Wait until we have 3 peers
	start := time.Now()
CHECK1:
	for _, s := range servers {
		peers, _ := s.raftPeers.Peers()
		if len(peers) != 3 {
			if time.Now().Sub(start) >= 2*time.Second {
				t.Fatalf("should have 3 peers")
			} else {
				time.Sleep(100 * time.Millisecond)
				goto CHECK1
			}
		}
	}

	// Kill any server
	servers[0].Shutdown()

	// Wait for failure detection
	time.Sleep(500 * time.Millisecond)

	// Force remove the non-leader (transition to left state)
	if err := servers[1].RemoveFailedNode(servers[0].config.NodeName); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Wait for intent propagation
	time.Sleep(500 * time.Millisecond)

	// Wait until we have 2 peers
	start = time.Now()
CHECK2:
	for _, s := range servers[1:] {
		peers, _ := s.raftPeers.Peers()
		if len(peers) != 2 {
			if time.Now().Sub(start) >= 2*time.Second {
				t.Fatalf("should have 2 peers")
			} else {
				time.Sleep(100 * time.Millisecond)
				goto CHECK2
			}
		}
	}
}

func TestLeader_MultiBootstrap(t *testing.T) {
	dir1, s1 := testServer(t)
	defer os.RemoveAll(dir1)
	defer s1.Shutdown()

	dir2, s2 := testServer(t)
	defer os.RemoveAll(dir2)
	defer s2.Shutdown()

	servers := []*Server{s1, s2}

	// Try to join
	addr := fmt.Sprintf("127.0.0.1:%d",
		s1.config.SerfLANConfig.MemberlistConfig.BindPort)
	if _, err := s2.JoinLAN([]string{addr}); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Wait until we have 2 peers
	start := time.Now()
CHECK1:
	for _, s := range servers {
		peers := s.serfLAN.Members()
		if len(peers) != 2 {
			if time.Now().Sub(start) >= 2*time.Second {
				t.Fatalf("should have 2 peers")
			} else {
				time.Sleep(100 * time.Millisecond)
				goto CHECK1
			}
		}
	}

	// Wait to ensure no peer is added
	time.Sleep(200 * time.Millisecond)

	// Ensure we don't have multiple raft peers
	for _, s := range servers {
		peers, _ := s.raftPeers.Peers()
		if len(peers) != 1 {
			t.Fatalf("should only have 1 raft peer!")
		}
	}
}
