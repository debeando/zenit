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

	client "github.com/influxdata/influxdb1-client/v2"
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

	_, err := url.Parse(config.File.Outputs.InfluxDB.URL)
	if err != nil {
		log.Error(fmt.Sprintf("Plugin - OutputIndluxDB:Parser - Error: %s", err))
		return
	}


	conf := client.HTTPConfig{
		Addr: config.File.Outputs.InfluxDB.URL,
		Username: config.File.Outputs.InfluxDB.Username,
		Password: config.File.Outputs.InfluxDB.Password,
	}
	con, err := client.NewHTTPClient(conf)
	if err != nil {
		log.Error(fmt.Sprintf("Plugin - OutputIndluxDB:Client - Error: %s", err))
		return
	}
	defer con.Close()

	_, ver, err := con.Ping(0)
	if err != nil {
		log.Error(fmt.Sprintf("Plugin - OutputIndluxDB:Ping - Error: %s", err))
		return
	}

	log.Debug(fmt.Sprintf("Plugin - OutputIndluxDB - Connected to InfluxDB V-%s", ver))

	var events = Normalize(metrics.Load())
	var database = config.File.Outputs.InfluxDB.Database

	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database: database,
		Precision: "s",
	})

	for k, l := range events {
		for _, m := range l {
			pt, err := client.NewPoint(
				k,
				m["tags"].(map[string]string),
				m["fields"].(map[string]interface{}),
				time.Now(),
			)
			if err != nil {
				fmt.Println("Error: ", err.Error())
			}
			bp.AddPoint(pt)
		}
	}

	if con.Write(bp) != nil {
		if strings.Contains(err.Error(), errDatabaseNotFound) {
			query := client.NewQuery(fmt.Sprintf(
					`CREATE DATABASE "%s"`,
					config.File.Outputs.InfluxDB.Database,
				), "", "",
			)

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
