package inputs

import (
	"zenit/config"
	"zenit/agent/plugins/lists/metrics"
)

// Input defines the interface that can interact with the registry
type Input interface {
	Collect(string, *config.Config, *metrics.Items)
}

// Creator lets us use a closure to get intsances of the Input struct
type Creator func() Input

// Inputs registry
var Inputs = map[string]Creator{}

// Add can be called from init() on a plugin in this package
// It will automatically be added to the Inputs map to be called externally
func Add(name string, creator Creator) {
	Inputs[name] = creator
}
