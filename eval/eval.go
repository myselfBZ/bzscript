package eval

import (
	"fmt"

	"github.com/myselfBZ/bzscript/ast"
	"github.com/myselfBZ/bzscript/object"
)

var(
	False = &object.Bool{Value: false}
	True = &object.Bool{Value: true}
)


func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Intiger:
		return &object.Intiger{Value: node.Value}
	case *ast.Float:
		return &object.Float{Value: node.Value}
	case *ast.String:
		return &object.String{Value: node.Value}
	case *ast.Bool:
		if node.Value {
			return True
		} else {
			return False
		}
	case *ast.InfixExpression:
		left := Eval(node, env)
		if isError(left) {
			return left
		}
		right := Eval(node, env)
		if isError(right) {
			return right
		}
		return evalInfix(left, right, node.Operator)
	default:
		return newError("unrecognized AST node: %T", node)
	}
}

func evalInfix(left object.Object, right object.Object, oprt string) object.Object {
	if left.Type() != right.Type() {
		return newError("cannot perform '%s' operator for %s and %s, mismatched types", oprt, left.Type(), right.Type())
	}
	switch left.Type() {
	case object.STRING:
		if oprt != "+" {
			return newError("cannot perform '%s' operator for strings", oprt)
		}
		leftStr := left.(*object.String).Value
		rightStr := right.(*object.String).Value
		return &object.String{Value: leftStr + rightStr}
	case object.FLOAT:
		return evalFloatInfix(left, right, oprt)


	// TODO 
	case object.INTIGER:
		return nil
	case object.BOOL:
		return nil
	default:
		return nil
	}
}

func evalFloatInfix(left, right object.Object, oprt string) object.Object {
	leftVal := left.(*object.Float).Value
	rightVal := right.(*object.Float).Value

	switch oprt {
	case "+":
		return &object.Float{Value: leftVal + rightVal}
	case "-":
		return &object.Float{Value: leftVal - rightVal}
	case "/":
		return &object.Float{Value: leftVal / rightVal}
	case "*":
		return &object.Float{Value: leftVal * rightVal}
	case "==":
		if leftVal == rightVal {
			return True
		}
		return False
	case "!=":
		if leftVal != rightVal {
			return True
		}
		return False
	case "<=":
		if leftVal <= rightVal {
			return True
		}
		return False
	case ">=":
		if leftVal >= rightVal {
			return True
		}
		return False
	case ">":
		if leftVal > rightVal {
			return True
		}
		return False
	case "<":
		if leftVal < rightVal {
			return True
		}
		return False
	default:
		return newError("invalid infix operator '%s'", oprt)
	}
}


func isError(obj object.Object) bool {
	_, ok := obj.(*object.Error)
	return ok
}
func newError(msg string, args...any) *object.Error {
	return &object.Error{Message: fmt.Sprintf(msg, args...)}
}
