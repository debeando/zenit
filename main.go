package main

import (
  "flag"
  "fmt"
  "os"
  "strings"
  "gitlab.com/swapbyt3s/zenit/config"
  "gitlab.com/swapbyt3s/zenit/collect"
  "gitlab.com/swapbyt3s/zenit/command"
)

const USAGE = "zenit (%s) written by Nicola Strappazzon C. <swapbyt3s@gmail.com>\nUsage: %s <command>\n"

func main() {
  fHelp    := flag.Bool("help",      false, "Show this help.")
  fVersion := flag.Bool("version",   false, "Show version.")
  fCollect := flag.String("collect",    "", "List of metrics to collect.")
  fRun     := flag.String("run",        "", "Run bash command and wait to finish to notify via slack.")
  // -collect-mysql {status,variables,slave status}
  // check is open port from any
  // -output-prometheus
  // -output-influxdb
  // -percona-skip-replication-error
  // -percona-eta-catchup
  // -collect-os {network,swap,disk,iops}
  // -collect-pgbouncer ?
  // -collect-postgresql ?
  // -collect-mongodb ?

  flag.Parse()

  if len(os.Args) == 1 {
    help()
  } else if len(os.Args) >= 2 {
    if *fHelp {
      help()
    } else if *fVersion {
      fmt.Printf("%s\n", config.VERSION)
    } else if len(*fCollect) > 0 {
      collect.Run(strings.Split(*fCollect, ","))
    } else if len(*fRun) > 0 {
      command.Run(*fRun)
    } else {
      fmt.Printf("%q is not valid command.\n", os.Args[1])
      help()
    }
  } else {
    help()
  }
}

func help() {
  fmt.Printf(USAGE, config.VERSION, os.Args[0])
  flag.PrintDefaults()
  os.Exit(1)
}
