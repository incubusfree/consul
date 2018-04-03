package proxy

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"

	"github.com/hashicorp/consul/lib/freeport"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"
)

// TestLocalBindAddrs returns n localhost address:port strings with free ports
// for binding test listeners to.
func TestLocalBindAddrs(t testing.T, n int) []string {
	ports := freeport.GetT(t, n)
	addrs := make([]string, n)
	for i, p := range ports {
		addrs[i] = fmt.Sprintf("localhost:%d", p)
	}
	return addrs
}

// TestTCPServer is a simple TCP echo server for use during tests.
type TestTCPServer struct {
	l                        net.Listener
	stopped                  int32
	accepted, closed, active int32
}

// NewTestTCPServer opens as a listening socket on the given address and returns
// a TestTCPServer serving requests to it. The server is already started and can
// be stopped by calling Close().
func NewTestTCPServer(t testing.T, addr string) (*TestTCPServer, error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	log.Printf("test tcp server listening on %s", addr)
	s := &TestTCPServer{
		l: l,
	}
	go s.accept()
	return s, nil
}

// Close stops the server
func (s *TestTCPServer) Close() {
	atomic.StoreInt32(&s.stopped, 1)
	if s.l != nil {
		s.l.Close()
	}
}

func (s *TestTCPServer) accept() error {
	for {
		conn, err := s.l.Accept()
		if err != nil {
			if atomic.LoadInt32(&s.stopped) == 1 {
				log.Printf("test tcp echo server %s stopped", s.l.Addr())
				return nil
			}
			log.Printf("test tcp echo server %s failed: %s", s.l.Addr(), err)
			return err
		}

		atomic.AddInt32(&s.accepted, 1)
		atomic.AddInt32(&s.active, 1)

		go func(c net.Conn) {
			io.Copy(c, c)
			atomic.AddInt32(&s.closed, 1)
			atomic.AddInt32(&s.active, -1)
		}(conn)
	}
}

// TestEchoConn attempts to write some bytes to conn and expects to read them
// back within a short timeout (10ms). If prefix is not empty we expect it to be
// poresent at the start of all echoed responses (for example to distinguish
// between multiple echo server instances).
func TestEchoConn(t testing.T, conn net.Conn, prefix string) {
	t.Helper()

	// Write some bytes and read them back
	n, err := conn.Write([]byte("Hello World"))
	require.Equal(t, 11, n)
	require.Nil(t, err)

	expectLen := 11 + len(prefix)

	buf := make([]byte, expectLen)
	// read until our buffer is full - it might be separate packets if prefix is
	// in use.
	got := 0
	for got < expectLen {
		n, err = conn.Read(buf[got:])
		require.Nilf(t, err, "err: %s", err)
		got += n
	}
	require.Equal(t, expectLen, got)
	require.Equal(t, prefix+"Hello World", string(buf[:]))
}
