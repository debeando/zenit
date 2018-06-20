package output

import (
  "fmt"
  "strings"
  "gitlab.com/swapbyt3s/zenit/common"
)

func Prometheus() {
  var a = LoadAccumulator()

  for i := range(a.Items) {
    fmt.Printf("%s", a.Items[i].Key)
    s := []string{}
    for t := range(a.Items[i].Tags) {
      k := a.Items[i].Tags[t].Name
      v := strings.ToLower(a.Items[i].Tags[t].Value)
      s = append(s, fmt.Sprintf("%s=\"%s\"", k, v))
    }
    t := strings.Join(s,",")
    fmt.Printf("{%s}", t)
    if common.IsIntegral(a.Items[i].Value) {
      fmt.Printf(" %d", uint64(a.Items[i].Value))
    } else {
      fmt.Printf(" %.2f", a.Items[i].Value)
    }
    fmt.Printf("\n")
  }
}
