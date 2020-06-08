// TODO: Add check for every X time to force send to CH and purge the buffer.

package clickhouse

import (
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
	log.Info("OutputClickHouse", map[string]interface{}{"dsn": config.File.Outputs.ClickHouse.DSN})

	if http.Post(config.File.Outputs.ClickHouse.DSN, "SELECT 1;", map[string]string{}) != 200 {
		log.Error("OutputClickHouse", map[string]interface{}{"error": "Impossible to connect."})
		return false
	}

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
			log.Debug("OutputClickHouse", map[string]interface{}{"type": e.Type, "values": e.Values})

			if len(e.Values) > 0 {
				sql := sql.Insert(e.Schema, e.Table, e.Wildcard, e.Values)
				e.Values = []map[string]string{}

				log.Debug("OutputClickHouse", map[string]interface{}{"type": e.Type, "values": sql})

				go http.Post(config.File.Outputs.ClickHouse.DSN, sql, nil)
			}
		case d := <-data:
			log.Debug("OutputClickHouse", map[string]interface{}{"type": e.Type, "values": d})

			e.Values = append(e.Values, d)
			if len(e.Values) == e.Size {
				sql := sql.Insert(e.Schema, e.Table, e.Wildcard, e.Values)
				e.Values = []map[string]string{}

				log.Debug("OutputClickHouse", map[string]interface{}{"type": e.Type, "query": sql})

				go http.Post(config.File.Outputs.ClickHouse.DSN, sql, nil)
			}
		}
	}
}
