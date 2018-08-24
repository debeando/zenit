package os

import (
	"regexp"
	"strings"

	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/common/file"
	"github.com/swapbyt3s/zenit/plugins/accumulator"
)

func Net() {
	reGroups := regexp.MustCompile(`(\d+)`)
	net := file.Read("/proc/net/dev")
	lines := strings.Split(net, "\n")

	for index, line := range lines {
		if index > 1 && len(line) > 0 {
			data := strings.Split(line, ":")
			dev := strings.TrimSpace(data[0])
			match := reGroups.FindAllString(data[1], -1)
			receive_bytes := common.StringToUInt64(match[0])
			transmit_bytes := common.StringToUInt64(match[8])

			accumulator.Load().AddItem(accumulator.Metric{
				Key: "os",
				Tags: []accumulator.Tag{{"system", "linux"},
					{"hardware", "net"},
					{"device", dev}},
				Values: []accumulator.Value{{"receive", receive_bytes},
					{"transmit", transmit_bytes}},
			})
		}
	}
}
