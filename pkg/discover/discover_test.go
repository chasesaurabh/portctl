package discover

import (
	"testing"
)

func TestParseSS(t *testing.T) {
	sample := `Netid  State      Recv-Q Send-Q Local Address:Port               Peer Address:Port  
LISTEN 0         128        127.0.0.1:8000                    *:*      users:(("python3",pid=1234,fd=3))
`
	res := parseSS([]byte(sample), "8000")
	if len(res) != 1 {
		t.Fatalf("expected 1 listener, got %d", len(res))
	}
	if res[0].PID != 1234 || res[0].Cmd != "python3" {
		t.Fatalf("unexpected parsed listener: %+v", res[0])
	}
}

func TestParseLSOF(t *testing.T) {
	sample := `COMMAND     PID   USER   FD   TYPE DEVICE SIZE/OFF NODE NAME
python3    2345   alice  3u  IPv4  0xabc      0t0  TCP 127.0.0.1:9000 (LISTEN)`
	res := parseLSOF([]byte(sample))
	if len(res) != 1 {
		t.Fatalf("expected 1 listener, got %d", len(res))
	}
	if res[0].PID != 2345 || res[0].User != "alice" || res[0].Cmd != "python3" {
		t.Fatalf("unexpected parsed listener: %+v", res[0])
	}
}

func TestParseNetstat(t *testing.T) {
	sample := `tcp        0      0 0.0.0.0:22            0.0.0.0:*               LISTEN      1234/sshd
`
	res := parseNetstat([]byte(sample), "22")
	if len(res) != 1 {
		t.Fatalf("expected 1 listener, got %d", len(res))
	}
	if res[0].PID != 1234 || res[0].Cmd != "sshd" {
		t.Fatalf("unexpected parsed listener: %+v", res[0])
	}
}
