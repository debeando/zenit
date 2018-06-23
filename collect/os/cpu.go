package os

import (
  "strings"
  "time"
  "gitlab.com/swapbyt3s/zenit/common"
  "gitlab.com/swapbyt3s/zenit/output"
)

type CPU struct {
  idle  uint64
  total uint64
}

func GatherCPU() {
  c := [2]CPU{}

  c[0].idle, c[0].total = getCPUSample()
  time.Sleep(300 * time.Millisecond)
  c[1].idle, c[1].total = getCPUSample()

  if c[0].total > 0 && c[1].total > 0 && c[0].idle > 0 && c[1].idle > 0 {
    total := c[1].total - c[0].total
    idle  := c[1].idle  - c[0].idle
    percentage := (float64(total) - float64(idle)) / float64(total) * 100.0

    output.Load().AddItem(output.Metric{
      Key: "os",
      Tags: []output.Tag{output.Tag{"system", "linux"},
                         output.Tag{"hardware", "cpu"}},
      Values: []output.Value{output.Value{"used_percent", percentage}},
    })
  }
}

func getCPUSample() (idle uint64, total uint64) {
  lines := common.ReadFile("/proc/stat")
  if len(lines) > 0 {
    fields := strings.Fields(lines[0])

      for i := 1; i < len(fields); i++ {
        total += common.StringToUInt64(fields[i])
        if i == 4 {
          idle = common.StringToUInt64(fields[i])
        }
    }
  }
  return
}
