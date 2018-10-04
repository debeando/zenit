package alerts_test

import (
	"testing"

	"github.com/swapbyt3s/zenit/plugins/lists/alerts"
)

func TestCheck(t *testing.T) {
	var checks = []struct{
		Key       int
		FirstSeen int
		LastSeen  int
		Status    uint8
		Duration  int
		Critical  uint64
		Warning   uint64
		Value     uint64
		Decide    bool
		ExpectedStatus bool
	}{
		{1, 1537705900, 1537705910, alerts.Normal, 10, 3, 2, 1, true,  false},
		{2, 1537705900, 1537705910, alerts.Normal, 10, 3, 2, 0, false, true},
		{3, 1537705900, 1537705910, alerts.Normal, 10, 3, 2, 2, true,  true},
		{4, 1537705900, 1537705910, alerts.Normal, 10, 3, 2, 3, true,  true},
		{5, 1537705900, 1537705905, alerts.Normal, 10, 3, 2, 2, true,  false},
		{6, 1537705900, 1537705905, alerts.Normal, 10, 3, 2, 3, true,  false},
		{7, 1537705900, 1537705910, alerts.Normal, 10, 2, 2, 2, true,  true},
		{8, 1537705900, 1537705905, alerts.Normal, 10, 2, 2, 2, true,  false},
		// {9, 1537705900, 1537705915, alerts.Normal, 10, 3, 2, 0, false, false},
	}

	for _, check := range checks {
		test := alerts.Check {
			FirstSeen: check.FirstSeen,
			LastSeen: check.LastSeen,
			Status: check.Status,
			Duration: check.Duration,
			Critical: check.Critical,
			Warning: check.Warning,
			Value: check.Value,
			Decide: check.Decide,
		}

		test.Evaluate()

		if check.ExpectedStatus != test.Notify() {
			t.Logf("\nCheck: %#v\nTest: %#v", check, test)
		}
	}
}
