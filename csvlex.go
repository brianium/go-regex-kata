package regexkata

import "regexp"
import "bytes"
import "strings"

const PlainFieldToken = 1
const QuotedFieldToken = 2
const FieldSeparatorToken = 3
const LineSeparatorToken = 4

//internal map of patterns to tokens
var tokens = map[string]int{
	`[^",\r\n]+`:                     PlainFieldToken,
	`"[^"\\\\]*(?:\\\\.[^"\\\\]*)*"`: QuotedFieldToken,
	`,`:        FieldSeparatorToken,
	`\r\n?|\n`: LineSeparatorToken,
}

type CsvLexer struct {
	Src     []byte
	Pattern *regexp.Regexp
}

type CsvToken struct {
	Token int
	Value string
}

func keys(mp map[string]int) []string {
	ks := make([]string, 0, len(mp))
	for k := range mp {
		ks = append(ks, k)
	}
	return ks
}

func compile(tokenMap map[string]int) *regexp.Regexp {
	buffer := new(bytes.Buffer)
	buffer.WriteString("((")
	patterns := keys(tokenMap)
	buffer.WriteString(strings.Join(patterns, ")|("))
	buffer.WriteString("))")
	return regexp.MustCompile(buffer.String())
}

func New(src []byte) *CsvLexer {
	return &CsvLexer{Src: src, Pattern: compile(tokens)}
}

func (l *CsvLexer) GetNext() *CsvToken {
	return &CsvToken{Token: PlainFieldToken, Value: "One"}
}
