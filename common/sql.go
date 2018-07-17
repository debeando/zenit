package common

import (
  "regexp"
  "strings"
)

func QueryNormalizer(s string) string {
  s = RemoveComments(s)
  s = ClearWhiteSpaces(s)
  s = ClearString(s)
  s = ClearNumber(s)
  return s
}

func RemoveComments(s string) string {
  re := regexp.MustCompile(`\/\*(.|[\r\n])*\*\/`)
  s   = re.ReplaceAllString(s, "")
  re  = regexp.MustCompile(`\#.*\n`)
  s   = re.ReplaceAllString(s, "")
  re  = regexp.MustCompile(`--.*\n`)
  s   = re.ReplaceAllString(s, "")
  return s
}

func ClearWhiteSpaces(s string) string {
  re := regexp.MustCompile(`\s+`)
  s   = re.ReplaceAllString(s, " ")
  return strings.Trim(s, " ")
}

func ClearString(s string) string {
  re := regexp.MustCompile(`=\s?'([^']*)'`)
  s = strings.Replace(s, "\"", "'", -1)
  s = strings.Replace(s, "`", "", -1)
  s = re.ReplaceAllString(s, "= '?'")

  return s
}

func ClearNumber(s string) string {
  re := regexp.MustCompile(`=\s?([\d.]+)`)
  s = re.ReplaceAllString(s, "= ?")

  return s
}
