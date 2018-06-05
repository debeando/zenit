// Package generic lib.
package lib

import(
  "crypto/md5"
  "fmt"
)

func MD5(s string) string {
  return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}
