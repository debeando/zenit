// TODO: Add check for every X time to force send to CH and purge the buffer.

package clickhouse

import (
	"time"

	"zenit/config"

	"github.com/debeando/go-common/http"
	"github.com/debeando/go-common/log"
	"github.com/debeando/go-common/mysql/sql"
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
	cnf := config.GetInstance()

	log.InfoWithFields("OutputClickHouse", log.Fields{"dsn": cnf.Outputs.ClickHouse.DSN})

	if http.Post(cnf.Outputs.ClickHouse.DSN, "SELECT 1;", map[string]string{}) != 200 {
		log.ErrorWithFields("OutputClickHouse", log.Fields{"error": "Impossible to connect."})
		return false
	}

	return true
}

func Send(e *Event, data <-chan map[string]string) {
	timeout := make(chan bool)
	ticker := time.NewTicker(time.Duration(e.Timeout) * time.Second)
	cnf := config.GetInstance()

	go func() {
		for range ticker.C {
			timeout <- true
		}
	}()

	for {
		select {
		case <-timeout:
			log.DebugWithFields("OutputClickHouse", log.Fields{"type": e.Type, "values": e.Values})

			if len(e.Values) > 0 {
				sql := sql.Insert(e.Schema, e.Table, e.Wildcard, e.Values)
				e.Values = []map[string]string{}

				log.DebugWithFields("OutputClickHouse", log.Fields{"type": e.Type, "values": sql})

				go http.Post(cnf.Outputs.ClickHouse.DSN, sql, nil)
			}
		case d := <-data:
			log.DebugWithFields("OutputClickHouse", log.Fields{"type": e.Type, "values": d})

			e.Values = append(e.Values, d)
			if len(e.Values) == e.Size {
				sql := sql.Insert(e.Schema, e.Table, e.Wildcard, e.Values)
				e.Values = []map[string]string{}

				log.DebugWithFields("OutputClickHouse", log.Fields{"type": e.Type, "query": sql})

				go http.Post(cnf.Outputs.ClickHouse.DSN, sql, nil)
			}
		}
	}
}
