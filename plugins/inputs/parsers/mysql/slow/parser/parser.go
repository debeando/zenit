package parser

import (
	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/common/mysql"
	"github.com/swapbyt3s/zenit/common/sql"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/outputs/clickhouse"
)

type parser struct {}

var instance *parser = nil

func New() *parser {
	if instance == nil {
		instance = &parser{}
	}
	return instance
}

func (p *parser) Start() {
	channel_tail := make(chan string)
	channel_parser := make(chan map[string]string)
	channel_data := make(chan map[string]string)

	event := &clickhouse.Event{
		Type:    "SlowLog",
		Schema:  "zenit",
		Table:   "mysql_slow_log",
		Size:    config.File.Parser.MySQL.SlowLog.BufferSize,
		Timeout: config.File.Parser.MySQL.SlowLog.BufferTimeOut,
		Wildcard: map[string]string{
			"_time":         "'%s'",
			"bytes_sent":    "%s",
			"host_ip":       "IPv4StringToNum('%s')",
			"host_name":     "'%s'",
			"killed":        "%s",
			"last_errno":    "%s",
			"lock_time":     "%s",
			"query":         "'%s'",
			"query_digest":  "'%s'",
			"query_time":    "%s",
			"rows_affected": "%s",
			"rows_examined": "%s",
			"rows_read":     "%s",
			"rows_sent":     "%s",
			"schema":        "'%s'",
			"thread_id":     "%s",
			"user_host":     "'%s'",
		},
		Values: []map[string]string{},
	}

	go common.Tail(config.File.Parser.MySQL.SlowLog.LogPath, channel_tail)
	go p.Parser(config.File.Parser.MySQL.SlowLog.LogPath, channel_tail, channel_parser)
	go clickhouse.Send(event, channel_data)

	go func() {
		for channel_event := range channel_parser {
			channel_data <- channel_event
		}
	}()
}

func (p *parser) Parser(path string, in <-chan string, out chan<- map[string]string) {
	channelTail := make(chan string)
	channelEvent := make(chan string)

	go slow.Event(channelTail, channelEvent)

	go func() {
		defer close(channelTail)
		for line := range in {
			channelTail <- line
		}
	}()

	go func() {
		defer close(channelEvent)
		for event := range channelEvent {
			result := slow.Properties(event)

			if common.KeyInMap("user_host", result) {
				result["user_host"] = mysql.ClearUser(result["user_host"])
			}

			if common.KeyInMap("query", result) {
				result["query_digest"] = sql.Digest(result["query"])
			}

			result["_time"] = result["timestamp"]
			result["host_ip"] = config.File.IPAddress
			result["host_name"] = config.File.General.Hostname
			result["query"] = common.Escape(result["query"])
			result["query_digest"] = common.Escape(result["query_digest"])

			// Remove unnused key:
			delete(result, "time")
			delete(result, "timestamp")
			delete(result, "qc_hit")

			out <- result
		}
	}()
}
