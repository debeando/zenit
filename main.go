package main

import (
  "flag"
  "fmt"
  "os"
  "gitlab.com/swapbyt3s/zenit/config"
  "gitlab.com/swapbyt3s/zenit/collect"
)

func main() {
  flg_help     := flag.Bool("help",             false, "Show this help.")
  flg_version  := flag.Bool("version",          false, "Show version.")
  flg_collect  := flag.Bool("collect",          false, "Collect all.")
  flg_mysql    := flag.Bool("collect-mysql",    false, "Stats from MySQL.")
  flg_os       := flag.Bool("collect-os",       false, "Info from Linux Operating System.")
  flg_percona  := flag.Bool("collect-percona",  false, "Info & Stats from Percona Toolkit.")
  flg_proxysql := flag.Bool("collect-proxysql", false, "Stats from ProxySQL.")
  // -collect-mysql {limit pk datatype,status,variables,slave status}
  // check is open port from any
  // -output-prometheus
  // -output-influxdb
  // -collect-os {cpu,cores,network,swap,disk[size,free],iops}

  flag.Parse()

  if len(os.Args) == 1 {
    help()
  } else if len(os.Args) >= 2 {
    if *flg_help {
      help()
    } else if *flg_version {
      fmt.Printf("%s\n", config.VERSION)
    } else if *flg_collect {
      collect.OS()
      collect.MySQL()
      collect.Percona()
      collect.ProxySQL()
    } else if *flg_mysql {
      collect.MySQL()
    } else if *flg_percona {
      collect.Percona()
    } else if *flg_proxysql {
      collect.ProxySQL()
    } else if *flg_os {
      collect.OS()
    } else {
      fmt.Printf("%q is not valid command.\n", os.Args[1])
      help()
    }
  } else {
    help()
  }
}

func help() {
  flag.PrintDefaults()
  os.Exit(1)
}
