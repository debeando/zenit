// TODO: Add check for every X time to force send to CH and purge the buffer.

package clickhouse

import (
	"fmt"
	"time"

	"github.com/swapbyt3s/zenit/common/http"
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/common/sql"
	"github.com/swapbyt3s/zenit/config"
)

type Event struct {
	Type     string
	Schema   string
	Table    string
	Size     int
	Timeout  int
	Wildcard map[string]string
	Values   []map[string]string
}

func Check() bool {
	log.Info(fmt.Sprintf("ClickHouse - DSN: %s", config.File.Outputs.ClickHouse.DSN))
	if http.Post(config.File.Outputs.ClickHouse.DSN, "SELECT 1;", map[string]string{}) != 200 {
		log.Error("ClickHouse - Impossible to connect.")
		return false
	}

	log.Info("ClickHouse - Connected successfully.")
	return true
}

func Send(e *Event, data <-chan map[string]string) {
	timeout := make(chan bool)
	ticker := time.NewTicker(time.Duration(e.Timeout) * time.Second)

	go func() {
		for range ticker.C {
			timeout <- true
		}
	}()

	for {
		select {
		case <-timeout:
			log.Debug(fmt.Sprintf("ClickHouse - Event timeout: %s - %#v", e.Type, e.Values))
			if len(e.Values) > 0 {
				sql := sql.Insert(e.Schema, e.Table, e.Wildcard, e.Values)
				e.Values = []map[string]string{}

				log.Debug(fmt.Sprintf("ClickHouse - Event insert: %s - %s", e.Type, sql))

				go http.Post(config.File.Outputs.ClickHouse.DSN, sql, nil)
			}
		case d := <-data:
			log.Debug(fmt.Sprintf("ClickHouse - Event capture: %s - %#v", e.Type, d))

			e.Values = append(e.Values, d)
			if len(e.Values) == e.Size {
				sql := sql.Insert(e.Schema, e.Table, e.Wildcard, e.Values)
				e.Values = []map[string]string{}

				log.Debug(fmt.Sprintf("ClickHouse - Event insert: %s - %s", e.Type, sql))

				go http.Post(config.File.Outputs.ClickHouse.DSN, sql, nil)
			}
		}
	}
}
