package ast

import (
	"bytes"
	"fmt"

	"github.com/myselfBZ/bzscript/token"
)

var(
	_ Expression = (*InfixExpression)(nil)
	_ Expression = (*Intiger)(nil)
	_ Expression = (*Bool)(nil)
	_ Expression = (*Float)(nil)
	_ Expression = (*String)(nil)
)

var(
	_ Statement = (*VarStatement)(nil)
	_ Statement = (*ExpressionStatement)(nil)
)

type Program struct {
	Statements []Statement
}

func (p *Program) AddStatement(s Statement) {
	p.Statements = append(p.Statements, s)
}

type Node interface {
	String() string
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type VarStatement struct {
	Token *token.Token
	Ident *Ident
	Value Expression
}
func (v *VarStatement) String() string {
	return fmt.Sprintf("%s", v.Value)
}
func (v *VarStatement) TokenLiteral() string {
	return v.Token.Literal
}
func (v *VarStatement) statementNode() {}


type Expression interface {
	Node
	expressionNode()
}

type PrefixExpression struct {
	Operator string
	Expression Expression
}
func (i *PrefixExpression) TokenLiteral() string {
	return i.Operator
}
func (i *PrefixExpression) expressionNode() {}
func (i *PrefixExpression) String() string {
	return i.Operator + " " + i.Expression.String()
}

type InfixExpression struct {
	Operator string
	Left Expression
	Right Expression
}

func (i *InfixExpression) TokenLiteral() string {
	return i.Operator
}
func (i *InfixExpression) expressionNode() {}
func (i *InfixExpression) String() string {
	buff := bytes.Buffer{}
	leftStr := i.Left.String()
	rightStr := i.Right.String()
	buff.WriteString(leftStr)
	buff.WriteString(" ")
	buff.WriteString(i.Operator)
	buff.WriteString(" ")
	buff.WriteString(rightStr)
	return buff.String()
}


type Bool struct {
	Token *token.Token
	Value bool
}
func (b *Bool) expressionNode() {}
func (b *Bool) String() string {
	return fmt.Sprintf("%v", b.Value)
}
func (b *Bool) TokenLiteral() string {
	return b.Token.Literal
}

type Float struct {
	Token *token.Token
	Value float64
}
func (f *Float) expressionNode() {}
func (f *Float) String() string {
	return fmt.Sprintf("%v", f.Value)
}
func (f *Float) TokenLiteral() string {
	return f.Token.Literal
}


type Intiger struct {
	Token *token.Token
	Value int64
}
func (i *Intiger) String() string {
	return fmt.Sprintf("%d", i.Value)
}
func (i *Intiger) TokenLiteral() string {
	return i.Token.Literal
}
func (i *Intiger) expressionNode() {}

type String struct {
	Token *token.Token
	Value string
}
func (s *String) String() string {
	return fmt.Sprintf("%s", s.Value)
}
func (s *String) TokenLiteral() string {
	return s.Token.Literal
}
func (s *String) expressionNode() {}


type Ident struct {
	Token *token.Token
	Value string
}
func (i *Ident) String() string {
	return fmt.Sprintf("%s", i.Value)
}
func (i *Ident) TokenLiteral() string {
	return i.Token.Literal
}
func (i *Ident) expressionNode() {}


type ExpressionStatement struct {
	Token *token.Token
	Expression Expression
}
func (e *ExpressionStatement) String() string {
	return e.Expression.String()
}
func (e *ExpressionStatement) TokenLiteral() string {
	return e.Token.Literal
}
func (e *ExpressionStatement) statementNode() {}
