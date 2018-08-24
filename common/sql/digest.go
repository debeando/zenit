// NOTES:
// - This a test version for normalize SQL without regex, because is very slowly
// to parse.
// - Digest converts an SQL statement into a normalized version (with all
// string/numeric literals replaced with ?).
// - It most definitely does not validate that a query is syntactically correct.

package sql

import (
	"unicode"
)

func Digest(s string) string {
	whitespace := false
	quote := rune(0)
	comment := false
	multiline := false
	number := false
	result := []rune("")
	sql := []rune(s)
	length := len(sql)
	endnumber := []rune{' ', ',', '+', '-', '*', '/', '^', '%'}

	IsNumber := func(r rune) bool {
		if unicode.IsNumber(r) || r == '.' {
			return true
		}
		return false
	}

	IsEndNumber := func(r rune) bool {
		for z := 0; z < len(endnumber); z++ {
			if r == endnumber[z] {
				return true
			}
		}
		return false
	}

	for x := 0; x < length; x++ {
		// Remove comments:
		if !comment && !multiline && sql[x] == '#' {
			comment = true
		} else if comment && !multiline && sql[x] == '\n' {
			comment = false
			continue
		}

		if !comment && !multiline && sql[x] == '-' && x+1 < length && sql[x+1] == '-' {
			comment = true
		} else if comment && !multiline && sql[x] == '\n' {
			comment = false
			continue
		}

		if !comment && sql[x] == '/' && x+1 < length && sql[x+1] == '*' {
			comment = true
			multiline = true
			// continue
		} else if comment && multiline && sql[x] == '*' && x+1 < length && sql[x+1] == '/' {
			x += 1
			comment = false
			multiline = false
			continue
		}

		if comment {
			continue
		}

		// Remove new lines:
		if sql[x] == '\n' || sql[x] == '\r' {
			sql[x] = ' '
			whitespace = true
			number = false
			// continue
		}

		// Remove whitespaces:
		if quote == 0 && sql[x] == ' ' {
			whitespace = true
			number = false
			// continue
		} else if quote == 0 {
			if whitespace {
				whitespace = false
				result = append(result, ' ')
			}
		}

		if whitespace {
			continue
		}

		// Remove backtick
		if quote == 0 && sql[x] == '`' {
			continue
		}

		// Remove string between quotes:
		if quote == 0 && (sql[x] == '"' || sql[x] == '\'') {
			quote = sql[x]
			result = append(result, '\'')
		} else if quote > 0 && sql[x] == '\\' && x+1 < length && sql[x+1] == quote {
			x += 1
		} else if sql[x] == quote {
			quote = 0
			result = append(result, '?')
			result = append(result, '\'')
			continue
		}

		if quote > 0 {
			continue
		}

		// Remove numbers:
		if !number && IsNumber(sql[x]) {
			number = true

			// Check to skip word composed with number and letter:
			for y := x; y >= 0; y-- {
				if IsEndNumber(sql[y]) {
					break
				} else {
					if unicode.IsLetter(sql[y]) {
						number = false
					}
				}
			}

			for y := 0; y < (length - x); y++ {
				if IsEndNumber(sql[x+y]) {
					break
				} else {
					if unicode.IsLetter(sql[x+y]) {
						number = false
					}
				}
			}

			// Add ? symbol to remove nombre:
			if number {
				result = append(result, '?')
			}
		}

		if number && !IsNumber(sql[x]) {
			number = false
		}

		if number {
			continue
		}

		// Add character:
		result = append(result, sql[x])
	}

	return string(result)
}
