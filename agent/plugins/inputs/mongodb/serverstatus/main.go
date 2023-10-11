package serverstatus

import (
	"fmt"
	"reflect"

	"zenit/config"
	"zenit/agent/plugins/inputs"
	"zenit/agent/plugins/lists/metrics"

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
			"hostname":  cnf.Inputs.MongoDB[host].Hostname,
			"enable":    cnf.Inputs.MongoDB[host].Enable,
			"variables": cnf.Inputs.MongoDB[host].ServerStatus,
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
		m.Connect()
		ss := m.GetServerStatus()
		data := reflect.ValueOf(ss)

		p.Name = name
		p.Hostname = cnf.Inputs.MongoDB[host].Hostname
		p.iterate("", data)

		mtc.Add(metrics.Metric{
			Key: "mongodb_serverstatus",
			Tags: []metrics.Tag{
				{Name: "hostname", Value: cnf.Inputs.MongoDB[host].Hostname},
			},
			Values: p.Values,
		})

		p.Values.Reset()
	}
}

func (p *Plugin) iterate(prefixKey string, data reflect.Value) {
	switch data.Kind() {
	case reflect.Map:
		for _, mapKey := range data.MapKeys() {
			mapValue := data.MapIndex(mapKey).Interface()

			if reflect.TypeOf(mapValue).Kind() == reflect.Map {
				p.iterate(mapKey.String(), reflect.ValueOf(mapValue))
			} else {
				if cast.InterfaceIsNumber(mapValue) {
					key := strings.ToCamel(fmt.Sprintf("%s %s", prefixKey, mapKey.String()))
					log.DebugWithFields(p.Name, log.Fields{
						"hostname": p.Hostname,
						key:        mapValue,
					})
					p.Values.Add(metrics.Value{Key: key, Value: mapValue})
				}
			}
		}
	}
}

func init() {
	inputs.Add("InputMongoDBServerStatus", func() inputs.Input { return &Plugin{} })
}
