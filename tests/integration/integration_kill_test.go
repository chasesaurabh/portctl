package integration

import (
	"net"
	"net/http"
	"os/exec"
	"strconv"
	"testing"
	"time"
)

// startTestHTTPServer starts a HTTP server in a goroutine and returns the listener, port, and a cleanup function
func startTestHTTPServer(t *testing.T) (listener net.Listener, port int, cleanup func()) {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	server := &http.Server{Handler: http.FileServer(http.Dir("."))}
	done := make(chan struct{})
	go func() {
		_ = server.Serve(listener)
		close(done)
	}()
	portStr := listener.Addr().(*net.TCPAddr).Port
	return listener, portStr, func() {
		_ = server.Close()
		listener.Close()
		select {
		case <-done:
		case <-time.After(2 * time.Second):
			t.Log("server did not shut down in time")
		}
	}
}

func TestPortctlKillIntegration(t *testing.T) {
	_, port, cleanup := startTestHTTPServer(t)
	defer cleanup()

	t.Logf("Started test HTTP server on port %d", port)

	// Run portctl in dry-run mode
	dryRunCmd := exec.Command("../../cmd/portctl/portctl", "-dry-run", strconv.Itoa(port))
	dryRunOut, dryRunErr := dryRunCmd.CombinedOutput()
	t.Logf("portctl dry-run output: %s", string(dryRunOut))
	if dryRunErr != nil {
		t.Errorf("portctl dry-run failed: %v", dryRunErr)
	}

	// Check if server is still running (should be)
	conn, err := net.DialTimeout("tcp", "127.0.0.1:"+strconv.Itoa(port), 500*time.Millisecond)
	if err != nil {
		t.Errorf("server not running after dry-run: %v", err)
	} else {
		_ = conn.Close()
	}

	// Run portctl in force mode
	forceCmd := exec.Command("../../cmd/portctl/portctl", "-force", "-timeout", "3s", "-kill-after", strconv.Itoa(port))
	forceOut, forceErr := forceCmd.CombinedOutput()
	t.Logf("portctl force output: %s", string(forceOut))
	if forceErr != nil {
		t.Errorf("portctl force failed: %v", forceErr)
	}

	// Check if server is terminated (should not be reachable)
	conn2, err2 := net.DialTimeout("tcp", "127.0.0.1:"+strconv.Itoa(port), 500*time.Millisecond)
	if err2 == nil {
		_ = conn2.Close()
		t.Errorf("server still running after force kill")
	}
}
