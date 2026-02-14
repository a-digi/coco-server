package log_test

import (
	"strings"
	"testing"
	"github.com/a-digi/coco-server/server"
)

func TestLogFileName(t *testing.T) {
	name := server.LogFileName("test")
	if !strings.HasPrefix(name, "test_") || !strings.HasSuffix(name, ".log") {
		t.Errorf("unexpected log file name: %s", name)
	}
}
