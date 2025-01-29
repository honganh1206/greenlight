package mailer

import (
	"net"
	"testing"
	"time"
)

var (
	testConn = &net.TCPConn{} // Implementation of Conn interface for TCP
)

type mockClient struct {
	t       *testing.T
	i       int
	want    []string
	addr    string
	timeout bool
}

// HELPER FUNCTIONS
func doTestSendMail(t *testing.T, d *Dialer, want []string, timeout bool) {
	testClient := &mockClient{
		t:       t,
		want:    want,
		addr:    addr(d.Host, d.Port),
		timeout: timeout,
	}

	// Stub for test, defined in smtp.go later
	netDialTimeout = func(network, address string, d time.Duration) (net.Conn, error) {
		if network != "tcp" {
			t.Errorf("Invalid network, got: %q, want tcp", network)
		}

		if address != testClient.addr {
			t.Errorf("Invalid address. got: %q, want: %q", address, testClient.addr)
		}

		return testConn, nil // Stub for testing
	}

	// smtpClient = func(conn net.Conn, host string) (smtpClient, error) {
	// 	if host != testHost {
	// 		t.Errorf("Invalid host. got: %q, want: %q", host, testHost)
	// 	}

	// 	return testClient, nil
	// }

}
