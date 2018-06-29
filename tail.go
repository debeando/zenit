package main

import (
  "fmt"
  "os"
  "gitlab.com/swapbyt3s/zenit/common"
)

func main() {
  filename := os.Args[1]
  out := make(chan string)

  go common.Tail(filename, out)

  for tail := range out {
    fmt.Print(tail)
  }

  close(out)
}
