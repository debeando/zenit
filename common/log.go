package common

import (
  "log"
  "os"
)

func LogInit(log_file string) {
  file, err := os.OpenFile(log_file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
  if err != nil {
    log.Fatalln("Failed to open log file:", err)
  }

  log.SetOutput(file)
}
