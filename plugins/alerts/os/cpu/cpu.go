package cpu

import (
	"fmt"

	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/alerts"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
	"github.com/swapbyt3s/zenit/plugins/lists/checks"
)

type OSCPU struct {}

func (l *OSCPU) Collect() {
	if ! config.File.OS.Alerts.CPU.Enable {
		log.Info("Require to enable OS CPU in config file.")
		return
	}

	var metrics = metrics.Load()
	var message string = ""
	var value = metrics.FetchOne("zenit_os", "name", "cpu")
	var percentage = uint64(common.InterfaceToFloat64(value))

	message += fmt.Sprintf("*CPU:* %d%%\n", percentage)

	checks.Load().Register(
		"cpu",
		"CPU",
		config.File.OS.Alerts.CPU.Duration,
		config.File.OS.Alerts.CPU.Warning,
		config.File.OS.Alerts.CPU.Critical,
		percentage,
		message,
	)
}

func init() {
	alerts.Add("AlertOSCPU", func() alerts.Alert { return &OSCPU{} })
}
