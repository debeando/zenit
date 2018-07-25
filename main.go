package main

import (
  "flag"
  "fmt"
  "os"
  "strings"
  "gitlab.com/swapbyt3s/zenit/collect"
  "gitlab.com/swapbyt3s/zenit/command"
  "gitlab.com/swapbyt3s/zenit/config"
  "gitlab.com/swapbyt3s/zenit/daemonize"
  "gitlab.com/swapbyt3s/zenit/status"
)

const USAGE = `zenit (%s) written by %s
Usage: %s [-version] [-help] <options> [args]
Options:

`

const HELP = `
Available arguments for collect:
  - mysql
  - mysql-overflow
  - mysql-slave
  - mysql-status
  - mysql-table
  - mysql-variables
  - os
  - os-cpu
  - os-disk
  - os-limits
  - os-mem
  - os-net
  - percona-process
  - proxysql

Environment variables:
  - export DSN_CLICKHOUSE="http://127.0.0.1:8123/?database=zenit"
  - export DSN_MYSQL="root@tcp(127.0.0.1:3306)/"
  - export DSN_PROXYSQL="radminuser:radminpass@tcp(127.0.0.1:6032)/"
  - export SLACK_CHANNEL="alerts"
  - export SLACK_TOKEN="XXXXXXXXX/YYYYYYYYY/ZZZZZZZZZZZZZZZZZZZZZZZZ"

Examples:
  - zenit -status
  - zenit -run="mysqldump -h 127.0.0.1 -u root > /tmp/mysql.dump"
  - zenit -collect=os,mysql
  - zenit -parser-format=slowlog -parser-file=/tmp/slow.log
  - zenit -parser-format=slowlog -parser-file=/tmp/slow.log -daemonize
  - zenit -parser-format=slowlog -parser-file=/tmp/slow.log -stop

`

var cmd string

func init() {
  cmd = os.Args[0]
}

func main() {
  fHelp         := flag.Bool("help", false, "Show this help.")
  fVersion      := flag.Bool("version", false, "Show version.")
  fCollect      := flag.String("collect", "", "List of metrics to collect.")
  fDaemonize    := flag.Bool("daemonize", false, "Fork to the background and detach from the shell.")
  fParserFile   := flag.String("parser-file", "", "File path to Tail to parse.")
  fParserFormat := flag.String("parser-format", "", "Parser log format.")
  fHostname     := flag.String("hostname", "", "Rename the hostname.")
  fRun          := flag.String("run", "", "Run bash command and wait to finish to notify via slack.")
  fStatus       := flag.Bool("status", false, "Status for each environment variable and own process.")
  fStop         := flag.Bool("stop", false, "Stop daemon.")

  flag.Parse()

  if len(os.Args) == 1 {
    help()
  } else if len(os.Args) >= 2 {
    if *fHelp {
      help()
    }

    if *fVersion {
      fmt.Printf("%s\n", config.VERSION)
      os.Exit(0)
    }

    if *fStatus {
      status.Run()
      os.Exit(0)
    }

    if len(*fRun) > 0 {
      command.Run(*fRun)
    }

    if *fDaemonize {
      daemonize.Start()
    } else if *fStop {
      daemonize.Stop()
    }

    if len(*fHostname) > 0 {
      config.HOSTNAME = *fHostname
    }

    if len(*fCollect) > 0 && len(*fParserFormat) == 0 && len(*fParserFile) == 0 {
      collect.Run(strings.Split(*fCollect, ","))
      os.Exit(0)
    }

    if len(*fCollect) == 0 && len(*fParserFormat) > 0 && len(*fParserFile) > 0 {
      collect.Parser(*fParserFormat, *fParserFile)
    }
  } else {
    help()
  }
}

func help() {
  fmt.Printf(USAGE, config.VERSION, config.AUTHOR, cmd)
  flag.PrintDefaults()
  fmt.Printf(HELP)
  os.Exit(1)
}
