package eval

import (
	"testing"

	"github.com/myselfBZ/bzscript/ast"
	"github.com/myselfBZ/bzscript/lexer"
	"github.com/myselfBZ/bzscript/object"
	"github.com/myselfBZ/bzscript/parser"
)

func TestPrefix(t *testing.T) {
	tests := []struct{
		input string
		expected any
	}{
		{input: "-100", expected: -100},
		{input: "-(-100)", expected: 100},
		{input: "-3.14", expected: -3.14},
		{input: "-(4.5) + (-(-4.5))", expected: 0.0},

		{input: "!true", expected: false},
		{input: "!!true", expected: true},
		{input: "!(!(!false))", expected: true},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		env := object.NewEnvironment()
		expression := program.Statements[0].(*ast.ExpressionStatement).Expression
		obj := Eval(expression, env)
		testObject(t, obj, tt.expected)
	}
}

func TestInfix(t *testing.T) {
	tests := []struct{
		input string
		expected any
	}{
		{input: "1+1", expected: int64(2)},
		{input: "1+2*4", expected: int64(9)},
		{input: "2.3 + 1.1", expected: 3.4},
		{input: "33 + 55", expected: 88},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		env := object.NewEnvironment()
		expression := program.Statements[0].(*ast.ExpressionStatement).Expression
		obj := Eval(expression, env)
		testObject(t, obj, tt.expected)
	}
}

func TestEvalTypes(t *testing.T) {
	tests := []struct{
		input string
		expected any
	}{
		{input: "1", expected: int64(1)},
		{input: "3.14", expected: float64(3.14)},
		{input: "true", expected: true},
		{input: "false", expected: false},
		{input: "\"hello, world\"", expected: "hello, world"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		env := object.NewEnvironment()
		expression := program.Statements[0].(*ast.ExpressionStatement).Expression
		obj := Eval(expression, env)
		testObject(t, obj, tt.expected)
	}
}

func testObject(t *testing.T, obj object.Object, expected any) {
	switch expectedType := expected.(type) {
	case int64:
		val, ok := obj.(*object.Intiger)
		if !ok {
			t.Errorf("expected Intiger, got %T", obj)
			return
		}
		if val.Value != expectedType {
			t.Errorf("expected %d, got %d", val.Value, expectedType)
			return
		}
	case string:
		val, ok := obj.(*object.String)
		if !ok {
			t.Errorf("expected String, got %T", obj)
			return
		}
		if val.Value != expectedType {
			t.Errorf("expected %s, got %s", val.Value, expectedType)
			return
		}
	case bool:
		val, ok := obj.(*object.Bool)
		if !ok {
			t.Errorf("expected Bool, got %T", obj)
			return
		}
		if val.Value != expectedType {
			t.Errorf("expected %v, got %v", val.Value, expectedType)
			return
		}
	case float64:
		val, ok := obj.(*object.Float)
		if !ok {
			t.Errorf("expected Float, got %T", obj)
			return
		}
		if val.Value != expectedType {
			t.Errorf("expected %v, got %v", val.Value, expectedType)
			return
		}
	}
}

