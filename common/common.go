package common

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
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
	_, exitcode := ExecCommand("/usr/bin/pgrep -x '" + cmd + "' > /dev/null")

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
	if flag.Lookup("test.v") != nil {
		return "localhost.test"
	}

	host, err := os.Hostname()
	if err != nil {
		return ""
	}

	return host
}

func IPAddress() string {
	if flag.Lookup("test.v") != nil {
		return "127.0.0.1"
	}

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
	return strings.Replace(text, "'", `\'`, -1)
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

func ComparteMapString(a, b map[string]string) bool {
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for k := range a {
		if a[k] != b[k] {
			return false
		}
	}

	return true
}

func SplitKeyAndValue(s *string) (key string, value string) {
	kv := strings.SplitN(*s, "=", 2)
	if len(kv) == 2 {
		return strings.TrimSpace(strings.ToLower(kv[0])), kv[1]
	}
	return "", ""
}

func Trim(value *string) string {
	*value = strings.TrimSpace(*value)
	*value = strings.TrimRight(*value, "\"")
	*value = strings.TrimLeft(*value, "\"")
	return *value
}

func Percentage(value float64, max float64) float64 {
	return (value / max) * 100
}

func InterfaceToInt(value interface{}) int {
	if v, ok := value.(float64); ok {
		return int(v)
	}
	return -1
}

func InterfaceToFloat64(value interface{}) float64 {
	if v, ok := value.(float64); ok {
		return v
	}
	return -1
}
