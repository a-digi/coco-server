package pid_test

import (
	"os"
	"testing"
	"github.com/a-digi/coco-server/server"
)

func TestWriteReadRemovePID(t *testing.T) {
	pidFile := "test.pid"
	defer os.Remove(pidFile)
	if err := server.WritePID(pidFile); err != nil {
		t.Fatalf("WritePID error: %v", err)
	}
	pid, err := server.ReadPID(pidFile)
	if err != nil {
		t.Fatalf("ReadPID error: %v", err)
	}
	if pid != os.Getpid() {
		t.Errorf("PID mismatch: got %d, want %d", pid, os.Getpid())
	}
	if err := server.RemovePID(pidFile); err != nil {
		t.Fatalf("RemovePID error: %v", err)
	}
	if _, err := os.Stat(pidFile); !os.IsNotExist(err) {
		t.Errorf("PID file still exists after RemovePID")
	}
}
