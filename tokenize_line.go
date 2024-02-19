package zone

import (
	"bytes"
	"unicode"
)

// tokenizeLine parses a line into a set of record tokens. A record token
// is a sequence of non-whitespace characters, or a sequence of quoted
// characters. Trailing comments are ignored.
func tokenizeLine(line []byte) [][]byte {
	line = stripComment(line)
	result := make([][]byte, 0)

	inQuote := false
	token := make([]byte, 0)
	for i, b := range line {
		if unicode.IsSpace(rune(b)) == true && inQuote == false {
			if len(token) > 0 {
				result = append(result, token)
				token = make([]byte, 0)
			}
			continue
		}

		if b == bracketOpenByte || b == bracketCloseByte {
			if inQuote == false {
				continue
			}
		}

		if b == quoteByte {
			switch inQuote {
			case false:
				inQuote = true
			case true:
				if line[i-1] != escapeByte {
					inQuote = false
				}
			}
		}

		token = append(token, b)
	}

	if len(token) > 0 {
		// If the line did not end with spaces, then the last token would not have
		// been appended to the result. Also, if the line was a single token then
		// the token would not have been appended, i.e. `len(result) == 0`.
		if len(result) == 0 || bytes.Equal(result[len(result)-1], token) == false {
			result = append(result, token)
		}
	}

	return result
}
