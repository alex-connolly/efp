package efp

import (
	"strings"
	"testing"
)

func TestLexerBasicOperators(t *testing.T) {
	SingleLexer(t, tknAssign, "=")
	SingleLexer(t, tknOpenBrace, "{")
	SingleLexer(t, tknCloseBrace, "}")
	SingleLexer(t, tknOpenSquare, "[")
	SingleLexer(t, tknCloseSquare, "]")
	SingleLexer(t, tknRequired, "!")
	SingleLexer(t, tknComma, ",")
	SingleLexer(t, tknOpenBracket, "(")
	SingleLexer(t, tknCloseBracket, ")")
}

func SingleLexer(t *testing.T, tkn tokenType, data string) {
	l := lex([]byte(data))
	if len(l.tokens) != 1 {
		t.Fail()
	} else if l.tokens[0].tkntype != tkn {
		t.Fail()
	}
}

func TestLexerDuplicateOperators(t *testing.T) {
	MultiLexer(t, tknAssign, "=", 3)
	MultiLexer(t, tknOpenBrace, "{", 2)
	MultiLexer(t, tknCloseBrace, "}", 5)
	MultiLexer(t, tknOpenSquare, "[", 3)
	MultiLexer(t, tknCloseSquare, "]", 8)
	MultiLexer(t, tknRequired, "!", 1)
	MultiLexer(t, tknComma, ",", 11)
	MultiLexer(t, tknOpenBracket, "(", 9)
	MultiLexer(t, tknCloseBracket, ")", 9)
}

func MultiLexer(t *testing.T, tkn tokenType, data string, times int) {
	l := lex([]byte(strings.Repeat(data, times)))
	if len(l.tokens) != times {
		t.Fail()
	}
	for i := 0; i < times; i++ {
		if l.tokens[i].tkntype != tkn {
			t.Fail()
		}
	}
}

func TestLexerValueLength(t *testing.T) {
	l := lex([]byte("hello this is dog"))
	if len(l.tokens) != 4 {
		t.Fail()
	}
	expected := []int{5, 4, 2, 3}
	for i, tk := range l.tokens {
		if tk.end-tk.start != expected[i] {
			t.Fail()
		}
	}
}

func TestLexerTokenString(t *testing.T) {
	l := lex([]byte("hello this is dog"))
	expected := []string{"hello", "this", "is", "dog"}
	for i, tk := range l.tokens {
		if l.tokenString(tk) != expected[i] {
			t.Fail()
		}
	}
}

func TestLexerTokenLengths(t *testing.T) {
	l := lex([]byte("alias x : y = 5"))
	if len(l.tokens) != 6 {
		t.Fail()
	}
}
