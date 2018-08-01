package config

import (
  "log"
  "os"
  "time"

  "github.com/go-ini/ini"

  "gitlab.com/swapbyt3s/zenit/common"
)

const (
  AUTHOR  string = "Nicola Strappazzon C. <swapbyt3s@gmail.com>"
  VERSION string = "1.0.0"
)

type GeneralConfig struct {
  Hostname string        `ini:"hostname"`
  Interval time.Duration `ini:"interval"`
  LogFile  string        `ini:"log_file"`
  PIDFile  string        `ini:"pid_file"`
  Debug    bool          `ini:"debug"`
}

type OSConfig struct {
  CPU    bool `ini:"cpu"`
  Disk   bool `ini:"disk"`
  Limits bool `ini:"limits"`
  Mem    bool `ini:"mem"`
  Net    bool `ini:"net"`
}

type MySQLConfig struct {
  DSN       string `ini:"dsn"`
  Overflow  bool   `ini:"overflow"`
  Slave     bool   `ini:"slave"`
  Status    bool   `ini:"status"`
  Tables    bool   `ini:"tables"`
  Variables bool   `ini:"variables"`
  AuditLog  bool   `ini:"auditlog"`
  SlowLog   bool   `ini:"slowlog"`
}

type MySQLSlowLogConfig struct {
  LogPath    string `ini:"log_path"`
  BufferSize int    `ini:"buffer_size"`
}

type MySQLAuditLogConfig struct {
  Format     string `ini:"format"`
  LogPath    string `ini:"log_path"`
  BufferSize int    `ini:"buffer_size"`
}

type ProxySQLConfig struct {
  DSN         string `ini:"dsn"`
  QueryDigest bool   `ini:"query_digest"`
}

type ClickHouseConfig struct {
  DSN string `ini:"dsn"`
}

type PrometheusConfig struct {
  TextFile string `ini:"textfile"`
}

type ProcessConfig struct {
  PerconaToolKitKill           bool `ini:"pt_kill"`
  PerconaToolKitDeadlockLogger bool `ini:"pt_deadlock_logger"`
  PerconaToolKitSlaveDelay     bool `ini:"pt_slave_delay"`
}

var (
  ConfigFile string = "/etc/zenit/zenit.ini"
  IpAddress  string = ""

  ClickHouse    = new(ClickHouseConfig)
  General       = new(GeneralConfig)
  MySQL         = new(MySQLConfig)
  MySQLAuditLog = new(MySQLAuditLogConfig)
  MySQLSlowLog  = new(MySQLSlowLogConfig)
  OS            = new(OSConfig)
  Process       = new(ProcessConfig)
  Prometheus    = new(PrometheusConfig)
  ProxySQL      = new(ProxySQLConfig)
)

func init() {
  IpAddress = common.IpAddress()
}

func Load() {
  cfg, err := ini.Load(ConfigFile)
  if err != nil {
    cfg, err = ini.Load("zenit.ini")
    if err != nil {
      log.Printf("Fail to read config file: %s or %s", ConfigFile, "./zenit.ini")
      os.Exit(1)
    }
  }

  cfg.Section("clickhouse").MapTo(ClickHouse)
  cfg.Section("general").MapTo(General)
  cfg.Section("mysql").MapTo(MySQL)
  cfg.Section("mysql-auditlog").MapTo(MySQLAuditLog)
  cfg.Section("mysql-slowlog").MapTo(MySQLSlowLog)
  cfg.Section("os").MapTo(OS)
  cfg.Section("process").MapTo(Process)
  cfg.Section("prometheus").MapTo(Prometheus)
  cfg.Section("proxysql").MapTo(ProxySQL)

  // Validate interval variable
  // > 5 && <= 86400 (24h)

  // TODO: Define default values:
  //
  //HOSTNAME = cfg.Section("general").Key("hostname").String()
  //if len(HOSTNAME) == 0 {
  //  HOSTNAME = common.Hostname()
  //}
  // LOG_FILE            string = "/var/log/zenit.log"
  // PID_FILE            string = "/var/run/zenit.pid"
  // PROMETHEUS_TEXTFILE string = "/var/tmp/zenit.prom"
  // DSN_CLICKHOUSE      string = "http://127.0.0.1:8123/?database=zenit"
  // DSN_MYSQL           string = "root@tcp(127.0.0.1:3306)/"
  // DSN_PROXYSQL        string = "radminuser:radminpass@tcp(127.0.0.1:6032)/"
  // HOSTNAME            string = ""
  // INTERVAL            int    = 0
}
