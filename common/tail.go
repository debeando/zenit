package common

import (
  "github.com/hpcloud/tail"
)

func Tail(filename string, out chan<- string) {
  defer close(out)

  t, _ := tail.TailFile(filename, tail.Config{Follow: true})
  for line := range t.Lines {
    out <- line.Text
  }
}
