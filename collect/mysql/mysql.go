package mysql

import (
  "strings"
)

func ClearUser(u string) string {
  index := strings.Index(u, "[")
  if index > 0 {
    return u[0:index]
  }
  return u
}
