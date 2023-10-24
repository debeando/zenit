package serverstatus

import (
	"encoding/json"
	"reflect"

	"zenit/agent/plugins/inputs"
	"zenit/agent/plugins/lists/metrics"
	"zenit/config"

	"github.com/debeando/go-common/cast"
	"github.com/debeando/go-common/log"
	"github.com/debeando/go-common/mongodb"
	"github.com/debeando/go-common/strings"
)

type Plugin struct {
	Name     string
	Hostname string
	Values   metrics.Values
}

func (p *Plugin) Collect(name string, cnf *config.Config, mtc *metrics.Items) {
	defer func() {
		if err := recover(); err != nil {
			log.ErrorWithFields(name, log.Fields{"message": err})
		}
	}()

	for host := range cnf.Inputs.MongoDB {
		log.DebugWithFields(name, log.Fields{
			"hostname":      cnf.Inputs.MongoDB[host].Hostname,
			"enable":        cnf.Inputs.MongoDB[host].Enable,
			"server_status": cnf.Inputs.MongoDB[host].ServerStatus,
		})

		if !cnf.Inputs.MongoDB[host].Enable {
			continue
		}

		if !cnf.Inputs.MongoDB[host].ServerStatus {
			continue
		}

		log.InfoWithFields(name, log.Fields{
			"hostname": cnf.Inputs.MongoDB[host].Hostname,
		})

		m := mongodb.New(cnf.Inputs.MongoDB[host].Hostname, cnf.Inputs.MongoDB[host].DSN)
		if err := m.Connect(); err != nil {
			continue
		}

		ss := m.ServerStatus()

		p.Name = name
		p.Hostname = cnf.Inputs.MongoDB[host].Hostname

		var obj map[string]interface{}
		t, _ := json.Marshal(ss)
		json.Unmarshal(t, &obj)

		entry := reflect.ValueOf(obj)
		p.Iterate([]string{""}, entry)

		mtc.Add(metrics.Metric{
			Key: "mongodb_serverstatus",
			Tags: []metrics.Tag{
				{Name: "hostname", Value: cnf.Inputs.MongoDB[host].Hostname},
			},
			Values: p.Values,
		})

		m.Close()
		p.Values.Reset()
	}
}

func (p *Plugin) Iterate(parent []string, data reflect.Value) {
	switch data.Kind() {
	case reflect.Map:
		for _, key := range data.MapKeys() {
			value := data.MapIndex(key)

			if value.IsZero() {
				continue
			}
			if value.IsNil() {
				continue
			}

			parent = append(parent, key.String())

			p.Iterate(parent, reflect.ValueOf(value.Interface()))

			parent = parent[:len(parent)-1]
		}
	default:
		if cast.InterfaceIsNumber(data.Interface()) {
			key := strings.ToCamel(strings.Join(parent, ""))
			log.DebugWithFields(p.Name, log.Fields{
				"hostname": p.Hostname,
				key:        data,
			})

			p.Values.Add(
				metrics.Value{
					Key:   key,
					Value: cast.InterfaceToFloat64(data.Interface()),
				},
			)
		}
	}
}

func init() {
	inputs.Add("InputMongoDBServerStatus", func() inputs.Input { return &Plugin{} })
}
