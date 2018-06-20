package output

import (
  "fmt"
  "strings"
)

func Prometheus() {
  var a = Load()

  for _, m := range *a {
    fmt.Printf("%s", m.Key)
    s := []string{}
    for t := range(m.Tags) {
      k := m.Tags[t].Name
      v := strings.ToLower(m.Tags[t].Value)
      s = append(s, fmt.Sprintf("%s=\"%s\"", k, v))
    }
    t := strings.Join(s,",")
    fmt.Printf("{%s}", t)
    switch value := m.Values.(type) {
    case uint64:
      fmt.Printf(" %d", value)
    case float64:
      fmt.Printf(" %.2f", value)
    default:
      fmt.Printf(" 0")
    }
    fmt.Printf("\n")
  }
}
