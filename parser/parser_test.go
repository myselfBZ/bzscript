package parser

import (
	"fmt"
	"testing"

	"github.com/myselfBZ/bzscript/ast"
	"github.com/myselfBZ/bzscript/lexer"
	"github.com/myselfBZ/bzscript/token"
)

func TestParseMalformedExpressionVar(t *testing.T) {
	input := "var x ="

	l := lexer.New(input)
	p := New(l)

	_ = p.ParseProgram()

	if len(p.Errors()) == 0 {
		t.Error("no parser errors")
	}

	expectedError := fmt.Errorf("invalid prefix expression token: EOF")
	if expectedError.Error() != p.Errors()[0].Error() {
		t.Errorf("expected error '%s', got '%s'", expectedError.Error(), p.Errors()[0].Error())
	}
}

func TestParseMalformedIdentVar(t *testing.T) {
	input := "var 123 = 12123"

	l := lexer.New(input)
	p := New(l)

	_ = p.ParseProgram()

	if len(p.Errors()) == 0 {
		t.Error("no parser errors")
	}

	expectedError := p.expectedError(token.IDENT, token.INT)
	if expectedError.Error() != p.Errors()[0].Error() {
		t.Errorf("expected error '%s', got '%s'", expectedError.Error(), p.Errors()[0].Error())
	}
}

func TestParseExpressionStament(t *testing.T) {
	tests := []string{"1 + 1", "2 - 3", "1 + 4 + 4 + 5 - 1 - 34"}

	for _, tt := range tests {
		p := New(lexer.New(tt))
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			t.Errorf("expected 0 errors, got '%d'", len(p.Errors()))
			return
		}

		if len(program.Statements) != 1 {
			t.Errorf("expected 1 statement, got %d", len(program.Statements))
		}

		s := program.Statements[0].String()

		if s != tt {
			t.Errorf("expected end result '%s', got '%s'", tt, s)
		}
	}
}

func TestParseVar(t *testing.T) {
	tests := []struct {
		input         string
		expectedIdent any
		expectedVal   expectedVal
	}{
		{input: "var x = 1", expectedIdent: "x", expectedVal: expectedVal{Type: token.INT, Val: int64(1)}},
		{input: "var x = 1.0", expectedIdent: "x", expectedVal: expectedVal{Type: token.FLOAT, Val: float64(1.0)}},
		{input: "var x = \"Hello, World\"", expectedIdent: "x", expectedVal: expectedVal{Type: token.STRING, Val: "Hello, World"}},
		{input: "var isTrue = true", expectedIdent: "isTrue", expectedVal: expectedVal{Type: token.TRUE, Val: true}},
		{input: "var isTrue = false", expectedIdent: "isTrue", expectedVal: expectedVal{Type: token.FALSE, Val: false}},
	}

	for _, tt := range tests {
		p := New(lexer.New(tt.input))
		program := p.ParseProgram()
		if len(program.Statements) != 1 {
			t.Errorf("expected 1 statement, got %d", len(program.Statements))
			return
		}
		varStmnt, ok := program.Statements[0].(*ast.VarStatement)
		if !ok {
			t.Errorf("expected *ast.VarStatement statement, got '%T'", program.Statements[0])
			return
		}
		if varStmnt.Ident.Value != tt.expectedIdent.(string) {
			t.Errorf("expected identifier '%s', got '%s'", tt.expectedIdent.(string), varStmnt.Ident.Value)
			return
		}
		testLiteralExpression(t, varStmnt.Value, &tt.expectedVal)
	}

}

func TestGrouped(t *testing.T) {
	tests := []struct {
		input    string
		expected ast.Expression
	}{
		{input: "(1)", expected: &ast.Intiger{Value: 1}},
		{input: "-(1)", expected: &ast.PrefixExpression{Expression: &ast.Intiger{Value: 1}, Operator: "-"}},
		{input: "-(1 + 3)", expected: &ast.PrefixExpression{Operator: "-", Expression: &ast.InfixExpression{Operator: "+", Left: &ast.Intiger{Value: 1}, Right: &ast.Intiger{Value: 3}}}},
		{
			input: "(1 + 3) * 2", 
			expected: &ast.InfixExpression{
				Operator: "*", 
				Left: &ast.InfixExpression{Operator: "+", Left: &ast.Intiger{Value: 1}, Right: &ast.Intiger{Value: 3}},
				Right: &ast.Intiger{Value: 2},
			},
		},
	}
	for _, tt := range tests {
		p := New(lexer.New(tt.input))
		program := p.ParseProgram()

		if len(p.Errors()) > 0 {
			printErrors(t, p.Errors())
			t.FailNow()
		}

		if len(program.Statements) != 1 {
			t.Errorf("expected 1 statement, got %d", len(program.Statements))
			return
		}

		exprStatement, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Errorf("expected *ast.ExpressionStatement, got '%T'", program.Statements[0])
			return
		}
		assertExpressionEqual(t, tt.expected, exprStatement.Expression)
	}
}

func TestPrefix(t *testing.T) {
	tests := []struct {
		input    string
		expected *ast.PrefixExpression
	}{
		{input: "-x", expected: &ast.PrefixExpression{Operator: "-", Expression: &ast.Ident{Value: "x"}}},
		{input: "!x", expected: &ast.PrefixExpression{Operator: "!", Expression: &ast.Ident{Value: "x"}}},
		{input: "-1", expected: &ast.PrefixExpression{Operator: "-", Expression: &ast.Intiger{Value: 1}}},
		{input: "!true", expected: &ast.PrefixExpression{Operator: "!", Expression: &ast.Bool{Value: true}}},
		{input: "!false", expected: &ast.PrefixExpression{Operator: "!", Expression: &ast.Bool{Value: false}}},
		{input: "!!false", expected: &ast.PrefixExpression{
			Operator: "!", Expression: &ast.PrefixExpression{Operator: "!", Expression: &ast.Bool{Value: false}},
		}},
		{
			input: "-----1",
			expected: &ast.PrefixExpression{
				Operator: "-",
				Expression: &ast.PrefixExpression{
					Operator: "-",
					Expression: &ast.PrefixExpression{
						Operator: "-",
						Expression: &ast.PrefixExpression{
							Operator: "-",
							Expression: &ast.PrefixExpression{
								Operator:   "-",
								Expression: &ast.Intiger{Value: 1},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		p := New(lexer.New(tt.input))
		program := p.ParseProgram()
		if len(p.Errors()) > 0 {
			printErrors(t, p.Errors())
			t.FailNow()
		}

		if len(program.Statements) != 1 {
			t.Errorf("expected 1 statement, got %d", len(program.Statements))
			return
		}

		exprStatement, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Errorf("expected *ast.ExpressionStatement, got '%T'", program.Statements[0])
			return
		}
		assertExpressionEqual(t, tt.expected, exprStatement.Expression)
	}
}

func TestInfix(t *testing.T) {
	tests := []struct {
		input    string
		expected *ast.InfixExpression
	}{
		{input: "1 + 1", expected: &ast.InfixExpression{Operator: "+", Left: &ast.Intiger{Value: 1}, Right: &ast.Intiger{Value: 1}}},
		{input: "x - y", expected: &ast.InfixExpression{Operator: "-", Left: &ast.Ident{Value: "x"}, Right: &ast.Ident{Value: "y"}}},
		{input: "x == y", expected: &ast.InfixExpression{Operator: "==", Left: &ast.Ident{Value: "x"}, Right: &ast.Ident{Value: "y"}}},
		{input: "x <= y", expected: &ast.InfixExpression{Operator: "<=", Left: &ast.Ident{Value: "x"}, Right: &ast.Ident{Value: "y"}}},
		{input: "x >= y", expected: &ast.InfixExpression{Operator: ">=", Left: &ast.Ident{Value: "x"}, Right: &ast.Ident{Value: "y"}}},
		{input: "x != y", expected: &ast.InfixExpression{Operator: "!=", Left: &ast.Ident{Value: "x"}, Right: &ast.Ident{Value: "y"}}},
		{input: "1 + 2 != y", expected: &ast.InfixExpression{Operator: "!=", 
		Left: &ast.InfixExpression{
			Operator: "+",
			Left: &ast.Intiger{Value: 1},
			Right: &ast.Intiger{Value: 2},
		}, 
		Right: &ast.Ident{Value: "y"}}},
		{
			input: "1 + 1 * 2",
			expected: &ast.InfixExpression{Operator: "+", Left: &ast.Intiger{Value: 1}, Right: &ast.InfixExpression{
				Operator: "*",
				Left:     &ast.Intiger{Value: 1},
				Right:    &ast.Intiger{Value: 2},
			}},
		},
		{
			input: "x + 3.14 / 4",
			expected: &ast.InfixExpression{Operator: "+", Left: &ast.Ident{Value: "x"}, Right: &ast.InfixExpression{
				Operator: "/",
				Left:     &ast.Float{Value: 3.14},
				Right:    &ast.Intiger{Value: 4},
			}},
		},
	}

	for _, tt := range tests {
		p := New(lexer.New(tt.input))
		program := p.ParseProgram()
		if len(p.Errors()) > 0 {
			printErrors(t, p.Errors())
			t.FailNow()
		}

		if len(program.Statements) != 1 {
			t.Errorf("expected 1 statement, got %d", len(program.Statements))
			return
		}
		stmnt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("expected *ast.ExpressionStatement statement, got '%T'", program.Statements[0])
			return
		}
		infix := assertExpressionNodeType[*ast.InfixExpression](t, stmnt.Expression)

		assertInfixEqual(t, tt.expected, infix)
	}

}

func printErrors(t *testing.T, errs []error) {
	for _, e := range errs {
		t.Logf("error: %v", e)
	}
}

func assertExpressionEqual(t *testing.T, expected, actual ast.Expression) {
	t.Helper()

	if expected == nil || actual == nil {
		if expected != actual {
			t.Errorf("expected expression %v, got %v", expected, actual)
		}
		return
	}

	switch exp := expected.(type) {
	case *ast.Intiger:
		act := assertExpressionNodeType[*ast.Intiger](t, actual)
		if exp.Value != act.Value {
			t.Errorf("expected int value to be %d, got %d", exp.Value, act.Value)
		}
	case *ast.String:
		act := assertExpressionNodeType[*ast.String](t, actual)
		if exp.Value != act.Value {
			t.Errorf("expected string value to be '%s', got '%s'", exp.Value, act.Value)
		}
	case *ast.Bool:
		act := assertExpressionNodeType[*ast.Bool](t, actual)
		if exp.Value != act.Value {
			t.Errorf("expected boolean value to be '%v', got '%v'", exp.Value, act.Value)
		}
	case *ast.Ident:
		act := assertExpressionNodeType[*ast.Ident](t, actual)
		if exp.Value != act.Value {
			t.Errorf("expected identifier name to be '%v', got '%v'", exp.Value, act.Value)
		}
	case *ast.Float:
		act := assertExpressionNodeType[*ast.Float](t, actual)
		if exp.Value != act.Value {
			t.Errorf("expected float value to be '%v', got '%v'", exp.Value, act.Value)
		}
	case *ast.InfixExpression:
		act := assertExpressionNodeType[*ast.InfixExpression](t, actual)
		assertInfixEqual(t, exp, act)
	case *ast.PrefixExpression:
		act := assertExpressionNodeType[*ast.PrefixExpression](t, actual)
		if act.Operator != exp.Operator {
			t.Errorf("expected prefix operator '%s', got '%s'", exp.Operator, act.Operator)
			return
		}
		assertExpressionEqual(t, exp.Expression, act.Expression)
	default:
		t.Fatalf("unhandled expression type: %T", expected)
	}
}

func assertInfixEqual(t *testing.T, expected *ast.InfixExpression, actual *ast.InfixExpression) {
	t.Helper()

	if expected.Operator != actual.Operator {
		t.Errorf("expected infix expression operator '%s', got '%s'", expected.Operator, actual.Operator)
		return
	}

	assertExpressionEqual(t, expected.Left, actual.Left)
	assertExpressionEqual(t, expected.Right, actual.Right)
}

func assertExpressionNodeType[T ast.Expression](t *testing.T, node ast.Expression) T {
	t.Helper()
	if node == nil {
		t.Fatalf("expected node type %T, got nil", (*T)(nil))
	}

	tt, ok := node.(T)
	if !ok {
		t.Errorf("expected node type %T, got %T", (*T)(nil), node)
	}

	return tt
}

type expectedVal struct {
	Type token.TokenType
	Val  any
}

func testLiteralExpression(t *testing.T, expr ast.Expression, expected *expectedVal) {
	t.Helper()
	switch expected.Type {
	case token.INT:
		val := expected.Val.(int64)
		testIntegerLiteral(t, expr, val)
	case token.FLOAT:
		val := expected.Val.(float64)
		testFloatLiteral(t, expr, val)
	case token.IDENT:
		val := expected.Val.(string)
		testIdent(t, expr, val)
	case token.TRUE, token.FALSE:
		val := expected.Val.(bool)
		testBooleanLiteral(t, expr, val)
	case token.STRING:
		val := expected.Val.(string)
		testStringLiteral(t, expr, val)
	default:
		t.Errorf("type of expr not handled. got=%T", expr)
	}
}

func testStringLiteral(t *testing.T, expr ast.Expression, value string) {
	s, ok := expr.(*ast.String)
	if !ok {
		t.Errorf("s not *ast.String. got=%T", expr)
	}
	if s.Value != value {
		t.Errorf("b.Value not %s. got=%s", value, s.Value)
	}
	if s.TokenLiteral() != fmt.Sprintf("%s", value) {
		t.Errorf("b.TokenLiteral() not %s. got=%s", value, s.TokenLiteral())
	}
}

func testBooleanLiteral(t *testing.T, expr ast.Expression, value bool) {
	b, ok := expr.(*ast.Bool)
	if !ok {
		t.Errorf("b not *ast.Boolean. got=%T", expr)
	}
	if b.Value != value {
		t.Errorf("b.Value not %t. got=%t", value, b.Value)
	}
	if b.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("b.TokenLiteral() not %t. got=%s", value, b.TokenLiteral())
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) {
	i, ok := il.(*ast.Intiger)
	if !ok {
		t.Fatalf("il not *ast.IntegerLiteral. got=%T", il)
	}

	if i.Value != value {
		t.Errorf("i.Value not %d. got=%d", value, i.Value)
		return
	}

	if i.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("i.TokenLiteral() not %d. got=%s", value, i.TokenLiteral())
	}
}

func testFloatLiteral(t *testing.T, expr ast.Expression, value float64) {
	fl, ok := expr.(*ast.Float)
	if !ok {
		t.Errorf("expr not *ast.FloatLiteral. got=%T", fl)
		return
	}

	if fl.Value != value {
		t.Errorf("fl.Value not %f. got=%f", value, fl.Value)
	}
}

func testIdent(t *testing.T, expr ast.Expression, value string) {
	ident, ok := expr.(*ast.Ident)
	if !ok {
		t.Errorf("expr not *ast.Ident. got=%T", expr)
		return
	}
	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return
	}
	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral() not %s. got=%s", value, ident.TokenLiteral())
		return
	}
}
