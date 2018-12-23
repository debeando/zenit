package metrics_test

import (
	"testing"

	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

var a = metrics.Load()

func TestCount(t *testing.T) {
	if a.Count() != 0 {
		t.Error("Expected 0, got ", 1)
	}
}

func TestTagsEquals(t *testing.T) {
	result := metrics.TagsEquals([]metrics.Tag{{"baz", "foo"}},
		[]metrics.Tag{{"baz", "foo"}})

	if !result == true {
		t.Error("Expected false")
	}

	result = metrics.TagsEquals([]metrics.Tag{{"bar", "foo1"}},
		[]metrics.Tag{{"bar", "foo1"},
			{"baz", "foo2"}})

	if !result == false {
		t.Error("Expected true")
	}
}

func TestAdd(t *testing.T) {
	a.Add(metrics.Metric{
		Key:    "test",
		Tags:   []metrics.Tag{{"foo", "bar"}},
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

	result := a.Unique(metrics.Metric{
		Key:    "test",
		Tags:   []metrics.Tag{{"foo", "bar"}},
		Values: 123,
	})

	if result == false {
		t.Error("Expected true")
	}

	result = a.Unique(metrics.Metric{
		Key:    "test",
		Tags:   []metrics.Tag{{"foo", "baz"}},
		Values: 123,
	})

	if result == true {
		t.Error("Expected false")
	}
}

func TestAccumulator(t *testing.T) {
	a.Add(metrics.Metric{
		Key:  "test_sum_values",
		Tags: []metrics.Tag{{"foo", "bar"}},
		Values: []metrics.Value{{"a", uint(1)},
			{"b", uint(1)}},
	})

	a.Add(metrics.Metric{
		Key:  "test_sum_values",
		Tags: []metrics.Tag{{"foo", "bar"}},
		Values: []metrics.Value{{"a", uint(2)},
			{"b", uint(2)}},
	})

	if a.Count() != 2 {
		t.Error("Expected 1")
	}
}

func TestFetchOne(t *testing.T) {
	a.Add(metrics.Metric{
		Key: "test_values",
		Tags: []metrics.Tag{
			{"name", "fulano"},
		},
		Values: 1,
	})

	a.Add(metrics.Metric{
		Key: "test_values",
		Tags: []metrics.Tag{
			{"name", "mengano"},
		},
		Values: 2,
	})

	a.Add(metrics.Metric{
		Key: "test_values",
		Tags: []metrics.Tag{
			{"name", "zutano"},
		},
		Values: 3,
	})

	value := a.FetchOne("test_values", "name", "mengano")

	t.Log(value)
}
