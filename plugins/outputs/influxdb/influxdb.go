package influxdb

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
	"github.com/swapbyt3s/zenit/plugins/outputs"

	client "github.com/influxdata/influxdb1-client"
)

const (
	defaultDatabase     = "zenit"
	errDatabaseNotFound = "database not found"
)

type OutputIndluxDB struct{}

func (l *OutputIndluxDB) Collect() {
	defer func() {
		if err := recover(); err != nil {
			log.Debug(fmt.Sprintf("Plugin - OutputIndluxDB - Panic (code %d) has been recover from somewhere.\n", err))
		}
	}()

	if !config.File.Outputs.InfluxDB.Enable {
		return
	}

	if config.File.Outputs.InfluxDB.Database == "" {
		config.File.Outputs.InfluxDB.Database = defaultDatabase
	}

	host, err := url.Parse(config.File.Outputs.InfluxDB.URL)
	if err != nil {
		log.Error(fmt.Sprintf("Plugin - OutputIndluxDB:Parser - Error: %s", err))
		return
	}

	conf := client.Config{
		URL:      *host,
		Username: config.File.Outputs.InfluxDB.Username,
		Password: config.File.Outputs.InfluxDB.Password,
	}

	con, err := client.NewClient(conf)
	if err != nil {
		log.Error(fmt.Sprintf("Plugin - OutputIndluxDB:Client - Error: %s", err))
		return
	}

	_, ver, err := con.Ping()
	if err != nil {
		log.Error(fmt.Sprintf("Plugin - OutputIndluxDB:Ping - Error: %s", err))
		return
	}

	log.Debug(fmt.Sprintf("Plugin - OutputIndluxDB - Connected to InfluxDB V-%s", ver))

	pts := make([]client.Point, 1000)

	events := metrics.Load()

	for i, m := range *events {
		tags   := map[string]string{}
		values := make(map[string]interface{})

		tags["hostname"] = config.File.General.Hostname

		switch v := m.Values.(type) {
			case []metrics.Value:
				for y := range v {
					log.Debug(fmt.Sprintf("Plugin - OutputIndluxDB - Metric: [%s:%#v]", v[y].Key, v[y].Value))

					values[v[y].Key] = v[y].Value
				}
			default:
				for t := range m.Tags {
					if m.Tags[t].Name == "name" {
						log.Debug(fmt.Sprintf("Plugin - OutputIndluxDB - Metric: [%s:%#v]", m.Tags[t].Value, v))
						values[m.Tags[t].Value] = v
					}
				}
		}

		pts[i] = client.Point{
			Measurement: m.Key,
			Tags:      tags,
			Fields:    values,
			Time:      time.Now(),
			Precision: "s",
		}
	}

	bps := client.BatchPoints{
		Points:          pts,
		Database:        config.File.Outputs.InfluxDB.Database,
	}

	_, err = con.Write(bps)
	if err != nil {
		if strings.Contains(err.Error(), errDatabaseNotFound) {
			query := client.Query{
				Command:  fmt.Sprintf(
					`CREATE DATABASE "%s"`,
					config.File.Outputs.InfluxDB.Database,
				),
			}

			log.Debug(fmt.Sprintf(
				"Plugin - OutputIndluxDB:CreateDatabase %s",
				config.File.Outputs.InfluxDB.Database,
			))

			if _, err := con.Query(query); err != nil {
				log.Error(fmt.Sprintf("Plugin - OutputIndluxDB:CreateDatabase - Error: %s", err))
				return
			}
		} else {
			log.Error(fmt.Sprintf("Plugin - OutputIndluxDB:Write - Error: %s", err))
		}
	}
}

func init() {
	outputs.Add("OutputIndluxDB", func() outputs.Output { return &OutputIndluxDB{} })
}
