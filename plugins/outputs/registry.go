package outputs

import (
	"os"

	"github.com/swapbyt3s/zenit/common/log"
)

// Output defines the interface that can interact with the registry
type Output interface {
	Collect()
}

// Creator lets us use a closure to get intsances of the Output struct
type Creator func() Output

// Outputs registry
var Outputs = map[string]Creator{}

// Add can be called from init() on a plugin in this package
// It will automatically be added to the Outputs map to be called externally
func Add(name string, creator Creator) {
	if len(os.Args) == 1 {
		log.Info("Load Plugin", map[string]interface{}{"name": name})
	}

	Outputs[name] = creator
}
