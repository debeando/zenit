package net

import (
	"regexp"
	"strings"

	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/common/file"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/loader"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

type InputOSNet struct {}

func (l *InputOSNet) Collect() {
	if ! config.File.OS.Inputs.Net {
		return
	}

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

			metrics.Load().Add(metrics.Metric{
				Key: "zenit_os",
				Tags: []metrics.Tag{
					{"name", "net"},
					{"device", dev},
				},
				Values: []metrics.Value{
					{"receive", uint(receive_bytes)},
					{"transmit", uint(transmit_bytes)},
				},
			})
		}
	}
}

func init() {
	loader.Add("InputOSNet", func() loader.Plugin { return &InputOSNet{} })
}
