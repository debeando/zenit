package alerts

import (
	"time"
)

// Check is a status for each alert.
type Check struct {
	Critical  uint64 // Critical value
	Duration  int    // Duration limit between FirstSeen and LastSeen to alert.
	FirstSeen int    // UnixTimeStamp when is register alert.
	Key       string // Key to identify own check.
	LastSeen  int    // What is the last seen (in UnixTimeStamp) value to rest with first seen and compare with duration.
	Message   string // Custom message.
	Name      string // Name human own check.
	Notified  uint8  // Has the same values from Status, and indicate what level is notificated.
	Status    uint8  // Current status for the check: normal, warning, critical or resolved.
	Value     uint64 // Value to evaluate.
	Warning   uint64 // Warning value.
	Decide    bool   // Use Evaluate function.
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

// Count all checks in accumulator.
func (l *Items) Count() int {
	return len(*l)
}

// Add new check.
func (l *Items) Add(key string, name string, duration int, warning uint64, critical uint64, value uint64, message string, decide bool) {
	if items.Exist(key) == nil {
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
			Decide: decide,
		}

		*l = append(*l, c)
	}
}

func (l *Items) Keys() []string {
	var keys = []string{}

	for i := 0; i < len(*l); i++ {
		keys = append(keys, (*l)[i].Key)
	}
	return keys
}

func (l *Items) Delete(key string) bool {
	for i := 0; i < len(*l); i++ {
		if (*l)[i].Key == key {
			*l = (*l)[:i+copy((*l)[i:], (*l)[i+1:])]
			return true
		}
	}
	return false
}

func (l *Items) Exist(key string) *Check {
	for i := 0; i < len(*l); i++ {
		if (*l)[i].Key == key {
			return &(*l)[i]
		}
	}
	return nil
}

func (c *Check) SetLastSeen(value uint64) {
	if c == nil {
		return
	}

	c.LastSeen = int(time.Now().Unix())
	c.Value = value
}

// IsAlert verify the last time and current time is a valid alert.
func (c *Check) Evaluate() {
	if c == nil {
		return
	}

	if ! ((c.LastSeen - c.FirstSeen) >= c.Duration) {
		return
	}

	if ! c.Decide {
		c.Status = Critical
		return
	}

	if c.Status == Warning || c.Status == Critical {
		if c.Value < c.Warning {
			c.Status = Resolved
		}
	} else if c.Warning == c.Critical && c.Value >= c.Critical {
		c.Status = Critical
	} else if c.Value >= c.Warning {
		if c.Value >= c.Critical {
			c.Status = Critical
		} else {
			c.Status = Warning
		}
	}
}

func (c *Check) Update(value uint64, message string) {
	c.SetLastSeen(value)
	c.Evaluate()
	c.Message = message
}

func (c *Check) Notify() bool {
	if c == nil {
		return false
	}

	if c.Status == Normal && c.Notified == 0 {
		return false
	}

	if c.Status == c.Notified {
		return false
	}

	c.Notified = c.Status
	return true
}
