package os

import (
  "bufio"
  "fmt"
  "os"
  "strings"
  "strconv"
  "time"
)

type CPU struct {
  idle  uint64
  total uint64
}

func GetCPU() {
  c := [2]CPU{}

  for i := 0; i < 2; i++ {
    file, err := os.Open("/proc/stat")
    defer file.Close()
    if err == nil {
      scanner := bufio.NewScanner(file)
      for scanner.Scan() {
        fields := strings.Split(scanner.Text(), " ")
        fields  = ArrayRemoveEmpty(fields)

        if len(fields) == 11 {
          key := strings.TrimSpace(fields[0])

          if key == "cpu" {
            user    := StringToUInt64(fields[1])
            nice    := StringToUInt64(fields[2])
            system  := StringToUInt64(fields[3])
            idle    := StringToUInt64(fields[4])
            iowait  := StringToUInt64(fields[5])
            irq     := StringToUInt64(fields[6])
            softirq := StringToUInt64(fields[7])
            steal   := StringToUInt64(fields[8])

            c[i].idle  = idle + iowait
            c[i].total = idle + iowait + user + nice + system + irq + softirq + steal

            if i == 0 {
              time.Sleep(1 * time.Second)
            }

            break
          }
        }
      }
    }
  }

  if c[0].total > 0 && c[1].total > 0 && c[0].idle > 0 && c[1].idle > 0 {
    total := c[1].total - c[0].total
    idle  := c[1].idle  - c[0].idle
    percentage := (float64(total) - float64(idle)) / float64(total) * 100.0

    fmt.Printf("os.linux_cpu.used_percent %.2f\n", percentage)
  }
}

func StringToUInt64(value string) uint64 {
  i, err := strconv.ParseUint(strings.TrimSpace(value), 10, 64)
  if err != nil {
    panic(err)
  }
  return i
}

func ArrayRemoveEmpty (s []string) []string {
    var r []string
    for _, str := range s {
        if str != "" {
            r = append(r, str)
        }
    }
    return r
}
