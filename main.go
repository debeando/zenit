package main

import (
  "flag"
  "fmt"
  "os"

  "gitlab.com/swapbyt3s/zenit/common"
  "gitlab.com/swapbyt3s/zenit/config"
  "gitlab.com/swapbyt3s/zenit/daemonize"
  "gitlab.com/swapbyt3s/zenit/plugins/inputs"
)

const USAGE = `zenit (%s) written by %s
Usage: %s [--help | --quiet | --start | --stop | --version]
Options:
  --help        Show this help.
  --quiet       Run in quiet mode.
  --start       Start daemon.
  --stop        Stop daemon.
  --version     Print version numbers.
`

func init() {
  config.Load()
  common.LogInit(config.General.LogFile)
}

func main() {
  fHelp    := flag.Bool("help", false, "Show this help.")
  fQuiet   := flag.Bool("quiet", false, "Run in quiet mode.")
  fStart   := flag.Bool("start", false, "Fork to the background and detach from the shell.")
  fStop    := flag.Bool("stop", false, "Stop daemon.")
  fVersion := flag.Bool("version", false, "Show version.")

  flag.Usage = func() { help(0) }
  flag.Parse()

  if len(os.Args) == 1 {
    help(0)
  }

  switch {
  case *fVersion:
    fmt.Printf("%s\n", config.VERSION)
    return
  case *fHelp:
    help(0)
  case *fStart:
    daemonize.Start()
  case *fStop:
    daemonize.Stop()
  case *fQuiet:
    inputs.Gather()
  }
}

func help(rc int) {
  fmt.Printf(USAGE, config.VERSION, config.AUTHOR, os.Args[0])
  os.Exit(rc)
}
