package common_test

import (
  "testing"
  "gitlab.com/swapbyt3s/zenit/common"
)

func TestStringToUInt64(t *testing.T) {
  expected := uint64(1234)
  result   := common.StringToUInt64("1234")

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

func TestMD5(t *testing.T) {
  expected := "098f6bcd4621d373cade4e832627b4f6"
  result   := common.MD5("test")

  if result != expected {
    t.Error("Expected: 098f6bcd4621d373cade4e832627b4f6")
  }
}

func TestStringInArray(t *testing.T) {
  list   := []string{"foo", "bar"}
  result := common.StringInArray("bar", list)

  if ! result {
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

//ParseDateTime
//GetEnv
//Hostname
//IpAddress: usar regex
