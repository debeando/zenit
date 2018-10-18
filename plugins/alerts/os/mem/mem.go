package mem

import (
	"fmt"

	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/accumulator"
	"github.com/swapbyt3s/zenit/plugins/lists/alerts"
)

func Check() {
	if ! config.File.OS.Alerts.MEM.Enable {
		log.Info("Require to enable OS MEM in config file.")
		return
	}

	var metrics = accumulator.Load()
	var message string = ""
	var value = metrics.FetchOne("os", "name", "mem")
	var percentage = int(common.InterfaceToFloat64(value))

	if percentage == -1 {
		return
	}

	message += fmt.Sprintf("*Memory:* %d%%\n", percentage)

	alerts.Load().Register(
		"mem",
		"MEM",
		config.File.OS.Alerts.MEM.Duration,
		config.File.OS.Alerts.MEM.Warning,
		config.File.OS.Alerts.MEM.Critical,
		percentage,
		message,
	)
}
