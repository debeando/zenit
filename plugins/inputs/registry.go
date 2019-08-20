package inputs

import (
	"fmt"
	"os"

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
	if len(os.Args) == 1 {
		log.Info(fmt.Sprintf("Load Plugin - %s", name))
	}

	Inputs[name] = creator
}
