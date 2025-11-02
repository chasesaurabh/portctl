package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/kumarsaurabh/killport/pkg/discover"
)

func sampleListeners() []discover.Listener {
	return []discover.Listener{
		{PID: 1234, User: "alice", Cmd: "python3", Address: "127.0.0.1:8000"},
		{PID: 2345, User: "bob", Cmd: "node", Address: "0.0.0.0:3000"},
	}
}

func TestText(t *testing.T) {
	var buf bytes.Buffer
	Text(&buf, sampleListeners())
	out := buf.String()
	if !strings.Contains(out, "PID=1234") || !strings.Contains(out, "CMD=node") {
		t.Fatalf("unexpected text output: %s", out)
	}
}

func TestJSON(t *testing.T) {
	var buf bytes.Buffer
	if err := JSON(&buf, sampleListeners()); err != nil {
		t.Fatalf("json encode failed: %v", err)
	}
	var arr []discover.Listener
	if err := json.Unmarshal(buf.Bytes(), &arr); err != nil {
		t.Fatalf("json decode failed: %v", err)
	}
	if len(arr) != 2 || arr[0].PID != 1234 {
		t.Fatalf("unexpected json content: %+v", arr)
	}
}

func TestCSV(t *testing.T) {
	var buf bytes.Buffer
	if err := CSV(&buf, sampleListeners()); err != nil {
		t.Fatalf("csv failed: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "pid,user,cmd,address") || !strings.Contains(out, "127.0.0.1:8000") {
		t.Fatalf("unexpected csv output: %s", out)
	}
}
