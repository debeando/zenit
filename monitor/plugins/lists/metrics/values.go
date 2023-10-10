package metrics

// Value is a collection for specific metric.
type Value struct {
	Key   string
	Value interface{}
}

// Values of a colletion of metric
type Values []Value

// Reset the metric metrics.
func (v *Values) Reset() {
	*v = (*v)[:0]
}

// Add is aggregator for value in metrics values.
func (v *Values) Add(m Value) {
	*v = append(*v, m)
}
