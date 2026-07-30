package lexer

import (
	"bytes"
	"unicode"

	"github.com/myselfBZ/bzscript/token"
)

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

type Lexer struct {
	input   string
	pos     int
	ch      byte
	readPos int
}

func isDigit(ch byte) bool {
	return unicode.IsDigit(rune(ch))
}

func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch))
}

func (l *Lexer) peek() byte {
	if l.readPos >= len(l.input) {
		return 0
	}
	return l.input[l.readPos]
}

func (l *Lexer) skipWhiteSpace() {
	for l.ch == ' ' || l.ch == '\n' || l.ch == '\t' {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.pos
	for isLetter(l.ch) || l.ch == '_' {
		l.readChar()
	}
	return l.input[position:l.pos]
}

func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPos]
	}
	l.pos = l.readPos
	l.readPos += 1
}


func (l *Lexer) readDigit() string {
	position := l.pos
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.pos]
}

func (l *Lexer) readNumberToken() token.Token {
	intPart := l.readDigit()
	if l.ch != '.' {
		return token.NewToken(token.INT, intPart)
	}
	l.readChar()
	fracPart := l.readDigit()
	return token.NewToken(token.FLOAT, intPart + "." + fracPart)
}

// TODO: handle malformed strings and edge cases
func (l *Lexer) readString() string {
	l.readChar()
	buff := bytes.NewBuffer([]byte{})
	for l.ch != '"' && l.ch != 0 {
		buff.WriteString(string(l.ch))
		l.readChar()
	}
	return buff.String()
}

func (l *Lexer) NextToken() *token.Token {
	l.skipWhiteSpace()
	var t token.Token
	switch l.ch {
	// TODO start of the float
	// case '.':
	// ...
	case '"':
		strContent := l.readString()
		if l.ch == '"' {
			t = token.NewToken(token.STRING, strContent)
		} else {
			t = token.NewToken(token.ILLEGAL, strContent)
		}
	case '=':
		if l.peek() == '=' {
			l.readChar()
			t = token.NewToken(token.EQ, string(l.ch)+"=")
		} else {
			t = token.NewToken(token.ASSIGN, string(l.ch))
		}
	case '-':
		t = token.NewToken(token.MINUS, string(l.ch))
	case '/':
		t = token.NewToken(token.DIVISION, string(l.ch))
	case '+':
		t = token.NewToken(token.PLUS, string(l.ch))
	case '(':
		t = token.NewToken(token.LPAREN, string(l.ch))
	case '!':
		if l.peek() == '=' {
			l.readChar()
			t = token.NewToken(token.NOT_EQ, "!=")
		} else {
			t = token.NewToken(token.BANG, "!")
		}
	case ')':
		t = token.NewToken(token.RPAREN, string(l.ch))
	case '{':
		t = token.NewToken(token.LBRACE, string(l.ch))
	case '}':
		t = token.NewToken(token.RBRACE, string(l.ch))
	case ':':
		t = token.NewToken(token.COLON, string(l.ch))
	case '[':
		t = token.NewToken(token.LBRACK, string(l.ch))
	case ']':
		t = token.NewToken(token.RBRACK, string(l.ch))
	case '>':
		if l.peek() == '=' {
			l.readChar()
			t = token.NewToken(token.GTOREQ, ">=")
		} else {
			t = token.NewToken(token.GT, string(l.ch))
		}
	case ',':
		t = token.NewToken(token.COMMA, string(l.ch))
	case ';':
		t = token.NewToken(token.SEMCOLON, string(l.ch))
	case '<':
		if l.peek() == '=' {
			l.readChar()
			t = token.NewToken(token.LTOREQ, "<=")
		} else {
			t = token.NewToken(token.LT, string(l.ch))
		}
	case 0:
		t.Literal = "EOF"
		t.Type = token.EOF
	case '*':
		t = token.NewToken(token.MULTIPLICATION, string(l.ch))
	case '%':
		t = token.NewToken(token.MODULO, string(l.ch))
	default:
		if isDigit(l.ch) {
			t = l.readNumberToken()
			return &t
		} else if isLetter(l.ch) {
			word := l.readIdentifier()
			if kind, ok := token.Keywords[word]; ok {
				t.Literal = word
				t.Type = token.TokenType(kind)
			} else {
				t.Type = token.IDENT
				t.Literal = word
			}
			return &t
		} else {
			t = token.NewToken(token.ILLEGAL, string(l.ch))
		}
	}

	l.readChar()
	return &t
}
