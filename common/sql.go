// NOTES:
// - This a test version for normalize SQL without regex, because is very slowly
// to parse.
// - NormalizeQuery converts an SQL statement into a normalized version (with all
// string/numeric literals replaced with ?).
// - It most definitely does not validate that a query is syntactically correct.

package common

func NormalizeQuery(s string) string {
  whitespace := false
  quote      := rune(0)
  comment    := false
  multiline  := false
  number     := false
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
      number     = false
      continue
    }

    // Remove whitespaces:
    if quote == 0 && sql[i] == ' ' {
      whitespace = true
      number     = false
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

    // Remove numbers:
    if number {
      if (sql[i] < '0' || sql[i] > '9') && sql[i] != '.' {
        number = false
        result = append(result, sql[i])
      }
      continue
    }

    if ((sql[i] >= '0' && sql[i] <= '9') || sql[i] == '.') {
      number = true
      result = append(result, '?')
      continue
    }

    // Add character:
    result = append(result, sql[i])
  }

  return string(result)
}
