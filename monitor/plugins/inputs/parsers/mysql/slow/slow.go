package slow

import (
	"sync"

	"zenit/config"
	"zenit/monitor/plugins/inputs"
	"zenit/monitor/plugins/lists/metrics"
	"zenit/monitor/plugins/outputs/clickhouse"

	"github.com/debeando/go-common/log"
	"github.com/debeando/go-common/maps"
	"github.com/debeando/go-common/mysql"
	"github.com/debeando/go-common/mysql/sql"
	"github.com/debeando/go-common/mysql/sql/parser/slow"
	"github.com/debeando/go-common/strings"
	"github.com/debeando/go-common/tail"
)

type Plugin struct {
	Config *config.Config
	Name   string
}

var (
	instance *Plugin
	once     sync.Once
)

func (p *Plugin) Collect(name string, cnf *config.Config, mtc *metrics.Items) {
	once.Do(func() {
		if instance == nil {
			instance = &Plugin{}

			p.Config = cnf
			p.Name = name
			p.Load()
		}
	})
}

func (p *Plugin) Load() {
	defer func() {
		if err := recover(); err != nil {
			log.ErrorWithFields(p.Name, log.Fields{"error": err})
		}
	}()

	if p.Config.Parser.MySQL.SlowLog.Enable {
		log.DebugWithFields(p.Name, log.Fields{"slow_log_path": p.Config.Parser.MySQL.SlowLog.LogPath})

		if !clickhouse.Check() {
			log.ErrorWithFields(p.Name, log.Fields{"error": "SlowLog require active connection to ClickHouse."})
		}

		channel_tail := make(chan string)
		channel_parser := make(chan map[string]string)
		channel_data := make(chan map[string]string)

		event := &clickhouse.Event{
			Type:    "SlowLog",
			Schema:  "zenit",
			Table:   "mysql_slow_log",
			Size:    p.Config.Parser.MySQL.SlowLog.BufferSize,
			Timeout: p.Config.Parser.MySQL.SlowLog.BufferTimeOut,
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

		go tail.Tail(p.Config.Parser.MySQL.SlowLog.LogPath, channel_tail)
		go p.Parser(channel_tail, channel_parser)
		go clickhouse.Send(event, channel_data)

		go func() {
			for channel_event := range channel_parser {
				channel_data <- channel_event
			}
		}()
	}
}

func (p *Plugin) Parser(in <-chan string, out chan<- map[string]string) {
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

			if maps.In("user_host", result) {
				result["user_host"] = mysql.ClearUser(result["user_host"])
			}

			if maps.In("query", result) {
				result["query_digest"] = sql.Digest(result["query"])
			}

			result["_time"] = result["timestamp"]
			result["host_ip"] = p.Config.IPAddress
			result["host_name"] = p.Config.General.Hostname
			result["query"] = strings.Escape(result["query"])
			result["query_digest"] = strings.Escape(result["query_digest"])

			// Remove unnused key:
			delete(result, "time")
			delete(result, "timestamp")
			delete(result, "qc_hit")

			out <- result
		}
	}()
}

func init() {
	inputs.Add("InputMySQLSlowLog", func() inputs.Input { return &Plugin{} })
}
