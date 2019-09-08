package evaluator

import (
	"fmt"

	"github.com/magicly/monkey-go/ast"
	"github.com/magicly/monkey-go/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node, env *object.Env) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.FunctionLiteral:
		return &object.Function{
			Parameters: node.Parameters,
			Body:       node.Body,
			Env:        env,
		}
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args)
	}
	return NULL
}

func applyFunction(function object.Object, args []object.Object) object.Object {
	switch fn := function.(type) {
	case *object.Function:
		env := bindArgs(args, fn)
		return evalStatements(fn.Body.Statements, env)
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %s", function.Type())
	}
}

func bindArgs(args []object.Object, fn *object.Function) *object.Env {
	env := object.NewEnclosedEnv(fn.Env)

	for idx, arg := range args {
		env.Set(fn.Parameters[idx].Value, arg)
	}

	return env
}

func evalExpressions(exps []ast.Expression, env *object.Env) []object.Object {
	var result []object.Object

	for _, e := range exps {
		val := Eval(e, env)
		if isError(val) {
			return []object.Object{val}
		}
		result = append(result, val)
	}

	return result
}

func evalIdentifier(node *ast.Identifier, env *object.Env) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: " + node.Value)
}

func isError(obj object.Object) bool {
	return obj != nil && obj.Type() == object.ERROR_OBJ
}

func evalIfExpression(ie *ast.IfExpression, env *object.Env) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if truthy(condition) {
		return evalStatements(ie.Consequence.Statements, env)
	} else if ie.Alternative != nil {
		return evalStatements(ie.Alternative.Statements, env)
	}
	return NULL
}

func truthy(obj object.Object) bool {
	switch obj.Type() {
	case object.BOOLEAN_OBJ:
		return obj.(*object.Boolean).Value
	case object.INTERGER_OBJ:
		v := obj.(*object.Integer).Value
		return v != 0
	case object.NULL_OBJ:
		return false
	default:
		return true
	}

}

func evalInfixExpression(op string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTERGER_OBJ && right.Type() == object.INTERGER_OBJ:
		return evalIntegerInfixExpression(op, left, right)
	case left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ:
		return evalBooleanInfixExpression(op, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(op, left, right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), op, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), op, right.Type())
	}
}
func evalStringInfixExpression(op string, left object.Object, right object.Object) object.Object {
	lv := left.(*object.String).Value
	rv := right.(*object.String).Value
	switch op {
	case "+":
		return &object.String{Value: lv + rv}
	default:
		return newError("unknown operator: %s %s %s", left.Type(), op, right.Type())
	}
}
func evalBooleanInfixExpression(op string, left object.Object, right object.Object) object.Object {
	lv := left.(*object.Boolean).Value
	rv := right.(*object.Boolean).Value
	switch op {
	case "==":
		return nativeBoolToBooleanObject(lv == rv)
	case "!=":
		return nativeBoolToBooleanObject(lv != rv)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), op, right.Type())
	}
}

func evalIntegerInfixExpression(op string, left object.Object, right object.Object) object.Object {
	lv := left.(*object.Integer).Value
	rv := right.(*object.Integer).Value
	switch op {
	case "+":
		return &object.Integer{Value: lv + rv}
	case "-":
		return &object.Integer{Value: lv - rv}
	case "*":
		return &object.Integer{Value: lv * rv}
	case "/":
		return &object.Integer{Value: lv / rv}
	case "<":
		return nativeBoolToBooleanObject(lv < rv)
	case ">":
		return nativeBoolToBooleanObject(lv > rv)
	case "==":
		return nativeBoolToBooleanObject(lv == rv)
	case "!=":
		return nativeBoolToBooleanObject(lv != rv)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), op, right.Type())
	}
}

func evalPrefixExpression(op string, right object.Object) object.Object {
	switch op {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusOperatorExpression(right)
	default:
		return newError("unknown operator: %s %s", op, right.Type())
	}
}

func evalMinusOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTERGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE

	default:
		switch right.(type) {
		case *object.Integer:
			r := right.(*object.Integer)
			if r.Value != 0 {
				return nativeBoolToBooleanObject(false)
			}
			return nativeBoolToBooleanObject(true)
		case *object.Boolean:
			r := right.(*object.Boolean)
			return nativeBoolToBooleanObject(!r.Value)
		default:
			return FALSE
		}
		// default:
		// 	return FALSE
	}
}

func nativeBoolToBooleanObject(input bool) object.Object {
	if input {
		return TRUE
	}
	return FALSE
}

func evalProgram(p *ast.Program, env *object.Env) object.Object {
	var result object.Object

	for _, stmt := range p.Statements {
		result = Eval(stmt, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalStatements(stmts []ast.Statement, env *object.Env) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt, env)

		if result != nil && (result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ) {
			return result
		}
	}

	return result
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}
