package outputs

import (
	"zenit/agent/plugins/lists/metrics"
	"zenit/config"
)

// Output defines the interface that can interact with the registry
type Output interface {
	Deliver(string, *config.Config, *metrics.Items)
}

// Creator lets us use a closure to get intsances of the Output struct
type Creator func() Output

// Outputs registry
var Outputs = map[string]Creator{}

// Add can be called from init() on a plugin in this package
// It will automatically be added to the Outputs map to be called externally
func Add(name string, creator Creator) {
	Outputs[name] = creator
}
