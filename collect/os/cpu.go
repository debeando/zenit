package os

import (
  "fmt"
  "strings"
  "time"
  "gitlab.com/swapbyt3s/zenit/lib"
)

type CPU struct {
  idle  uint64
  total uint64
}

func GetCPU() {
  c := [2]CPU{}

  c[0].idle, c[0].total = getCPUSample()
  time.Sleep(300 * time.Millisecond)
  c[1].idle, c[1].total = getCPUSample()

  if c[0].total > 0 && c[1].total > 0 && c[0].idle > 0 && c[1].idle > 0 {
    total := c[1].total - c[0].total
    idle  := c[1].idle  - c[0].idle
    percentage := (float64(total) - float64(idle)) / float64(total) * 100.0

    fmt.Printf("os.linux_cpu.used_percent %.2f\n", percentage)
  }
}

func getCPUSample() (idle uint64, total uint64) {
  lines := lib.ReadFile("/proc/stat")
  if len(lines) > 0 {
    fields := strings.Fields(lines[0])

      for i := 1; i < len(fields); i++ {
        total += lib.StringToUInt64(fields[i])
        if i == 4 {
          idle = lib.StringToUInt64(fields[i])
        }
    }
  }
  return
}
