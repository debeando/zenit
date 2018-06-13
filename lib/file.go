package lib

import (
  "io/ioutil"
  "os"
  "strconv"
)

func GetValueFromFile(path string) uint64 {
  if _, err := os.Stat(path); err == nil {
    content, err := ioutil.ReadFile(path)
    if err != nil {
      panic(err)
    }

    value, err := strconv.ParseUint(string(content), 10, 64)
    if err != nil {
      panic(err)
    }
    return value
  }

  return 0
}
