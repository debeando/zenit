package os

import (
  "bufio"
  "fmt"
  "os"
  "strconv"
  "strings"
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

func GetMem() {
  file, err := os.Open("/proc/meminfo")
  defer file.Close()
  if err == nil {
    mem := Mem{}
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
      fields := strings.Split(scanner.Text(), ":")
      if len(fields) != 2 {
        continue
      }

      key   := strings.TrimSpace(fields[0])
      value := strings.TrimSpace(fields[1])
      value  = strings.Replace(value, " kB", "", -1)

      t, err := strconv.ParseUint(value, 10, 64)
      if err != nil {
        panic(err)
      }

      switch key {
      case "MemTotal":
        mem.total = t * 1024
      case "MemFree":
        mem.free = t * 1024
      case "MemAvailable":
        mem.available = t * 1024
      case "Buffers":
        mem.buffers = t * 1024
      case "Cached":
        mem.cached = t * 1024
      }
    }

    mem.used    = mem.total - mem.free - mem.buffers - mem.cached
    mem.percent = float64(mem.used) / float64(mem.total) * 100.0

    fmt.Printf("os.linux_mem.total %d\n", mem.total)
    fmt.Printf("os.linux_mem.free %d\n", mem.free)
    fmt.Printf("os.linux_mem.available %d\n", mem.available)
    fmt.Printf("os.linux_mem.buffers %d\n", mem.buffers)
    fmt.Printf("os.linux_mem.cached %d\n", mem.cached)
    fmt.Printf("os.linux_mem.used %d\n", mem.used)
    fmt.Printf("os.linux_mem.used_percent %.2f\n", mem.percent)
  }
}
