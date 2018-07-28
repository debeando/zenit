package slow_test

import (
  "regexp"
  "testing"
  "gitlab.com/swapbyt3s/zenit/common"
  "gitlab.com/swapbyt3s/zenit/collect/mysql/slow"
)

func TestRow(t *testing.T) {
  row := "SELECT * FROM foo; " +
         "# Time: 180625 15:25:03 " +
         "# User@Host: test[test] @ [127.0.0.1] " +
         "# Thread_id: 123456  Schema: test  Last_errno: 0  Killed: 0 " +
         "# Query_time: 0.792864  Lock_time: 0.000160  Rows_sent: 1  Rows_examined: 100  Rows_affected: 0  Rows_read: 100 " +
         "# Bytes_sent: 60 " +
         "SET timestamp=1529940303; " +
         "SELECT count(*) AS total FROM foo; " +
         "# Time: 180625 15:25:03 "

  expected := "# Time: 180625 15:25:03 # User@Host: test[test] @ [127.0.0.1] # Thread_id: 123456  Schema: test  Last_errno: 0  Killed: 0 # Query_time: 0.792864  Lock_time: 0.000160  Rows_sent: 1  Rows_examined: 100  Rows_affected: 0  Rows_read: 100 # Bytes_sent: 60 SET timestamp=1529940303; SELECT count(*) AS total FROM foo;"
  reRow    := regexp.MustCompile(slow.ROW)
  result   := reRow.FindString(row)

  if result != expected {
    t.Error("Expected: " + expected)
  }

  row = "SELECT count(*) AS total FROM foo; " +
        "# User@Host: test[test] @ [127.0.0.1] " +
        "# Thread_id: 123456  Schema: test  Last_errno: 0  Killed: 0 " +
        "# Query_time: 0.792864  Lock_time: 0.000160  Rows_sent: 1  Rows_examined: 100  Rows_affected: 0  Rows_read: 100 " +
        "# Bytes_sent: 60 " +
        "SET timestamp=1529940303; " +
        "UPDATE foo SET bar = NOW() WHERE id = 1; " +
        "# Time: 180625 15:25:03 "

  expected = "# User@Host: test[test] @ [127.0.0.1] # Thread_id: 123456  Schema: test  Last_errno: 0  Killed: 0 # Query_time: 0.792864  Lock_time: 0.000160  Rows_sent: 1  Rows_examined: 100  Rows_affected: 0  Rows_read: 100 # Bytes_sent: 60 SET timestamp=1529940303; UPDATE foo SET bar = NOW() WHERE id = 1;"
  reRow    = regexp.MustCompile(slow.ROW)
  result   = reRow.FindString(row)

  if result != expected {
    t.Error("Expected: " + expected)
  }

  row = "# User@Host: test[test] @ [127.0.0.1] " +
        "# Thread_id: 0  Schema: test  Last_errno: 0  Killed: 0 " +
        "# Query_time: 0.0  Lock_time: 0.0  Rows_sent: 1  Rows_examined: 0  Rows_affected: 0  Rows_read: 0 " +
        "# Bytes_sent: 0 " +
        "UPDATE foo SET bar = NOW() WHERE id = 1; " +
        "# Time: 000000 00:00:00 "

  reRow  = regexp.MustCompile(slow.ROW)
  result = reRow.FindString(row)

  if result != "" {
    t.Error("Expected nothing.")
  }
}

func TestKeysAndValues(t *testing.T) {
  row := "# Time: 180625 15:25:03 " +
         "# User@Host: test[test] @ [127.0.0.1] " +
         "# Thread_id: 123456  Schema: test  Last_errno: 0  Killed: 0 " +
         "# Query_time: 0.792864  Lock_time: 0.000160  Rows_sent: 1  Rows_examined: 100  Rows_affected: 0  Rows_read: 100 " +
         "# Bytes_sent: 60 " +
         "SET timestamp=1529940303; " +
         "SELECT count(*) AS total FROM foo; "

  reKV    := regexp.MustCompile(slow.KV)
  result  := common.RegexpGetGroups(reKV, row)

  if value, ok := result["time"]; ok {
    if value != "180625 15:25:03" {
      t.Error("Expected: 180625 15:25:03")
    }
  }

  if value, ok := result["user_host"]; ok {
    if value != "test[test] @ [127.0.0.1]" {
      t.Error("Expected: test[test] @ [127.0.0.1]")
    }
  }

  if value, ok := result["Thread_id"]; ok {
    if value != "123456" {
      t.Error("Expected: 123456")
    }
  }

  if value, ok := result["Query_time"]; ok {
    if value != "0.792864" {
      t.Error("Expected: 0.792864")
    }
  }
}
