package main

import (
  "flag"
  "fmt"
  "os"
  "github.com/swapbyt3s/zenit/collect"
  "github.com/swapbyt3s/zenit/config"
  "github.com/swapbyt3s/zenit/setup"
)

func main() {
  helpPtr            := flag.Bool("help",             false, "Show this help.")
  versionPtr         := flag.Bool("version",          false, "Show version.")
  collectProxySQLPtr := flag.Bool("collect-proxysql", false, "Collect all stats from ProxySQL.")
  setupPtr           := flag.Bool("setup",            false, "Create schemas on MySQL.")

  flag.Parse()

  if len(os.Args) == 1 {
    help()
  } else if len(os.Args) >= 2 {
    if *helpPtr {
      help()
    } else if *versionPtr {
      fmt.Printf("%s\n", config.VERSION)
    } else if *collectProxySQLPtr {
      collect.Prepare()
      collect.Run()
    } else if *setupPtr {
      setup.Run()
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
