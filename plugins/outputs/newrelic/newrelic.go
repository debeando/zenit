package newrelic

import (
	"encoding/json"
	"fmt"

	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/common/http"
	"github.com/swapbyt3s/zenit/plugins/outputs"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

type OutputNewrelicInsights struct {}

var events = make(map[string]map[string]interface{})

func (l *OutputNewrelicInsights) Collect() {
	defer func () {
		if err := recover(); err != nil {
			log.Debug(fmt.Sprintf("Plugin - OutputNewrelicInsights - Panic (code %d) has been recover from somewhere.\n", err))
		}
	}()

	if ! config.File.Newrelic.Insight.Enable {
		return
	}

	for _, m := range *metrics.Load() {
		switch v := m.Values.(type) {
		case int, uint, uint64, float64:
			if _, ok := events[m.Key]; ! ok {
				events[m.Key] = make(map[string]interface{})
				events[m.Key]["host"] = config.File.General.Hostname
				events[m.Key]["eventType"] = common.ToCamel(m.Key)
			}

			for t := range m.Tags {
				if m.Tags[t].Name == "name" {
					events[m.Key][m.Tags[t].Value] = v
				} else {
					events[m.Key][m.Tags[t].Name] = m.Tags[t].Value
				}
			}
		case []metrics.Value:
			if _, ok := events[m.Key]; ! ok {
				events[m.Key] = make(map[string]interface{})
				events[m.Key]["host"] = config.File.General.Hostname
				events[m.Key]["eventType"] = common.ToCamel(m.Key)
			}

			for y := range v {
				events[m.Key][v[y].Key] = v[y].Value
			}
		}
	}

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
