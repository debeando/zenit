package slow_test

import (
  "regexp"
  "testing"

  "github.com/swapbyt3s/zenit/common"
  "github.com/swapbyt3s/zenit/common/sql"
  "github.com/swapbyt3s/zenit/config"
  "github.com/swapbyt3s/zenit/plugins/inputs/mysql/slow"
)

var slowLog = []struct{ID, Input, Expected string}{
  {"row_case_1",
   "SELECT * FROM foo; " +
   "# Time: 180625 15:25:03 " +
   "# User@Host: test[test] @ [127.0.0.1] " +
   "# Thread_id: 123456  Schema: test  Last_errno: 0  Killed: 0 " +
   "# Query_time: 0.792864  Lock_time: 0.000160  Rows_sent: 1  Rows_examined: 100  Rows_affected: 0  Rows_read: 100 " +
   "# Bytes_sent: 60 " +
   "SET timestamp=1529940303; " +
   "SELECT count(*) AS total FROM foo WHERE att = 'bar'; " +
   "# Time: 180625 15:25:03 ",
   "# Time: 180625 15:25:03 # User@Host: test[test] @ [127.0.0.1] # Thread_id: 123456  Schema: test  Last_errno: 0  Killed: 0 # Query_time: 0.792864  Lock_time: 0.000160  Rows_sent: 1  Rows_examined: 100  Rows_affected: 0  Rows_read: 100 # Bytes_sent: 60 SET timestamp=1529940303; SELECT count(*) AS total FROM foo WHERE att = 'bar';"},
  {"row_case_2",
   "SELECT count(*) AS total FROM foo; " +
   "# User@Host: test[test] @ [127.0.0.1] " +
   "# Thread_id: 123456  Schema: test  Last_errno: 0  Killed: 0 " +
   "# Query_time: 0.792864  Lock_time: 0.000160  Rows_sent: 1  Rows_examined: 100  Rows_affected: 0  Rows_read: 100 " +
   "# Bytes_sent: 60 " +
   "SET timestamp=1529940303; " +
   "UPDATE foo SET bar = NOW() WHERE id = 1; " +
   "# Time: 180625 15:25:03 ",
   "# User@Host: test[test] @ [127.0.0.1] # Thread_id: 123456  Schema: test  Last_errno: 0  Killed: 0 # Query_time: 0.792864  Lock_time: 0.000160  Rows_sent: 1  Rows_examined: 100  Rows_affected: 0  Rows_read: 100 # Bytes_sent: 60 SET timestamp=1529940303; UPDATE foo SET bar = NOW() WHERE id = 1;"},
  {"row_case_3",
   "# User@Host: test[test] @ [127.0.0.1] " +
   "# Thread_id: 0  Schema: test  Last_errno: 0  Killed: 0 " +
   "# Query_time: 0.0  Lock_time: 0.0  Rows_sent: 1  Rows_examined: 0  Rows_affected: 0  Rows_read: 0 " +
   "# Bytes_sent: 0 " +
   "UPDATE foo SET bar = NOW() WHERE id = 1; " +
   "# Time: 000000 00:00:00 ",
   ""},
}

var keysValues = []struct{Key, Value string}{
  {"time",
   "180625 15:25:03"},
  {"user_host",
   "test[test] @ [127.0.0.1]"},
  {"Thread_id",
   "123456"},
  {"Query_time",
   "0.792864"},
  {"query",
   "SELECT count(*) AS total FROM foo WHERE att = 'bar'"},
  {"query_digest",
   "SELECT count(*) AS total FROM foo WHERE att = '?'"},
}

func TestRow(t *testing.T) {
  re := regexp.MustCompile(slow.ROW)

  for _, test := range slowLog {
    actual := re.FindString(test.Input)

    if test.Expected != actual {
      t.Error("test '" + test.ID + "' failed. actual = " + actual)
    }
  }
}

func TestKeysValues(t *testing.T) {
  re     := regexp.MustCompile(slow.KV)
  result := common.RegexpGetGroups(re, slowLog[0].Input)

  for _, test := range keysValues {
    switch test.Key {
    case "query_digest":
      result["query_digest"] = sql.Digest(result["query"])
    }

    if value, ok := result[test.Key]; ok {
      if value != test.Value {
        t.Error("test '" + test.Key + "' failed. actual = " + test.Value)
      }
    }
  }
}

func TestParser(t *testing.T) {
  channelTail   := make(chan string)
  channelParser := make(chan map[string]string)

  defer close(channelTail)

  expected := map[string]string{
    "_time":"1529940303",
    "bytes_sent":"6",
    "host_ip":"127.0.0.1",
    "host_name":"",
    "killed":"0",
    "last_errno":"0",
    "lock_time":"0.000160",
    "query":"SELECT count(*) AS total FROM foo WHERE att = 'bar';",
    "query_digest":"SELECT count(*) AS total FROM foo WHERE att = '?';",
    "query_time":"0.792864",
    "rows_affected":"0",
    "rows_examined":"100",
    "rows_read":"100",
    "rows_sent":"1",
    "schema":"test",
    "thread_id":"123456",
    "user_host":"test",
  }

  lines := []string{
    "# Time: 180625 15:25:03",
    "# User@Host: test[test] @ [127.0.0.1]",
    "# Thread_id: 123456  Schema: test  Last_errno: 0  Killed: 0",
    "# Query_time: 0.792864  Lock_time: 0.000160  Rows_sent: 1  Rows_examined: 100  Rows_affected: 0  Rows_read: 100",
    "# Bytes_sent: 60",
    "SET timestamp=1529940303;",
    "SELECT count(*) AS total FROM foo WHERE att = 'bar';",
  }

  go slow.Parser(config.MySQLSlowLog.LogPath, channelTail, channelParser)

  for _, line := range lines {
    channelTail <- line
  }

  found := <-channelParser

  for key, value := range found {
    if expected[key] != value {
      t.Errorf("Expected key '%s' value %s, found %s", key, expected[key], value)
    }
  }
}
