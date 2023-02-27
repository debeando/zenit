package influxdb

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/debeando/zenit/common/log"
	"github.com/debeando/zenit/config"
	"github.com/debeando/zenit/plugins/lists/metrics"
	"github.com/debeando/zenit/plugins/outputs"

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
			log.Error("OutputIndluxDB", map[string]interface{}{"error": err})
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
		log.Error("OutputIndluxDB", map[string]interface{}{"step": "parser", "error": err})
		return
	}

	conf := client.HTTPConfig{
		Addr:     config.File.Outputs.InfluxDB.URL,
		Username: config.File.Outputs.InfluxDB.Username,
		Password: config.File.Outputs.InfluxDB.Password,
	}
	con, err := client.NewHTTPClient(conf)
	if err != nil {
		log.Error("OutputIndluxDB", map[string]interface{}{"step": "client", "error": err})
		return
	}
	defer con.Close()

	_, ver, err := con.Ping(0)
	if err != nil {
		log.Error("OutputIndluxDB", map[string]interface{}{"step": "ping", "error": err})
		return
	}

	log.Debug("OutputIndluxDB", map[string]interface{}{"version": ver})

	var events = Normalize(metrics.Load())
	var database = config.File.Outputs.InfluxDB.Database

	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  database,
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

			log.Debug(fmt.Sprintf("OutputIndluxDB\t%s", k), m["fields"].(map[string]interface{}))

			if err != nil {
				log.Error("OutputIndluxDB", map[string]interface{}{"step": "event", "error": err})
			}
			bp.AddPoint(pt)
		}
	}

	err = con.Write(bp)
	if err != nil {
		if strings.Contains(err.Error(), errDatabaseNotFound) {
			query := client.NewQuery(fmt.Sprintf(
				`CREATE DATABASE "%s"`,
				config.File.Outputs.InfluxDB.Database,
			), "", "",
			)

			log.Debug("OutputIndluxDB", map[string]interface{}{"database": config.File.Outputs.InfluxDB.Database})

			if _, err := con.Query(query); err != nil {
				log.Error("OutputIndluxDB", map[string]interface{}{"step": "create", "error": err})
				return
			}
		} else {
			log.Error("OutputIndluxDB", map[string]interface{}{"error": err})
		}
	}
}

func init() {
	outputs.Add("OutputIndluxDB", func() outputs.Output { return &OutputIndluxDB{} })
}
