package net

import (
	"regexp"
	"strings"

	"zenit/config"
	"zenit/monitor/plugins/inputs"
	"zenit/monitor/plugins/lists/metrics"

	"github.com/debeando/go-common/cast"
	"github.com/debeando/go-common/file"
	"github.com/debeando/go-common/log"
)

type Plugin struct{}

func (p *Plugin) Collect(name string, cnf *config.Config, mtc *metrics.Items) {
	defer func() {
		if err := recover(); err != nil {
			log.ErrorWithFields(name, log.Fields{"error": err})
		}
	}()

	if !cnf.Inputs.OS.Net {
		return
	}

	reGroups := regexp.MustCompile(`(\d+)`)
	net := file.Read("/proc/net/dev")
	lines := strings.Split(string(net), "\n")

	for index, line := range lines {
		if index > 1 && len(line) > 0 {
			data := strings.Split(line, ":")
			dev := strings.TrimSpace(data[0])
			match := reGroups.FindAllString(data[1], -1)
			receive_bytes := cast.StringToInt64(match[0])
			transmit_bytes := cast.StringToInt64(match[8])

			log.DebugWithFields(name, log.Fields{
				"device":   dev,
				"receive":  receive_bytes,
				"transmit": transmit_bytes,
			})

			mtc.Add(metrics.Metric{
				Key: "os_net",
				Tags: []metrics.Tag{
					{Name: "hostname", Value: cnf.General.Hostname},
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
	inputs.Add("InputOSNet", func() inputs.Input { return &Plugin{} })
}
