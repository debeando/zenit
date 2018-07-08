package os

import (
  "regexp"
  "strings"
  "gitlab.com/swapbyt3s/zenit/common"
  "gitlab.com/swapbyt3s/zenit/accumulator"
)

func GatherDisk() {
  reGroups := regexp.MustCompile(`^(.+)\s+(\w+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+%)\s+(.+)$`)
  stdout, _ := common.ExecCommand("df -Tl")
  lines := strings.Split(stdout, "\n")

  for index, line := range(lines) {
    if index > 1 && len(line) > 0 {
      match := reGroups.FindStringSubmatch(line)
      percentage := common.StringToUInt64(strings.TrimRight(strings.TrimSpace(match[6]), "%"))

      accumulator.Load().AddItem(accumulator.Metric{
        Key: "os",
        Tags: []accumulator.Tag{accumulator.Tag{"system", "linux"},
                                accumulator.Tag{"hardware", "disk"},
                                accumulator.Tag{"filesystem", strings.TrimSpace(match[1])},
                                accumulator.Tag{"type", strings.TrimSpace(match[2])},
                                accumulator.Tag{"mounted", strings.TrimSpace(match[7])}},
        Values: []accumulator.Value{accumulator.Value{"used_percent", percentage}},
      })
    }
  }
}
