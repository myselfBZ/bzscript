package lexer

import (
	"testing"

	"github.com/myselfBZ/bzscript/token"
)



func TestLexer(t *testing.T) {
	input := []struct{
		input string
		expected token.Token
	}{
		{ input: "     var    ", expected: token.NewToken(token.VAR, "var")},
		{ input: "     fun    ", expected: token.NewToken(token.FUNCTION, "fun")},
		{ input: "     if    \n", expected: token.NewToken(token.IF, "if")},
		{ input: "+", expected: token.NewToken(token.PLUS, "+")},
		{ input: "-", expected: token.NewToken(token.MINUS, "-")},
		{ input: "map", expected: token.NewToken(token.MAP, "map")},
		{ input: "else", expected: token.NewToken(token.ELSE, "else")},
		{ input: "{", expected: token.NewToken(token.LBRACE, "{")},
		{ input: "}", expected: token.NewToken(token.RBRACE, "}")},
		{ input: "(", expected: token.NewToken(token.LPAREN, "(")},
		{ input: ")", expected: token.NewToken(token.RPAREN, ")")},
		{ input: "1.", expected: token.NewToken(token.FLOAT, "1.")},
		{ input: "0.1", expected: token.NewToken(token.FLOAT, "0.1")},
		{ input: "3.14", expected: token.NewToken(token.FLOAT, "3.14")},
		{ input: "\"Hello\"", expected: token.NewToken(token.STRING, "Hello")},
		{ input: "\"Hello", expected: token.NewToken(token.ILLEGAL, "Hello")},
	}

	for _, tt := range input {
		l := New(tt.input) 
		tok := l.NextToken()

		if tok.Type != tt.expected.Type {
			t.Errorf("expected token type: %s, got: %s", tt.expected.Type, tok.Type)
		}

		if tok.Literal != tt.expected.Literal {
			t.Errorf("expected token literal: %s, got literal: %s", tt.expected.Literal, tok.Literal)
		}
	}
}
