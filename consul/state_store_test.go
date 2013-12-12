package consul

import (
	"sort"
	"testing"
)

func TestEnsureNode(t *testing.T) {
	store, err := NewStateStore()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer store.Close()

	if err := store.EnsureNode("foo", "127.0.0.1"); err != nil {
		t.Fatalf("err: %v")
	}

	found, addr := store.GetNode("foo")
	if !found || addr != "127.0.0.1" {
		t.Fatalf("Bad: %v %v", found, addr)
	}

	if err := store.EnsureNode("foo", "127.0.0.2"); err != nil {
		t.Fatalf("err: %v")
	}

	found, addr = store.GetNode("foo")
	if !found || addr != "127.0.0.2" {
		t.Fatalf("Bad: %v %v", found, addr)
	}
}

func TestGetNodes(t *testing.T) {
	store, err := NewStateStore()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer store.Close()

	if err := store.EnsureNode("foo", "127.0.0.1"); err != nil {
		t.Fatalf("err: %v")
	}

	if err := store.EnsureNode("bar", "127.0.0.2"); err != nil {
		t.Fatalf("err: %v")
	}

	nodes := store.Nodes()
	if len(nodes) != 4 {
		t.Fatalf("Bad: %v", nodes)
	}
	if nodes[0] != "foo" && nodes[2] != "bar" {
		t.Fatalf("Bad: %v", nodes)
	}
}

func TestEnsureService(t *testing.T) {
	store, err := NewStateStore()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer store.Close()

	if err := store.EnsureNode("foo", "127.0.0.1"); err != nil {
		t.Fatalf("err: %v")
	}

	if err := store.EnsureService("foo", "api", "", 5000); err != nil {
		t.Fatalf("err: %v")
	}

	if err := store.EnsureService("foo", "api", "", 5001); err != nil {
		t.Fatalf("err: %v")
	}

	if err := store.EnsureService("foo", "db", "master", 8000); err != nil {
		t.Fatalf("err: %v")
	}

	services := store.NodeServices("foo")

	entry, ok := services["api"]
	if !ok {
		t.Fatalf("missing api: %#v", services)
	}
	if entry.Tag != "" || entry.Port != 5001 {
		t.Fatalf("Bad entry: %#v", entry)
	}

	entry, ok = services["db"]
	if !ok {
		t.Fatalf("missing db: %#v", services)
	}
	if entry.Tag != "master" || entry.Port != 8000 {
		t.Fatalf("Bad entry: %#v", entry)
	}
}

func TestDeleteNodeService(t *testing.T) {
	store, err := NewStateStore()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer store.Close()

	if err := store.EnsureNode("foo", "127.0.0.1"); err != nil {
		t.Fatalf("err: %v")
	}

	if err := store.EnsureService("foo", "api", "", 5000); err != nil {
		t.Fatalf("err: %v")
	}

	if err := store.DeleteNodeService("foo", "api"); err != nil {
		t.Fatalf("err: %v")
	}

	services := store.NodeServices("foo")
	_, ok := services["api"]
	if ok {
		t.Fatalf("has api: %#v", services)
	}
}

func TestDeleteNode(t *testing.T) {
	store, err := NewStateStore()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer store.Close()

	if err := store.EnsureNode("foo", "127.0.0.1"); err != nil {
		t.Fatalf("err: %v")
	}

	if err := store.EnsureService("foo", "api", "", 5000); err != nil {
		t.Fatalf("err: %v")
	}

	if err := store.DeleteNode("foo"); err != nil {
		t.Fatalf("err: %v")
	}

	services := store.NodeServices("foo")
	_, ok := services["api"]
	if ok {
		t.Fatalf("has api: %#v", services)
	}

	found, _ := store.GetNode("foo")
	if found {
		t.Fatalf("found node")
	}
}

func TestGetServices(t *testing.T) {
	store, err := NewStateStore()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer store.Close()

	if err := store.EnsureNode("foo", "127.0.0.1"); err != nil {
		t.Fatalf("err: %v")
	}

	if err := store.EnsureNode("bar", "127.0.0.2"); err != nil {
		t.Fatalf("err: %v")
	}

	if err := store.EnsureService("foo", "api", "", 5000); err != nil {
		t.Fatalf("err: %v")
	}

	if err := store.EnsureService("foo", "db", "master", 8000); err != nil {
		t.Fatalf("err: %v")
	}

	if err := store.EnsureService("bar", "db", "slave", 8000); err != nil {
		t.Fatalf("err: %v")
	}

	services := store.Services()

	tags, ok := services["api"]
	if !ok {
		t.Fatalf("missing api: %#v", services)
	}
	if len(tags) != 1 || tags[0] != "" {
		t.Fatalf("Bad entry: %#v", tags)
	}

	tags, ok = services["db"]
	sort.Strings(tags)
	if !ok {
		t.Fatalf("missing db: %#v", services)
	}
	if len(tags) != 2 || tags[0] != "master" || tags[1] != "slave" {
		t.Fatalf("Bad entry: %#v", tags)
	}
}
