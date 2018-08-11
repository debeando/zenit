package common

import (
  "bufio"
  "log"
  "os"
  "os/exec"
)

func Tail(path string, channel chan<- string) {
  if _, err := os.Stat(path); err != nil {
    log.Printf("E! - Tail - File not exist: %s\n", path)
    os.Exit(1)
  }

  cmd       := exec.Command("/usr/bin/tail", "-n", "0", "-f", path)
  stdout, _ := cmd.StdoutPipe()
  scanner   := bufio.NewScanner(stdout)

  go func() {
    defer close(channel)

    for scanner.Scan() {
      channel <- scanner.Text()
    }
  }()

  cmd.Run()
}
