package slow_test

import (
	//"regexp"
	//"strings"
	"testing"

	//"github.com/swapbyt3s/zenit/common"
	//"github.com/swapbyt3s/zenit/common/sql"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/inputs/mysql/slow"
)

var slowLog = []struct{
	ID         string
	Input    []string
	Row        string
	KeyValue []struct{ Key, Value string }
	Parsed   []struct{ Key, Value string }
}{
	{
		"row_case_1",
		[]string{
			"# Time: 180625 15:25:03",
			"# User@Host: test[test] @ [127.0.0.1]",
			"# Thread_id: 123456  Schema: test  Last_errno: 0  Killed: 0",
			"# Query_time: 0.12345  Lock_time: 0.000123  Rows_sent: 1  Rows_examined: 100  Rows_affected: 0  Rows_read: 100",
			"# Bytes_sent: 60 ",
			"SET timestamp=1529940303;",
			"SELECT count(*) AS total FROM foo WHERE att = 'bar';",
		},
		"# Time: 180625 15:25:03 # User@Host: test[test] @ [127.0.0.1] # Thread_id: 123456  Schema: test  Last_errno: 0  Killed: 0 # Query_time: 0.12345  Lock_time: 0.000123  Rows_sent: 1  Rows_examined: 100  Rows_affected: 0  Rows_read: 100 # Bytes_sent: 60  SET timestamp=1529940303; SELECT count(*) AS total FROM foo WHERE att = 'bar';",
		[]struct{ Key, Value string }{
			{"time",         "180625 15:25:03"},
			{"user_host",    "test[test] @ [127.0.0.1]"},
			{"Thread_id",    "123456"},
			{"Query_time",   "0.12345"},
			{"Bytes_sent",   "60"},
			{"query",        "SELECT count(*) AS total FROM foo WHERE att = 'bar'"},
			{"query_digest", "SELECT count(*) AS total FROM foo WHERE att = '?'"},
		},
		[]struct{ Key, Value string }{
			{"_time",        "1529940303"},
			{"bytes_sent",   "60"},
			{"host_ip",      "127.0.0.1"},
			{"host_name",    ""},
			{"killed",       "0"},
			{"last_errno",   "0"},
			{"lock_time",    "0.000123"},
			{"query",        `SELECT count(*) AS total FROM foo WHERE att = \'bar\';`},
			{"query_digest", `SELECT count(*) AS total FROM foo WHERE att = \'?\';`},
			{"query_time",   "0.12345"},
			{"rows_read",    "100"},
			{"rows_sent",    "1"},
			{"schema",       "test"},
			{"thread_id",    "123456"},
			{"user_host",    "test"},
		},
	},{
		"row_case_2",
		[]string{
			"SELECT * FROM foo;",
			"# User@Host: test[test] @ [127.0.0.1]",
			"# Thread_id: 123456  Schema: test  Last_errno: 0  Killed: 0",
			"# query_time: 0.12345  Lock_time: 0.000123  Rows_sent: 1  Rows_examined: 100  Rows_affected: 0  Rows_read: 100",
			"# Bytes_sent: 64 ",
			"SET timestamp=1529940303;",
			"UPDATE foo SET bar = NOW() WHERE id = 1;",
		},
		"# User@Host: test[test] @ [127.0.0.1] # Thread_id: 123456  Schema: test  Last_errno: 0  Killed: 0 # Query_time: 0.12345  Lock_time: 0.000123  Rows_sent: 1  Rows_examined: 100  Rows_affected: 0  Rows_read: 100 # Bytes_sent: 64  SET timestamp=1529940303; UPDATE foo SET bar = NOW() WHERE id = 1;",
		[]struct{ Key, Value string }{
			{"Rows_sent",     "1"},
			{"Rows_examined", "100"},
		},
		[]struct{ Key, Value string }{
			{"_time",        "1529940303"},
			{"bytes_sent",   "64"},
			{"host_ip",      "127.0.0.1"},
			{"host_name",    ""},
			{"killed",       "0"},
			{"last_errno",   "0"},
			{"lock_time",    "0.000123"},
			{"query",        `UPDATE foo SET bar = NOW() WHERE id = 1;`},
			{"query_digest", `UPDATE foo SET bar = NOW() WHERE id = ?;`},
			{"query_time",   "0.12345"},
			{"rows_read",    "100"},
			{"rows_sent",    "1"},
			{"schema",       "test"},
			{"thread_id",    "123456"},
			{"user_host",    "test"},
		},
	},
}

//func TestRow(t *testing.T) {
//	re := regexp.MustCompile(slow.ROW)
//
//	for _, test := range slowLog {
//		actual := re.FindString(strings.Join(test.Input, " "))
//
//		if test.Row != actual {
//			t.Error("test '" + test.ID + "' failed. actual = " + actual)
//		}
//	}
//}

//func TestKeysValues(t *testing.T) {
//	re := regexp.MustCompile(slow.KV)
//	result := common.RegexpGetGroups(re, strings.Join(slowLog[0].Input, " "))
//
//	for _, test := range slowLog[0].KeyValue {
//		switch test.Key {
//		case "query_digest":
//			result["query_digest"] = sql.Digest(result["query"])
//		}
//
//		if value, ok := result[test.Key]; ok {
//			if value != test.Value {
//				t.Error("test '" + test.Key + "' failed. " + value + " != " + test.Value)
//			}
//		}
//	}
//}

func TestParser(t *testing.T) {
	channelTail := make(chan string)
	channelParser := make(chan map[string]string)

	defer close(channelTail)

	go slow.Parser(config.MySQLSlowLog.LogPath, channelTail, channelParser)

	for _, test := range slowLog {
		t.Logf("--> TEST: %#v\n", test)

		//for _, line := range test.Input {
		//	//t.Logf("--> LINE: %#v\n", line)
		//	channelTail <- line
		//}

		//result := <-channelParser

		//t.Logf("--> RES: %#v\n", result)

		//for parseKey, parseValue := range result {
		//	t.Logf("--> %s, %#v\n", parseKey, parseValue)
		//	for _, testParsed := range test.Parsed {
		//		if parseKey == testParsed.Key && parseValue != testParsed.Value {
		//			t.Errorf("Expected key '%s' bad value %s, found %s", parseKey, parseValue, testParsed.Value)
		//		}
		//	}
		//}
	}
}
