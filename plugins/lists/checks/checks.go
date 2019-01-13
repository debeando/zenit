package checks

import (
	"fmt"
	"time"

	"github.com/swapbyt3s/zenit/common/log"
)

// Struct to manage each alerts.
type Check struct {
	// Critical value from config file.
	Critical uint64
	// Duration limit between FirstSeen and LastSeen to alert.
	Duration int
	// UnixTimeStamp when is register alert or recovered or in normal state.
	FirstSeen int
	// Key to identify own check.
	Key string
	// What is the last seen (in UnixTimeStamp) value to rest with first seen and compare with duration.
	LastSeen int
	// Custom message from alert plugin.
	Message string
	// Name human own check from alert plugin.
	Name string
	// Current status for the check: normal, warning, critical or Recovered.
	Status uint8
	// Value to evaluate become input plugin.
	Value uint64
	// Last evaluated value.
	BeforeValue uint64
	// Warning value from config file..
	Warning uint64
}

// Alert levels
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

// Return all keys from the list.
func (l *Items) Keys() []string {
	var keys = []string{}

	for _, i := range *l {
		keys = append(keys, i.Key)
	}
	return keys
}

// Check de alert by key name exist on the list and return the item.
func (l *Items) Exist(key string) *Check {
	for i := 0; i < len(*l); i++ {
		if (*l)[i].Key == key {
			return &(*l)[i]
		}
	}
	return nil
}

// Verify the item on list is a valid alert.
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

// Calculate difference between dates
func (c *Check) Delay() int {
	return c.LastSeen - c.FirstSeen
}

// Check the value is between range of warning and critical.
func (c *Check) Between(value uint64) uint8 {
	if value >= c.Warning {
		if value >= c.Critical {
			return Critical
		}
		return Warning
	}

	return Normal
}
