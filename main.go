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
  flg_proxysql := flag.Bool("collect-proxysql", false, "Stats from ProxySQL.")
  flg_os       := flag.Bool("collect-os",       false, "Info from Linux Operating System.")
  // -collect-mysql
  // -output-prometheus
  // -output-influxdb

  flag.Parse()

  if len(os.Args) == 1 {
    help()
  } else if len(os.Args) >= 2 {
    if *flg_help {
      help()
    } else if *flg_version {
      fmt.Printf("%s\n", config.VERSION)
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
