// This a test version for normalize SQL without regex because is very slowly to parse.
package main

import (
  "fmt"
)

var sql = `
-- Disable AC.
SET autocommit = 0;
/* Comment type 1 */
SET @foo = 1;
/*
 * Comment type 2
 */
# Comment type 3
-- Comment type 4
SELECT /* id */ user_id
     , COUNT(*) # AS count
     , id AS 'Number'
     , IF(foo = "3", 1, 2) AS "test"
FROM foo    -- , bar
WHERE id           = 123
  AND email        = 'abc@def.aaa'
  AND text         = "<foo='test'/>Don't bar</foo>"
  AND text         = '<foo="test"/>Don\'t bar</foo>'
  AND text         = '<foo=\'test\'/>Don\'t bar</foo>'
  AND created_at  >= "2015-06-19"
  AND modified_at <> "2015-06-19 00:00:00"
  AND amount       > 0.10
LIMIT 1,10;
`

func main() {
  sql = `AND text = "<foo='test'/>Don't bar</foo>"`
  fmt.Printf(">%s.\n", QueryNormalizer(sql))

  sql = `AND text = '<foo="test"/>Don\'t bar</foo>'`
  fmt.Printf(">%s.\n", QueryNormalizer(sql))

  sql = `AND modified_at <> "2015-06-19 00:00:00"`
  fmt.Printf(">%s.\n", QueryNormalizer(sql))

  sql = `AND modified_at <> "2015-06-19"`
  fmt.Printf(">%s.\n", QueryNormalizer(sql))

  sql = `IF(foo = "3", 1, 2) AS "test"`
  fmt.Printf(">%s.\n", QueryNormalizer(sql))
}

func QueryNormalizer(s string) string {
  whitespace := false
  quote      := rune(0)
  comment    := false
  multiline  := false
  result     := []rune("")
  sql        := []rune(s)

  for i := 0; i < len(sql); i++ {
    // Remove comments:
    if ! comment && ! multiline && sql[i] == '#' {
      comment = true
    } else if comment && ! multiline && sql[i] == '\n' {
      comment = false
      continue
    }

    if ! comment && ! multiline && sql[i] == '-' && sql[i + 1] == '-' {
      comment = true
    } else if comment && ! multiline && sql[i] == '\n' {
      comment = false
      continue
    }

    if ! comment && sql[i] == '/' && sql[i + 1] == '*' {
      comment = true
      multiline = true
      // continue
    } else if comment && multiline && sql[i - 1] == '*' && sql[i] == '/' {
      comment = false
      multiline = false
      continue
    }

    if comment {
      continue
    }

    // Remove new lines:
    if sql[i] == '\n' || sql[i] == '\r' {
      sql[i] = ' '
      whitespace = true
      continue
    }

    // Remove whitespaces:
    if quote == 0 && sql[i] == ' ' {
      whitespace = true
      continue
    } else if quote == 0 {
      if whitespace {
        whitespace = false
        result = append(result, ' ')
      }
    }

    // Remove string between quotes:
    if quote == 0 && ( sql[i] == '"' || sql[i] == '\'' ) {
      quote = sql[i]
      result = append(result, '\'')
    } else if quote > 0 && sql[i] == '\\' && sql[i + 1] == quote {
      i += 1
    } else if sql[i] == quote {
      quote = 0
      result = append(result, '?')
      result = append(result, '\'')
      continue
    }

    if quote > 0 {
      continue
    }

    result = append(result, sql[i])
  }

  return string(result)
}

// OK SELECT "<foo='test'/>Don't bar</foo>"
// OK SELECT "<foo=\"test\"/>Don\'t bar</foo>"
// KO SELECT "<foo=\"test\"/>Don"t bar</foo>";
// OK SELECT '<foo="test"/>Don\'t bar</foo>'
// OK SELECT '<foo=\'test\'/>Don\'t bar</foo>'
// KO SELECT '<foo=\'test\'/>Don't bar</foo>'
// 
// guarda el quote " para compararlo con el final.
// verifica que no sea un scape quote \' or \"
// ignora todo el recorrido hasta que la posicion x sea igual al quote ".
// añade '?'
// 
// '<foo=\'test\'/>Don\'t bar</foo>'
