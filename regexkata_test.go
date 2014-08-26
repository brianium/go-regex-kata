package regexkata

import (
	"regexp"
	"testing"
	"unicode/utf8"
)

//Match pattern against byte array
func TestMatch(t *testing.T) {
	pattern := "^(B|b)rian$"
	if match, err := regexp.Match(pattern, []byte("Brian")); match != true {
		t.Errorf("Brian did not match %q %v", pattern, err)
	}

	if match, err := regexp.Match(pattern, []byte("brian")); match != true {
		t.Errorf("brian did not match %q %v", pattern, err)
	}
}

//MatchReader against a rune reader
type TestRuneReader struct{}
type RuneError struct{}

func (e *RuneError) Error() string {
	return "Rune error!!!"
}
func (reader *TestRuneReader) ReadRune() (r rune, size int, e error) {
	r = '\a'
	size = utf8.RuneLen(r)
	if r == '\v' {
		return 0, 0, &RuneError{}
	}
	return
}
func TestMatchReader(t *testing.T) {
	pattern, reader := "\a", &TestRuneReader{}
	run := '\a'
	reader.ReadRune()

	if match, err := regexp.MatchReader(pattern, reader); match != true {
		t.Errorf("MatchReader did not match %v %v", run, err)
	}
}

//MatchString against a string literal
func TestMatchString(t *testing.T) {
	pattern, upper, lower := "^(B|b)rian$", "Brian", "brian"

	if match, err := regexp.MatchString(pattern, upper); match != true {
		t.Errorf("MatchString did not match %q %v", upper, err)
	}

	if match, err := regexp.MatchString(pattern, lower); match != true {
		t.Errorf("MatchString did not match %q %v", lower, err)
	}
}

//MatchString using regexp producted by regexp.QuoteMeta
func TestMatchStringWithQuoteMeta(t *testing.T) {
	pattern, str := regexp.QuoteMeta("[foo]"), "[foo]"

	if match, err := regexp.MatchString(pattern, str); match != true {
		t.Errorf("MatchString did not match %q %v", str, err)
	}
}

//Compiling a regular expression
func TestCompileRegexp(t *testing.T) {
	regex, err := regexp.Compile("^(B|b)rian$")

	if regex == nil {
		t.Errorf("Regex did not compile %v", err)
	}
}

//Compiled regex are safe for access from multiple go routines
func TestCompiledRegexInGoRoutine(t *testing.T) {
	regex, err := regexp.Compile("^(B|b)rian$")
	if err != nil {
		t.Errorf("Regex did not compile %v", err)
	}

	ch := make(chan bool)
	tests := []string{"brian", "Brian"}
	for _, test := range tests {
		go func(t string) {
			ch <- regex.MatchString(t)
		}(test)
	}

	first, second := <-ch, <-ch

	if first && second {
		return
	}

	t.Error("String did not match")
}

//Expand matches in a template
func TestExpand(t *testing.T) {
	regex, err := regexp.Compile(`^(Brian)[\s]+Scaturro$`)
	if err != nil {
		t.Error("Regexp could not compile")
	}
	subject, template := []byte("Brian Scaturro"), []byte("Your name is $1")
	dest := make([]byte, 0)
	expanded := regex.Expand(dest, template, subject, regex.FindSubmatchIndex(subject))
	destStr := string(expanded)
	if destStr != "Your name is Brian" {
		t.Errorf("Expected 'Your name is Brian' got '%q'", destStr)
	}
}

//ExpandString matches using strings
func TestExpandString(t *testing.T) {
	regex, err := regexp.CompilePOSIX("Set|SetValue")
	if err != nil {
		t.Error("Regexp could not compile")
	}
	subject, template := "SetValue", "POSIX match was $0"
	dest := make([]byte, 0)
	expanded := regex.ExpandString(dest, template, subject, regex.FindStringSubmatchIndex(subject))
	destStr := string(expanded)
	if destStr != "POSIX match was SetValue" {
		t.Errorf("Expected 'POSIX match was SetValue' got '%q'", destStr)
	}
}

//Find should return leftmost match
func TestFind(t *testing.T) {
	regex, err := regexp.Compile("Set|SetValue")
	if err != nil {
		t.Error("Regexp could not compile")
	}
	subject := []byte{'S', 'e', 't', 'V', 'a', 'l', 'u', 'e'}
	match := regex.Find(subject)
	if match == nil {
		t.Error("Could not find match for regexp.Find")
		return
	}
	str := string(match)
	if str != "Set" {
		t.Errorf("Expected 'Set', got %q", str)
	}
}

//FindAll should return all matches
func TestFindAll(t *testing.T) {
	regex, err := regexp.Compile("(brian|bryce|jason)")
	if err != nil {
		t.Error("Regexp could not compile")
	}
	subject := []byte("brian and bryce and jason are rad dudes")
	expected := []string{"brian", "bryce", "jason"}

	matches := regex.FindAll(subject, 3)

	if matches == nil {
		t.Error("Could not find matches for regexp.FindAll")
	}

	for i, match := range matches {
		strMatch := string(match)
		if strMatch != expected[i] {
			t.Errorf("Expected '%q', got '%q'", expected[i], strMatch)
		}
	}
}

//start using a helper function with dynamic types for Assertion
func Assert(expected interface{}, actual interface{}, t *testing.T) {
	if expected != actual {
		t.Errorf("Expected '%v', got '%v'", expected, actual)
	}
}

//FindIndex should return location of match - now enforcing compilation to save on error checks
func TestFindIndex(t *testing.T) {
	regex := regexp.MustCompile("Brian")
	subject := []byte("My name is Brian, pleased to meet you")
	index := regex.FindIndex(subject)
	Assert(index[0], 11, t)
	Assert(index[1], 16, t)
}

//FindAllIndex should return a collection of indexes
func TestFindAllIndex(t *testing.T) {
	regex := regexp.MustCompile("Brian")
	subject := []byte("Brian, your name is Brian right?")
	indexes := regex.FindAllIndex(subject, 2)
	Assert(indexes[0][0], 0, t)
	Assert(indexes[0][1], 5, t)
	Assert(indexes[1][0], 20, t)
	Assert(indexes[1][1], 25, t)
}

//FindString should return a matched string
func TestFindString(t *testing.T) {
	regex := regexp.MustCompile("Brian")
	subject := "Hello Brian"
	match := regex.FindString(subject)
	Assert("Brian", match, t)
}

//FindAllString should return all matching strings
func TestFindAllString(t *testing.T) {
	regex := regexp.MustCompile("Brian|Austin")
	subject := "Hello Austin, this is Brian"
	matches := regex.FindAllString(subject, 2)
	Assert("Austin", matches[0], t)
	Assert("Brian", matches[1], t)
}

//FindStringIndex should return index of string
func TestFindStringIndex(t *testing.T) {
	regex := regexp.MustCompile("Brian")
	subject := "Hello Brian"
	index := regex.FindStringIndex(subject)
	Assert(6, index[0], t)
	Assert(11, index[1], t)
}

//FindAllStringIndex should return multiple indexes
func TestFindAllStringIndex(t *testing.T) {
	regex := regexp.MustCompile("Brian")
	subject := "Brian. Meet Brian"
	indexes := regex.FindAllStringIndex(subject, 2)
	Assert(0, indexes[0][0], t)
	Assert(5, indexes[0][1], t)
	Assert(12, indexes[1][0], t)
	Assert(17, indexes[1][1], t)
}

//FindStringSubmatch should return an array of submathces
func TestFindStringSubmatch(t *testing.T) {
	regex := regexp.MustCompile("Hello.*(world)")
	subject := "Hello brave new world"
	matches := regex.FindStringSubmatch(subject)
	Assert("world", matches[1], t)
}

func TestFindAllStringSubmatch(t *testing.T) {
	regex := regexp.MustCompile("a(x*)b")
	matches := regex.FindAllStringSubmatch("-axxb-ab-", -1)
	Assert("xx", matches[0][1], t)
	Assert("", matches[1][1], t)
}
