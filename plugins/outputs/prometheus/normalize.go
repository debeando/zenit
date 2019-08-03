package prometheus

import (
	"fmt"
	"strings"

	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

func Normalize(a *metrics.Items) string {
	var s string
	var e []string

	for _, m := range *a {
		switch m.Values.(type) {
		case int, int64, float64:
			s = fmt.Sprintf("%s{%s} %s", m.Key, getTags(m.Tags), getValue(m.Values))

			if config.File.General.Debug {
				log.Debug("Prometheus - " + s)
			}

			e = append(e, s)
		case []metrics.Value:
			for _, i := range m.Values.([]metrics.Value) {
				s = fmt.Sprintf("%s{%s,type=\"%s\"} %s", m.Key, getTags(m.Tags), i.Key, getValue(i.Value))

				if config.File.General.Debug {
					log.Debug("Prometheus - " + s)
				}

				e = append(e, s)
			}
		}
	}

	return strings.Join(e, "\n") + "\n"
}

func getTags(tags []metrics.Tag) string {
	s := []string{}
	for t := range tags {
		k := tags[t].Name
		v := strings.ToLower(tags[t].Value)
		s = append(s, fmt.Sprintf("%s=\"%s\"", k, v))
	}
	return strings.Join(s, ",")
}

func getValue(value interface{}) string {
	switch v := value.(type) {
	case int, int64:
		return fmt.Sprintf("%d", v)
	case float64:
		return fmt.Sprintf("%.2f", v)
	}

	return "0"
}
