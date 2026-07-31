package eval

import (
	"fmt"

	"github.com/myselfBZ/bzscript/ast"
	"github.com/myselfBZ/bzscript/object"
)

var(
	False = &object.Bool{Value: false}
	True = &object.Bool{Value: true}
	
	// Temp 
	Nothing = &object.String{Value: "nothing"}
)


func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
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
	case *ast.Ident:
		return evalIdent(env, node.Value)
	case *ast.VarStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Ident.Value, val)
		return Nothing
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.Block:
		enclosedEnv := object.NewEnclosedEnvironment(env)
		return evalBlockStatement(node, enclosedEnv, )
	case *ast.IfStatement:
		condition := Eval(node.Condition, env)
		if isError(condition) {
			return condition
		}
		boolVal, ok := condition.(*object.Bool)
		if !ok {
			return newError("non-boolean expression: %T", condition)
		}
		var result object.Object
		enclosedEnv := object.NewEnclosedEnvironment(env)
		if boolVal.Value {
			result = evalBlockStatement(node.Consequence, enclosedEnv)
			return result
		} else if node.Alternative != nil {
			result = evalBlockStatement(node.Alternative, enclosedEnv)
			return result
		}
		return Nothing
	case *ast.ReturnStatement:
		obj := &object.ReturnValue{}
		obj.Value = Eval(node.Value, env)
		if isError(obj.Value) {
			return obj.Value
		}
		return obj
	case *ast.AssignStatement:
		switch lhs := node.LHS.(type) {
		case *ast.Ident:
			_, scope, ok := env.Get(lhs.Value)
			if !ok {
				return newError("identifier not found %s", lhs.Value)
			}
			rhs := Eval(node.RHS, env)
			if isError(rhs) {
				return rhs 
			}
			scope.Set(lhs.Value, rhs)
		case *ast.StructMemberAccess:
			left := Eval(lhs.Lhs, env)
			if isError(left) {
				return left
			}
			structInstance, ok := left.(*object.StructInstance)
			if !ok {
				return newError("field access on a non-struct object: %s", left.Type())
			}
			rhs := lhs.Rhs.(*ast.Ident)

			val := Eval(node.RHS, env)
			if isError(val) {
				return val
			}
			if _, ok := structInstance.Fields[rhs.Value]; !ok {
				return newError("field not found on a struct: %s", rhs.Value)
			}
			structInstance.Fields[rhs.Value] = val
		default:
			return newError("cannot assign to %s, not addressable", node.LHS.TokenLiteral())

		}
		return Nothing
	case *ast.StructLiteral:
		obj := &object.Struct{Fields: make(map[string]bool)}
		for _, f := range node.Fields {
			obj.Fields[f.Name] = true
		}
		env.Set(node.Name.Value, obj)
		return Nothing
	case *ast.StructMemberAccess:
		left := Eval(node.Lhs, env)
		if isError(left){
			return left
		}
		s, ok := left.(*object.StructInstance)
		if !ok {
			return newError("field access on a non-struct object: %s", left.Type())
		}
		fieldName := node.Rhs.(*ast.Ident)
		val, ok := s.Fields[fieldName.Value]
		if !ok {
			return newError("field not found on a struct: %s", fieldName)
		}
		return val
	case *ast.WhileLoop:
		for {
			enclosedEnv := object.NewEnclosedEnvironment(env)
			condition := Eval(node.Condition, env)
			if isError(condition) {
				return condition
			}
			boolVal, ok := condition.(*object.Bool)
			if !ok {
				return newError("non-boolean expression: %T", condition)
			}
			if !boolVal.Value {
				return Nothing
			}
			obj := evalWhileLoopBody(node.Body, enclosedEnv)
			switch obj.(type) {
			case *object.Error, *object.ReturnValue:
				return obj
			case *object.Break:
				return Nothing
			}
		}
	case *ast.AnonymousFuncLiteral:
		obj := &object.Function{Body: node.Body, Params: node.Params, Env: env}
		return obj
	case *ast.FunctionLiteral:
		obj := &object.Function{Body: node.Body, Params: node.Params, Env: env}
		env.Set(node.Ident.Value, obj)
		return obj
	case *ast.FunctionCall:
		f := Eval(node.Function, env)
		if isError(f) {
			return f
		}
		args := evalExpressions(env, node.Args)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(f, args)
	case *ast.PrefixExpression:
		obj := Eval(node.Expression, env)
		if isError(obj) {
			return obj
		}
		return evalPrefix(obj, node.Operator)
	case *ast.ArrayLiteral:
		elements := evalExpressions(env, node.Elements)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	case *ast.IndexOperator:
		obj := Eval(node.Left, env)
		if isError(obj) {
			return obj
		}
		idx := Eval(node.Index, env)
		if isError(idx) {
			return idx
		}
		return evalArrayIdx(obj, idx)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfix(left, right, node.Operator)
	case *ast.BreakStatement:
		return &object.Break{}
	case *ast.ContinueStatement:
		return &object.Continue{}
	default:
		return newError("invalid statement: %s", node.TokenLiteral())
	}
}

func evalProgram(node *ast.Program, env *object.Environment) object.Object {
	var result object.Object
	for _, s := range node.Statements {
		result = Eval(s, env)
		switch obj := result.(type) {
		case *object.Error:
			return obj
		case *object.ReturnValue:
			return obj.Value
		case *object.Break: 
			return newError("break outside of while loop")
		case *object.Continue:
			return newError("continue outside of while loop")
		}
	}
	return result
}

func evalBlockStatement(block *ast.Block, env *object.Environment) object.Object {
	var result object.Object 
	for _, s := range block.Statements {
		obj := Eval(s, env)
		if isError(obj) || 
		obj.Type() == object.RETURN || 
		obj.Type() == object.CONTINUE || 
		obj.Type() == object.BREAK {
			return obj
		}
		result = obj
	}
	return result
}

func evalWhileLoopBody(block *ast.Block, env *object.Environment) object.Object {
	var result object.Object 
	for _, s := range block.Statements {
		obj := Eval(s, env)
		switch obj.(type) {
		case *object.ReturnValue, *object.Error, *object.Break, *object.Continue:
			return obj
		}
		result = obj
	}
	return result
}

func evalArrayIdx(array object.Object, idx object.Object) object.Object {
	validArr, ok := array.(*object.Array)
	if !ok {
		return newError("not an array: %s", array.Type())
	}
	validIdx, ok := idx.(*object.Intiger) 
	if !ok {
		return newError("invalid index value: %s", idx.Type())
	}
	if validIdx.Value < 0 || validIdx.Value > int64(len(validArr.Elements)) {
		return newError("index %d out of bounds array length: %d", validIdx.Value, len(validArr.Elements))
	}
	return validArr.Elements[validIdx.Value]
}

func evalIdent(env *object.Environment, name string) object.Object {
	if f, ok := builtIns[name]; ok {
		return f
	}
	obj, _, ok := env.Get(name)
	if !ok {
		return newError("identifier not found %s", name)
	}
	return obj
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
	case object.INTIGER:
		return evalIntigerInfix(left, right, oprt)
	case object.BOOL:
		return evalBoolInfix(left, right, oprt)
	}
	return newError("invalid type inside an infix expression: %s", left.Type())	
}

func evalBoolInfix(left, right object.Object, oprt string) object.Object {
	leftVal := left.(*object.Bool).Value
	rightVal := right.(*object.Bool).Value

	switch oprt {
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
	default:
		return newError("cannot perform '%s' operator on booleans", oprt)
	}
}

func evalIntigerInfix(left, right object.Object, oprt string) object.Object {
	leftVal := left.(*object.Intiger).Value
	rightVal := right.(*object.Intiger).Value

	switch oprt {
	case "+":
		return &object.Intiger{Value: leftVal + rightVal}
	case "-":
		return &object.Intiger{Value: leftVal - rightVal}
	case "/":
		return &object.Intiger{Value: leftVal / rightVal}
	case "*":
		return &object.Intiger{Value: leftVal * rightVal}
	case "%":
		return &object.Intiger{Value:leftVal % rightVal}
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

func evalPrefix(obj object.Object, oprtr string) object.Object {
	switch oprtr {
	case "!":
		return evalBang(obj)
	case "-":
		return evalMinus(obj)
	default:
		return newError("invalid prefix operator: %s", oprtr)
	}
}

func evalMinus(obj object.Object) object.Object {
	switch o := obj.(type) {
	case *object.Intiger:
		return &object.Intiger{Value: -o.Value}
	case *object.Float:
		return &object.Float{Value: -o.Value}
	default:
		return newError("cannot perform '-' operator on %T", obj)
	}
}

func evalBang(obj object.Object) object.Object {
	switch obj {
	case True:
		return False
	case False:
		return True
	default:
		return newError("cannot perform '!' operator on %T", obj)
	}
}

func unwrapReturnVal(v object.Object) object.Object {
	if r, ok := v.(*object.ReturnValue); ok {
		return r.Value
	}
	return v
}

func extendFunctionEnv( fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)
	for paramIdx, param := range fn.Params {
		env.Set(param.Value, args[paramIdx])
	}
	return env
}

func evalExpressions(env *object.Environment, exprs []ast.Expression) []object.Object {
	result := make([]object.Object, len(exprs))
	for i, e := range exprs {
		r := Eval(e, env)
		if isError(r) {
			return []object.Object{r}
		}
		result[i] = r
	}
	return result
}

func applyFunction(f object.Object, args []object.Object) object.Object {
	switch function := f.(type) {
	case *object.Function:
		if len(function.Params) != len(args) {
			return newError("function call missing arguments")
		}
		extendedEnv := extendFunctionEnv(function, args)
		evaluated := Eval(function.Body, extendedEnv)
		return unwrapReturnVal(evaluated)
	case *object.BuiltIn:
		return function.Fn(args...)
	default:
		return newError("not a function: %s", f.Type())
	}
}

func isError(obj object.Object) bool {
	_, ok := obj.(*object.Error)
	return ok
}
func newError(msg string, args...any) *object.Error {
	return &object.Error{Message: fmt.Sprintf(msg, args...)}
}
