package alerts_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/swapbyt3s/zenit/plugins/lists/alerts"
)

func TestBetween(t *testing.T) {
	var checks = []struct{
		Value uint64
		Status uint8
	}{
		{ 10, alerts.Normal   },
		{ 50, alerts.Warning  },
		{ 90, alerts.Critical },
		{ 49, alerts.Normal   },
		{ 89, alerts.Warning  },
		{ 99, alerts.Critical },
	}

	alerts.Load().Register(
		"TestRange",
		"TestRange",
		0,
		50, // Warning
		90, // Critical
		0,
		"TestRange",
	)

	alert := alerts.Load().Exist("TestRange")

	for _, check := range checks {
		r := alert.Between(check.Value)

		if ! (r == check.Status) {
			t.Errorf("Expected: '%d', got: '%d'", check.Status, r)
		}
	}
}

func TestNotify(t *testing.T) {
// ASCII Graph for generic alert:
//
// 100 ┤-------------------------------------------------------------------------------------------------------
//  90 ┤                                                                        ╭───────╮       ╭─────────c─╮
//  80 ┤                                                                        │       │       │           │
//  70 ┤                                                                        │       │       │           │
//  60 ┤        ╭───╮   ╭───────╮   ╭───────────╮   ╭─────────────w─╮       ╭───╯       │   ╭───╯           │
//  50 ┤        │   │   │       │   │           │   │               │       │           │   │               │
//  40 ┤        │   │   │       │   │           │   │               │       │           │   │               │
//  30 ┤    ╭───╯   │   │       │   │           │   │               │   ╭───╯           │   │               │
//  20 ┤    │       │   │       │   │           │   │               │   │               │   │               │
//  10 ┤────╯       ╰───╯       ╰───╯           ╰───╯               ╰─r─╯               ╰───╯               ╰─r
//   0 ---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
//       01  02  03  04  05  06  07  08  09  10  11  12  13  14  15  16  17  18  19  20  21  22  23  24  25  26
//
// 100 ┤-------------------------------------------------------------------------------------------------------
//  90 ┤                                                           ╭─────────c───────╮
//  80 ┤                                                           │                 │
//  70 ┤                                                           │                 │
//  60 ┤        ╭─────────────w─────────╮        ╭─────────────w───╯                 │
//  50 ┤        │                       │        │                                   │
//  40 ┤        │                       │        │                                   │
//  30 ┤    ╭───╯                       │    ╭───╯                                   │
//  20 ┤    │                           │    │                                       │
//  10 ┤────╯                           ╰─r──╯                                       ╰r─
//   0 ---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
//       27  28  29  30  31  32  33  34  35 (            No test implemented            )
//
// Legend:
// - X = Time (s)
// - Y = Percentage (%)
// - w = Notify Warning
// - c = Notify Critical
// - r = Notify Recovered
//

	var checks = []struct{
		Value uint64
		Status uint8
		Notify bool
	}{
		{  0, alerts.Normal   , false }, // 0s
		{ 10, alerts.Normal   , false }, // 1s
		{ 30, alerts.Normal   , false }, // 2s
		{ 60, alerts.Normal   , false }, // 3s
		{ 10, alerts.Normal   , false }, // 4s
		{ 60, alerts.Normal   , false }, // 5s
		{ 60, alerts.Normal   , false }, // 6s
		{ 10, alerts.Normal   , false }, // 7s
		{ 60, alerts.Normal   , false }, // 8s
		{ 60, alerts.Normal   , false }, // 9s
		{ 60, alerts.Normal   , false }, // 10s
		{ 10, alerts.Normal   , false }, // 11s
		{ 60, alerts.Normal   , false }, // 12s
		{ 60, alerts.Normal   , false }, // 13s
		{ 60, alerts.Normal   , false }, // 14s
		{ 60, alerts.Notified , true  }, // 15s
		{ 10, alerts.Recovered, true  }, // 16s
		{ 30, alerts.Normal   , false }, // 17s
		{ 60, alerts.Normal   , false }, // 18s
		{ 90, alerts.Normal   , false }, // 19s
		{ 90, alerts.Normal   , false }, // 20s
		{ 10, alerts.Normal   , false }, // 21s
		{ 60, alerts.Normal   , false }, // 22s
		{ 90, alerts.Normal   , false }, // 23s
		{ 90, alerts.Normal   , false }, // 24s
		{ 90, alerts.Notified , true  }, // 25s
		{ 10, alerts.Recovered, true  }, // 26s
		{ 10, alerts.Normal   , false }, // 27s
		{ 30, alerts.Normal   , false }, // 28s
		{ 60, alerts.Normal   , false }, // 29s
		{ 60, alerts.Normal   , false }, // 30s
		{ 60, alerts.Normal   , false }, // 31s
		{ 60, alerts.Notified , true  }, // 32s
		{ 60, alerts.Notified , false }, // 33s
		{ 60, alerts.Notified , false }, // 34s
		{ 10, alerts.Recovered, true  }, // 35s
	}

	fmt.Printf("  * Second\tLastSeen\tFirstSeen\tDelay\tMaxTime\tNotify\tValue\tBetween\tStatus\n")

	for second, check := range checks {
		time.Sleep(time.Second)

		alerts.Load().Register(
			"test",
			"Test",
			4,
			50, // Warning
			90, // Critical
			check.Value,
			"",
		)

		alert := alerts.Load().Exist("test")
		delay := alert.Delay()
		notify := alert.Notify()

		fmt.Printf("  - %d\t\t%d\t%d\t%d\t%d\t%t\t%d\t%d\t%d\n",
			second,
			alert.LastSeen,
			alert.FirstSeen,
			delay,          // Delay
			alert.Duration, // MaxTime
			notify,         // Notify
			alert.Value,
			alert.Between(alert.Value),
			alert.Status,
		)

		if ! (alert.Status == check.Status && check.Notify == notify) {
			t.Errorf("Second: %d, Value: %d, Evaluated: %t, Expected: '%d', Got: '%d'.",
				second,
				check.Value,
				notify,
				check.Status,
				alert.Status,
			)
		}
	}
}
