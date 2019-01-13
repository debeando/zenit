package alerts

import (
  "fmt"

  "github.com/swapbyt3s/zenit/common/log"
)

// Alert defines the interface that can interact with the registry
type Alert interface {
  Collect()
}

// Creator lets us use a closure to get intsances of the Alert struct
type Creator func() Alert

// Alerts registry
var Alerts = map[string]Creator{}

// Add can be called from init() on a plugin in this package
// It will automatically be added to the Alerts map to be called externally
func Add(name string, creator Creator) {
  log.Info(fmt.Sprintf("Plugin Alerts - %s", name))

  Alerts[name] = creator
}
