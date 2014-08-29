package regexkata

import (
	"testing"
)

func AssertNotNull(obj interface{}, t *testing.T) {
	if obj == nil {
		t.Error("Expected non nil value, got nil")
	}
}

func AssertToken(tokInt int, val string, tok *CsvToken, t *testing.T) {
	Assert(tokInt, tok.Token, t)
	Assert(val, string(tok.Value), t)
}

func TestNew(t *testing.T) {
	lex := New([]byte("src"))
	AssertNotNull(lex, t)
}

func TestGetNext(t *testing.T) {
	lex := New([]byte("One,Two,Three"))
	token := lex.GetNext()
	Assert(PlainFieldToken, token.Token, t)
	Assert("One", string(token.Value), t)
}

func TestGetSequentialTokens(t *testing.T) {
	lex := New([]byte("One,Two,Three"))
	var token *CsvToken
	for i := 0; i < 2; i++ {
		token = lex.GetNext()
	}
	Assert(FieldSeparatorToken, token.Token, t)
	Assert(",", string(token.Value), t)
}

func TestGetNextTokenReturnsNilWhenNoMatch(t *testing.T) {
	lex := New([]byte("One"))
	var token *CsvToken
	for i := 0; i < 2; i++ {
		token = lex.GetNext()
	}
	if token != nil {
		t.Error("Expected token to be nil")
	}
}

func TestLexAll(t *testing.T) {
	lex := New([]byte(`One,Two` + "\n" + `"Three","Four"`))
	tokens := lex.LexAll()

	AssertToken(PlainFieldToken, "One", tokens[0], t)
	AssertToken(FieldSeparatorToken, ",", tokens[1], t)
	AssertToken(PlainFieldToken, "Two", tokens[2], t)
	AssertToken(LineSeparatorToken, "\n", tokens[3], t)
	AssertToken(QuotedFieldToken, `"Three"`, tokens[4], t)
	AssertToken(FieldSeparatorToken, ",", tokens[5], t)
	AssertToken(QuotedFieldToken, `"Four"`, tokens[6], t)
}
