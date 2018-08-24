package daemonize_test

import (
	"testing"

	"github.com/swapbyt3s/zenit/daemonize"
)

func TestRun(t *testing.T) {
	cmd := "echo 'test'"
	expected := 0
	result := daemonize.Run(cmd)

	if result == expected {
		t.Error("Expected: pid > 0")
	}
}
