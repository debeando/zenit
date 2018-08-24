package accumulator

type Tag struct {
	Name  string
	Value string
}

type Metric struct {
	Key    string
	Tags   []Tag
	Values interface{}
}

type Value struct {
	Key   string
	Value interface{}
}

type Items []Metric

var items *Items

func Load() *Items {
	if items == nil {
		items = &Items{}
	}
	return items
}

func (l *Items) Reset() {
	*l = (*l)[:0]
}

func (l *Items) Count() int {
	return len(*l)
}

func (l *Items) AddItem(m Metric) {
	if !items.Unique(m) {
		*l = append(*l, m)
	} else {
		items.SumValues(m)
	}
}

func (l *Items) Unique(m Metric) bool {
	for _, i := range *l {
		if i.Key == m.Key && TagsEquals(i.Tags, m.Tags) {
			return true
		}
	}
	return false
}

func (l *Items) SumValues(m Metric) {
	for item_index := 0; item_index < len(*l); item_index++ {
		if (*l)[item_index].Key == m.Key && TagsEquals((*l)[item_index].Tags, m.Tags) == true {
			for item_value_index, item_value := range (*l)[item_index].Values.([]Value) {
				for _, metric_value := range m.Values.([]Value) {
					if item_value.Key == metric_value.Key {
						sum_value := metric_value.Value.(uint)
						old_value := (*l)[item_index].Values.([]Value)[item_value_index].Value.(uint)
						new_value := old_value + sum_value

						(*l)[item_index].Values.([]Value)[item_value_index].Value = new_value

						break
					}
				}
			}
		}
	}
}

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
