package mysql_test

import (
  "testing"

  "gitlab.com/swapbyt3s/zenit/plugins/inputs/mysql"
)

func TestClearUser(t *testing.T) {
  user     := "test[test] @ [127.0.0.1]"
  expected := "test"
  result   := mysql.ClearUser(user)

  if result != expected {
    t.Error("Expected: " + expected)
  }
}
