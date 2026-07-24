package ast

import (
	"bytes"
	"fmt"

	"github.com/myselfBZ/bzscript/token"
)

var (
	_ Expression = (*InfixExpression)(nil)
	_ Expression = (*Intiger)(nil)
	_ Expression = (*Bool)(nil)
	_ Expression = (*Float)(nil)
	_ Expression = (*String)(nil)
	_ Expression = (*Ident)(nil)
)

var (
	_ Statement = (*VarStatement)(nil)
	_ Statement = (*ExpressionStatement)(nil)
	_ Statement = (*Block)(nil)
	_ Statement = (*IfStatement)(nil)
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
	return fmt.Sprintf("var %s = %s", v.Ident.String(), v.Value.String())
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
	Operator   string
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
	Left     Expression
	Right    Expression
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
	return fmt.Sprintf("\"%s\"", s.Value)
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
	Token      *token.Token
	Expression Expression
}

func (e *ExpressionStatement) String() string {
	return e.Expression.String()
}
func (e *ExpressionStatement) TokenLiteral() string {
	return e.Token.Literal
}
func (e *ExpressionStatement) statementNode() {}

type Block struct {
	Token      *token.Token
	Statements []Statement
}

func (b *Block) String() string {
	buff := bytes.Buffer{}
	buff.WriteString("{\n")
	for _, s := range b.Statements {
		buff.WriteString(s.String())
		buff.WriteString("\n")
	}
	buff.WriteString("}")
	return buff.String()
}
func (b *Block) TokenLiteral() string {
	return b.Token.Literal
}
func (b *Block) statementNode() {}

type IfStatement struct {
	Token       *token.Token
	Condition   Expression
	Consequence *Block
	Alternative *Block
}

func (i *IfStatement) String() string {
	buff := bytes.Buffer{}
	buff.WriteString("if")
	buff.WriteString(" ")
	buff.WriteString(i.Condition.String())
	buff.WriteString(" ")
	buff.WriteString(i.Consequence.String())

	if i.Alternative != nil {
		buff.WriteString(" ")
		buff.WriteString("else")
		buff.WriteString(" ")
		buff.WriteString(i.Alternative.String())
	}

	return buff.String()
}
func (i *IfStatement) TokenLiteral() string {
	return i.Token.Literal
}
func (i *IfStatement) statementNode() {}

type FunctionLiteral struct {
	Token  *token.Token
	Params []*Ident
	Ident  *Ident
	Body   *Block
}

func (f *FunctionLiteral) String() string {
	buff := &bytes.Buffer{}
	buff.WriteString("fun")
	buff.WriteString(" ")
	buff.WriteString(f.Ident.String())
	buff.WriteString("(")
	for i, ident := range f.Params {
		buff.WriteString(ident.String())
		if i != len(f.Params) - 1 {
			buff.WriteString(", ")
		}
	}
	buff.WriteString(")")
	buff.WriteString(" ")
	buff.WriteString(f.Body.String())
	return buff.String()
}
func (f *FunctionLiteral) TokenLiteral() string {
	return f.Token.Literal
}
func (f *FunctionLiteral) expressionNode() {}

type AnonymousFuncLiteral struct {
	Token  *token.Token
	Params []*Ident
	Body   *Block
}

func (f *AnonymousFuncLiteral) String() string {
	buff := &bytes.Buffer{}
	buff.WriteString("fun")
	buff.WriteString(" ")
	buff.WriteString("(")
	for i, ident := range f.Params {
		buff.WriteString(ident.String())
		if i != len(f.Params) - 1 {
			buff.WriteString(", ")
		}
	}
	buff.WriteString(")")
	buff.WriteString(" ")
	buff.WriteString(f.Body.String())
	return buff.String()
}
func (f *AnonymousFuncLiteral) TokenLiteral() string {
	return f.Token.Literal
}
func (f *AnonymousFuncLiteral) expressionNode() {}
