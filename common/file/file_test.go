package file_test

import (
  "os"
  "path/filepath"
  "testing"

  "github.com/swapbyt3s/zenit/common/file"
)

var (
 twd string
 tf  string
)

func init() {
  ex, _ := os.Executable()
  twd = filepath.Dir(ex)
  tf  = twd + "/zenit.txt"
}

func TestCreate(t *testing.T) {
  if ! file.Create(tf) {
    t.Error("Problem to create file.")
  }

  if _, err := os.Stat(tf); os.IsNotExist(err) {
    t.Error("File not exist in: zenit.txt")
  }
}

func TestWrite(t *testing.T) {
  if ! file.Write(tf, "Test 1\nTest 2") {
    t.Error("Problem to write in file.")
  }
}

func TestRead(t *testing.T) {
  result   := file.Read(tf)
  expected := "Test 1\nTest 2"

  if result != expected {
    t.Error("Expected: " + expected)
  }
}

func TestTruncate(t *testing.T) {
  if ! file.Truncate(tf) {
    t.Error("Problem to truncate file.")
  }

  if len(file.Read(tf)) != 0 {
    t.Error("Is not truncated file.")
  }
}
