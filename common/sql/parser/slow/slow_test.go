package slow_test

import (
        "fmt"
        "strings"
        "testing"

        "github.com/swapbyt3s/zenit/common"
        "github.com/swapbyt3s/zenit/common/sql/parser/slow"
)

var event = []string{
                "# Time: 180625 15:25:03",
                "# User@Host: test[test] @ [127.0.0.1]",
                "# Thread_id: 123456  Schema: test  Last_errno: 0  Killed: 0",
                "# Query_time: 0.1  Lock_time: 0.1  Rows_sent: 1",
                "# Bytes_sent: 60",
                "SET timestamp=1529940303;",
                "SELECT count(*) FROM foo;",
        }

func TestEvent(t *testing.T) {
        channelTail := make(chan string)
        channelEvent := make(chan string)

        defer close(channelTail)
        defer close(channelEvent)

        go slow.Event(channelTail, channelEvent)

        for _, e := range event {
                channelTail <- e
        }

        result := <- channelEvent
        expected := strings.Join(event[0:7], "\n")

        if result != expected {
                t.Error("Expected: " + expected)
        }
}

func TestProperty(t *testing.T) {
        result := slow.Properties(strings.Join(event, "\n"))
        expected := map[string]string{
                "thread_id": "123456",
                "schema": "test",
                "killed": "0",
                "bytes_sent": "60",
                "timestamp": "1529940303",
                "time": "180625 15:25:03",
                "user_host": "test[test] _ [127.0.0.1]",
                "last_errno": "0",
                "query_time": "0.1",
                "lock_time": "0.1",
                "rows_sent": "1",
                "query": "SELECT count(*) FROM foo;",
        }

        if ! common.ComparteMapString(result, expected) {
                t.Error("Expected:", fmt.Sprintf("%#v", expected))
        }
}
