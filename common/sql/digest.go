package sql

import (
	"strings"
	"unicode"
)

func Digest(s string) string {
	comment    := false
	endnumber  := []rune{' ', ',', '+', '-', '*', '/', '^', '%', '(', ')'}
	list       := false
	listValues := false
	multiline  := false
	number     := false
	quote      := rune(0)
	result     := []rune("")
	sql        := []rune(strings.ToLower(s))
	whitespace := false
	length     := len(sql)

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

		// Remove literals inside of list "IN":
		if sql[x] == 'i' && x+1 < length && sql[x+1] == 'n' {
			// fmt.Printf("==> %t, %d\n", list, (length - x))

			for y := 0; y < (length - x); y++ {
				if sql[x+y] == '(' {
					// fmt.Printf("-- --> S %s, %d\n", string(sql[x+y]), x+y)
					list = true
					listValues = false
					break
				}
			}
		}

		if list && ! listValues && sql[x] == '(' {
			listValues = true
		} else if list && listValues && sql[x] == ')' {
			// fmt.Printf("-- --> E %s, %d\n", string(sql[x]), x)
			list = false
			listValues = false
			result = append(result, '?')
		} else if list && listValues {
			continue
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
