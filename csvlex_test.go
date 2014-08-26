package regexkata

import (
	"testing"
)

func AssertNotNull(obj interface{}, t *testing.T) {
	if obj == nil {
		t.Error("Expected non nil value, got nil")
	}
}

func TestNew(t *testing.T) {
	lex := New("src")
	AssertNotNull(lex, t)
}

func TestGetNext(t *testing.T) {
	lex := New(`One,Two,Three`)
	token := lex.GetNext()
	Assert(PlainFieldToken, token.Token, t)
	Assert("One", token.Value, t)
}
