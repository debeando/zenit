package config

import (
  "log"
  "os"
  "time"

  "github.com/go-ini/ini"

  "github.com/swapbyt3s/zenit/common"
)

const (
  AUTHOR  string = "Nicola Strappazzon C. <swapbyt3s@gmail.com>"
  VERSION string = "1.0.1"
)

// Struct to save in memory config settings from [general] section.
type GeneralConfig struct {
  Hostname string        `ini:"hostname"`
  Interval time.Duration `ini:"interval"`
  LogFile  string        `ini:"log_file"`
  PIDFile  string        `ini:"pid_file"`
  Debug    bool          `ini:"debug"`
}

// Struct to save in memory config settings from [os] section.
type OSConfig struct {
  CPU    bool `ini:"cpu"`
  Disk   bool `ini:"disk"`
  Limits bool `ini:"limits"`
  Mem    bool `ini:"mem"`
  Net    bool `ini:"net"`
}

// Struct to save in memory config settings from [mysql] section.
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

// Struct to save in memory config settings from [mysql-audit] section.
type MySQLAuditLogConfig struct {
  Format     string `ini:"format"`
  LogPath    string `ini:"log_path"`
  BufferSize int    `ini:"buffer_size"`
}

// Struct to save in memory config settings from [mysql-slowlog] section.
type MySQLSlowLogConfig struct {
  LogPath    string `ini:"log_path"`
  BufferSize int    `ini:"buffer_size"`
}

// Struct to save in memory config settings from [proxysql] section.
type ProxySQLConfig struct {
  DSN         string `ini:"dsn"`
  QueryDigest bool   `ini:"query_digest"`
}

// Struct to save in memory config settings from [clickhouse] section.
type ClickHouseConfig struct {
  DSN string `ini:"dsn"`
}

// Struct to save in memory config settings from [prometheus] section.
type PrometheusConfig struct {
  TextFile string `ini:"textfile"`
}

// Struct to save in memory config settings from [process] section.
type ProcessConfig struct {
  PerconaToolKitKill           bool `ini:"pt_kill"`
  PerconaToolKitDeadlockLogger bool `ini:"pt_deadlock_logger"`
  PerconaToolKitSlaveDelay     bool `ini:"pt_slave_delay"`
}

// Define default variables and initialize structs.
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

// Init does any initialization necessary for the module.
func init() {
  IpAddress = common.IpAddress()
}

// Loading settings from config file and set into struct.
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
}

// Check minimun config settings and set default values to start.
func SanityCheck() {
  if General.Interval < 5 {
    log.Println("W! Config - general.Interval: Use positive value, and minimun start from 5 seconds.")
    log.Println("W! Config - general.Interval: Using default 30 seconds.")
    General.Interval = 30
  }

  if len(General.Hostname) == 0 {
    log.Println("W! Config - general.Hostname: Custom value is not set, using current.")
    General.Hostname = common.Hostname()
  }

  if len(General.LogFile) == 0 {
    log.Println("W! Config - general.LogFile: Custom value is not set, using default /var/log/zenit.log")
    General.LogFile = "/var/log/zenit.log"
  }

  if len(General.PIDFile) == 0 {
    log.Println("W! Config - general.PIDFile: Custom value is not set, using default /var/run/zenit.pid")
    General.PIDFile = "/var/run/zenit.pid"
  }
}
