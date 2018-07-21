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
)

const USAGE = `zenit (%s) written by %s
Usage: %s [-version] [-help] <options> [args]
Options:

`

const HELP = `
Environment variables:
  - export DSN_CLICKHOUSE="http://127.0.0.1:8123/?database=zenit"
  - export DSN_MYSQL="root@tcp(127.0.0.1:3306)/"
  - export DSN_PROXYSQL="radminuser:radminpass@tcp(127.0.0.1:6032)/"
  - export SLACK_CHANNEL="alerts"
  - export SLACK_TOKEN="XXXXXXXXX/YYYYYYYYY/ZZZZZZZZZZZZZZZZZZZZZZZZ"
Examples:
  - zenit -run="mysqldump -h 127.0.0.1 -u root > /tmp/mysql.dump"
  - zenit -collect=os,mysql
  - zenit -parser-format=slowlog -parser-file=/tmp/slow.log
  - zenit -parser-format=slowlog -parser-file=/tmp/slow.log -daemonize
  - zenit -parser-format=slowlog -parser-file=/tmp/slow.log -stop
`

const COLLECT = `zenit (%s) written by %s
Usage: %s -collect=[argument1,argument2,...]
Available arguments:
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
Examples:
  - zenit -collect=mysql,os,proxysql
  - zenit -collect=os-cpu,os-disk,mysql-slave
`

var cmd string

func init() {
  cmd = os.Args[0]
}

func main() {
  fHelp         := flag.Bool("help", false, "Show this help.")
  fVersion      := flag.Bool("version", false, "Show version.")
  // status: check env variable conn and each daemonize process.
  fDaemonize    := flag.Bool("daemonize", false, "Fork to the background and detach from the shell.")
  fCollect      := flag.String("collect", "", "List of metrics to collect.")
  fParserFormat := flag.String("parser-format", "", "Parser log format.")
  fParserFile   := flag.String("parser-file", "", "File path to Tail to parse.")
  fStop         := flag.Bool("stop", false, "Stop daemon.")
  fRun          := flag.String("run", "", "Run bash command and wait to finish to notify via slack.")

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

    if len(*fRun) > 0 {
      command.Run(*fRun)
    }

    if *fDaemonize {
      daemonize.Start()
    } else if *fStop {
      daemonize.Stop()
    }

    if *fCollect == "" {
      fmt.Printf(COLLECT, config.VERSION, config.AUTHOR, cmd)
    } else if len(*fCollect) > 0 {
      collect.Run(strings.Split(*fCollect, ","))
      os.Exit(0)
    }

    if len(*fParserFormat) > 0 && len(*fParserFile) > 0 {
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
