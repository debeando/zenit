package alerts

import (
	"fmt"
	"time"

	"github.com/swapbyt3s/zenit/common/log"
)

// Check is a status for each alert.
type Check struct {
	Critical  int    // Critical value
	Duration  int    // Duration limit between FirstSeen and LastSeen to alert.
	FirstSeen int    // UnixTimeStamp when is register alert.
	Key       string // Key to identify own check.
	LastSeen  int    // What is the last seen (in UnixTimeStamp) value to rest with first seen and compare with duration.
	Message   string // Custom message.
	Name      string // Name human own check.
	Notified  uint8  // Has the same values from Status, and indicate what level is notificated.
	Status    uint8  // Current status for the check: normal, warning, critical or resolved.
	Value     int    // Value to evaluate.
	Warning   int    // Warning value.
}

const (
	Normal   = 0
	Warning  = 1
	Critical = 2
	Resolved = 3
)

// Items is a collection of checks
type Items []Check

var items *Items

// Load is a singleton method to return same object.
func Load() *Items {
	if items == nil {
		items = &Items{}
	}
	return items
}

// Add new check or update last status from existed check.
func (l *Items) Register(key string, name string, duration int, warning int, critical int, value int, message string) {
	check := items.Exist(key)

	if check == nil {
		c := Check{
			Critical: critical,
			Duration: duration,
			FirstSeen: int(time.Now().Unix()),
			Key: key,
			Name: name,
			Status: Normal,
			Warning: warning,
			Value: value,
			Message: message,
		}

		*l = append(*l, c)

		log.Debug(fmt.Sprintf("Alert - Insert: %#v", c))
	} else {
		check.LastSeen = int(time.Now().Unix())
		check.Value = value
		check.Message = message

		log.Debug(fmt.Sprintf("Alert - Update: %#v", check))
	}
}

func (l *Items) Keys() []string {
	var keys = []string{}

	for _, i := range *l {
		keys = append(keys, i.Key)
	}
	return keys
}

func (l *Items) Exist(key string) *Check {
	for i := 0; i < len(*l); i++ {
		if (*l)[i].Key == key {
			return &(*l)[i]
		}
	}
	return nil
}

// IsAlert verify the last time and current time is a valid alert.
func (c *Check) Evaluate() bool {
	if c == nil {
		return false
	}

	if ! ((c.LastSeen - c.FirstSeen) >= c.Duration) {
		return false
	}

	if (c.Status == Warning || c.Status == Critical) && c.Value < c.Warning {
		c.Status = Resolved
	} else if c.Warning == c.Critical && c.Value >= c.Critical {
		c.Status = Critical
	} else if c.Value >= c.Warning && c.Value >= c.Critical {
		c.Status = Critical
	} else if c.Value >= c.Warning && c.Value < c.Critical {
		c.Status = Warning
	}

	if c.Notified == 0 && c.Status == Normal {
		return false
	}

	if c.Notified == c.Status {
		return false
	}

	c.Notified = c.Status

	log.Debug(fmt.Sprintf("Alert - Evaluated: %#v", c))

	return true
}
