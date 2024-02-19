package zone

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_tokenizeLine(t *testing.T) {
	line := []byte("        IN      NS      dns1.example.com.")
	expected := [][]byte{[]byte("IN"), []byte("NS"), []byte("dns1.example.com.")}
	found := tokenizeLine(line)
	assert.Equal(t, expected, found)

	line = []byte("        IN      MX      10      mail.example.com.")
	expected = [][]byte{[]byte("IN"), []byte("MX"), []byte("10"), []byte("mail.example.com.")}
	found = tokenizeLine(line)
	assert.Equal(t, expected, found)

	line = []byte("mail    IN      CNAME   server1")
	expected = [][]byte{[]byte("mail"), []byte("IN"), []byte("CNAME"), []byte("server1")}
	found = tokenizeLine(line)
	assert.Equal(t, expected, found)

	// tokenizes a line that ends with space characters
	line = []byte("    foo\tbar\v\r\n")
	expected = [][]byte{[]byte("foo"), []byte("bar")}
	found = tokenizeLine(line)
	assert.Equal(t, expected, found)

	// honor quoted spaces
	line = []byte("\"foo bar\"")
	expected = [][]byte{[]byte("\"foo bar\"")}
	found = tokenizeLine(line)
	assert.Equal(t, expected, found)

	// handle escaped quotes within a quote block ("foo \"bar baz\"")
	line = []byte("\"foo \\\"bar baz\\\"\"")
	expected = [][]byte{[]byte("\"foo \\\"bar baz\\\"\"")}
	found = tokenizeLine(line)
	assert.Equal(t, expected, found)

	// a rfc4408§3.1.3 style record
	line = []byte("@ txt \"foo \" \"bar \" \"baz\"")
	expected = [][]byte{[]byte("@"), []byte("txt"), []byte("\"foo \""), []byte("\"bar \""), []byte("\"baz\"")}
	found = tokenizeLine(line)
	assert.Equal(t, expected, found)

	// ignores line comment
	line = []byte("in ns dns1 ; a comment")
	expected = [][]byte{[]byte("in"), []byte("ns"), []byte("dns1")}
	found = tokenizeLine(line)
	assert.Equal(t, expected, found)

	// skips brackets
	line = []byte("foo ( bar baz )")
	expected = [][]byte{[]byte("foo"), []byte("bar"), []byte("baz")}
	found = tokenizeLine(line)
	assert.Equal(t, expected, found)

	// does not skip quoted brackets
	line = []byte("foo \"( bar baz )\"")
	expected = [][]byte{[]byte("foo"), []byte("\"( bar baz )\"")}
	found = tokenizeLine(line)
	assert.Equal(t, expected, found)
}
