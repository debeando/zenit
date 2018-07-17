package common

import (
  "regexp"
  "strings"
)

var (
  reRemoveCommentCase1 = regexp.MustCompile(`\/\*(.|[\r\n])*\*\/`)
  reRemoveCommentCase2 = regexp.MustCompile(`\#.*\n`)
  reRemoveCommentCase3 = regexp.MustCompile(`--.*\n`)
  reClearWhiteSpaces   = regexp.MustCompile(`\s+`)
  reClearString        = regexp.MustCompile(`=\s?'([^']*)'`)
  reClearNumber        = regexp.MustCompile(`=\s?([\d.]+)`)
)

func QueryNormalizer(s string) string {
  s = RemoveComments(s)
  s = ClearWhiteSpaces(s)
  s = ClearString(s)
  s = ClearNumber(s)
  return s
}

func RemoveComments(s string) string {
  s = reRemoveCommentCase1.ReplaceAllString(s, "")
  s = reRemoveCommentCase2.ReplaceAllString(s, "")
  s = reRemoveCommentCase3.ReplaceAllString(s, "")
  return s
}

func ClearWhiteSpaces(s string) string {
  s = reClearWhiteSpaces.ReplaceAllString(s, " ")
  return strings.Trim(s, " ")
}

func ClearString(s string) string {
  s = strings.Replace(s, "\"", "'", -1)
  s = strings.Replace(s, "`", "", -1)
  s = reClearString.ReplaceAllString(s, "= '?'")

  return s
}

func ClearNumber(s string) string {
  s = reClearNumber.ReplaceAllString(s, "= ?")

  return s
}
