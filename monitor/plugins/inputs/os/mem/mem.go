package mem

import (
	"fmt"

	"zenit/config"
	"zenit/monitor/plugins/inputs"
	"zenit/monitor/plugins/lists/metrics"

	"github.com/debeando/go-common/log"
	"github.com/shirou/gopsutil/mem"
)

type Plugin struct{}

func (p *Plugin) Collect(name string, cnf *config.Config, mtc *metrics.Items) {
	defer func() {
		if err := recover(); err != nil {
			log.ErrorWithFields(name, log.Fields{"error": err})
		}
	}()

	if !cnf.Inputs.OS.Mem {
		return
	}

	vmStat, err := mem.VirtualMemory()
	if err == nil {
		mtc.Add(metrics.Metric{
			Key: "os_mem",
			Tags: []metrics.Tag{
				{Name: "hostname", Value: cnf.General.Hostname},
			},
			Values: []metrics.Value{
				{Key: "percentage", Value: vmStat.UsedPercent},
			},
		})

		log.DebugWithFields(name, log.Fields{"value": fmt.Sprintf("%.2f", vmStat.UsedPercent)})
	}
}

func init() {
	inputs.Add("InputOSMem", func() inputs.Input { return &Plugin{} })
}
