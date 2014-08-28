package regexkata

import (
	"bytes"
	"regexp"
	"strings"
)

const PlainFieldToken = 1
const QuotedFieldToken = 2
const FieldSeparatorToken = 3
const LineSeparatorToken = 4

//internal map of patterns to tokens
var patterns = []string{
	`[^",\r\n]+`,
	`"[^"\\\\]*(?:\\\\.[^"\\\\]*)*"`,
	`,`,
	`\r\n?|\n`,
}

var tokens = []int{
	PlainFieldToken,
	QuotedFieldToken,
	FieldSeparatorToken,
	LineSeparatorToken,
}

type CsvLexer struct {
	Src     []byte
	Pattern *regexp.Regexp
	offset  int
}

type CsvToken struct {
	Token int
	Value []byte
}

func compile() *regexp.Regexp {
	buffer := new(bytes.Buffer)
	buffer.WriteString("(")
	buffer.WriteString(strings.Join(patterns, ")|("))
	buffer.WriteString(")")
	return regexp.MustCompile(buffer.String())
}

func New(src []byte) *CsvLexer {
	return &CsvLexer{Src: src, Pattern: compile()}
}

func (l *CsvLexer) GetNext() (token *CsvToken) {
	subject := l.Src[l.offset:]
	indexes := l.Pattern.FindSubmatchIndex(subject)
	indexes = indexes[2:]
	var index int
	for i := 0; i < len(indexes); i++ {
		if indexes[i] > 0 {
			index = i
			break
		}
	}
	token = &CsvToken{Token: tokens[index/2], Value: subject[0:indexes[index]]}
	l.offset = indexes[index]
	return
}
