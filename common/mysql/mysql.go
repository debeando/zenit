package mysql

import (
	"bytes"
	"database/sql"
	"log"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func Connect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, err
}

func Check(dsn string, name string) bool {
	log.Printf("I! - %s - DSN: %s\n", name, dsn)
	conn, err := Connect(dsn)
	if err != nil {
		log.Printf("E! - %s - Impossible to connect: %s\n", name, err)
		return false
	}

	log.Printf("I! - %s - Connected successfully.\n", name)
	conn.Close()
	return true
}

func ParseValue(value sql.RawBytes) (uint64, bool) {
	if bytes.EqualFold(value, []byte("YES")) || bytes.Compare(value, []byte("ON")) == 0 {
		return 1, true
	}

	if bytes.EqualFold(value, []byte("NO")) || bytes.Compare(value, []byte("OFF")) == 0 {
		return 0, true
	}

	if val, err := strconv.ParseUint(string(value), 10, 64); err == nil {
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

func YesOrNo(v int) string {
	if v == 1 {
		return "Yes"
	}
	return "No"
}
