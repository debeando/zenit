package common

import (
  "crypto/md5"
  "encoding/hex"
  "io/ioutil"
  "net"
  "os"
  "os/exec"
  "strconv"
  "strings"
  "syscall"
  "time"
)

// Puede que el .Output() este de más.
// Se podria reutilizar la funcion ExecCommand
func PGrep(cmd string) int {
  _, err := exec.Command("/bin/bash", "-c", "/usr/bin/pgrep -x '" + cmd +"' > /dev/null").Output()
  if err != nil {
    if exitError, ok := err.(*exec.ExitError); ok {
      ws := exitError.Sys().(syscall.WaitStatus)
      return ws.ExitStatus()
    }
  }
  return 0
}

func ReadFile(path string) (lines []string) {
  if _, err := os.Stat(path); err == nil {
    contents, err := ioutil.ReadFile(path)
    if err != nil {
      panic(err)
    }

    lines = strings.Split(string(contents), "\n")
  }
  return
}

func GetUInt64FromFile(path string) uint64 {
  lines := ReadFile(path)
  if len(lines) > 0 {
    return StringToUInt64(lines[0])
  }
  return 0
}

func StringToUInt64(value string) uint64 {
  i, err := strconv.ParseUint(strings.TrimSpace(value), 10, 64)
  if err != nil {
    panic(err)
  }
  return i
}

func StringInArray(key string, list []string) bool {
  for _, l := range list {
    if l == key {
      return true
    }
  }
  return false
}

func GetEnv(key string, default_value string) string {
  val, ok := os.LookupEnv(key)
  if !ok {
    if len(default_value) > 0 {
      return default_value
    }
  }
  return val
}

// Retorna "" si hay err.
func Hostname() string {
  host, err := os.Hostname()
  if err != nil {
    panic(err)
  }

  return host
}

func IpAddress() string {
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

// The timestamp is represent in ISO 8601 (UTC) -> RFC 3339
func ToDateTime(timestamp string) string {
  layout := "2006-01-02T15:04:05 UTC"
  t, err := time.Parse(layout, timestamp)
  if err != nil {
    return ""
  }
  return t.Format("2006-01-02 15:04:05")
}

// Creo que esta ya no se usa.
func ISO8601V2toRFC3339(timestamp string) string {
  layout := "20060102 15:04:05"
  t, err := time.Parse(layout, "20" + timestamp)
  if err != nil {
    return ""
  }
  return t.Format("2006-01-02 15:04:05")
}

// Adding slash by quotes in strings:
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
