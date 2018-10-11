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
		Value     int
		Warning   int
		Critical  int
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
		LastSeen: 1537705910,
		Status: alerts.Normal,
		Duration: 10,
		Value: 60,
		Warning: 60,
		Critical: 80,
	}

	t.Logf("Check Evaluate: %t\n", t1.Evaluate())
	t.Logf("Check Status: %d\n", t1.Status)

	t1.LastSeen = 1537705920
	t1.Value = 80

	t.Logf("Check Evaluate: %t\n", t1.Evaluate())
	t.Logf("Check Status: %d\n", t1.Status)

	t1.LastSeen = 1537705930
	t1.Value = 70

	t.Logf("Check Evaluate: %t\n", t1.Evaluate())
	t.Logf("Check Status: %d\n", t1.Status)

	t1.LastSeen = 1537705940
	t1.Value = 75

	t.Logf("Check Evaluate: %t\n", t1.Evaluate())
	t.Logf("Check Status: %d\n", t1.Status)

	t1.LastSeen = 1537705950
	t1.Value = 10

	t.Logf("Check Evaluate: %t\n", t1.Evaluate())
	t.Logf("Check Status: %d\n", t1.Status)
}
