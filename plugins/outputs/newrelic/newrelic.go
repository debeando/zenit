package newrelic

import (
	"fmt"
	"encoding/json"

	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/common/http"
	"github.com/swapbyt3s/zenit/plugins/outputs"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

type OutputNewrelicInsights struct {}

func (l *OutputNewrelicInsights) Collect() {
	defer func () {
		if err := recover(); err != nil {
			fmt.Printf("Plugin - InputProxySQLCommands - Panic (code %d) has been recover from somewhere.\n", err)
		}
	}()

	if ! config.File.Newrelic.Insight.Enable {
		return
	}

	var events []interface{}

	for _, m := range *metrics.Load() {
		switch v := m.Values.(type) {
                case int, uint, uint64, float64:
                	event := make(map[string]interface{})
			event["eventType"] = common.ToCamel(m.Key)
			event["host"] = config.File.General.Hostname

			for t := range m.Tags {
				event[common.ToCamel(m.Tags[t].Value)] = v
			}

			events = append(events, event)
                case []metrics.Value:
                	for y := range v {
	                	event := make(map[string]interface{})
				event["eventType"] = common.ToCamel(m.Key)
				event["host"] = config.File.General.Hostname

				for t := range m.Tags {
					k := fmt.Sprintf("%s_%s", m.Tags[t].Value, v[y].Key)
					k = common.ToCamel(k)
					event[k] = v[y].Value
				}

				events = append(events, event)

                	}
                }
	}

        events_json, err := json.Marshal(events)
        if err != nil {
                log.Error(fmt.Sprintf("Fail parsing to JSON: %s", err))
                return;
        }

        log.Debug(fmt.Sprintf("OutputNewrelicInsights - JSON Event: %s", string(events_json)))

        headers := make(map[string]string)
        headers["Content-Type"] = "application/json"
        headers["X-Insert-Key"] = config.File.Newrelic.Insight.InsertKey

        url := fmt.Sprintf("https://insights-collector.newrelic.com/v1/accounts/%d/events", config.File.Newrelic.Insight.AccountID)

        http.Post(url, string(events_json), headers)
}

func init() {
	outputs.Add("OutputNewrelicInsights", func() outputs.Output { return &OutputNewrelicInsights{} })
}
