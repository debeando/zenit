package mem_test

import (
  "testing"
  "time"

  "github.com/swapbyt3s/zenit/config"
  "github.com/swapbyt3s/zenit/plugins/alerts/os/mem"
  "github.com/swapbyt3s/zenit/plugins/lists/checks"
  "github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

func TestMain(m *testing.M) {
  // Configure:
  config.File.OS.Alerts.MEM.Enable   = true
  config.File.OS.Alerts.MEM.Duration = 4
  config.File.OS.Alerts.MEM.Warning  = 85
  config.File.OS.Alerts.MEM.Critical = 90

  // Run Tests:
  m.Run()
}

func TestCPU(t *testing.T) {
  var histogram = []struct{
    Value uint64
    Status uint8
    Notify bool
  }{
    { 80, checks.Normal   , false }, // 1s
    { 80, checks.Normal   , false }, // 2s
    { 88, checks.Normal   , false }, // 3s
    { 80, checks.Normal   , false }, // 4s
    { 80, checks.Normal   , false }, // 5s
    { 86, checks.Normal   , false }, // 6s
    { 85, checks.Normal   , false }, // 7s
    { 87, checks.Normal   , false }, // 8s
    { 85, checks.Notified , true  }, // 9s
    { 95, checks.Normal   , false }, // 10s
    { 80, checks.Recovered, true  }, // 11s
    { 80, checks.Normal   , false }, // 12s
  }

  for second, variable := range histogram {
    // Add test value on metrics:
    metrics.Load().Reset()
    metrics.Load().Add(metrics.Metric{
      Key: "zenit_os",
      Tags: []metrics.Tag{
        {"name", "mem"},
      },
      Values: variable.Value,
    })

    // Register alert:
    var c mem.OSMEM
    c.Collect()

    // Evaluate alert status
    check := checks.Load().Exist("mem")
    notify := check.Notify()

    if ! (check.Status == variable.Status && variable.Notify == notify) {
      t.Errorf("Second: %d, Value: %d, Evaluated: %t, Expected: %d, Got: %d.",
        second,
        variable.Value,
        notify,
        variable.Status,
        check.Status,
      )
    }

    // Wait:
    time.Sleep(1 * time.Second)
  }
}
