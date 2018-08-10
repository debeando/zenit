package daemonize_test

import (
  "testing"

  "github.com/swapbyt3s/zenit/config"
  "github.com/swapbyt3s/zenit/daemonize"
)

func TestArgs(t *testing.T) {
  args     := []string{"zenit", "-start"}
  expected := ""
  result   := daemonize.Args(args)

  if result != expected {
    t.Error("Expected.")
  }
}

func TestBuild(t *testing.T) {
  args     := "-start"
  cmd      := "/usr/local/bin/zenit"
  expected := "/usr/local/bin/zenit -start"
  result   := daemonize.Build(cmd, args)

  if result != expected {
    t.Error("Expected: /usr/local/bin/zenit -start")
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
  expected := "/var/run/zenit.pid"
  result   := config.General.PIDFile

  if result != expected {
    t.Error("Expected: /var/run/zenit.pid")
  }
}

func TestPIDFileExist(t *testing.T) {
  file     := "/var/run/zenit.pid"
  expected := false
  result   := config.General.PIDFile

  if ! result == expected {
    t.Error("Expected: /var/run/zenit.pid")
  }
}
