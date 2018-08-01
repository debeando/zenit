// TODO:
// - Convert this into module/package called "collect" because use for inputs and parsers.
// - Find any way to simplify this to make more dinamyc.

package inputs

import (
  "log"
  "time"
  "sync"

  "gitlab.com/swapbyt3s/zenit/common"
  "gitlab.com/swapbyt3s/zenit/config"
  "gitlab.com/swapbyt3s/zenit/plugins/accumulator"
  "gitlab.com/swapbyt3s/zenit/plugins/inputs/mysql"
  "gitlab.com/swapbyt3s/zenit/plugins/inputs/mysql/audit"
  "gitlab.com/swapbyt3s/zenit/plugins/inputs/mysql/slow"
  "gitlab.com/swapbyt3s/zenit/plugins/inputs/os"
  "gitlab.com/swapbyt3s/zenit/plugins/inputs/process"
  "gitlab.com/swapbyt3s/zenit/plugins/inputs/proxysql"
  "gitlab.com/swapbyt3s/zenit/plugins/outputs/clickhouse"
  "gitlab.com/swapbyt3s/zenit/plugins/outputs/prometheus"
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
    if ! clickhouse.Check() {
      log.Println("E! - AuditLog require active connection to ClickHouse.")
    }

    if config.MySQLAuditLog.Format == "xml-old" {
      channel_tail   := make(chan string)
      channel_parser := make(chan map[string]string)
      channel_event  := make(chan map[string]string)

      go common.Tail(config.MySQLAuditLog.LogPath, channel_tail)
      go audit.Parser(config.MySQLAuditLog.LogPath, channel_tail, channel_parser)
      go clickhouse.SendMySQLAuditLog(channel_event)

      for event := range channel_parser {
        channel_event <- event
      }
    }
  }

  if config.MySQL.SlowLog {
    if ! clickhouse.Check() {
      log.Println("E! - SlowLog require active connection to ClickHouse.")
    }

    channel_tail   := make(chan string)
    channel_parser := make(chan map[string]string)
    channel_event  := make(chan map[string]string)

    go common.Tail(config.MySQLSlowLog.LogPath, channel_tail)
    go slow.Parser(config.MySQLSlowLog.LogPath, channel_tail, channel_parser)
    go clickhouse.SendMySQLSlowLog(channel_event)

    for event := range channel_parser {
      channel_event <- event
    }
  }
}
