package alerts_test

import (
	"testing"

	"github.com/swapbyt3s/zenit/plugins/lists/alerts"
)

func TestEvaluate(t *testing.T) {
	var checks = []struct{
		Key       int
		FirstSeen int
		LastSeen  int
		Status    uint8
		Duration  int
		Value     uint64
		Warning   uint64
		Critical  uint64
		Notify    bool
	}{
		{ 1, 1537705900, 1537705910, alerts.Normal , 10, 1, 2, 3, false},
		{ 3, 1537705900, 1537705910, alerts.Normal , 10, 2, 2, 3, true },
		{ 4, 1537705900, 1537705910, alerts.Normal , 10, 3, 2, 3, true },
		{ 5, 1537705900, 1537705905, alerts.Normal , 10, 2, 2, 3, false},
		{ 6, 1537705900, 1537705905, alerts.Normal , 10, 3, 2, 3, false},
		{ 7, 1537705900, 1537705910, alerts.Normal , 10, 2, 2, 2, true },
		{ 8, 1537705900, 1537705905, alerts.Normal , 10, 2, 2, 2, false},
		{ 9, 1537705900, 1537705910, alerts.Normal , 10, 0, 1, 1, false},
		{10, 1537705900, 1537705910, alerts.Normal , 10, 1, 1, 1, true },
		{11, 1537705900, 1537705910, alerts.Warning, 10, 2, 2, 3, true },
		{12, 1537705900, 1537705910, alerts.Warning, 10, 3, 2, 3, true },
	}

	for _, check := range checks {
		test := alerts.Check {
			FirstSeen: check.FirstSeen,
			LastSeen: check.LastSeen,
			Status: check.Status,
			Duration: check.Duration,
			Value: check.Value,
			Warning: check.Warning,
			Critical: check.Critical,
		}

		if test.Evaluate() != check.Notify {
			t.Logf("\nCheck: %#v\nTest: %#v", check, test)
		}
	}

	// Check live cicle:
	t1 := alerts.Check {
		FirstSeen: 1537705900,
		Status: alerts.Normal,
		Duration: 10,
		Warning: 60,
		Critical: 80,
	}

	var states = []struct{
		Key      int
		LastSeen int
		Value    uint64
		Status   uint8
		Expected bool
	}{
		{1, 1537705910, 60, alerts.Warning , true},
		{2, 1537705920, 80, alerts.Critical, true},
		{3, 1537705925, 70, alerts.Critical, false},
		{4, 1537705940, 75, alerts.Warning , true},
		{5, 1537705945, 70, alerts.Warning , false},
		{6, 1537705960, 10, alerts.Resolved, true},
	}

	for _, s := range states {
		t1.LastSeen = s.LastSeen
		t1.Value = s.Value

		evaluated := t1.Evaluate()

		if evaluated != s.Expected {
			t.Errorf("Test %d - Expected: '%t', got: '%t'", s.Key, s.Expected, evaluated)
		}

		if t1.Status != s.Status {
			t.Errorf("Test %d - Expected: '%d', got: '%d'", s.Key, s.Status, t1.Status)
		}
	}
}
