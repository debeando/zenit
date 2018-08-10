// TODO: Add more tests;
// - example of process use.
// - accumulator.Value by another data type.
// - test for Reset().

package accumulator_test

import (
  "testing"

  "github.com/swapbyt3s/zenit/plugins/accumulator"
)

var a = accumulator.Load()

func TestCount(t *testing.T) {
  if a.Count() != 0 {
    t.Error("Expected 0, got ", 1)
  }
}

func TestTagsEquals(t *testing.T) {
  result := accumulator.TagsEquals([]accumulator.Tag{accumulator.Tag{"baz", "foo"}},
                                   []accumulator.Tag{accumulator.Tag{"baz", "foo"}})

  if ! result == true {
    t.Error("Expected false")
  }

  result = accumulator.TagsEquals([]accumulator.Tag{accumulator.Tag{"bar", "foo1"}},
                                  []accumulator.Tag{accumulator.Tag{"bar", "foo1"},
                                                    accumulator.Tag{"baz", "foo2"}})

  if ! result == false {
    t.Error("Expected true")
  }
}

func TestAddItem(t *testing.T) {
  a.AddItem(accumulator.Metric{
    Key: "test",
    Tags: []accumulator.Tag{accumulator.Tag{"foo", "bar"}},
    Values: 123,
  })

  if a.Count() != 1 {
    t.Error("Expected 1")
  }
}

func TestUnique(t *testing.T) {
  if a.Count() == 0 {
    t.Error("Expected > 0")
  }

  result := a.Unique(accumulator.Metric{
    Key: "test",
    Tags: []accumulator.Tag{accumulator.Tag{"foo", "bar"}},
    Values: 123,
  })

  if result == false {
    t.Error("Expected true")
  }

  result = a.Unique(accumulator.Metric{
    Key: "test",
    Tags: []accumulator.Tag{accumulator.Tag{"foo", "baz"}},
    Values: 123,
  })

  if result == true {
    t.Error("Expected false")
  }
}

func TestSumValues(t *testing.T) {
  a.AddItem(accumulator.Metric{
    Key: "test_sum_values",
    Tags: []accumulator.Tag{accumulator.Tag{"foo", "bar"}},
    Values: []accumulator.Value{accumulator.Value{"a", uint(1)},
                                accumulator.Value{"b", uint(1)}},
  })

  a.AddItem(accumulator.Metric{
    Key: "test_sum_values",
    Tags: []accumulator.Tag{accumulator.Tag{"foo", "bar"}},
    Values: []accumulator.Value{accumulator.Value{"a", uint(2)},
                                accumulator.Value{"b", uint(2)}},
  })

  if a.Count() != 2 {
    t.Error("Expected 1")
  }
}
