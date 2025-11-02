package kill

import (
	"testing"
	"time"
)

func TestSignalToNumber(t *testing.T) {
	n, err := signalToNumber("TERM")
	if err != nil || n <= 0 {
		t.Fatalf("expected TERM to map to a positive signal number, got %v %v", n, err)
	}
	n2, err := signalToNumber("SIGKILL")
	if err != nil || n2 <= 0 {
		t.Fatalf("expected SIGKILL to map to a positive signal number, got %v %v", n2, err)
	}
}

func TestProcessExistsFalse(t *testing.T) {
	// pick a high pid that likely doesn't exist
	if processExists(999999) {
		t.Fatalf("expected pid 999999 to not exist")
	}
}

func TestWaitForExitTimeout(t *testing.T) {
	// use a pid that doesn't exist so WaitForExit returns immediately
	err := WaitForExit(999999, 500*time.Millisecond)
	if err != nil {
		t.Fatalf("expected no error for non-existent pid, got %v", err)
	}
}
