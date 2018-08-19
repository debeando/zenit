package common

import (
  "crypto/md5"
  "encoding/hex"
  "net"
  "os"
  "os/exec"
  "sort"
  "strconv"
  "strings"
  "syscall"
  "time"

  "github.com/swapbyt3s/zenit/common/file"
)

func PGrep(cmd string) int {
  _, exitcode := ExecCommand("/usr/bin/pgrep -x '" + cmd +"' > /dev/null")

  return exitcode
}

func GetUInt64FromFile(path string) uint64 {
  lines := file.Read(path)
  if len(lines) > 0 {
    return StringToUInt64(lines)
  }
  return 0
}

func GetIntFromFile(path string) int {
  lines := file.Read(path)
  if len(lines) > 0 {
    return StringToInt(lines)
  }
  return 0
}

func StringToInt(value string) int {
  i, err := strconv.Atoi(strings.TrimSpace(value))
  if err != nil {
    return 0
  }
  return i
}

func IntToString(value int) string {
  return strconv.Itoa(value)
}

func StringToUInt64(value string) uint64 {
  i, err := strconv.ParseUint(strings.TrimSpace(value), 10, 64)
  if err != nil {
    return 0
  }
  return i
}

func KeyInMap(key string, list map[string]string) bool {
  if _, ok := list[key]; ok {
    return true
  }
  return false
}

func KeyOfMaps(v []map[string]string) (keys []string) {
  if len(v) > 0 {
    for key := range v[0] {
      keys = append(keys, key)
    }
    sort.Strings(keys)
  }
  return
}

func StringInArray(key string, list []string) bool {
  for _, l := range list {
    if l == key {
      return true
    }
  }
  return false
}

func Hostname() string {
  host, err := os.Hostname()
  if err != nil {
    return ""
  }

  return host
}

func IPAddress() string {
  addrs, _ := net.InterfaceAddrs()

  for _, a := range addrs {
    if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
      if ipnet.IP.To4() != nil {
        return ipnet.IP.String()
      }
    }
  }

  return ""
}

func ToDateTime(timestamp string, layout string) string {
  t, err := time.Parse(layout, timestamp)
  if err != nil {
    return ""
  }
  return t.Format("2006-01-02 15:04:05")
}

func Escape(text string) string {
  return strings.Replace(text, "'", "\\'", -1)
}

func ExecCommand(cmd string) (stdout string, exitcode int) {
  out, err := exec.Command("/bin/bash", "-c", cmd).Output()
  if err != nil {
    if exitError, ok := err.(*exec.ExitError); ok {
      ws := exitError.Sys().(syscall.WaitStatus)
      exitcode = ws.ExitStatus()
    }
  }
  stdout = string(out[:])
  return
}

func MD5(s string) string {
  hash := md5.Sum([]byte(s))
  return hex.EncodeToString(hash[:])
}
