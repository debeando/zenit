package os

import (
  "math"
  "strings"

  "github.com/swapbyt3s/zenit/common"
  "github.com/swapbyt3s/zenit/common/file"
  "github.com/swapbyt3s/zenit/plugins/accumulator"
)

type MemMetric struct {
  total     uint64
  free      uint64
  available uint64
  buffers   uint64
  cached    uint64
  used      uint64
  percent   float64
}

func Mem() {
  meminfo := file.Read("/proc/meminfo")
  lines   := strings.Split(meminfo, "\n")
  mem     := MemMetric{}

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

  accumulator.Load().AddItem(accumulator.Metric{
    Key: "os",
    Tags: []accumulator.Tag{accumulator.Tag{"system", "linux"},
                            accumulator.Tag{"hardware", "mem"}},
    Values: []accumulator.Value{accumulator.Value{"total", mem.total},
                                accumulator.Value{"free", mem.free},
                                accumulator.Value{"available", mem.available},
                                accumulator.Value{"buffers", mem.buffers},
                                accumulator.Value{"cached", mem.cached},
                                accumulator.Value{"used", mem.used},
                                accumulator.Value{"used_percent", mem.percent}},
  })
}
