package zone

import (
	"bytes"
	"unicode"
)

// compactWhiteSpace reduces any leading and trailing whitespace characters
// in a line to a single empty space character (0x20).
func compactWhiteSpace(line []byte) []byte {
	startIdx := bytes.IndexFunc(line, isNotWhiteSpace)
	endIdx := bytes.LastIndexFunc(line, isNotWhiteSpace)
	result := []byte{0x20}
	result = append(result, line[startIdx:endIdx+1]...)
	return append(result, 0x20)
}

func isContinuedLine(line []byte) bool {
	openByte := byte('(')
	closeByte := byte(')')
	openIdx := indexNonEscapedByte(line, openByte)
	endIdx := lastIndexNonEscapedByte(line, closeByte)

	result := false
	if openIdx > 0 {
		result = true
	}
	if endIdx > 0 {
		commentIdx := indexNonEscapedByte(line, commentStartBytes[0])
		if commentIdx > 0 && commentIdx < endIdx {
			result = true
		} else {
			result = false
		}
	}

	return result
}

// indexNonEscapedByte finds the first position of the given needle in the
// haystack if that needle is not preceded by the escape character
// U+005C (reverse solidus).
func indexNonEscapedByte(haystack []byte, needle byte) int {
	idx := bytes.IndexByte(haystack, needle)
	if idx < 1 {
		return idx
	}

	if haystack[idx-1] == escapeByte {
		return -1
	}

	return idx
}

// isNotWhiteSpace determines if rune r constitutes a "visible" character.
// See [unicode.IsSpace] for the list of runes considered to be invisible
// characters.
func isNotWhiteSpace(r rune) bool {
	if unicode.IsSpace(r) == true {
		return false
	}
	return true
}

// isTtl determines if the provided string is a time-to-live number or not.
// Basically, it's a simple check for the string being all digits.
func isTtl(input string) bool {
	return isTtlToken.Match([]byte(input))
}

// lastIndexNonEscapedByte finds the last position of the given needle in the
// haystack if that needle is not preceded by the escape character
// U+005C (reverse solidus).
func lastIndexNonEscapedByte(haystack []byte, needle byte) int {
	idx := bytes.LastIndexByte(haystack, needle)
	if idx < 1 {
		return idx
	}

	if haystack[idx-1] == escapeByte {
		return -1
	}

	return idx
}

// stripComment removes any text following a semicolon, including the
// semicolon, in a slice of bytes. Does not strip text inside of quote
// blocks.
//
//	result := stripComment([]byte("foo \" bar ; baz \" ; remove me"))
//	fmt.Println(string(result)) // `foo " bar ; baz "`
func stripComment(line []byte) []byte {
	inQuote := false
	idx := -1
	for i, j := range line {
		if j == quoteByte && i == 0 {
			inQuote = true
			continue
		}
		if j == quoteByte && line[i-1] != escapeByte && inQuote == false {
			inQuote = true
			continue
		}
		if j == quoteByte && line[i-1] != escapeByte && inQuote == true {
			inQuote = false
			continue
		}
		if j == commentStartByte && inQuote == false {
			idx = i
			break
		}
	}

	if idx == -1 {
		return line
	}
	return line[0:idx]
}
