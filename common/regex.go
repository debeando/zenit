package common

import (
  // "fmt"
  "regexp"
)

func RegexpGetGroups(re *regexp.Regexp, s string) map[string]string {
  result := make(map[string]string)
  kv := re.FindStringSubmatch(s)

  // fmt.Printf("--> %#v\n", kv)

  if len(kv) > 0 {
    for i, name := range re.SubexpNames() {
      if i != 0 && name != "" {
        result[name] = kv[i]
      }
    }
  }

  return result
}
