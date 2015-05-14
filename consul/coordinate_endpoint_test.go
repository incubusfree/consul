package consul

import (
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/consul/consul/structs"
	"github.com/hashicorp/consul/testutil"
	"github.com/hashicorp/serf/coordinate"
)

// getRandomCoordinate generates a random coordinate.
func getRandomCoordinate() *coordinate.Coordinate {
	config := coordinate.DefaultConfig()
	// Randomly apply updates between n clients
	n := 5
	clients := make([]*coordinate.Client, n)
	for i := 0; i < n; i++ {
		clients[i], _ = coordinate.NewClient(config)
	}

	for i := 0; i < n*100; i++ {
		k1 := rand.Intn(n)
		k2 := rand.Intn(n)
		if k1 == k2 {
			continue
		}
		clients[k1].Update(clients[k2].GetCoordinate(), time.Duration(rand.Int63())*time.Microsecond)
	}
	return clients[rand.Intn(n)].GetCoordinate()
}

func coordinatesEqual(a, b *coordinate.Coordinate) bool {
	return reflect.DeepEqual(a, b)
}

func TestCoordinate_Update(t *testing.T) {
	name := fmt.Sprintf("Node %d", getPort())
	dir1, config1 := testServerConfig(t, name)
	config1.CoordinateUpdatePeriod = 1000 * time.Millisecond
	s1, err := NewServer(config1)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir1)
	defer s1.Shutdown()

	client := rpcClient(t, s1)
	defer client.Close()

	testutil.WaitForLeader(t, client.Call, "dc1")

	arg1 := structs.CoordinateUpdateRequest{
		Datacenter: "dc1",
		Node:       "node1",
		Op:         structs.CoordinateSet,
		Coord:      getRandomCoordinate(),
	}

	arg2 := structs.CoordinateUpdateRequest{
		Datacenter: "dc1",
		Node:       "node2",
		Op:         structs.CoordinateSet,
		Coord:      getRandomCoordinate(),
	}

	var out struct{}
	if err := client.Call("Coordinate.Update", &arg1, &out); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Verify
	state := s1.fsm.State()
	_, d, err := state.CoordinateGet("node1")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if d != nil {
		t.Fatalf("should be nil because the update should be batched")
	}

	// Wait a while and send another update; this time the updates should be sent
	time.Sleep(2 * s1.config.CoordinateUpdatePeriod)
	if err := client.Call("Coordinate.Update", &arg2, &out); err != nil {
		t.Fatalf("err: %v", err)
	}
	// Yield the current goroutine to allow the goroutine that sends the updates to run
	time.Sleep(100 * time.Millisecond)

	_, d, err = state.CoordinateGet("node1")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if d == nil {
		t.Fatalf("should return a coordinate but it's nil")
	}
	if !coordinatesEqual(d.Coord, arg1.Coord) {
		t.Fatalf("should be equal\n%v\n%v", d.Coord, arg1.Coord)
	}

	_, d, err = state.CoordinateGet("node2")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if d == nil {
		t.Fatalf("should return a coordinate but it's nil")
	}
	if !coordinatesEqual(d.Coord, arg2.Coord) {
		t.Fatalf("should be equal\n%v\n%v", d.Coord, arg2.Coord)
	}
}

func TestCoordinate_GetLAN(t *testing.T) {
	dir1, s1 := testServer(t)
	defer os.RemoveAll(dir1)
	defer s1.Shutdown()

	client := rpcClient(t, s1)
	defer client.Close()
	testutil.WaitForLeader(t, client.Call, "dc1")

	arg := structs.CoordinateUpdateRequest{
		Datacenter: "dc1",
		Node:       "node1",
		Op:         structs.CoordinateSet,
		Coord:      getRandomCoordinate(),
	}

	var out struct{}
	if err := client.Call("Coordinate.Update", &arg, &out); err != nil {
		t.Fatalf("err: %v", err)
	}
	// Yield the current goroutine to allow the goroutine that sends the updates to run
	time.Sleep(100 * time.Millisecond)

	// Get via RPC
	out2 := structs.IndexedCoordinate{}
	arg2 := structs.NodeSpecificRequest{
		Datacenter: "dc1",
		Node:       "node1",
	}
	if err := client.Call("Coordinate.GetLAN", &arg2, &out2); err != nil {
		t.Fatalf("err: %v", err)
	}
	if !coordinatesEqual(out2.Coord, arg.Coord) {
		t.Fatalf("should be equal\n%v\n%v", out2.Coord, arg.Coord)
	}

	// Now let's override the original coordinate; Coordinate.Get should return
	// the latest coordinate
	arg.Coord = getRandomCoordinate()
	if err := client.Call("Coordinate.Update", &arg, &out); err != nil {
		t.Fatalf("err: %v", err)
	}
	// Yield the current goroutine to allow the goroutine that sends the updates to run
	time.Sleep(100 * time.Millisecond)

	if err := client.Call("Coordinate.GetLAN", &arg2, &out2); err != nil {
		t.Fatalf("err: %v", err)
	}
	if !coordinatesEqual(out2.Coord, arg.Coord) {
		t.Fatalf("should be equal\n%v\n%v", out2.Coord, arg.Coord)
	}
}

func TestCoordinate_GetWAN(t *testing.T) {
	// Create 1 server in dc1, 2 servers in dc2
	dir1, s1 := testServerDC(t, "dc1")
	defer os.RemoveAll(dir1)
	defer s1.Shutdown()

	dir2, s2 := testServerDC(t, "dc2")
	defer os.RemoveAll(dir2)
	defer s2.Shutdown()

	dir3, s3 := testServerDC(t, "dc2")
	defer os.RemoveAll(dir3)
	defer s3.Shutdown()

	client := rpcClient(t, s1)
	defer client.Close()
	testutil.WaitForLeader(t, client.Call, "dc1")

	// Try to join
	addr := fmt.Sprintf("127.0.0.1:%d",
		s1.config.SerfWANConfig.MemberlistConfig.BindPort)
	if _, err := s2.JoinWAN([]string{addr}); err != nil {
		t.Fatalf("err: %v", err)
	}
	if _, err := s3.JoinWAN([]string{addr}); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Check the members
	testutil.WaitForResult(func() (bool, error) {
		return len(s1.WANMembers()) == 3, nil
	}, func(err error) {
		t.Fatalf("bad len")
	})

	// Wait for coordinates to be exchanged
	time.Sleep(s1.config.SerfWANConfig.MemberlistConfig.ProbeInterval * 2)

	var coords []*coordinate.Coordinate
	arg := structs.DCSpecificRequest{
		Datacenter: "dc1",
	}
	if err := client.Call("Coordinate.GetWAN", &arg, &coords); err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(coords) != 1 {
		t.Fatalf("there is 1 server in dc1")
	}

	arg = structs.DCSpecificRequest{
		Datacenter: "dc2",
	}
	if err := client.Call("Coordinate.GetWAN", &arg, &coords); err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(coords) != 2 {
		t.Fatalf("there are 2 servers in dc2")
	}
}
