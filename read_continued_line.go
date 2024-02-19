package zone

import (
	"bufio"
	"io"
)

// readContinuedLine reads from the reader until the end of a continued line
// and returns a single line of bytes. A continued line is one in which an
// opening parentheses is found with no closing parentheses on the same line.
// For example, if the data stream contains `foo (\n bar\n baz)` then the
// currentLine would be `foo (` and the result of this function will be
// `foo ( bar baz)`.
func readContinuedLine(reader io.Reader, currentLine []byte) ([]byte, error) {
	r := bufio.NewReader(reader)
	currentLine = compactWhiteSpace(stripComment(currentLine))
	for {
		line, err := r.ReadBytes('\n')
		if err != nil {
			return nil, err
		}

		endIdx := lastIndexNonEscapedByte(line, byte(')'))
		line = compactWhiteSpace(stripComment(line))
		currentLine = append(currentLine, line...)
		if endIdx > 0 {
			break
		}
	}
	return currentLine, nil
}
