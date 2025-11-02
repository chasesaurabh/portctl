package discover

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// Listener holds information about a listening socket/process
type Listener struct {
	PID     int
	User    string
	Cmd     string
	Address string
}

// RunnerFunc allows injection of command execution for testing
type RunnerFunc func(name string, args ...string) ([]byte, error)

// defaultRunner runs real commands
func defaultRunner(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	// set a timeout to avoid hanging
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	done := make(chan error)
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		if err != nil {
			return out.Bytes(), err
		}
		return out.Bytes(), nil
	case <-time.After(3 * time.Second):
		// kill process
		_ = cmd.Process.Kill()
		return out.Bytes(), errors.New("command timeout")
	}
}

// Discover finds listeners on the given TCP port. It chooses the best available tool.
func Discover(port string) ([]Listener, error) {
	return DiscoverWithRunner(defaultRunner, port)
}

// DiscoverWithRunner is injectable for testing
func DiscoverWithRunner(runner RunnerFunc, port string) ([]Listener, error) {
	// try ss
	if out, err := runner("ss", "-ltnp"); err == nil && len(out) > 0 {
		listeners := parseSS(out, port)
		if len(listeners) > 0 {
			return listeners, nil
		}
	}
	// try lsof
	if out, err := runner("lsof", "-nP", fmt.Sprintf("-iTCP:%s", port), "-sTCP:LISTEN"); err == nil && len(out) > 0 {
		listeners := parseLSOF(out)
		if len(listeners) > 0 {
			return listeners, nil
		}
	}
	// try netstat
	if out, err := runner("netstat", "-plnt"); err == nil && len(out) > 0 {
		listeners := parseNetstat(out, port)
		if len(listeners) > 0 {
			return listeners, nil
		}
	}
	return nil, errors.New("no listeners found or no supported tools available")
}

// parseSS parses `ss -ltnp` output (Linux)
func parseSS(output []byte, port string) []Listener {
	var res []Listener
	s := bufio.NewScanner(bytes.NewReader(output))
	// Example line: LISTEN   0      128        127.0.0.1:8000       *:*    users:(("python3",pid=1234,fd=3))
	re := regexp.MustCompile(`LISTEN\s+\S+\s+\S+\s+([^\s]+:` + regexp.QuoteMeta(port) + `)\s+[^\s]+\s+users:\(\("([^"]+)",pid=(\d+),`) 
	for s.Scan() {
		line := s.Text()
		if matches := re.FindStringSubmatch(line); len(matches) == 4 {
			addr := matches[1]
			cmd := matches[2]
			pidStr := matches[3]
			pid := atoiSafe(pidStr)
			res = append(res, Listener{PID: pid, User: "", Cmd: cmd, Address: addr})
		}
	}
	return res
}

// parseLSOF parses `lsof -iTCP:PORT -sTCP:LISTEN` output
func parseLSOF(output []byte) []Listener {
	var res []Listener
	s := bufio.NewScanner(bytes.NewReader(output))
	// lsof header: COMMAND PID USER FD TYPE DEVICE SIZE/OFF NODE NAME
	first := true
	for s.Scan() {
		line := s.Text()
		if first {
			first = false
			continue
		}
		// split by spaces (multiple) - conservative approach
		parts := fieldsN(line, 9)
		if len(parts) < 9 {
			continue
		}
		cmd := parts[0]
		pid := atoiSafe(parts[1])
		user := parts[2]
		name := parts[8]
		res = append(res, Listener{PID: pid, User: user, Cmd: cmd, Address: name})
	}
	return res
}

// parseNetstat parses `netstat -plnt` output
func parseNetstat(output []byte, port string) []Listener {
	var res []Listener
	s := bufio.NewScanner(bytes.NewReader(output))
	// Example: tcp        0      0 0.0.0.0:22            0.0.0.0:*               LISTEN      1234/sshd
	re := regexp.MustCompile(`\S+\s+\S+\s+\S+\s+([^\s]+:` + regexp.QuoteMeta(port) + `)\s+\S+\s+LISTEN\s+(\d+)/(\S+)`)
	for s.Scan() {
		line := s.Text()
		if matches := re.FindStringSubmatch(line); len(matches) == 4 {
			addr := matches[1]
			pid := atoiSafe(matches[2])
			cmd := matches[3]
			res = append(res, Listener{PID: pid, Cmd: cmd, Address: addr})
		}
	}
	return res
}

// fieldsN splits a line into at most n fields separated by whitespace
func fieldsN(s string, n int) []string {
	f := regexp.MustCompile(`\s+`).Split(strings.TrimSpace(s), n)
	return f
}

func atoiSafe(s string) int {
	var v int
	fmt.Sscanf(s, "%d", &v)
	return v
}
