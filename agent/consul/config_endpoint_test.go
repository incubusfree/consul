package consul

import (
	"os"
	"testing"

	"github.com/hashicorp/consul/acl"
	"github.com/hashicorp/consul/agent/structs"
	"github.com/hashicorp/consul/testrpc"
	msgpackrpc "github.com/hashicorp/net-rpc-msgpackrpc"
	"github.com/stretchr/testify/require"
)

func TestConfigEntry_Apply(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	dir1, s1 := testServer(t)
	defer os.RemoveAll(dir1)
	defer s1.Shutdown()
	codec := rpcClient(t, s1)
	defer codec.Close()

	args := structs.ConfigEntryRequest{
		Datacenter: "dc1",
		Entry: &structs.ServiceConfigEntry{
			Name: "foo",
		},
	}
	var out struct{}
	require.NoError(msgpackrpc.CallWithCodec(codec, "ConfigEntry.Apply", &args, &out))

	state := s1.fsm.State()
	_, entry, err := state.ConfigEntry(nil, structs.ServiceDefaults, "foo")
	require.NoError(err)

	serviceConf, ok := entry.(*structs.ServiceConfigEntry)
	require.True(ok)
	require.Equal("foo", serviceConf.Name)
	require.Equal(structs.ServiceDefaults, serviceConf.Kind)
}

func TestConfigEntry_Apply_ACLDeny(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	dir1, s1 := testServerWithConfig(t, func(c *Config) {
		c.ACLDatacenter = "dc1"
		c.ACLsEnabled = true
		c.ACLMasterToken = "root"
		c.ACLDefaultPolicy = "deny"
	})
	defer os.RemoveAll(dir1)
	defer s1.Shutdown()
	testrpc.WaitForTestAgent(t, s1.RPC, "dc1")
	codec := rpcClient(t, s1)
	defer codec.Close()

	// Create the ACL.
	arg := structs.ACLRequest{
		Datacenter: "dc1",
		Op:         structs.ACLSet,
		ACL: structs.ACL{
			Name: "User token",
			Type: structs.ACLTokenTypeClient,
			Rules: `
service "foo" {
	policy = "write"
}
operator = "write"
`,
		},
		WriteRequest: structs.WriteRequest{Token: "root"},
	}
	var id string
	if err := msgpackrpc.CallWithCodec(codec, "ACL.Apply", &arg, &id); err != nil {
		t.Fatalf("err: %v", err)
	}

	// This should fail since we don't have write perms for the "db" service.
	args := structs.ConfigEntryRequest{
		Datacenter: "dc1",
		Entry: &structs.ServiceConfigEntry{
			Name: "db",
		},
		WriteRequest: structs.WriteRequest{Token: id},
	}
	var out struct{}
	err := msgpackrpc.CallWithCodec(codec, "ConfigEntry.Apply", &args, &out)
	if !acl.IsErrPermissionDenied(err) {
		t.Fatalf("err: %v", err)
	}

	// The "foo" service should work.
	args.Entry = &structs.ServiceConfigEntry{
		Name: "foo",
	}
	err = msgpackrpc.CallWithCodec(codec, "ConfigEntry.Apply", &args, &out)
	require.NoError(err)

	state := s1.fsm.State()
	_, entry, err := state.ConfigEntry(nil, structs.ServiceDefaults, "foo")
	require.NoError(err)

	serviceConf, ok := entry.(*structs.ServiceConfigEntry)
	require.True(ok)
	require.Equal("foo", serviceConf.Name)
	require.Equal(structs.ServiceDefaults, serviceConf.Kind)

	// Try to update the global proxy args with the anonymous token - this should fail.
	proxyArgs := structs.ConfigEntryRequest{
		Datacenter: "dc1",
		Entry: &structs.ProxyConfigEntry{
			Config: map[string]interface{}{
				"foo": 1,
			},
		},
	}
	err = msgpackrpc.CallWithCodec(codec, "ConfigEntry.Apply", &proxyArgs, &out)
	if !acl.IsErrPermissionDenied(err) {
		t.Fatalf("err: %v", err)
	}

	// Now with the privileged token.
	proxyArgs.WriteRequest.Token = id
	err = msgpackrpc.CallWithCodec(codec, "ConfigEntry.Apply", &proxyArgs, &out)
	require.NoError(err)
}

func TestConfigEntry_Get(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	dir1, s1 := testServer(t)
	defer os.RemoveAll(dir1)
	defer s1.Shutdown()
	codec := rpcClient(t, s1)
	defer codec.Close()

	// Create a dummy service in the state store to look up.
	entry := &structs.ServiceConfigEntry{
		Kind: structs.ServiceDefaults,
		Name: "foo",
	}
	state := s1.fsm.State()
	require.NoError(state.EnsureConfigEntry(1, entry))

	args := structs.ConfigEntryQuery{
		Kind:       structs.ServiceDefaults,
		Name:       "foo",
		Datacenter: s1.config.Datacenter,
	}
	var out structs.IndexedConfigEntries
	require.NoError(msgpackrpc.CallWithCodec(codec, "ConfigEntry.Get", &args, &out))

	serviceConf, ok := out.Entries[0].(*structs.ServiceConfigEntry)
	require.True(ok)
	require.Equal("foo", serviceConf.Name)
	require.Equal(structs.ServiceDefaults, serviceConf.Kind)
}

func TestConfigEntry_Get_ACLDeny(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	dir1, s1 := testServerWithConfig(t, func(c *Config) {
		c.ACLDatacenter = "dc1"
		c.ACLsEnabled = true
		c.ACLMasterToken = "root"
		c.ACLDefaultPolicy = "deny"
	})
	defer os.RemoveAll(dir1)
	defer s1.Shutdown()
	testrpc.WaitForTestAgent(t, s1.RPC, "dc1")
	codec := rpcClient(t, s1)
	defer codec.Close()

	// Create the ACL.
	arg := structs.ACLRequest{
		Datacenter: "dc1",
		Op:         structs.ACLSet,
		ACL: structs.ACL{
			Name: "User token",
			Type: structs.ACLTokenTypeClient,
			Rules: `
service "foo" {
	policy = "read"
}
operator = "read"
`,
		},
		WriteRequest: structs.WriteRequest{Token: "root"},
	}
	var id string
	if err := msgpackrpc.CallWithCodec(codec, "ACL.Apply", &arg, &id); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Create a dummy service in the state store to look up.
	// Create some dummy service/proxy configs to be looked up.
	state := s1.fsm.State()
	require.NoError(state.EnsureConfigEntry(1, &structs.ProxyConfigEntry{
		Kind: structs.ProxyDefaults,
		Name: structs.ProxyConfigGlobal,
	}))
	require.NoError(state.EnsureConfigEntry(2, &structs.ServiceConfigEntry{
		Kind: structs.ServiceDefaults,
		Name: "foo",
	}))

	// This should fail since we don't have write perms for the "db" service.
	args := structs.ConfigEntryQuery{
		Kind:         structs.ServiceDefaults,
		Name:         "db",
		Datacenter:   s1.config.Datacenter,
		QueryOptions: structs.QueryOptions{Token: id},
	}
	var out structs.IndexedConfigEntries
	err := msgpackrpc.CallWithCodec(codec, "ConfigEntry.Get", &args, &out)
	if !acl.IsErrPermissionDenied(err) {
		t.Fatalf("err: %v", err)
	}

	// The "foo" service should work.
	args.Name = "foo"
	require.NoError(msgpackrpc.CallWithCodec(codec, "ConfigEntry.Get", &args, &out))

	serviceConf, ok := out.Entries[0].(*structs.ServiceConfigEntry)
	require.True(ok)
	require.Equal("foo", serviceConf.Name)
	require.Equal(structs.ServiceDefaults, serviceConf.Kind)

	// Try to look up the proxy config with no token.
	args.Kind = structs.ProxyDefaults
	args.Name = structs.ProxyConfigGlobal
	args.QueryOptions.Token = ""
	err = msgpackrpc.CallWithCodec(codec, "ConfigEntry.Get", &args, &out)
	if !acl.IsErrPermissionDenied(err) {
		t.Fatalf("err: %v", err)
	}

	args.QueryOptions.Token = id
	require.NoError(msgpackrpc.CallWithCodec(codec, "ConfigEntry.Get", &args, &out))
}

func TestConfigEntry_List(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	dir1, s1 := testServer(t)
	defer os.RemoveAll(dir1)
	defer s1.Shutdown()
	codec := rpcClient(t, s1)
	defer codec.Close()

	// Create some dummy services in the state store to look up.
	state := s1.fsm.State()
	expected := structs.IndexedConfigEntries{
		Entries: []structs.ConfigEntry{
			&structs.ServiceConfigEntry{
				Kind: structs.ServiceDefaults,
				Name: "bar",
			},
			&structs.ServiceConfigEntry{
				Kind: structs.ServiceDefaults,
				Name: "foo",
			},
		},
	}
	require.NoError(state.EnsureConfigEntry(1, expected.Entries[0]))
	require.NoError(state.EnsureConfigEntry(2, expected.Entries[1]))

	args := structs.ConfigEntryQuery{
		Kind:       structs.ServiceDefaults,
		Datacenter: "dc1",
	}
	var out structs.IndexedConfigEntries
	require.NoError(msgpackrpc.CallWithCodec(codec, "ConfigEntry.List", &args, &out))

	expected.QueryMeta = out.QueryMeta
	require.Equal(expected, out)
}

func TestConfigEntry_List_ACLDeny(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	dir1, s1 := testServerWithConfig(t, func(c *Config) {
		c.ACLDatacenter = "dc1"
		c.ACLsEnabled = true
		c.ACLMasterToken = "root"
		c.ACLDefaultPolicy = "deny"
	})
	defer os.RemoveAll(dir1)
	defer s1.Shutdown()
	testrpc.WaitForTestAgent(t, s1.RPC, "dc1")
	codec := rpcClient(t, s1)
	defer codec.Close()

	// Create the ACL.
	arg := structs.ACLRequest{
		Datacenter: "dc1",
		Op:         structs.ACLSet,
		ACL: structs.ACL{
			Name: "User token",
			Type: structs.ACLTokenTypeClient,
			Rules: `
service "foo" {
	policy = "read"
}
operator = "read"
`,
		},
		WriteRequest: structs.WriteRequest{Token: "root"},
	}
	var id string
	if err := msgpackrpc.CallWithCodec(codec, "ACL.Apply", &arg, &id); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Create some dummy service/proxy configs to be looked up.
	state := s1.fsm.State()
	require.NoError(state.EnsureConfigEntry(1, &structs.ProxyConfigEntry{
		Kind: structs.ProxyDefaults,
		Name: structs.ProxyConfigGlobal,
	}))
	require.NoError(state.EnsureConfigEntry(2, &structs.ServiceConfigEntry{
		Kind: structs.ServiceDefaults,
		Name: "foo",
	}))
	require.NoError(state.EnsureConfigEntry(3, &structs.ServiceConfigEntry{
		Kind: structs.ServiceDefaults,
		Name: "db",
	}))

	// This should filter out the "db" service since we don't have permissions for it.
	args := structs.ConfigEntryQuery{
		Kind:         structs.ServiceDefaults,
		Datacenter:   s1.config.Datacenter,
		QueryOptions: structs.QueryOptions{Token: id},
	}
	var out structs.IndexedConfigEntries
	err := msgpackrpc.CallWithCodec(codec, "ConfigEntry.List", &args, &out)
	require.NoError(err)

	serviceConf, ok := out.Entries[0].(*structs.ServiceConfigEntry)
	require.Len(out.Entries, 1)
	require.True(ok)
	require.Equal("foo", serviceConf.Name)
	require.Equal(structs.ServiceDefaults, serviceConf.Kind)

	// Get the global proxy config.
	args.Kind = structs.ProxyDefaults
	err = msgpackrpc.CallWithCodec(codec, "ConfigEntry.List", &args, &out)
	require.NoError(err)

	proxyConf, ok := out.Entries[0].(*structs.ProxyConfigEntry)
	require.Len(out.Entries, 1)
	require.True(ok)
	require.Equal(structs.ProxyConfigGlobal, proxyConf.Name)
	require.Equal(structs.ProxyDefaults, proxyConf.Kind)
}

func TestConfigEntry_Delete(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	dir1, s1 := testServer(t)
	defer os.RemoveAll(dir1)
	defer s1.Shutdown()
	codec := rpcClient(t, s1)
	defer codec.Close()

	// Create a dummy service in the state store to look up.
	entry := &structs.ServiceConfigEntry{
		Kind: structs.ServiceDefaults,
		Name: "foo",
	}
	state := s1.fsm.State()
	require.NoError(state.EnsureConfigEntry(1, entry))

	// Verify it's there.
	_, existing, err := state.ConfigEntry(nil, structs.ServiceDefaults, "foo")
	require.NoError(err)

	serviceConf, ok := existing.(*structs.ServiceConfigEntry)
	require.True(ok)
	require.Equal("foo", serviceConf.Name)
	require.Equal(structs.ServiceDefaults, serviceConf.Kind)

	args := structs.ConfigEntryRequest{
		Datacenter: "dc1",
	}
	args.Entry = entry
	var out struct{}
	require.NoError(msgpackrpc.CallWithCodec(codec, "ConfigEntry.Delete", &args, &out))

	// Verify the entry was deleted.
	_, existing, err = state.ConfigEntry(nil, structs.ServiceDefaults, "foo")
	require.NoError(err)
	require.Nil(existing)
}

func TestConfigEntry_Delete_ACLDeny(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	dir1, s1 := testServerWithConfig(t, func(c *Config) {
		c.ACLDatacenter = "dc1"
		c.ACLsEnabled = true
		c.ACLMasterToken = "root"
		c.ACLDefaultPolicy = "deny"
	})
	defer os.RemoveAll(dir1)
	defer s1.Shutdown()
	testrpc.WaitForTestAgent(t, s1.RPC, "dc1")
	codec := rpcClient(t, s1)
	defer codec.Close()

	// Create the ACL.
	arg := structs.ACLRequest{
		Datacenter: "dc1",
		Op:         structs.ACLSet,
		ACL: structs.ACL{
			Name: "User token",
			Type: structs.ACLTokenTypeClient,
			Rules: `
service "foo" {
	policy = "write"
}
operator = "write"
`,
		},
		WriteRequest: structs.WriteRequest{Token: "root"},
	}
	var id string
	if err := msgpackrpc.CallWithCodec(codec, "ACL.Apply", &arg, &id); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Create some dummy service/proxy configs to be looked up.
	state := s1.fsm.State()
	require.NoError(state.EnsureConfigEntry(1, &structs.ProxyConfigEntry{
		Kind: structs.ProxyDefaults,
		Name: structs.ProxyConfigGlobal,
	}))
	require.NoError(state.EnsureConfigEntry(2, &structs.ServiceConfigEntry{
		Kind: structs.ServiceDefaults,
		Name: "foo",
	}))

	// This should fail since we don't have write perms for the "db" service.
	args := structs.ConfigEntryRequest{
		Datacenter: s1.config.Datacenter,
		Entry: &structs.ServiceConfigEntry{
			Name: "db",
		},
		WriteRequest: structs.WriteRequest{Token: id},
	}
	var out struct{}
	err := msgpackrpc.CallWithCodec(codec, "ConfigEntry.Delete", &args, &out)
	if !acl.IsErrPermissionDenied(err) {
		t.Fatalf("err: %v", err)
	}

	// The "foo" service should work.
	args.Entry = &structs.ServiceConfigEntry{
		Name: "foo",
	}
	require.NoError(msgpackrpc.CallWithCodec(codec, "ConfigEntry.Delete", &args, &out))

	// Verify the entry was deleted.
	_, existing, err := state.ConfigEntry(nil, structs.ServiceDefaults, "foo")
	require.NoError(err)
	require.Nil(existing)

	// Try to delete the global proxy config without a token.
	args = structs.ConfigEntryRequest{
		Datacenter: s1.config.Datacenter,
		Entry: &structs.ProxyConfigEntry{
			Name: structs.ProxyConfigGlobal,
		},
	}
	err = msgpackrpc.CallWithCodec(codec, "ConfigEntry.Delete", &args, &out)
	if !acl.IsErrPermissionDenied(err) {
		t.Fatalf("err: %v", err)
	}

	// Now delete with a valid token.
	args.WriteRequest.Token = id
	require.NoError(msgpackrpc.CallWithCodec(codec, "ConfigEntry.Delete", &args, &out))

	_, existing, err = state.ConfigEntry(nil, structs.ServiceDefaults, "foo")
	require.NoError(err)
	require.Nil(existing)
}

func TestConfigEntry_ResolveServiceConfig(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	dir1, s1 := testServer(t)
	defer os.RemoveAll(dir1)
	defer s1.Shutdown()
	codec := rpcClient(t, s1)
	defer codec.Close()

	// Create a dummy proxy/service config in the state store to look up.
	state := s1.fsm.State()
	require.NoError(state.EnsureConfigEntry(1, &structs.ProxyConfigEntry{
		Kind: structs.ProxyDefaults,
		Name: structs.ProxyConfigGlobal,
		Config: map[string]interface{}{
			"foo": "bar",
		},
	}))
	require.NoError(state.EnsureConfigEntry(2, &structs.ServiceConfigEntry{
		Kind: structs.ServiceDefaults,
		Name: "foo",
	}))

	args := structs.ServiceConfigRequest{
		Name:       "foo",
		Datacenter: s1.config.Datacenter,
	}
	var out structs.ServiceConfigResponse
	require.NoError(msgpackrpc.CallWithCodec(codec, "ConfigEntry.ResolveServiceConfig", &args, &out))

	expected := structs.ServiceDefinition{
		Name: "foo",
		Proxy: &structs.ConnectProxyConfig{
			Config: map[string]interface{}{
				"foo": "bar",
			},
		},
	}
	out.Definition.Proxy.Config["foo"] = structs.Uint8ToString(out.Definition.Proxy.Config["foo"].([]uint8))
	require.Equal(expected, out.Definition)
}

func TestConfigEntry_ResolveServiceConfig_ACLDeny(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	dir1, s1 := testServerWithConfig(t, func(c *Config) {
		c.ACLDatacenter = "dc1"
		c.ACLsEnabled = true
		c.ACLMasterToken = "root"
		c.ACLDefaultPolicy = "deny"
	})
	defer os.RemoveAll(dir1)
	defer s1.Shutdown()
	testrpc.WaitForTestAgent(t, s1.RPC, "dc1")
	codec := rpcClient(t, s1)
	defer codec.Close()

	// Create the ACL.
	arg := structs.ACLRequest{
		Datacenter: "dc1",
		Op:         structs.ACLSet,
		ACL: structs.ACL{
			Name: "User token",
			Type: structs.ACLTokenTypeClient,
			Rules: `
service "foo" {
	policy = "write"
}
operator = "write"
`,
		},
		WriteRequest: structs.WriteRequest{Token: "root"},
	}
	var id string
	if err := msgpackrpc.CallWithCodec(codec, "ACL.Apply", &arg, &id); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Create some dummy service/proxy configs to be looked up.
	state := s1.fsm.State()
	require.NoError(state.EnsureConfigEntry(1, &structs.ProxyConfigEntry{
		Kind: structs.ProxyDefaults,
		Name: structs.ProxyConfigGlobal,
	}))
	require.NoError(state.EnsureConfigEntry(2, &structs.ServiceConfigEntry{
		Kind: structs.ServiceDefaults,
		Name: "foo",
	}))
	require.NoError(state.EnsureConfigEntry(3, &structs.ServiceConfigEntry{
		Kind: structs.ServiceDefaults,
		Name: "db",
	}))

	// This should fail since we don't have write perms for the "db" service.
	args := structs.ServiceConfigRequest{
		Name:         "db",
		Datacenter:   s1.config.Datacenter,
		QueryOptions: structs.QueryOptions{Token: id},
	}
	var out struct{}
	err := msgpackrpc.CallWithCodec(codec, "ConfigEntry.ResolveServiceConfig", &args, &out)
	if !acl.IsErrPermissionDenied(err) {
		t.Fatalf("err: %v", err)
	}

	// The "foo" service should work.
	args.Name = "foo"
	require.NoError(msgpackrpc.CallWithCodec(codec, "ConfigEntry.ResolveServiceConfig", &args, &out))

}
