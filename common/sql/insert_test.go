package sql_test

import (
	"testing"

	"github.com/debeando/zenit/common/sql"
)

func TestEmpty(t *testing.T) {
	wildcard := map[string]string{"FieldA": "'%s'", "FieldN": "%s"}
	values := []map[string]string{{}, {}}
	result := sql.Insert("zenit", "test", wildcard, values)

	if len(result) > 0 {
		t.Error("Expected empty result.")
	}
}

func TestInsert(t *testing.T) {
	wildcard := map[string]string{"FieldA": "'%s'", "FieldN": "%s"}
	values := []map[string]string{{"FieldA": "c", "FieldN": "3"},
		{"FieldA": "b", "FieldN": "2"},
		{"FieldA": "a", "FieldN": "1"}}
	result := sql.Insert("zenit", "test", wildcard, values)
	expected := "INSERT INTO zenit.test (FieldA,FieldN) VALUES ('a',1),('b',2),('c',3);"

	if result != expected {
		t.Errorf("Expected: '%#v', got: '%#v'.", expected, result)
	}
}
