package os

import (
  "regexp"
  "strings"
  "gitlab.com/swapbyt3s/zenit/common"
  "gitlab.com/swapbyt3s/zenit/accumulator"
)

func GatherNet() {
  reGroups := regexp.MustCompile(`(\d+)`)
  lines  := common.ReadFile("/proc/net/dev")

  for index, line := range(lines) {
    if index > 1 && len(line) > 0 {
      data := strings.Split(line, ":")
      dev := strings.TrimSpace(data[0])
      match := reGroups.FindAllString(data[1], -1)
      receive_bytes := common.StringToUInt64(match[0])
      transmit_bytes := common.StringToUInt64(match[8])

      accumulator.Load().AddItem(accumulator.Metric{
        Key: "os",
        Tags: []accumulator.Tag{accumulator.Tag{"system", "linux"},
                                accumulator.Tag{"hardware", "net"},
                                accumulator.Tag{"device", dev}},
        Values: []accumulator.Value{accumulator.Value{"receive", receive_bytes},
                                    accumulator.Value{"transmit", transmit_bytes}},
      })
    }
  }
}
