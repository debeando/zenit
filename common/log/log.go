package log

import (
  "log"

  "github.com/swapbyt3s/zenit/config"
)

func Info(m string, a ...interface{}) {
  log.Printf("I! - %s\n", m, a)
}

func Error(m string, a ...interface{}){
  log.Printf("E! - %s\n", m, a)
}

func Debug(m string, a ...interface{}) {
  if config.General.Debug {
    log.Printf("D! - %s\n", m, a)
  }
}
