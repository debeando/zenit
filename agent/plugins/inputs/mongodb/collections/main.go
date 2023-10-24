package collections

import (
	"zenit/agent/plugins/inputs"
	"zenit/agent/plugins/lists/metrics"
	"zenit/config"

	"github.com/debeando/go-common/log"
	"github.com/debeando/go-common/mongodb"
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
			"hostname":           cnf.Inputs.MongoDB[host].Hostname,
			"enable":             cnf.Inputs.MongoDB[host].Enable,
			"server_collections": cnf.Inputs.MongoDB[host].Collections,
		})

		if !cnf.Inputs.MongoDB[host].Enable {
			continue
		}

		if !cnf.Inputs.MongoDB[host].Collections {
			continue
		}

		log.InfoWithFields(name, log.Fields{
			"hostname": cnf.Inputs.MongoDB[host].Hostname,
		})

		m := mongodb.New(cnf.Inputs.MongoDB[host].Hostname, cnf.Inputs.MongoDB[host].DSN)
		if err := m.Connect(); err != nil {
			continue
		}

		databases := m.Databases()
		for _, database := range databases.Databases {
			collections := m.Collections(database.Name)
			for _, collection := range collections {
				colStats := m.CollectionStats(database.Name, collection)

				log.DebugWithFields(p.Name, log.Fields{
					"hostname":         cnf.Inputs.MongoDB[host].Hostname,
					"name":             colStats.Collection,
					"count":            colStats.Count,
					"size":             colStats.Size,
					"storage_size":     colStats.StorageSize,
					"total_index_size": colStats.TotalIndexSize,
				})

				mtc.Add(metrics.Metric{
					Key: "mongodb_collections",
					Tags: []metrics.Tag{
						{Name: "hostname", Value: cnf.Inputs.MongoDB[host].Hostname},
						{Name: "collection", Value: colStats.Collection},
					},
					Values: []metrics.Value{
						{Key: "count", Value: colStats.Count},
						{Key: "size", Value: colStats.Size},
						{Key: "storage_size", Value: colStats.StorageSize},
						{Key: "total_index_size", Value: colStats.TotalIndexSize},
					},
				})
			}
		}

		p.Values.Reset()
	}
}

func init() {
	inputs.Add("InputMongoDBCollections", func() inputs.Input { return &Plugin{} })
}
