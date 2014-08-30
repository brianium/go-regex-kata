package regexkata

import (
	"bytes"
	"fmt"
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
	"\r\n?|\n",
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

type LexError struct {
	Offset int
	Bytes  []byte
}

func NewError(offset int, bytes []byte) error {
	return &LexError{Offset: offset, Bytes: bytes}
}

func (e *LexError) Error() string {
	return fmt.Sprintf("Invalid sequence %q at offset %d", string(e.Bytes), e.Offset)
}

func (l *CsvLexer) GetNext() (token *CsvToken, e error) {
	subject := l.Src[l.offset:]
	if len(subject) == 0 {
		return nil, nil
	}
	indexes := l.Pattern.FindSubmatchIndex(subject)
	if indexes == nil {
		return nil, NewError(l.offset, subject)
	}
	indexes = indexes[2:]
	var index int
	for i := range indexes {
		if indexes[i] > 0 {
			index = i
			break
		}
	}
	token = &CsvToken{Token: tokens[index/2], Value: subject[0:indexes[index]]}
	l.offset = l.offset + indexes[index]
	return token, nil
}

func (l *CsvLexer) LexAll() ([]*CsvToken, error) {
	tokens := make([]*CsvToken, 0)
	token, err := l.GetNext()
	if err != nil {
		return nil, err
	}
	for token != nil {
		tokens = append(tokens, token)
		token, err = l.GetNext()
		if err != nil {
			return nil, err
		}
	}
	return tokens, nil
}
