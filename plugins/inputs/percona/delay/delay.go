package process

import (
	"fmt"

	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/loader"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

type InputsPerconaDelay struct {}

func (l *InputsPerconaDelay) Collect() {
	if ! config.File.Process.Inputs.PerconaToolKitSlaveDelay {
		return
	}

	var pid = common.PGrep("pt-slave-delay")
	var value = 0

	if pid > 0 {
		value = 1
	}

	metrics.Load().Add(metrics.Metric{
		Key: "zenit_os",
		Tags: []metrics.Tag{
			{"system", "linux"},
			{"process", "pt_slave_delay"},
		},
		Values: value,
	})

	log.Debug(fmt.Sprintf("Plugin - InputsPerconaDelay - %d", value))
}

func init() {
	loader.Add("InputsPerconaDelay", func() loader.Plugin { return &InputsPerconaDelay{} })
}
