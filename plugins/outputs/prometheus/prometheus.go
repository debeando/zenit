package prometheus

import (
	"github.com/swapbyt3s/zenit/common/file"
	"github.com/swapbyt3s/zenit/common/prometheus"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
	"github.com/swapbyt3s/zenit/plugins/outputs"
)

type OutputPrometheus struct {}

func (l *OutputPrometheus) Collect() {
	file.Create(config.File.Prometheus.TextFile)
	file.Truncate(config.File.Prometheus.TextFile)
	file.Write(config.File.Prometheus.TextFile, prometheus.Normalize(metrics.Load()))
}

func init() {
	outputs.Add("OutputPrometheus", func() outputs.Output { return &OutputPrometheus{} })
}
