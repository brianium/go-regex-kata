package regexkata

const PlainFieldToken = 1
const QuotedFieldToken = 2
const FieldSeparatorToken = 3
const LineSeparatorToken = 4

//internal map of patterns to tokens
var tokens = map[string]int{
	`[^",\r\n]+~`:                    PlainFieldToken,
	`"[^"\\\\]*(?:\\\\.[^"\\\\]*)*"`: QuotedFieldToken,
	`,`:        FieldSeparatorToken,
	`\r\n?|\n`: LineSeparatorToken,
}

type CsvLexer struct {
	Src string
}

type CsvToken struct {
	Token int
	Value string
}

func New(src string) *CsvLexer {
	return &CsvLexer{Src: src}
}

func (l *CsvLexer) GetNext() *CsvToken {
	return &CsvToken{Token: PlainFieldToken, Value: "One"}
}
