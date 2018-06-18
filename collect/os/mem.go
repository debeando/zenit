package os

import (
  "fmt"
  "math"
  "strings"
  "gitlab.com/swapbyt3s/zenit/common"
)

type Mem struct {
  total     uint64
  free      uint64
  available uint64
  buffers   uint64
  cached    uint64
  used      uint64
  percent   float64
}

func GatherMem() {
  lines  := common.ReadFile("/proc/meminfo")
  mem := Mem{}
  for _, line := range(lines) {
    fields := strings.Split(line, ":")

    if len(fields) != 2 {
      continue
    }

    key   := strings.TrimSpace(fields[0])
    value := common.StringToUInt64(strings.Replace(fields[1], " kB", "", -1)) * 1024

    switch key {
    case "MemTotal":
      mem.total = value
    case "MemFree":
      mem.free = value
    case "MemAvailable":
      mem.available = value
    case "Buffers":
      mem.buffers = value
    case "Cached":
      mem.cached = value
    }
  }

  mem.used    = mem.total - mem.free - mem.buffers - mem.cached
  mem.percent = float64(mem.used) / float64(mem.total) * 100.0

  if math.IsNaN(mem.percent) {
    mem.percent = 0
  }

  fmt.Printf("os.linux_mem.total %d\n", mem.total)
  fmt.Printf("os.linux_mem.free %d\n", mem.free)
  fmt.Printf("os.linux_mem.available %d\n", mem.available)
  fmt.Printf("os.linux_mem.buffers %d\n", mem.buffers)
  fmt.Printf("os.linux_mem.cached %d\n", mem.cached)
  fmt.Printf("os.linux_mem.used %d\n", mem.used)
  fmt.Printf("os.linux_mem.used_percent %.2f\n", mem.percent)
}
