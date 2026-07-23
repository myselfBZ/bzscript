package token

type TokenType string


var Keywords = map[string]string{
	"fun":FUNCTION,
	"true":TRUE,
	"false":FALSE,
	"var":VAR,
	"while":WHILE,
	"if":IF,
	"else":ELSE,
	"map":MAP,
}

type Token struct {
	Type    TokenType
	Literal string
}

func NewToken(t TokenType, literal string) Token {
	return Token{
		Type: t,
		Literal: literal,
	}
}

const (
	ILLEGAL        = "ILLEGAL"
	EOF            = "EOF"
	IDENT          = "IDENT"
	INT            = "INT"
	FLOAT		   = "FLOAT"
	ASSIGN         = "="
	PLUS           = "+"
	COMMA          = ","
	LPAREN         = "("
	RPAREN         = ")"
	LBRACE         = "{"
	RBRACE         = "}"
	FUNCTION       = "FUNCTION"
	VAR            = "VAR"
	WHILE          = "WHILE"
	MAP            = "MAP"
	MINUS          = "-"
	DIVISION       = "/"
	MULTIPLICATION = "*"
	LT             = "<"
	GT             = ">"
	RETURN         = "RETURN"
	IF             = "IF"
	ELSE           = "ELSE"
	TRUE           = "TRUE"
	FALSE          = "FALSE"
	EQ             = "=="
	NOT_EQ         = "!="
	GTOREQ         = ">="
	LTOREQ         = "<="
	BANG           = "!"
	STRING         = "STRING"
)
