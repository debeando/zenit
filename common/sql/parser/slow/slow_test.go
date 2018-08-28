package slow_test

import (
        "strings"
        "testing"

        "github.com/swapbyt3s/zenit/common/sql/parser/slow"
)

var event = []string{
                "SELECT count(*) FROM bar;",
                "# Time: 180625 15:25:03",
                "# User@Host: test[test] @ [127.0.0.1]",
                "# Thread_id: 123456  Schema: test  Last_errno: 0  Killed: 0",
                "# Query_time: 0.12345  Lock_time: 0.000123  Rows_sent: 1  Rows_examined: 100  Rows_affected: 0  Rows_read: 100",
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
        // expected := strings.Join(event[1:8], "\n")

        //if result != expected {
        //        t.Error("Result: " + result)
        //        t.Error("Expected: " + expected)
        //}

        // Case two:
        //for _, e := range event[1:8] {
        //        channelTail <- e
        //}

        //result = <- channelEvent
//        //expected = strings.Join(event[1:8], "\n")
//
        //if result != expected {
        //        t.Error("Result: " + result)
        //        t.Error("Expected: " + expected)
        //}


        t.Logf("--> Event: %#v\n", result)
}

func TestProperty(t *testing.T) {
        p := slow.Properties(strings.Join(event, "\n"))

        t.Logf("%#v\n", p)
}
