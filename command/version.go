package command

import (
	"fmt"
)

// Version is a const to have the latest version number for this code.
const Number string = "1.4.4"

// Build date and time when building.
var BuildTime string

// Return version number and build time.
func Version() string {
	return fmt.Sprintf("%s (%s)", Number, BuildTime)
}
