package output

type Tag struct {
  Name  string
  Value string
}

type Metric struct {
  Key     string
  Tags  []Tag
  Value   float64
}

type Accumulator struct {
  Items []Metric
}

var accumulator *Accumulator

func LoadAccumulator() *Accumulator {
  if accumulator == nil {
    accumulator = &Accumulator{}
  }
  return accumulator
}

func (a *Accumulator) AddItem(item Metric) []Metric {
  a.Items = append(a.Items, item)
  return a.Items
}
