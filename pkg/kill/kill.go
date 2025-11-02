package kill

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// RunnerFunc allows injecting commands for testing (not used heavily here)
type RunnerFunc func(name string, args ...string) ([]byte, error)

// SendSignal sends the given signal (name or number) to pid
func SendSignal(pid int, sig string) error {
	// resolve signal name to number if necessary
	num, err := signalToNumber(sig)
	if err != nil {
		return err
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	// On unix, use syscall.Kill
	if err := proc.Signal(syscall.Signal(num)); err != nil {
		return err
	}
	return nil
}

// WaitForExit waits until pid is gone or timeout expires
func WaitForExit(pid int, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if !processExists(pid) {
			return nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	if processExists(pid) {
		return errors.New("timeout waiting for process to exit")
	}
	return nil
}

func processExists(pid int) bool {
	// On unix, sending signal 0 checks existence
	err := syscall.Kill(pid, 0)
	if err == nil {
		return true
	}
	if err == syscall.ESRCH {
		return false
	}
	// other errors (EPERM) mean it exists but not accessible
	return true
}

func signalToNumber(sig string) (int, error) {
	// accept numeric
	if n, err := strconv.Atoi(sig); err == nil {
		return n, nil
	}
	// strip SIG prefix
	s := strings.ToUpper(sig)
	// Normalize both SIGTERM and TERM forms by trimming a leading "SIG" if present
	s = strings.TrimPrefix(s, "SIG")
	switch s {
	case "TERM":
		return int(syscall.SIGTERM), nil
	case "KILL":
		return int(syscall.SIGKILL), nil
	case "INT":
		return int(syscall.SIGINT), nil
	case "HUP":
		return int(syscall.SIGHUP), nil
	default:
		return 0, fmt.Errorf("unsupported signal: %s", sig)
	}
}

// Helper to run a check command; used by tests
func RunCommand(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	return cmd.CombinedOutput()
}
