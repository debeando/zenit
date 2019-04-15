package prometheus

import (
	"fmt"

	"github.com/swapbyt3s/zenit/common/file"
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
	"github.com/swapbyt3s/zenit/plugins/outputs"
)

type OutputPrometheus struct{}

func (l *OutputPrometheus) Collect() {
	defer func() {
		if err := recover(); err != nil {
			log.Debug(fmt.Sprintf("Plugin - OutputPrometheus - Panic (code %d) has been recover from somewhere.\n", err))
		}
	}()

	if !config.File.Outputs.Prometheus.Enable {
		return
	}

	file.Create(config.File.Outputs.Prometheus.TextFile)
	file.Truncate(config.File.Outputs.Prometheus.TextFile)
	file.Write(config.File.Outputs.Prometheus.TextFile, Normalize(metrics.Load()))
}

func init() {
	outputs.Add("OutputPrometheus", func() outputs.Output { return &OutputPrometheus{} })
}
