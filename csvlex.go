package regexkata

import (
	"fmt"
	"regexp"
)
import "bytes"
import "strings"

const PlainFieldToken = 1
const QuotedFieldToken = 2
const FieldSeparatorToken = 3
const LineSeparatorToken = 4

//internal map of patterns to tokens
var tokenMap = map[string]int{
	`[^",\r\n]+`:                     PlainFieldToken,
	`"[^"\\\\]*(?:\\\\.[^"\\\\]*)*"`: QuotedFieldToken,
	`,`:        FieldSeparatorToken,
	`\r\n?|\n`: LineSeparatorToken,
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

func keys(mp map[string]int) []string {
	ks := make([]string, 0, len(mp))
	for k := range mp {
		ks = append(ks, k)
	}
	return ks
}

func values(mp map[string]int) []int {
	vals := make([]int, 0, len(mp))
	for _, v := range mp {
		vals = append(vals, v)
	}
	return vals
}

var patterns = keys(tokenMap)
var tokens = values(tokenMap)

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

	value := subject[0:indexes[index]]
	fmt.Printf("Index = %d, Indexes = %v\n", index, indexes)
	tokenType := tokens[index/2]

	fmt.Printf("%v %v\n", string(value), tokenType)
	token = &CsvToken{Token: tokenType, Value: value}
	l.offset = indexes[index]
	return
}
