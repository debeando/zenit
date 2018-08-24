package os

import (
	"regexp"
	"strings"

	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/plugins/accumulator"
)

func Disk() {
	reGroups := regexp.MustCompile(`^(.+)\s+(\w+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+%)\s+(.+)$`)
	stdout, _ := common.ExecCommand("df -Tl")
	lines := strings.Split(stdout, "\n")

	for index, line := range lines {
		if index > 1 && len(line) > 0 {
			match := reGroups.FindStringSubmatch(line)
			percentage := common.StringToUInt64(strings.TrimRight(strings.TrimSpace(match[6]), "%"))

			accumulator.Load().AddItem(accumulator.Metric{
				Key: "os",
				Tags: []accumulator.Tag{{"system", "linux"},
					{"hardware", "disk"},
					{"filesystem", strings.TrimSpace(match[1])},
					{"type", strings.TrimSpace(match[2])},
					{"mounted", strings.TrimSpace(match[7])}},
				Values: []accumulator.Value{{"used_percent", percentage}},
			})
		}
	}
}
