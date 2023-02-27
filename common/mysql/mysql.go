package mysql

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/debeando/zenit/common/log"

	_ "github.com/go-sql-driver/mysql"
)

type singleton struct {
	Connection *sql.DB
	Name       string
}

var instance = make(map[string]*singleton)

func GetInstance(name string) *singleton {
	if instance[name] == nil {
		instance[name] = &singleton{}
		instance[name].Name = name
	}
	return instance[name]
}

func (s *singleton) Connect(dsn string) {
	if s.Connection == nil {
		conn, err := sql.Open("mysql", dsn)
		if err != nil {
			log.Error("MySQL Client", map[string]interface{}{"name": s.Name, "error": err})
		}

		if err := conn.Ping(); err != nil {
			log.Error("MySQL Client", map[string]interface{}{"name": s.Name, "error": err})
		}

		s.Connection = conn
	}
}

func (s *singleton) Query(query string) map[int]map[string]string {
	if err := s.Connection.Ping(); err != nil {
		log.Error("MySQL Client", map[string]interface{}{"name": s.Name, "error": err})
		return nil
	}

	// Execute the query
	rows, err := s.Connection.Query(query)
	if err != nil {
		log.Error("MySQL Client", map[string]interface{}{"name": s.Name, "error": err})
	}
	defer rows.Close()

	// Get column names
	cols, _ := rows.Columns()
	if err != nil {
		log.Error("MySQL Client", map[string]interface{}{"name": s.Name, "error": err})
	}

	dataset := make(map[int]map[string]string)
	row_id := 0
	columns := make([]sql.RawBytes, len(cols))
	columnPointers := make([]interface{}, len(cols))

	for i := range cols {
		columnPointers[i] = &columns[i]
	}

	for rows.Next() {
		err = rows.Scan(columnPointers...)
		if err != nil {
			log.Error("MySQL Client", map[string]interface{}{"name": s.Name, "error": err})
		}

		row := make(map[string]string)

		for columnIndex, columnValue := range columns {
			row[cols[columnIndex]] = string(columnValue)
			dataset[row_id] = row
		}

		row_id++
	}

	return dataset
}

func (s *singleton) Close() {
	if s.Connection != nil {
		s.Connection.Close()
	}
}

func ParseValue(value string) (int64, bool) {
	value = strings.ToLower(value)

	if value == "yes" || value == "on" {
		return 1, true
	}

	if value == "no" || value == "off" {
		return 0, true
	}

	if val, err := strconv.ParseInt(value, 10, 64); err == nil {
		return val, true
	}

	return 0, false
}

func ClearUser(u string) string {
	index := strings.Index(u, "[")
	if index > 0 {
		return u[0:index]
	}
	return u
}

func MaximumValueSigned(dataType string) uint64 {
	switch dataType {
	case "tinyint":
		return 127
	case "smallint":
		return 32767
	case "mediumint":
		return 8388607
	case "int":
		return 2147483647
	case "bigint":
		return 9223372036854775807
	}
	return 0
}

func MaximumValueUnsigned(dataType string) uint64 {
	switch dataType {
	case "tinyint":
		return 255
	case "smallint":
		return 65535
	case "mediumint":
		return 16777215
	case "int":
		return 4294967295
	case "bigint":
		return 18446744073709551615
	}
	return 0
}
