// TODO:
// - Convert this into module/package called "collect" because use for inputs and parsers.
// - If not set any option, ignore and no enter in infinite loop.

package inputs

import (
	"sync"
	"time"

	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/inputs/mysql"
	"github.com/swapbyt3s/zenit/plugins/inputs/mysql/audit"
	"github.com/swapbyt3s/zenit/plugins/inputs/mysql/indexes"
	"github.com/swapbyt3s/zenit/plugins/inputs/mysql/overflow"
	"github.com/swapbyt3s/zenit/plugins/inputs/mysql/slave"
	"github.com/swapbyt3s/zenit/plugins/inputs/mysql/slow"
	"github.com/swapbyt3s/zenit/plugins/inputs/mysql/status"
	"github.com/swapbyt3s/zenit/plugins/inputs/mysql/tables"
	"github.com/swapbyt3s/zenit/plugins/inputs/mysql/variables"
	"github.com/swapbyt3s/zenit/plugins/inputs/os/cpu"
	"github.com/swapbyt3s/zenit/plugins/inputs/os/disk"
	"github.com/swapbyt3s/zenit/plugins/inputs/os/mem"
	"github.com/swapbyt3s/zenit/plugins/inputs/os/net"
	"github.com/swapbyt3s/zenit/plugins/inputs/os/sys"
	"github.com/swapbyt3s/zenit/plugins/inputs/process"
	"github.com/swapbyt3s/zenit/plugins/inputs/proxysql"
	"github.com/swapbyt3s/zenit/plugins/inputs/proxysql/commands"
	"github.com/swapbyt3s/zenit/plugins/inputs/proxysql/pool"
	"github.com/swapbyt3s/zenit/plugins/inputs/proxysql/queries"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
	"github.com/swapbyt3s/zenit/plugins/outputs/clickhouse"
	"github.com/swapbyt3s/zenit/plugins/outputs/prometheus"
	"github.com/swapbyt3s/zenit/plugins/outputs/slack"
)

func Plugins(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		// Flush old metrics:
		metrics.Load().Reset()

		// Collect by OS:
		if config.File.OS.Inputs.CPU {
			cpu.Collect()
		}
		if config.File.OS.Inputs.Disk {
			disk.Collect()
		}
		if config.File.OS.Inputs.Mem {
			mem.Collect()
		}
		if config.File.OS.Inputs.Net {
			net.Collect()
		}
		if config.File.OS.Inputs.Limits {
			sys.Collect()
		}

		// Collect by MySQL:
		if mysql.Check() {
			if config.File.MySQL.Inputs.Indexes {
				indexes.Collect()
			}
			if config.File.MySQL.Inputs.Overflow {
				overflow.Collect()
			}
			if config.File.MySQL.Inputs.Slave {
				slave.Collect()
			}
			if config.File.MySQL.Inputs.Status {
				status.Collect()
			}
			if config.File.MySQL.Inputs.Tables {
				tables.Collect()
			}
			if config.File.MySQL.Inputs.Variables {
				variables.Collect()
			}
		}

		// Collect by ProxySQL:
		if proxysql.Check() {
			if config.File.ProxySQL.Inputs.Commands {
				commands.Collect()
			}
			if config.File.ProxySQL.Inputs.Pool {
				pool.Collect()
			}
			if config.File.ProxySQL.Inputs.Queries {
				queries.Collect()
			}
		}

		// Collect by process:
		if config.File.Process.Inputs.PerconaToolKitKill {
			process.PerconaToolKitKill()
		}
		if config.File.Process.Inputs.PerconaToolKitDeadlockLogger {
			process.PerconaToolKitDeadlockLogger()
		}
		if config.File.Process.Inputs.PerconaToolKitSlaveDelay {
			process.PerconaToolKitSlaveDelay()
		}
		if config.File.Process.Inputs.PerconaToolKitOnlineSchemaChange {
			process.PerconaToolKitOnlineSchemaChange()
		}
		if config.File.Process.Inputs.PerconaXtraBackup {
			process.PerconaXtraBackup()
		}
		if config.File.Prometheus.Enable {
			prometheus.Run()
		}
		if config.File.Slack.Enable {
			slack.Run()
		}

		// Wait loop:
		time.Sleep(config.File.General.Interval * time.Second)
	}
}

func Parsers(wg *sync.WaitGroup) {
	defer wg.Done()

	if config.File.MySQL.Inputs.AuditLog.Enable {
		if config.File.General.Debug {
			log.Debug("Load MySQL AuditLog")
			log.Debug("Read MySQL AuditLog: " + config.File.MySQL.Inputs.AuditLog.LogPath)
		}

		if !clickhouse.Check() {
			log.Error("AuditLog require active connection to ClickHouse.")
		}

		if config.File.MySQL.Inputs.AuditLog.Format == "xml-old" {
			channel_tail := make(chan string)
			channel_parser := make(chan map[string]string)
			channel_data := make(chan map[string]string)

			event := &clickhouse.Event{
				Type:    "AuditLog",
				Schema:  "zenit",
				Table:   "mysql_audit_log",
				Size:    config.File.MySQL.Inputs.AuditLog.BufferSize,
				Timeout: config.File.MySQL.Inputs.AuditLog.BufferTimeOut,
				Wildcard: map[string]string{
					"_time":          "'%s'",
					"command_class":  "'%s'",
					"connection_id":  "%s",
					"db":             "'%s'",
					"host":           "'%s'",
					"host_ip":        "IPv4StringToNum('%s')",
					"host_name":      "'%s'",
					"ip":             "'%s'",
					"name":           "'%s'",
					"os_login":       "'%s'",
					"os_user":        "'%s'",
					"priv_user":      "'%s'",
					"proxy_user":     "'%s'",
					"record":         "'%s'",
					"sqltext":        "'%s'",
					"sqltext_digest": "'%s'",
					"status":         "%s",
					"user":           "'%s'",
				},
				Values: []map[string]string{},
			}

			go common.Tail(config.File.MySQL.Inputs.AuditLog.LogPath, channel_tail)
			go audit.Parser(config.File.MySQL.Inputs.AuditLog.LogPath, channel_tail, channel_parser)
			go clickhouse.Send(event, channel_data)

			go func() {
				for channel_event := range channel_parser {
					channel_data <- channel_event
				}
			}()
		}
	}

	if config.File.MySQL.Inputs.SlowLog.Enable {
		if config.File.General.Debug {
			log.Debug("Load MySQL SlowLog")
			log.Debug("Read MySQL SlowLog: " + config.File.MySQL.Inputs.SlowLog.LogPath)
		}

		if !clickhouse.Check() {
			log.Error("SlowLog require active connection to ClickHouse.")
		}

		channel_tail := make(chan string)
		channel_parser := make(chan map[string]string)
		channel_data := make(chan map[string]string)

		event := &clickhouse.Event{
			Type:    "SlowLog",
			Schema:  "zenit",
			Table:   "mysql_slow_log",
			Size:    config.File.MySQL.Inputs.SlowLog.BufferSize,
			Timeout: config.File.MySQL.Inputs.SlowLog.BufferTimeOut,
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

		go common.Tail(config.File.MySQL.Inputs.SlowLog.LogPath, channel_tail)
		go slow.Parser(config.File.MySQL.Inputs.SlowLog.LogPath, channel_tail, channel_parser)
		go clickhouse.Send(event, channel_data)

		go func() {
			for channel_event := range channel_parser {
				channel_data <- channel_event
			}
		}()
	}
}
