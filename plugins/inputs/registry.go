package inputs

import (
	"github.com/swapbyt3s/zenit/common/log"
)

// Input defines the interface that can interact with the registry
type Input interface {
	Collect()
}

// Creator lets us use a closure to get intsances of the Input struct
type Creator func() Input

// Inputs registry
var Inputs = map[string]Creator{}

// Add can be called from init() on a plugin in this package
// It will automatically be added to the Inputs map to be called externally
func Add(name string, creator Creator) {
	log.Info("Load Plugin", map[string]interface{}{"name": name})

	Inputs[name] = creator
}
