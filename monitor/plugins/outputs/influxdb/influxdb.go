package influxdb

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"zenit/config"
	"zenit/monitor/plugins/lists/metrics"
	"zenit/monitor/plugins/outputs"

	"github.com/debeando/go-common/log"
	client "github.com/influxdata/influxdb1-client/v2"
)

const (
	defaultDatabase     = "zenit"
	errDatabaseNotFound = "database not found"
)

type Plugin struct{}

func (p *Plugin) Deliver(name string, cnf *config.Config, mtc *metrics.Items) {
	defer func() {
		if err := recover(); err != nil {
			log.ErrorWithFields(name, log.Fields{"error": err})
		}
	}()

	if !cnf.Outputs.InfluxDB.Enable {
		return
	}

	if cnf.Outputs.InfluxDB.Database == "" {
		cnf.Outputs.InfluxDB.Database = defaultDatabase
	}

	_, err := url.Parse(cnf.Outputs.InfluxDB.URL)
	if err != nil {
		log.ErrorWithFields(name, log.Fields{"step": "parser", "error": err})
		return
	}

	conf := client.HTTPConfig{
		Addr:     cnf.Outputs.InfluxDB.URL,
		Username: cnf.Outputs.InfluxDB.Username,
		Password: cnf.Outputs.InfluxDB.Password,
	}
	con, err := client.NewHTTPClient(conf)
	if err != nil {
		log.ErrorWithFields(name, log.Fields{"step": "client", "error": err})
		return
	}
	defer con.Close()

	_, ver, err := con.Ping(0)
	if err != nil {
		log.ErrorWithFields(name, log.Fields{"step": "ping", "error": err})
		return
	}

	log.DebugWithFields(name, log.Fields{"version": ver})

	var events = Normalize(mtc)
	var database = cnf.Outputs.InfluxDB.Database

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

			log.DebugWithFields(fmt.Sprintf("OutputIndluxDB\t%s", k), m["fields"].(log.Fields))

			if err != nil {
				log.ErrorWithFields(name, log.Fields{"step": "event", "error": err})
			}
			bp.AddPoint(pt)
		}
	}

	err = con.Write(bp)
	if err != nil {
		if strings.Contains(err.Error(), errDatabaseNotFound) {
			query := client.NewQuery(
				fmt.Sprintf(
					`CREATE DATABASE "%s"`,
					cnf.Outputs.InfluxDB.Database,
				), "", "",
			)

			log.DebugWithFields(name, log.Fields{"database": cnf.Outputs.InfluxDB.Database})

			if _, err := con.Query(query); err != nil {
				log.ErrorWithFields(name, log.Fields{"step": "create", "error": err})
				return
			}
		} else {
			log.ErrorWithFields(name, log.Fields{"error": err})
		}
	}
}

func init() {
	outputs.Add("OutputIndluxDB", func() outputs.Output { return &Plugin{} })
}
