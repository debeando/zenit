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

const USAGE = "zenit (%s) written by Nicola Strappazzon C. <swapbyt3s@gmail.com>\nUsage: %s <command>\n"

func main() {
  fHelp         := flag.Bool("help", false, "Show this help.")
  fVersion      := flag.Bool("version", false, "Show version.")
  fDaemonize    := flag.Bool("daemonize", false, "Fork to the background and detach from the shell.")
  fCollect      := flag.String("collect", "", "List of metrics to collect.")
  fParserFormat := flag.String("parser-format", "", "Parser log format.")
  fParserFile   := flag.String("parser-file", "", "Fail path to to parse.")
  fStop         := flag.Bool("stop", false, "Stop daemon.")
  fRun          := flag.String("run", "", "Run bash command and wait to finish to notify via slack.")

  flag.Parse()

  if len(os.Args) == 1 {
    help()
  } else if len(os.Args) >= 2 {
    if *fHelp {
      help()
    }

    if *fDaemonize {
      daemonize.Start()
    } else if *fStop {
      daemonize.Stop()
    }

    if *fVersion {
      fmt.Printf("%s\n", config.VERSION)
    } else if len(*fCollect) > 0 {
      collect.Run(strings.Split(*fCollect, ","))
    } else if len(*fParserFormat) > 0 && len(*fParserFile) > 0 {
      collect.Parser(*fParserFormat, *fParserFile)
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
