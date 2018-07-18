// This a test version for normalize SQL without regex because is very slowly to parse.
package main

import (
  "fmt"
  // "strings"
)

var newline    bool
var whitespace bool
var quote      bool
var comment    bool
var multiline  bool

var rv []rune
var sql = `
SET autocommit = 0;
/* Comment type 1 */
SET @foo = 1;
/*
 * Comment type 2
 */
# Comment type 3
-- Comment type 4
SELECT COUNT(*) # AS count
     , id AS 'Number'
     , IF(foo = "3", 1, 2) AS "test"
FROM foo    -- , bar
WHERE id           = 123
  AND email        = 'abc@def.aaa'
  AND created_at  >= "2015-06-19"
  AND modified_at <> "2015-06-19 00:00:00"
  AND amount       > 0.10
LIMIT 1,10;
`

func main() {
  sql_runes := []rune(sql)

  for i := 0; i < len(sql_runes); i++ {
    // fmt.Printf("%d: %c\n", i, sql_runes[i])

    // Remove comments:
    if ! comment && ! multiline && sql_runes[i] == '#' {
      comment = true
    } else if comment && ! multiline && sql_runes[i] == '\n' {
      comment = false
      continue
    }

    if ! comment && sql_runes[i] == '/' && sql_runes[i + 1] == '*' {
      comment = true
      multiline = true
      // continue
    } else if comment && multiline && sql_runes[i - 1] == '*' && sql_runes[i] == '/' {
      comment = false
      multiline = false
      continue
    }

    if comment {
      continue
    }

    // Remove new lines:
    if sql_runes[i] == '\n' || sql_runes[i] == '\r' {
      whitespace = true
      continue
    }

    // Remove whitespaces:
    if sql_runes[i] == ' ' {
      whitespace = true
      continue
    } else {
      if whitespace {
        whitespace = false
        rv = append(rv, ' ')
      }
    }

    // Replace duble quotes to single quotes
    if sql_runes[i] == '"' {
      quote = true
      rv = append(rv, '\'')
      continue
    } else {
      if quote {
        quote = false
      }
    }

    rv = append(rv, sql_runes[i])
  }

  fmt.Printf(">%s.\n", string(rv))
}
