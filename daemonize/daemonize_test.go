package daemonize_test

import (
  "testing"

  "gitlab.com/swapbyt3s/zenit/daemonize"
)

func TestArgs(t *testing.T) {
  args     := []string{"zenit", "-parser-format=slowlog", "-parser-file=/tmp/test_slow.log", "-daemonize"}
  expected := "-parser-format=slowlog -parser-file=/tmp/test_slow.log"
  result   := daemonize.Args(args)

  if result != expected {
    t.Error("Expected: -parser-format=slowlog -parser-file=/tmp/test_slow.log")
  }
}

func TestBuild(t *testing.T) {
  args     := "-collect=mysql"
  cmd      := "/usr/local/bin/zenit"
  expected := "/usr/local/bin/zenit -collect=mysql"
  result   := daemonize.Build(cmd, args)

  if result != expected {
    t.Error("Expected: /usr/local/bin/zenit -collect=mysql")
  }
}

func TestRun(t *testing.T) {
  cmd      := "echo 'test'"
  expected := 0
  result   := daemonize.Run(cmd)

  if result == expected {
    t.Error("Expected: pid > 0")
  }
}

func TestGetPIDFileName(t *testing.T) {
  args     := "-collect=mysql"
  expected := "/var/run/zenit-8791b0340bf9d02e19aa92a1ea6d07bc.pid"
  result   := daemonize.GetPIDFileName(args)

  if result != expected {
    t.Error("Expected: /var/run/zenit-8791b0340bf9d02e19aa92a1ea6d07bc.pid")
  }
}

func TestPIDFileExist(t *testing.T) {
  file     := "/var/run/zenit-098f6bcd4621d373cade4e832627b4f6.pid"
  expected := false
  result   := daemonize.PIDFileExist(file)

  if ! result == expected {
    t.Error("Expected: /var/run/zenit-098f6bcd4621d373cade4e832627b4f6.pid")
  }
}
