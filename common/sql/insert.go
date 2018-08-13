package sql

import (
  "fmt"
  "strings"

  "github.com/swapbyt3s/zenit/common"
)

func Insert(schema string, table string, wildcard map[string]string, values []map[string]string) string {
  fields := common.KeyOfMaps(values)
  rows   := []string{}

  for _, value := range values {
    v := []string{}
    for _, field := range fields {
      v = append(v, fmt.Sprintf(wildcard[field], value[field]))
    }
    rows = append(rows, fmt.Sprintf("(%s)", strings.Join(v, ",")))
  }

  return fmt.Sprintf("INSERT INTO %s.%s (%s) VALUES %s;", schema, table, strings.Join(fields, ","), strings.Join(rows, ","))
}
