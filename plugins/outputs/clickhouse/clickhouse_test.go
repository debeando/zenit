package clickhouse_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/swapbyt3s/zenit/common/sql"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/outputs/clickhouse"

	test "github.com/swapbyt3s/zenit/common/http"
)

var ts *httptest.Server
var event = &clickhouse.Event{
	Type:    "SlowLog",
	Schema:  "zenit",
	Table:   "mysql_slow_log",
	Size:    2,
	Timeout: 1,
	Wildcard: map[string]string{
		"_time":         "'%s'",
		"bytes_sent":    "%s",
		"host_ip":       "IPv4StringToNum('%s')",
		"host_name":     "'%s'",
		"killed":        "%s",
		"last_errno":    "%s",
		"lock_time":     "%s",
		"query":         "'%s'",
		"query_digest":  "'%s'",
		"query_time":    "%s",
		"rows_affected": "%s",
		"rows_examined": "%s",
		"rows_read":     "%s",
		"rows_sent":     "%s",
		"schema":        "'%s'",
		"thread_id":     "%s",
		"user_host":     "'%s'",
	},
	Values: []map[string]string{},
}

func HandlerResponse() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "database=zenit" {
			http.Error(w, "ERROR_PAGE_NOT_FOUND", http.StatusNotFound)
			return
		}

		if r.Method != "POST" {
			http.Error(w, "ERROR_BAD_REQUEST", http.StatusBadRequest)
			return
		}

		body, _ := ioutil.ReadAll(r.Body)
		query := string(body)

		if len(query) >= 6 {
			switch query[0:6] {
			case "SELECT":
				if query != "SELECT 1;" {
					http.Error(w, "ERROR_SQL_SELECT", http.StatusInternalServerError)
					return
				}
			case "INSERT":
				sql_insert := "INSERT INTO zenit.mysql_slow_log (" +
					"_time,bytes_sent,host_ip,host_name,killed,last_errno,lock_time,query,query_digest,query_time," +
					"rows_affected,rows_examined,rows_read,rows_sent,schema,thread_id,user_host" +
					") VALUES (" +
					"'1529940303',60,IPv4StringToNum('127.0.0.1'),'localhost',0,0,0.000160," +
					"'SELECT count(*) FROM tabletest WHERE att = 'foo' AND deleted_at IS NULL;'," +
					"'SELECT count(*) FROM tabletest WHERE att = '?' AND deleted_at IS NULL;'," +
					"0.792864,0,27997,27997,1,'test',55883795,'test');"

				fmt.Printf("          > %s\n", query)

				if query != sql_insert {
					http.Error(w, "ERROR_SQL_INSERT", http.StatusInternalServerError)
					return
				}
			default:
				http.Error(w, "ERROR_SQL", http.StatusBadRequest)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
	}

	return http.HandlerFunc(fn)
}

func setup() {
	ts = httptest.NewServer(HandlerResponse())

	config.ClickHouse.DSN = fmt.Sprintf("%s/?database=zenit", ts.URL)
}

func shutdown() {
	ts.Close()
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func TestCheck(t *testing.T) {
	if !clickhouse.Check() {
		t.Errorf("Expected response status to be 200.")
	}
}

func TestSQLInsert(t *testing.T) {
	event.Values = append(event.Values, map[string]string{
		"_time":         "1529940303",
		"bytes_sent":    "60",
		"host_ip":       "127.0.0.1",
		"host_name":     "localhost",
		"killed":        "0",
		"last_errno":    "0",
		"lock_time":     "0.000160",
		"query":         "SELECT count(*) FROM tabletest WHERE att = 'foo' AND deleted_at IS NULL;",
		"query_digest":  "SELECT count(*) FROM tabletest WHERE att = '?' AND deleted_at IS NULL;",
		"query_time":    "0.792864",
		"rows_affected": "0",
		"rows_examined": "27997",
		"rows_read":     "27997",
		"rows_sent":     "1",
		"schema":        "test",
		"thread_id":     "55883795",
		"user_host":     "test",
	})

	sql_insert := sql.Insert(event.Schema, event.Table, event.Wildcard, event.Values)
	status_code := test.Post(config.ClickHouse.DSN, sql_insert)

	if status_code != 200 {
		t.Errorf("Expected response status to be 200 got %v.", status_code)
	}
}

func TestSendTimeout(t *testing.T) {
	timeout := time.After(time.Second * 3)
	tick := time.Tick(time.Second * 2)
	channel := make(chan map[string]string)

	defer close(channel)

	go clickhouse.Send(event, channel)

	for {
		select {
		case <-timeout:
			return
		case <-tick:
			channel <- map[string]string{
				"_time":         "1529940303",
				"bytes_sent":    "60",
				"host_ip":       "127.0.0.1",
				"host_name":     "localhost",
				"killed":        "0",
				"last_errno":    "0",
				"lock_time":     "0.000160",
				"query":         "SELECT count(*) FROM tabletest WHERE att = 'foo' AND deleted_at IS NULL;",
				"query_digest":  "SELECT count(*) FROM tabletest WHERE att = '?' AND deleted_at IS NULL;",
				"query_time":    "0.792864",
				"rows_affected": "0",
				"rows_examined": "27997",
				"rows_read":     "27997",
				"rows_sent":     "1",
				"schema":        "test",
				"thread_id":     "55883795",
				"user_host":     "test",
			}
		}
	}
}
