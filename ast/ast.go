package ast

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/myselfBZ/bzscript/token"
)

var (
	_ Expression = (*InfixExpression)(nil)
	_ Expression = (*Intiger)(nil)
	_ Expression = (*Bool)(nil)
	_ Expression = (*Float)(nil)
	_ Expression = (*String)(nil)
	_ Expression = (*Ident)(nil)
	_ Expression = (*FunctionLiteral)(nil)
	_ Expression = (*AnonymousFuncLiteral)(nil)
	_ Expression = (*ArrayLiteral)(nil)
	_ Expression = (*StructMemberAccess)(nil)
)

var (
	_ Statement = (*VarStatement)(nil)
	_ Statement = (*ExpressionStatement)(nil)
	_ Statement = (*Block)(nil)
	_ Statement = (*IfStatement)(nil)
	_ Statement = (*ReturnStatement)(nil)
	_ Statement = (*BreakStatement)(nil)
	_ Statement = (*ContinueStatement)(nil)
	_ Statement = (*AssignStatement)(nil)
)

type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, v := range p.Statements {
		out.WriteString(v.String())
		out.WriteString("\n")
	}
	return out.String()
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
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
		if i != len(f.Params)-1 {
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
		if i != len(f.Params)-1 {
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

type FunctionCall struct {
	Token    *token.Token
	Function Expression
	Args     []Expression
}

func (f *FunctionCall) String() string {
	buff := &bytes.Buffer{}
	buff.WriteString(f.Function.String())
	buff.WriteString("(")
	for i, e := range f.Args {
		buff.WriteString(e.String())
		if i != len(f.Args)-1 {
			buff.WriteString(", ")
		}
	}
	buff.WriteString(")")
	return buff.String()
}
func (f *FunctionCall) TokenLiteral() string {
	return f.Token.Literal
}
func (f *FunctionCall) expressionNode() {}

type ReturnStatement struct {
	Token *token.Token
	Value Expression
}

func (f *ReturnStatement) String() string {
	buff := &bytes.Buffer{}
	buff.WriteString(f.Token.Literal)
	buff.WriteString(" ")
	buff.WriteString(f.Value.String())
	return buff.String()
}
func (f *ReturnStatement) TokenLiteral() string {
	return f.Token.Literal
}
func (f *ReturnStatement) statementNode() {}

type WhileLoop struct {
	Token     *token.Token
	Condition Expression
	Body      *Block
}

func (w *WhileLoop) String() string {
	buff := &bytes.Buffer{}
	buff.WriteString("while")
	buff.WriteString(" ")
	buff.WriteString(w.Condition.String())
	buff.WriteString(" ")
	buff.WriteString(w.Body.String())
	return buff.String()
}
func (w *WhileLoop) TokenLiteral() string {
	return w.Token.Literal
}
func (w *WhileLoop) statementNode() {}

type BreakStatement struct {
	Token *token.Token
}

func (b *BreakStatement) String() string {
	return "break"
}
func (b *BreakStatement) TokenLiteral() string {
	return b.Token.Literal
}
func (b *BreakStatement) statementNode() {}

type ContinueStatement struct {
	Token *token.Token
}

func (c *ContinueStatement) String() string {
	return "continue"
}
func (c *ContinueStatement) TokenLiteral() string {
	return c.Token.Literal
}
func (c *ContinueStatement) statementNode() {}

type AssignStatement struct {
	Token    *token.Token
	LHS, RHS Expression
}

func (a *AssignStatement) String() string {
	buff := &bytes.Buffer{}
	buff.WriteString(a.LHS.String())
	buff.WriteString(" ")
	buff.WriteString("=")
	buff.WriteString(" ")
	buff.WriteString(a.RHS.String())
	return buff.String()
}
func (a *AssignStatement) TokenLiteral() string {
	return a.Token.Literal
}
func (a *AssignStatement) statementNode() {}

type ArrayLiteral struct {
	Token    *token.Token
	Elements []Expression
}

func (a *ArrayLiteral) String() string {
	buff := &bytes.Buffer{}
	buff.WriteString("[")
	elements := []string{}
	for _, el := range a.Elements {
		elements = append(elements, el.String())
	}
	buff.WriteString(strings.Join(elements, ","))
	buff.WriteString("]")
	return buff.String()
}
func (a *ArrayLiteral) TokenLiteral() string {
	return a.Token.Literal
}
func (a *ArrayLiteral) expressionNode() {}

type IndexOperator struct {
	Token *token.Token
	Index Expression
	Left  Expression
}

func (i *IndexOperator) String() string {
	buff := &bytes.Buffer{}
	buff.WriteString(i.Left.String())
	buff.WriteString("[")
	buff.WriteString(i.Index.String())
	buff.WriteString("]")
	return buff.String()
}
func (i *IndexOperator) TokenLiteral() string {
	return i.Token.Literal
}
func (i *IndexOperator) expressionNode() {}

type MapLiteral struct {
	Token *token.Token
	Kv    map[Expression]Expression
}

func (m *MapLiteral) String() string {
	buff := &bytes.Buffer{}
	buff.WriteString("map")
	buff.WriteString("{")
	for _, k := range m.Kv {
		buff.WriteString(k.String())
		buff.WriteString(":")
		buff.WriteString(m.Kv[k].String())
		buff.WriteString(",")
	}
	buff.WriteString("}")
	return buff.String()
}
func (m *MapLiteral) TokenLiteral() string {
	return m.Token.Literal
}
func (m *MapLiteral) expressionNode() {}

type StructLiteral struct {
	Token  *token.Token
	Name   *Ident
	Fields []*FieldLiteral
}

func (s *StructLiteral) statementNode() {}
func (s *StructLiteral) String() string {
	buff := &bytes.Buffer{}
	buff.WriteString("struct")
	buff.WriteString(" ")
	buff.WriteString(s.Name.String())
	buff.WriteString(" ")
	buff.WriteString("{")
	fields := []string{}
	for _, f := range s.Fields {
		fields = append(fields, f.Name)
	}
	buff.WriteString(strings.Join(fields, ","))
	buff.WriteString("}")
	return buff.String()
}
func (s *StructLiteral) TokenLiteral() string {
	return s.Token.Literal
}

type FieldLiteral struct {
	Token *token.Token
	Name  string
}

func (f *FieldLiteral) statementNode() {}
func (f *FieldLiteral) String() string {
	return f.Name
}
func (f *FieldLiteral) TokenLiteral() string {
	return f.Token.Literal
}

type StructMemberAccess struct {
 	Token *token.Token // the . operator
	Rhs, Lhs Expression
}
func (s *StructMemberAccess) String() string {
	buff := &bytes.Buffer{}
	buff.WriteString(s.Rhs.String())
	buff.WriteString(".")
	buff.WriteString(s.Lhs.String())
	return buff.String()
}
func (s *StructMemberAccess) TokenLiteral() string {
	return s.Token.Literal
}
func (s *StructMemberAccess) expressionNode() {} 
