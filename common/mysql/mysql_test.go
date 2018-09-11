package mysql_test

import (
	"testing"

	"github.com/swapbyt3s/zenit/common/mysql"
)

func TestParseValue(t *testing.T) {
	if value, ok := mysql.ParseValue([]byte("YES")); !ok || value != 1 {
		t.Error("Expected: Found Parse and value = 1.")
	}

	if value, ok := mysql.ParseValue([]byte("NO")); !ok || value != 0 {
		t.Error("Expected: Found Parse and value = 0.")
	}

	if value, ok := mysql.ParseValue([]byte("ON")); !ok || value != 1 {
		t.Error("Expected: Found Parse and value = 1.")
	}

	if value, ok := mysql.ParseValue([]byte("OFF")); !ok || value != 0 {
		t.Error("Expected: Found Parse and value = 0.")
	}

	if value, ok := mysql.ParseValue([]byte("true")); ok && value == 0 {
		t.Error("Expected: Imposible Parse.")
	}

	if value, ok := mysql.ParseValue([]byte("1234567890")); !ok || value != 1234567890 {
		t.Error("Expected: Found Parse and value = 1234567890.")
	}
}

func TestClearUser(t *testing.T) {
	user := "test[test] @ [127.0.0.1]"
	expected := "test"
	result := mysql.ClearUser(user)

	if result != expected {
		t.Error("Expected: " + expected)
	}
}
