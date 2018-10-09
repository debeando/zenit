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
  var check = alerts.Load().Exist("mem")
  var message string = ""
  var value = metrics.FetchOne("os", "name", "mem")
  var percentage = common.InterfaceToInt(value)

  message += fmt.Sprintf("*Memory:* %d\n", percentage)

  if check == nil {
    alerts.Load().Add(
      "mem",
      "MEM",
      config.File.OS.Alerts.MEM.Duration,
      config.File.OS.Alerts.MEM.Warning,
      config.File.OS.Alerts.MEM.Critical,
      percentage,
      message,
      true,
    )
  } else {
    check.Update(percentage, message)
  }
}
