package metrics

import (
	"reflect"
)

// Tag for metric.
type Tag struct {
	Name  string
	Value string
}

// Metric is a collection with many Tags and Values.
type Metric struct {
	Key    string
	Tags   []Tag
	Values Values
}

// Items is a collection of metrics
type Items []Metric

var items *Items

// Load is a singleton method to return same object.
func Load() *Items {
	if items == nil {
		items = &Items{}
	}
	return items
}

// Reset the metric metrics.
func (l *Items) Reset() {
	*l = (*l)[:0]
}

// Count all metrics in metrics.
func (l *Items) Count() int {
	return len(*l)
}

// Add is aggregator for metric in metrics.
func (l *Items) Add(m Metric) {
	if !items.Unique(m) {
		*l = append(*l, m)
	} else {
		items.Update(m)
	}
}

// FetchOne and return specific metric.
func (l *Items) FetchOne(key string, tagName string, tagValue string) interface{} {
	for _, metric := range *l {
		if metric.Key == key {
			for _, tag := range metric.Tags {
				if tag.Name == tagName && tag.Value == tagValue {
					return metric.Values
				}
			}
		}
	}

	return nil
}

// Unique is a check to verify the metric key is one in the metrics.
func (l *Items) Unique(m Metric) bool {
	for _, i := range *l {
		if i.Key == m.Key && TagsEquals(i.Tags, m.Tags) {
			return true
		}
	}
	return false
}

// metrics sum values when we have the same key.
func (l *Items) Update(m Metric) {
	for itemIndex := 0; itemIndex < len(*l); itemIndex++ {
		if (*l)[itemIndex].Key == m.Key && TagsEquals((*l)[itemIndex].Tags, m.Tags) == true {
			if reflect.TypeOf((*l)[itemIndex].Values) == reflect.TypeOf([]Value{}) {
				for itemValueIndex, itemValue := range (*l)[itemIndex].Values {
					for _, metricValue := range m.Values {
						if itemValue.Key == metricValue.Key {
							switch v := metricValue.Value.(type) {
							case int64:
								oldValue := (*l)[itemIndex].Values[itemValueIndex].Value.(int64)
								newValue := oldValue + v
								(*l)[itemIndex].Values[itemValueIndex].Value = newValue
							}

							break
						}
					}
				}
			}
		}
	}
}

// TagsEquals verify two Tags are equals.
func TagsEquals(a, b []Tag) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
