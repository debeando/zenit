package cpu

import (
	"fmt"

	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
	"github.com/swapbyt3s/zenit/plugins/lists/loader"
	"github.com/swapbyt3s/zenit/plugins/lists/alerts"
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
	var percentage = int(common.InterfaceToFloat64(value))

	if percentage == -1 {
		return
	}

	message += fmt.Sprintf("*CPU:* %d%%\n", percentage)

	alerts.Load().Register(
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
	loader.Add("OSCPU", func() loader.Plugin { return &OSCPU{} })
}
