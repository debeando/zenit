// TODO:
// - Convert this into module/package called "collect" because use for inputs and parsers.
// - Find any way to simplify this to make more dinamyc.
// - If not set any option, ignore and no enter in infinite loop.

package inputs

import (
  "log"
  "time"
  "sync"

  "github.com/swapbyt3s/zenit/common"
  "github.com/swapbyt3s/zenit/config"
  "github.com/swapbyt3s/zenit/plugins/accumulator"
  "github.com/swapbyt3s/zenit/plugins/inputs/mysql"
  "github.com/swapbyt3s/zenit/plugins/inputs/mysql/audit"
  "github.com/swapbyt3s/zenit/plugins/inputs/mysql/slow"
  "github.com/swapbyt3s/zenit/plugins/inputs/os"
  "github.com/swapbyt3s/zenit/plugins/inputs/process"
  "github.com/swapbyt3s/zenit/plugins/inputs/proxysql"
  "github.com/swapbyt3s/zenit/plugins/outputs/clickhouse"
  "github.com/swapbyt3s/zenit/plugins/outputs/prometheus"
)

func Gather() {
  var wg sync.WaitGroup

  wg.Add(2)

  go doCollectPlugins(&wg)
  go doCollectParsers(&wg)

  wg.Wait()
}

func doCollectPlugins(wg *sync.WaitGroup) {
  defer wg.Done()

  for {
    if config.OS.CPU {
      os.CPU()
    }
    if config.OS.Disk {
      os.Disk()
    }
    if config.OS.Mem {
      os.Mem()
    }
    if config.OS.Net {
      os.Net()
    }
    if config.OS.Limits {
      os.SysLimits()
    }
    if config.MySQL.Overflow && mysql.Check() {
      mysql.Overflow()
    }
    if config.MySQL.Slave && mysql.Check() {
      mysql.Slave()
    }
    if config.MySQL.Status && mysql.Check() {
      mysql.Status()
    }
    if config.MySQL.Tables && mysql.Check() {
      mysql.Tables()
    }
    if config.MySQL.Variables && mysql.Check() {
      mysql.Variables()
    }
    if config.ProxySQL.QueryDigest && proxysql.Check() {
      proxysql.QueryDigest()
    }
    if config.Process.PerconaToolKitKill {
      process.PerconaToolKitKill()
    }
    if config.Process.PerconaToolKitDeadlockLogger {
      process.PerconaToolKitDeadlockLogger()
    }
    if config.Process.PerconaToolKitSlaveDelay {
      process.PerconaToolKitSlaveDelay()
    }
    prometheus.Run()
    accumulator.Load().Reset()
    time.Sleep(config.General.Interval * time.Second)
  }
}

func doCollectParsers(wg *sync.WaitGroup) {
  defer wg.Done()

  if config.MySQL.AuditLog {
    if config.General.Debug {
      log.Println("D! - Load MySQL AuditLog")
    }

    if ! clickhouse.Check() {
      log.Println("E! - AuditLog require active connection to ClickHouse.")
    }

    if config.MySQLAuditLog.Format == "xml-old" {
      channel_tail   := make(chan string)
      channel_parser := make(chan map[string]string)
      channel_data   := make(chan map[string]string)

      event := clickhouse.Event{
        Type: "AuditLog",
        Schema: "zenit",
        Table: "mysql_audit_log",
        Size: config.MySQLAuditLog.BufferSize,
        Timeout: config.MySQLAuditLog.BufferTimeOut,
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

      go common.Tail(config.MySQLAuditLog.LogPath, channel_tail)
      go audit.Parser(config.MySQLAuditLog.LogPath, channel_tail, channel_parser)
      go clickhouse.Run(event, channel_data, config.MySQLSlowLog.BufferTimeOut)

      go func() {
        for event := range channel_parser {
          channel_data <- event
        }
      }()
    }
  }

  if config.MySQL.SlowLog {
    if config.General.Debug {
      log.Println("D! - Load MySQL SlowLog")
    }

    if ! clickhouse.Check() {
      log.Println("E! - SlowLog require active connection to ClickHouse.")
    }

    channel_tail   := make(chan string)
    channel_parser := make(chan map[string]string)
    channel_data   := make(chan map[string]string)

    event := clickhouse.Event{
      Type: "SlowLog",
      Schema: "zenit",
      Table: "mysql_slow_log",
      Size: config.MySQLSlowLog.BufferSize,
      Timeout: config.MySQLSlowLog.BufferTimeOut,
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

    go common.Tail(config.MySQLSlowLog.LogPath, channel_tail)
    go slow.Parser(config.MySQLSlowLog.LogPath, channel_tail, channel_parser)
    go clickhouse.Run(event, channel_data, config.MySQLSlowLog.BufferTimeOut)

    go func() {
      for event := range channel_parser {
        channel_data <- event
      }
    }()
  }
}
