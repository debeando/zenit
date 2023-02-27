package sql

import (
	"fmt"
	"strings"

	"github.com/debeando/zenit/common"
)

func Insert(schema string, table string, wildcard map[string]string, values []map[string]string) string {
	sql := ""
	fields := common.KeyOfMaps(values)
	rows := []string{}

	for v := len(values) - 1; v >= 0; v-- {
		t := []string{}
		for f := 0; f < len(fields); f++ {
			t = append(t, fmt.Sprintf(wildcard[fields[f]], values[v][fields[f]]))
		}
		rows = append(rows, fmt.Sprintf("(%s)", strings.Join(t, ",")))
	}

	if len(fields) > 0 && len(rows) > 0 {
		sql = fmt.Sprintf("INSERT INTO %s.%s (%s) VALUES %s;", schema, table, strings.Join(fields, ","), strings.Join(rows, ","))
	}

	return sql
}
