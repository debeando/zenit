package os

import (
  "math"
  "strings"
  "gitlab.com/swapbyt3s/zenit/common"
  "gitlab.com/swapbyt3s/zenit/output"
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

  output.Load().AddItem(output.Metric{
    Key: "os",
    Tags: []output.Tag{output.Tag{"system", "linux"},
                       output.Tag{"hardware", "mem"}},
    Values: []output.Value{output.Value{"total", mem.total},
                           output.Value{"free", mem.free},
                           output.Value{"available", mem.available},
                           output.Value{"buffers", mem.buffers},
                           output.Value{"cached", mem.cached},
                           output.Value{"used", mem.used},
                           output.Value{"used_percent", mem.percent}},
  })
}
