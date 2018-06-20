package output

type Tag struct {
  Name  string
  Value string
}

type Metric struct {
  Key      string
  Tags   []Tag
  Values   interface{}
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

func (l *Items) AddItem(m Metric) {
  if ! items.Unique(m) {
    *l = append(*l, m)
  } else {
    // Sum...
  }
}

func (l *Items) Unique(m Metric) bool {
  for _, i := range *l {
    if i.Key == m.Key {
      if len(i.Tags) != len(m.Tags) {
        return false
      }
      for idx := range m.Tags {
        if i.Tags[idx] != m.Tags[idx] {
          return false
        }
      }
      return true
    }
  }
  return false
}
