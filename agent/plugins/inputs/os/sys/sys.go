package sys

import (
	"zenit/config"
	"zenit/agent/plugins/inputs"
	"zenit/agent/plugins/lists/metrics"

	"github.com/debeando/go-common/file"
	"github.com/debeando/go-common/log"
)

const NR_OPEN string = "/proc/sys/fs/nr_open"
const FILE_MAX string = "/proc/sys/fs/file-max"

type Plugin struct{}

func (p *Plugin) Collect(name string, cnf *config.Config, mtc *metrics.Items) {
	defer func() {
		if err := recover(); err != nil {
			log.ErrorWithFields(name, log.Fields{"message": err})
		}
	}()

	if !cnf.Inputs.OS.Limits {
		return
	}

	log.DebugWithFields(name, log.Fields{
		"nr_open":  file.GetInt64(NR_OPEN),
		"file_max": file.GetInt64(FILE_MAX),
	})

	mtc.Add(metrics.Metric{
		Key: "os_sys",
		Tags: []metrics.Tag{
			{Name: "hostname", Value: cnf.General.Hostname},
		},
		Values: []metrics.Value{
			{Key: "nr_open", Value: file.GetInt64(NR_OPEN)},
			{Key: "file_max", Value: file.GetInt64(FILE_MAX)},
		},
	})
}

func init() {
	inputs.Add("InputOSLimits", func() inputs.Input { return &Plugin{} })
}
