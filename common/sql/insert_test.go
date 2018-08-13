package sql_test

import (
  "testing"

  "github.com/swapbyt3s/zenit/common/sql"
)

func TestInsert(t *testing.T) {
  wildcard := map[string]string{"foo": "'%s'", "bar": "%s"}
  values   := []map[string]string{{"foo": "a", "bar": "1"},{"foo": "b", "bar": "2"}}
  result   := sql.Insert("zenit", "test", wildcard, values)
  expected := "INSERT INTO zenit.test (foo,bar) VALUES ('a',1),('b',2);"

  if result != expected {
    t.Error("Expected: " + expected)
  }
}
