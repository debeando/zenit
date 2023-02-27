package net

import (
	"regexp"
	"strings"

	"github.com/debeando/zenit/common"
	"github.com/debeando/zenit/common/file"
	"github.com/debeando/zenit/common/log"
	"github.com/debeando/zenit/config"
	"github.com/debeando/zenit/plugins/inputs"
	"github.com/debeando/zenit/plugins/lists/metrics"
)

type InputOSNet struct{}

func (l *InputOSNet) Collect() {
	defer func() {
		if err := recover(); err != nil {
			log.Error("InputOSNet", map[string]interface{}{"error": err})
		}
	}()

	if !config.File.Inputs.OS.Net {
		return
	}

	var a = metrics.Load()

	reGroups := regexp.MustCompile(`(\d+)`)
	net := file.Read("/proc/net/dev")
	lines := strings.Split(net, "\n")

	for index, line := range lines {
		if index > 1 && len(line) > 0 {
			data := strings.Split(line, ":")
			dev := strings.TrimSpace(data[0])
			match := reGroups.FindAllString(data[1], -1)
			receive_bytes := common.StringToInt64(match[0])
			transmit_bytes := common.StringToInt64(match[8])

			log.Debug("InputOSNet", map[string]interface{}{
				"device":   dev,
				"receive":  receive_bytes,
				"transmit": transmit_bytes,
			})

			a.Add(metrics.Metric{
				Key: "os_net",
				Tags: []metrics.Tag{
					{Name: "hostname", Value: config.File.General.Hostname},
					{Name: "device", Value: dev},
				},
				Values: []metrics.Value{
					{Key: "receive", Value: receive_bytes},
					{Key: "transmit", Value: transmit_bytes},
				},
			})
		}
	}
}

func init() {
	inputs.Add("InputOSNet", func() inputs.Input { return &InputOSNet{} })
}
