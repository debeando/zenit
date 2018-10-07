package mem

import (
  "fmt"
  "log"

  "github.com/swapbyt3s/zenit/config"
  "github.com/swapbyt3s/zenit/plugins/accumulator"
  "github.com/swapbyt3s/zenit/plugins/lists/alerts"
)

func Check() {
  if ! config.File.OS.Alerts.MEM.Enable {
    log.Printf("W! - Require to enable OS MEM in config file.\n")
    return
  }

  var metrics = accumulator.Load()
  var check = alerts.Load().Exist("mem")
  var message string = ""
  var value = metrics.Find("os", "name", "mem")
  var percentage = Float64ToInt(value)

  message += fmt.Sprintf("*Memory:* %d\n", percentage)

  if check == nil {
    log.Printf("D! - Alert:OS:MEM - Adding\n")
    log.Printf("D! - Alert:OS:MEM - Message: %s\n", message)
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
    log.Printf("D! - Alert:OS:MEM - Message: %s\n", message)
    log.Printf("D! - Alert:OS:MEM - Updateing\n")
    check.Update(percentage, message)
  }
}

func Float64ToInt(value interface{}) int {
  if v, ok := value.(float64); ok {
    return int(v)
  }
  return -1
}
