package mysql_test

import (
	"testing"

	"github.com/swapbyt3s/zenit/common/mysql"
)

func TestParseValue(t *testing.T) {
	if value, ok := mysql.ParseValue("yes"); !ok || value != 1 {
		t.Error("Expected: Found Parse and value = 1.")
	}

	if value, ok := mysql.ParseValue("Yes"); !ok || value != 1 {
		t.Error("Expected: Found Parse and value = 1.")
	}

	if value, ok := mysql.ParseValue("YES"); !ok || value != 1 {
		t.Error("Expected: Found Parse and value = 1.")
	}

	if value, ok := mysql.ParseValue("no"); !ok || value != 0 {
		t.Error("Expected: Found Parse and value = 0.")
	}

	if value, ok := mysql.ParseValue("No"); !ok || value != 0 {
		t.Error("Expected: Found Parse and value = 0.")
	}

	if value, ok := mysql.ParseValue("NO"); !ok || value != 0 {
		t.Error("Expected: Found Parse and value = 0.")
	}

	if value, ok := mysql.ParseValue("ON"); !ok || value != 1 {
		t.Error("Expected: Found Parse and value = 1.")
	}

	if value, ok := mysql.ParseValue("OFF"); !ok || value != 0 {
		t.Error("Expected: Found Parse and value = 0.")
	}

	if value, ok := mysql.ParseValue("true"); ok && value == 0 {
		t.Error("Expected: Imposible Parse.")
	}

	if value, ok := mysql.ParseValue("1234567890"); !ok || value != 1234567890 {
		t.Error("Expected: Found Parse and value = 1234567890.")
	}
}

func TestClearUser(t *testing.T) {
	user := "test[test] @ [127.0.0.1]"
	expected := "test"
	result := mysql.ClearUser(user)

	if result != expected {
		t.Errorf("Expected: '%s', got: '%s'.", expected, result)
	}
}
