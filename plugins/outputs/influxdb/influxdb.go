package influxdb

import (
	"fmt"
	"net/url"
	"strings"

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

	var points = make([]client.Point, 1000)
	var events = Normalize(metrics.Load())
	var i = 0

	for k, l := range events {
		for _, m := range l {
			points[i] = client.Point{
				Measurement: k,
				Tags:      m["tags"].(map[string]string),
				Fields:    m["fields"].(map[string]interface{}),
				Precision: "s",
			}
			i++
		}
	}

	bps := client.BatchPoints{
		Points:          points,
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
