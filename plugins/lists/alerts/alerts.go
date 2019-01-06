package alerts

import (
	"fmt"
	"time"

	"github.com/swapbyt3s/zenit/common/log"
)

// Check is a status for each alert.
type Check struct {
	Critical    uint64 // Critical value
	Duration    int    // Duration limit between FirstSeen and LastSeen to alert.
	FirstSeen   int    // UnixTimeStamp when is register alert.
	Key         string // Key to identify own check.
	LastSeen    int    // What is the last seen (in UnixTimeStamp) value to rest with first seen and compare with duration.
	Message     string // Custom message.
	Name        string // Name human own check.
	Status      uint8  // Current status for the check: normal, warning, critical or Recovered.
	Value       uint64 // Value to evaluate.
	BeforeValue uint64 // Last evaluated value.
	Warning     uint64 // Warning value.
}

const (
	Normal    = 0
	Warning   = 1
	Critical  = 2
	Notified  = 3
	Recovered = 4
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
func (l *Items) Register(key string, name string, duration int, warning uint64, critical uint64, value uint64, message string) {
	check := items.Exist(key)

	if check == nil {
		c := Check{
			Critical: critical,
			Duration: duration,
			FirstSeen: int(time.Now().Unix()),
			Key: key,
			LastSeen: int(time.Now().Unix()),
			Message: message,
			Name: name,
			BeforeValue: value,
			Status: Normal,
			Value: value,
			Warning: warning,
		}

		*l = append(*l, c)

		log.Debug(fmt.Sprintf("Alert - Insert: %#v", c))
	} else if l != nil {
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

func (c *Check) Notify() bool {
	if c == nil {
		return false
	}

	switch status := c.Status; status {
	case Normal:
		if c.Between(c.Value) == Normal {
			c.FirstSeen = c.LastSeen
		} else if c.Between(c.Value) == Warning || c.Between(c.Value) == Critical {
			if c.Between(c.BeforeValue) == Warning || c.Between(c.BeforeValue) == Critical {
				if c.Delay() >= c.Duration {
					c.Status = Notified
					return true
				}
			}
		}
	case Notified:
		if c.Between(c.Value) == Normal {
			c.Status = Recovered
			return true
		}
	case Recovered:
		if c.Between(c.Value) == Normal {
			c.FirstSeen = c.LastSeen
			c.Status = Normal
		}
	}

	c.BeforeValue = c.Value

	return false
}

func (c *Check) Delay() int {
	return c.LastSeen - c.FirstSeen
}

func (c *Check) Between(value uint64) uint8 {
	if value >= c.Warning {
		if value >= c.Critical {
			return Critical
		}
		return Warning
	}

	return Normal
}
