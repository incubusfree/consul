package consul

import (
	"database/sql"
	"fmt"
	"github.com/hashicorp/consul/rpc"
	_ "github.com/mattn/go-sqlite3"
	"sync/atomic"
)

// nextDBIndex is used to generate a new ID
// using sync/atomic to ensure it is safe
var nextDBIndex uint32 = 0

type namedQuery uint8

const (
	queryEnsureNode namedQuery = iota
	queryNode
	queryNodes
	queryEnsureService
	queryNodeServices
	queryDeleteNodeService
	queryDeleteNode
	queryServices
	queryServiceNodes
	queryServiceTagNodes
)

// The StateStore is responsible for maintaining all the Consul
// state. It is manipulated by the FSM which maintains consistency
// through the use of Raft. The goals of the StateStore are to provide
// high concurrency for read operations without blocking writes, and
// to provide write availability in the face of reads. The current
// implementation uses an in-memory SQLite database. This reduced the
// GC pressure on Go, and also gives us Multi-Version Concurrency Control
// for "free".
type StateStore struct {
	db       *sql.DB
	prepared map[namedQuery]*sql.Stmt
}

// NewStateStore is used to create a new state store
func NewStateStore() (*StateStore, error) {
	// Get the DB ID
	id := atomic.AddUint32(&nextDBIndex, 1)
	path := fmt.Sprintf("file:StateStore-%d?mode=memory&cache=shared", id)

	// Open the db
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %v", err)
	}

	s := &StateStore{
		db:       db,
		prepared: make(map[namedQuery]*sql.Stmt),
	}

	// Ensure we can initialize
	if err := s.initialize(); err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}

// Close is used to safely shutdown the state store
func (s *StateStore) Close() error {
	return s.db.Close()
}

// initialize is used to setup the sqlite store for use
func (s *StateStore) initialize() error {
	// Set the pragma first
	pragmas := []string{
		"pragma journal_mode=memory;",
		"pragma foreign_keys=ON;",
	}
	for _, p := range pragmas {
		if _, err := s.db.Exec(p); err != nil {
			return fmt.Errorf("Failed to set '%s': %v", p, err)
		}
	}

	// Create the tables
	tables := []string{
		`CREATE TABLE nodes (name text unique, address text);`,
		`CREATE TABLE services (node text REFERENCES nodes(name) ON DELETE CASCADE, service text, tag text, port integer);`,
		`CREATE INDEX servName ON services(service, tag);`,
		`CREATE INDEX nodeName ON services(node);`,
	}
	for _, t := range tables {
		if _, err := s.db.Exec(t); err != nil {
			return fmt.Errorf("Failed to call '%s': %v", t, err)
		}
	}

	// Prepare the queries
	queries := map[namedQuery]string{
		queryEnsureNode:        "INSERT OR REPLACE INTO nodes (name, address) VALUES (?, ?)",
		queryNode:              "SELECT address FROM nodes where name=?",
		queryNodes:             "SELECT * FROM nodes",
		queryEnsureService:     "INSERT OR REPLACE INTO services (node, service, tag, port) VALUES (?, ?, ?, ?)",
		queryNodeServices:      "SELECT service, tag, port from services where node=?",
		queryDeleteNodeService: "DELETE FROM services WHERE node=? AND service=?",
		queryDeleteNode:        "DELETE FROM nodes WHERE name=?",
		queryServices:          "SELECT DISTINCT service, tag FROM services",
		queryServiceNodes:      "SELECT n.name, n.address, s.tag, s.port from nodes n, services s WHERE s.service=? AND s.node=n.name",
		queryServiceTagNodes:   "SELECT n.name, n.address, s.tag, s.port from nodes n, services s WHERE s.service=? AND s.tag=? AND s.node=n.name",
	}
	for name, query := range queries {
		stmt, err := s.db.Prepare(query)
		if err != nil {
			return fmt.Errorf("Failed to prepare '%s': %v", query, err)
		}
		s.prepared[name] = stmt
	}
	return nil
}

func (s *StateStore) checkSet(res sql.Result, err error) error {
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return fmt.Errorf("Failed to set row")
	}
	return nil
}

func (s *StateStore) checkDelete(res sql.Result, err error) error {
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}

// EnsureNode is used to ensure a given node exists, with the provided address
func (s *StateStore) EnsureNode(name string, address string) error {
	stmt := s.prepared[queryEnsureNode]
	return s.checkSet(stmt.Exec(name, address))
}

// GetNode returns all the address of the known and if it was found
func (s *StateStore) GetNode(name string) (bool, string) {
	stmt := s.prepared[queryNode]
	row := stmt.QueryRow(name)

	var addr string
	if err := row.Scan(&addr); err != nil {
		if err == sql.ErrNoRows {
			return false, addr
		} else {
			panic(fmt.Errorf("Failed to get node: %v", err))
		}
	}
	return true, addr
}

// GetNodes returns all the known nodes, the slice alternates between
// the node name and address
func (s *StateStore) Nodes() []string {
	stmt := s.prepared[queryNodes]
	return parseNodes(stmt.Query())
}

// parseNodes parses the result of a queryNodes statement
func parseNodes(rows *sql.Rows, err error) []string {
	if err != nil {
		panic(fmt.Errorf("Failed to get nodes: %v", err))
	}
	data := make([]string, 0, 32)
	var name, address string
	for rows.Next() {
		if err := rows.Scan(&name, &address); err != nil {
			panic(fmt.Errorf("Failed to get nodes: %v", err))
		}
		data = append(data, name, address)
	}
	return data
}

// EnsureService is used to ensure a given node exposes a service
func (s *StateStore) EnsureService(name, service, tag string, port int) error {
	stmt := s.prepared[queryEnsureService]
	return s.checkSet(stmt.Exec(name, service, tag, port))
}

// NodeServices is used to return all the services of a given node
func (s *StateStore) NodeServices(name string) rpc.NodeServices {
	stmt := s.prepared[queryNodeServices]
	return parseNodeServices(stmt.Query(name))
}

// parseNodeServices is used to parse the results of a queryNodeServices
func parseNodeServices(rows *sql.Rows, err error) rpc.NodeServices {
	if err != nil {
		panic(fmt.Errorf("Failed to get node services: %v", err))
	}

	services := rpc.NodeServices(make(map[string]rpc.NodeService))
	var service string
	var entry rpc.NodeService
	for rows.Next() {
		if err := rows.Scan(&service, &entry.Tag, &entry.Port); err != nil {
			panic(fmt.Errorf("Failed to get node services: %v", err))
		}
		services[service] = entry
	}
	return services
}

// DeleteNodeService is used to delete a node service
func (s *StateStore) DeleteNodeService(node, service string) error {
	stmt := s.prepared[queryDeleteNodeService]
	return s.checkDelete(stmt.Exec(node, service))
}

// DeleteNode is used to delete a node and all it's services
func (s *StateStore) DeleteNode(node string) error {
	stmt := s.prepared[queryDeleteNode]
	return s.checkDelete(stmt.Exec(node))
}

// Services is used to return all the services with a list of associated tags
func (s *StateStore) Services() map[string][]string {
	stmt := s.prepared[queryServices]
	rows, err := stmt.Query()
	if err != nil {
		panic(fmt.Errorf("Failed to get services: %v", err))
	}

	services := make(map[string][]string)
	var service, tag string
	for rows.Next() {
		if err := rows.Scan(&service, &tag); err != nil {
			panic(fmt.Errorf("Failed to get services: %v", err))
		}

		tags := services[service]
		tags = append(tags, tag)
		services[service] = tags
	}

	return services
}

// ServiceNodes returns the nodes associated with a given service
func (s *StateStore) ServiceNodes(service string) rpc.ServiceNodes {
	stmt := s.prepared[queryServiceNodes]
	return parseServiceNodes(stmt.Query(service))
}

// ServiceTagNodes returns the nodes associated with a given service matching a tag
func (s *StateStore) ServiceTagNodes(service, tag string) rpc.ServiceNodes {
	stmt := s.prepared[queryServiceTagNodes]
	return parseServiceNodes(stmt.Query(service, tag))
}

// parseServiceNodes parses results from the queryServiceNodes / queryServiceTagNodes query
func parseServiceNodes(rows *sql.Rows, err error) rpc.ServiceNodes {
	if err != nil {
		panic(fmt.Errorf("Failed to get service nodes: %v", err))
	}
	var nodes rpc.ServiceNodes
	var node rpc.ServiceNode
	for rows.Next() {
		if err := rows.Scan(&node.Node, &node.Address, &node.ServiceTag, &node.ServicePort); err != nil {
			panic(fmt.Errorf("Failed to get services: %v", err))
		}
		nodes = append(nodes, node)
	}
	return nodes
}
