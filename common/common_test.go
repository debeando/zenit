package common_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/swapbyt3s/zenit/common"
)

var wd string

func TestMain(m *testing.M) {
	wd, _ = os.Getwd()
}

func TestPGrep(t *testing.T) {

}

func TestGetUInt64FromFile(t *testing.T) {
	expected := uint64(1234567890)
	result := common.GetUInt64FromFile(wd + "/../assets/tests/uint64.txt")

	if result != expected {
		t.Error("Expected: uint64(1234567890)")
	}

	expected = uint64(0)
	result = common.GetUInt64FromFile(wd + "/../assets/tests/uint64.log")

	if result != expected {
		t.Error("Expected: uint64(0)")
	}
}

func TestStringToUInt64(t *testing.T) {
	expected := uint64(1234)
	result := common.StringToUInt64("1234")

	if result != expected {
		t.Error("Expected: uint64(1234)")
	}

	result = common.StringToUInt64("abc")

	if result != 0 {
		t.Error("Expected: 0")
	}

	result = common.StringToUInt64("")

	if result != 0 {
		t.Error("Expected: 0")
	}
}

func TestGetIntFromFile(t *testing.T) {
	expected := 1234567890
	result := common.GetIntFromFile(wd + "/../assets/tests/int.txt")

	if result != expected {
		t.Error("Expected: 1234567890")
	}
}

func TestStringToInt(t *testing.T) {
	expected := 1234567890
	result := common.StringToInt("1234567890")

	if result != expected {
		t.Error("Expected: 1234567890")
	}
}

func TestMD5(t *testing.T) {
	expected := "098f6bcd4621d373cade4e832627b4f6"
	result := common.MD5("test")

	if result != expected {
		t.Error("Expected: 098f6bcd4621d373cade4e832627b4f6")
	}
}

func TestStringInArray(t *testing.T) {
	list := []string{"foo", "bar"}
	result := common.StringInArray("bar", list)

	if !result {
		t.Error("Expected: false")
	}

	result = common.StringInArray("test", list)

	if result {
		t.Error("Expected: true")
	}

	result = common.StringInArray("", list)

	if result {
		t.Error("Expected: false")
	}
}

func TestKeyInMap(t *testing.T) {
	expected := make(map[string]string)
	expected["test"] = "test"

	if !common.KeyInMap("test", expected) {
		t.Error("Expected: true")
	}

	if common.KeyInMap("foo", expected) {
		t.Error("Expected: false")
	}
}

func TestKeyOfMaps(t *testing.T) {
	result := common.KeyOfMaps([]map[string]string{{"foo": "a", "bar": "1"}, {"foo": "b", "bar": "2"}})
	expected := []string{"foo", "bar"}

	if !reflect.DeepEqual(result, expected) {
		t.Error("Expected: []string{\"foo\", \"bar\"}")
	}
}

func TestHostname(t *testing.T) {

}

func TestIpAddress(t *testing.T) {

}

func TestToDateTime(t *testing.T) {
	expected := "2018-12-31 15:04:05"
	result := common.ToDateTime("2018-12-31T15:04:05 UTC", "2006-01-02T15:04:05 UTC")

	if result != expected {
		t.Error("Expected: 2018-12-31 15:04:05")
	}
}

func TestExecCommand(t *testing.T) {

}

func TestEscape(t *testing.T) {
	expected := "<abc=\\'abc\\'>foo</abc>"
	result := common.Escape("<abc='abc'>foo</abc>")

	if result != expected {
		t.Error("Expected: <abc=\\'abc\\'>foo</abc>")
	}
}

func TestComparteMapString(t *testing.T) {
}
