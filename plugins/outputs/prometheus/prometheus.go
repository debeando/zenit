// TODO: Write in file.

package prometheus

import (
	"fmt"
	"log"
	"strings"

	"github.com/swapbyt3s/zenit/common/file"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/accumulator"
)

func Run() {
	file.Create(config.File.Prometheus.TextFile)
	file.Truncate(config.File.Prometheus.TextFile)

	var a = accumulator.Load()
	var s string
	var e []string

	for _, m := range *a {
		switch m.Values.(type) {
		case int, uint, uint64, float64:
			s = fmt.Sprintf("%s{%s} %s", m.Key, getTags(m.Tags), getValue(m.Values))
		case []accumulator.Value:
			for _, i := range m.Values.([]accumulator.Value) {
				s = fmt.Sprintf("%s{%s,type=\"%s\"} %s", m.Key, getTags(m.Tags), i.Key, getValue(i.Value))

				if config.File.General.Debug {
					log.Printf("D! - Prometheus - %s\n", s)
				}

				e = append(e, s)
			}
		}
	}

	file.Write(config.File.Prometheus.TextFile, strings.Join(e, "\n"))
}

func getTags(tags []accumulator.Tag) string {
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
	case int, uint, uint64:
		return fmt.Sprintf("%d", v)
	case float64:
		return fmt.Sprintf("%.2f", v)
	}

	return "0"
}
