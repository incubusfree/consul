package agent

import (
	"errors"
	"fmt"
	"github.com/hashicorp/consul/testutil"
	"github.com/hashicorp/serf/serf"
	"io"
	"net"
	"os"
	"strings"
	"testing"
	"time"
)

type rpcParts struct {
	dir    string
	client *RPCClient
	agent  *Agent
	rpc    *AgentRPC
}

func (r *rpcParts) Close() {
	r.client.Close()
	r.rpc.Shutdown()
	r.agent.Shutdown()
	os.RemoveAll(r.dir)
}

// testRPCClient returns an RPCClient connected to an RPC server that
// serves only this connection.
func testRPCClient(t *testing.T) *rpcParts {
	return testRPCClientWithConfig(t, nil)
}

func testRPCClientWithConfig(t *testing.T, c *Config) *rpcParts {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	lw := NewLogWriter(512)
	mult := io.MultiWriter(os.Stderr, lw)

	conf := nextConfig()
	if c != nil {
		conf = MergeConfig(conf, c)
	}

	dir, agent := makeAgentLog(t, conf, mult)
	rpc := NewAgentRPC(agent, l, mult, lw)

	rpcClient, err := NewRPCClient(l.Addr().String())
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	return &rpcParts{
		dir:    dir,
		client: rpcClient,
		agent:  agent,
		rpc:    rpc,
	}
}

func TestRPCClientForceLeave(t *testing.T) {
	p1 := testRPCClient(t)
	p2 := testRPCClient(t)
	defer p1.Close()
	defer p2.Close()

	s2Addr := fmt.Sprintf("127.0.0.1:%d", p2.agent.config.Ports.SerfLan)
	if _, err := p1.agent.JoinLAN([]string{s2Addr}); err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := p2.agent.Shutdown(); err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := p1.client.ForceLeave(p2.agent.config.NodeName); err != nil {
		t.Fatalf("err: %s", err)
	}

	m := p1.agent.LANMembers()
	if len(m) != 2 {
		t.Fatalf("should have 2 members: %#v", m)
	}

	testutil.WaitForResult(func() (bool, error) {
		m := p1.agent.LANMembers()
		success := m[1].Status == serf.StatusLeft
		return success, errors.New(m[1].Status.String())
	}, func(err error) {
		t.Fatalf("member status is %v, should be left", err)
	})
}

func TestRPCClientJoinLAN(t *testing.T) {
	p1 := testRPCClient(t)
	p2 := testRPCClient(t)
	defer p1.Close()
	defer p2.Close()

	s2Addr := fmt.Sprintf("127.0.0.1:%d", p2.agent.config.Ports.SerfLan)
	n, err := p1.client.Join([]string{s2Addr}, false)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if n != 1 {
		t.Fatalf("n != 1: %d", n)
	}
}

func TestRPCClientJoinWAN(t *testing.T) {
	p1 := testRPCClient(t)
	p2 := testRPCClient(t)
	defer p1.Close()
	defer p2.Close()

	s2Addr := fmt.Sprintf("127.0.0.1:%d", p2.agent.config.Ports.SerfWan)
	n, err := p1.client.Join([]string{s2Addr}, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if n != 1 {
		t.Fatalf("n != 1: %d", n)
	}
}

func TestRPCClientLANMembers(t *testing.T) {
	p1 := testRPCClient(t)
	p2 := testRPCClient(t)
	defer p1.Close()
	defer p2.Close()

	mem, err := p1.client.LANMembers()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if len(mem) != 1 {
		t.Fatalf("bad: %#v", mem)
	}

	s2Addr := fmt.Sprintf("127.0.0.1:%d", p2.agent.config.Ports.SerfLan)
	_, err = p1.client.Join([]string{s2Addr}, false)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	mem, err = p1.client.LANMembers()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if len(mem) != 2 {
		t.Fatalf("bad: %#v", mem)
	}
}

func TestRPCClientWANMembers(t *testing.T) {
	p1 := testRPCClient(t)
	p2 := testRPCClient(t)
	defer p1.Close()
	defer p2.Close()

	mem, err := p1.client.WANMembers()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if len(mem) != 1 {
		t.Fatalf("bad: %#v", mem)
	}

	s2Addr := fmt.Sprintf("127.0.0.1:%d", p2.agent.config.Ports.SerfWan)
	_, err = p1.client.Join([]string{s2Addr}, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	mem, err = p1.client.WANMembers()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if len(mem) != 2 {
		t.Fatalf("bad: %#v", mem)
	}
}

func TestRPCClientStats(t *testing.T) {
	p1 := testRPCClient(t)
	defer p1.Close()

	stats, err := p1.client.Stats()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if _, ok := stats["agent"]; !ok {
		t.Fatalf("bad: %#v", stats)
	}

	if _, ok := stats["consul"]; !ok {
		t.Fatalf("bad: %#v", stats)
	}
}

func TestRPCClientLeave(t *testing.T) {
	p1 := testRPCClient(t)
	defer p1.Close()

	if err := p1.client.Leave(); err != nil {
		t.Fatalf("err: %s", err)
	}

	time.Sleep(1 * time.Second)

	select {
	case <-p1.agent.ShutdownCh():
	default:
		t.Fatalf("agent should be shutdown!")
	}
}

func TestRPCClientMonitor(t *testing.T) {
	p1 := testRPCClient(t)
	defer p1.Close()

	eventCh := make(chan string, 64)
	if handle, err := p1.client.Monitor("debug", eventCh); err != nil {
		t.Fatalf("err: %s", err)
	} else {
		defer p1.client.Stop(handle)
	}

	found := false
OUTER1:
	for i := 0; ; i++ {
		select {
		case e := <-eventCh:
			if strings.Contains(e, "Accepted client") {
				found = true
				break OUTER1
			}
		default:
			if i > 100 {
				break OUTER1
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
	if !found {
		t.Fatalf("should log client accept")
	}

	// Join a bad thing to generate more events
	p1.agent.JoinLAN(nil)

	found = false
OUTER2:
	for i := 0; ; i++ {
		select {
		case e := <-eventCh:
			if strings.Contains(e, "joining") {
				found = true
				break OUTER2
			}
		default:
			if i > 100 {
				break OUTER2
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
	if !found {
		t.Fatalf("should log joining")
	}
}

func TestRPCClientListKeys(t *testing.T) {
	key1 := "tbLJg26ZJyJ9pK3qhc9jig=="
	conf := Config{EncryptKey: key1, Datacenter: "dc1"}
	p1 := testRPCClientWithConfig(t, &conf)
	defer p1.Close()

	// Key is initially installed to both wan/lan
	keys := listKeys(t, p1.client)
	if _, ok := keys["dc1"][key1]; !ok {
		t.Fatalf("bad: %#v", keys)
	}
	if _, ok := keys["wan"][key1]; !ok {
		t.Fatalf("bad: %#v", keys)
	}
}

func TestRPCClientInstallKey(t *testing.T) {
	key1 := "tbLJg26ZJyJ9pK3qhc9jig=="
	key2 := "xAEZ3uVHRMZD9GcYMZaRQw=="
	conf := Config{EncryptKey: key1}
	p1 := testRPCClientWithConfig(t, &conf)
	defer p1.Close()

	// key2 is not installed yet
	keys := listKeys(t, p1.client)
	if num, ok := keys["dc1"][key2]; ok || num != 0 {
		t.Fatalf("bad: %#v", keys)
	}
	if num, ok := keys["wan"][key2]; ok || num != 0 {
		t.Fatalf("bad: %#v", keys)
	}

	// install key2
	if _, err := p1.client.InstallKey(key2); err != nil {
		t.Fatalf("err: %s", err)
	}

	// key2 should now be installed
	keys = listKeys(t, p1.client)
	if num, ok := keys["dc1"][key2]; !ok || num != 1 {
		t.Fatalf("bad: %#v", keys)
	}
	if num, ok := keys["wan"][key2]; !ok || num != 1 {
		t.Fatalf("bad: %#v", keys)
	}
}

func TestRPCClientUseKey(t *testing.T) {
	key1 := "tbLJg26ZJyJ9pK3qhc9jig=="
	key2 := "xAEZ3uVHRMZD9GcYMZaRQw=="
	conf := Config{EncryptKey: key1}
	p1 := testRPCClientWithConfig(t, &conf)
	defer p1.Close()

	// add a second key to the ring
	if _, err := p1.client.InstallKey(key2); err != nil {
		t.Fatalf("err: %s", err)
	}

	// key2 is installed
	keys := listKeys(t, p1.client)
	if num, ok := keys["dc1"][key2]; !ok || num != 1 {
		t.Fatalf("bad: %#v", keys)
	}
	if num, ok := keys["wan"][key2]; !ok || num != 1 {
		t.Fatalf("bad: %#v", keys)
	}

	// can't remove key1 yet
	if _, err := p1.client.RemoveKey(key1); err == nil {
		t.Fatalf("expected error removing primary key")
	}

	// change primary key
	if _, err := p1.client.UseKey(key2); err != nil {
		t.Fatalf("err: %s", err)
	}

	// can remove key1 now
	if _, err := p1.client.RemoveKey(key1); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestRPCClientKeyOperation_encryptionDisabled(t *testing.T) {
	p1 := testRPCClient(t)
	defer p1.Close()

	r, err := p1.client.ListKeys()
	if err == nil {
		t.Fatalf("no error listing keys with encryption disabled")
	}

	if len(r.Messages) != 1 {
		t.Fatalf("bad: %#v", r)
	}
}

func listKeys(t *testing.T, c *RPCClient) map[string]map[string]int {
	resp, err := c.ListKeys()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	out := make(map[string]map[string]int)
	for _, k := range resp.Keys {
		respID := k.Datacenter
		if k.Pool == "WAN" {
			respID = "wan"
		}
		out[respID] = map[string]int{k.Key: k.Count}
	}
	return out
}
