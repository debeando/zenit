package newrelic

import (
	"encoding/json"
	"fmt"

	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/common/http"
	"github.com/swapbyt3s/zenit/plugins/outputs"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

type OutputNewrelicInsights struct {}

func (l *OutputNewrelicInsights) Collect() {
	defer func () {
		if err := recover(); err != nil {
			log.Debug(fmt.Sprintf("Plugin - OutputNewrelicInsights - Panic (code %d) has been recover from somewhere.\n", err))
		}
	}()

	if ! config.File.Newrelic.Insight.Enable {
		return
	}

	events := Normalize(metrics.Load())

	for j := range events {
		events_json, err := json.Marshal(events[j])
		if err != nil {
			log.Error(fmt.Sprintf("Fail parsing to JSON: %s", err))
			return;
		}

		log.Debug(fmt.Sprintf("Plugin - OutputNewrelicInsights - JSON Event: %s", string(events_json)))

		headers := make(map[string]string)
		headers["Content-Type"] = "application/json"
		headers["X-Insert-Key"] = config.File.Newrelic.Insight.InsertKey

		url := fmt.Sprintf("https://insights-collector.newrelic.com/v1/accounts/%s/events", config.File.Newrelic.Insight.AccountID)

		http.Post(url, string(events_json), headers)
	}
}

func init() {
	outputs.Add("OutputNewrelicInsights", func() outputs.Output { return &OutputNewrelicInsights{} })
}
