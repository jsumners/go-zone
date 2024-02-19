package zone

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_compactWhiteSpace(t *testing.T) {
	expected := []byte(" foo ")
	tests := []string{
		"\t \tfoo    ",
		" foo\n",
		"\vfoo\r",
	}

	for _, test := range tests {
		found := compactWhiteSpace([]byte(test))
		assert.Equal(t, expected, found)
	}
}

func Test_isContinuedLine(t *testing.T) {
	tests := [][]any{
		{"foo ( bar baz )", false},
		{"foo \\(bar", false},
		{"foo \\(bar\\)", false},
		{"foo ( ; comment bar)", true},
		{"foo ( bar", true},
	}

	for _, test := range tests {
		found := isContinuedLine([]byte(test[0].(string)))
		t.Log(test[0].(string))
		assert.Equal(t, test[1].(bool), found)
	}
}

func Test_indexNonEscapedByte(t *testing.T) {
	expected := -1
	found := indexNonEscapedByte([]byte("foobar"), byte('('))
	assert.Equal(t, expected, found)

	expected = -1
	found = indexNonEscapedByte([]byte("foo \\("), byte('('))
	assert.Equal(t, expected, found)

	expected = 4
	found = indexNonEscapedByte([]byte("foo ("), byte('('))
	assert.Equal(t, expected, found)
}

func Test_isNotWhiteSpace(t *testing.T) {
	found := isNotWhiteSpace(' ')
	assert.Equal(t, false, found)

	found = isNotWhiteSpace('a')
	assert.Equal(t, true, found)
}

func Test_lastIndexNonEscapedByte(t *testing.T) {
	expected := -1
	found := lastIndexNonEscapedByte([]byte("foobar"), byte(')'))
	assert.Equal(t, expected, found)

	expected = -1
	found = lastIndexNonEscapedByte([]byte("foobar \\)"), byte(')'))
	assert.Equal(t, expected, found)

	expected = 4
	found = lastIndexNonEscapedByte([]byte("foo )"), byte(')'))
	assert.Equal(t, expected, found)
}

func Test_stripComment(t *testing.T) {
	tests := [][]string{
		{"foo; bar", "foo"},
		{"foo ; bar", "foo "},
		{"\"foo ; keep me\" ; remove me", "\"foo ; keep me\" "},
	}

	for _, test := range tests {
		found := stripComment([]byte(test[0]))
		assert.Equal(t, []byte(test[1]), found)
	}
}
