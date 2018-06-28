package common

import "regexp"

type ExtRegexp struct {
  *regexp.Regexp
}

func (r *ExtRegexp) FindStringSubmatchMap(s string) (map[string]string) {
  match := r.FindStringSubmatch(s)
  if match == nil {
    return nil
  }

  captures := make(map[string]string)
  for i, name := range r.SubexpNames() {
    if i == 0 {
      continue
    }
    if name != "" {
      captures[name] = match[i]
    }
  }

  return captures
}
