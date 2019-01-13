package checks_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/swapbyt3s/zenit/plugins/lists/checks"
)

func TestBetween(t *testing.T) {
	var histogram = []struct{
		Value uint64
		Status uint8
	}{
		{ 10, checks.Normal   },
		{ 50, checks.Warning  },
		{ 90, checks.Critical },
		{ 49, checks.Normal   },
		{ 89, checks.Warning  },
		{ 99, checks.Critical },
	}

	checks.Load().Register(
		"TestRange",
		"TestRange",
		0,
		50, // Warning
		90, // Critical
		0,
		"TestRange",
	)

	alert := checks.Load().Exist("TestRange")

	for _, variable := range histogram {
		r := alert.Between(variable.Value)

		if ! (r == variable.Status) {
			t.Errorf("Expected: '%d', got: '%d'", variable.Status, r)
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

	var histogram = []struct{
		Value uint64
		Status uint8
		Notify bool
	}{
		{  0, checks.Normal   , false }, // 0s
		{ 10, checks.Normal   , false }, // 1s
		{ 30, checks.Normal   , false }, // 2s
		{ 60, checks.Normal   , false }, // 3s
		{ 10, checks.Normal   , false }, // 4s
		{ 60, checks.Normal   , false }, // 5s
		{ 60, checks.Normal   , false }, // 6s
		{ 10, checks.Normal   , false }, // 7s
		{ 60, checks.Normal   , false }, // 8s
		{ 60, checks.Normal   , false }, // 9s
		{ 60, checks.Normal   , false }, // 10s
		{ 10, checks.Normal   , false }, // 11s
		{ 60, checks.Normal   , false }, // 12s
		{ 60, checks.Normal   , false }, // 13s
		{ 60, checks.Normal   , false }, // 14s
		{ 60, checks.Notified , true  }, // 15s
		{ 10, checks.Recovered, true  }, // 16s
		{ 30, checks.Normal   , false }, // 17s
		{ 60, checks.Normal   , false }, // 18s
		{ 90, checks.Normal   , false }, // 19s
		{ 90, checks.Normal   , false }, // 20s
		{ 10, checks.Normal   , false }, // 21s
		{ 60, checks.Normal   , false }, // 22s
		{ 90, checks.Normal   , false }, // 23s
		{ 90, checks.Normal   , false }, // 24s
		{ 90, checks.Notified , true  }, // 25s
		{ 10, checks.Recovered, true  }, // 26s
		{ 10, checks.Normal   , false }, // 27s
		{ 30, checks.Normal   , false }, // 28s
		{ 60, checks.Normal   , false }, // 29s
		{ 60, checks.Normal   , false }, // 30s
		{ 60, checks.Normal   , false }, // 31s
		{ 60, checks.Notified , true  }, // 32s
		{ 60, checks.Notified , false }, // 33s
		{ 60, checks.Notified , false }, // 34s
		{ 10, checks.Recovered, true  }, // 35s
	}

	fmt.Printf("  * Second\tLastSeen\tFirstSeen\tDelay\tMaxTime\tNotify\tValue\tBetween\tStatus\n")

	for second, variable := range histogram {
		time.Sleep(time.Second)

		checks.Load().Register(
			"test",
			"Test",
			4,
			50, // Warning
			90, // Critical
			variable.Value,
			"",
		)

		check := checks.Load().Exist("test")
		delay := check.Delay()
		notify := check.Notify()

		fmt.Printf("  - %d\t\t%d\t%d\t%d\t%d\t%t\t%d\t%d\t%d\n",
			second,
			check.LastSeen,
			check.FirstSeen,
			delay,          // Delay
			check.Duration, // MaxTime
			notify,         // Notify
			check.Value,
			check.Between(check.Value),
			check.Status,
		)

		if ! (check.Status == variable.Status && variable.Notify == notify) {
			t.Errorf("Second: %d, Value: %d, Evaluated: %t, Expected: '%d', Got: '%d'.",
				second,
				variable.Value,
				notify,
				variable.Status,
				check.Status,
			)
		}
	}
}
