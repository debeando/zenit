package influxdb

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"zenit/agent/plugins/lists/metrics"
	"zenit/agent/plugins/outputs"
	"zenit/config"

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
			log.ErrorWithFields(name, log.Fields{"message": err})
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
		log.ErrorWithFields(name, log.Fields{"step": "parser", "message": err})
		return
	}

	conf := client.HTTPConfig{
		Addr:     cnf.Outputs.InfluxDB.URL,
		Username: cnf.Outputs.InfluxDB.Username,
		Password: cnf.Outputs.InfluxDB.Password,
	}
	con, err := client.NewHTTPClient(conf)
	if err != nil {
		log.ErrorWithFields(name, log.Fields{"step": "client", "message": err})
		return
	}
	defer con.Close()

	_, ver, err := con.Ping(0)
	if err != nil {
		log.ErrorWithFields(name, log.Fields{"step": "ping", "message": err})
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

			if err != nil {
				log.ErrorWithFields(name, log.Fields{
					"step":    "event",
					"message": err,
					"tags":    m["tags"],
					"fields":  m["fields"],
				})
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
				log.ErrorWithFields(name, log.Fields{"step": "create", "message": err})
				return
			}
		} else {
			log.ErrorWithFields(name, log.Fields{"message": err})
		}
	}
}

func init() {
	outputs.Add("OutputIndluxDB", func() outputs.Output { return &Plugin{} })
}
