package evaluate

import (
	"github.com/maiyama18/dog/ast"
	"github.com/maiyama18/dog/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func booleanObject(b bool) object.Object {
	if b {
		return TRUE
	} else {
		return FALSE
	}
}

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(node.Operator, left, right)
	case *ast.IntegerLiteral:
		return object.NewInteger(node.Value)
	case *ast.BooleanLiteral:
		return booleanObject(node.Value)
	default:
		return NULL
	}
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object
	for _, s := range stmts {
		result = Eval(s)
	}
	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangExpression(right)
	case "-":
		return evalMinusExpression(right)
	default:
		return NULL
	}
}

func evalBangExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE, NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusExpression(right object.Object) object.Object {
	integer, ok := right.(*object.Integer)
	if !ok {
		return NULL
	}
	return object.NewInteger(-integer.Value)
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.IntegerType && right.Type() == object.IntegerType:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.BooleanType && right.Type() == object.BooleanType:
		return evalBooleanInfixExpression(operator, left, right)
	default:
		return NULL
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	intLeft, ok1 := left.(*object.Integer)
	intRight, ok2 := right.(*object.Integer)
	if !ok1 || !ok2 {
		return NULL
	}

	switch operator {
	case "+":
		return object.NewInteger(intLeft.Value + intRight.Value)
	case "-":
		return object.NewInteger(intLeft.Value - intRight.Value)
	case "*":
		return object.NewInteger(intLeft.Value * intRight.Value)
	case "/":
		return object.NewInteger(intLeft.Value / intRight.Value)
	case "==":
		return booleanObject(intLeft.Value == intRight.Value)
	case "!=":
		return booleanObject(intLeft.Value != intRight.Value)
	case ">":
		return booleanObject(intLeft.Value > intRight.Value)
	case "<":
		return booleanObject(intLeft.Value < intRight.Value)
	default:
		return NULL
	}
}

func evalBooleanInfixExpression(operator string, left, right object.Object) object.Object {
	boolLeft, ok1 := left.(*object.Boolean)
	boolRight, ok2 := right.(*object.Boolean)
	if !ok1 || !ok2 {
		return NULL
	}

	switch operator {
	case "==":
		return booleanObject(boolLeft == boolRight)
	case "!=":
		return booleanObject(boolLeft != boolRight)
	default:
		return NULL
	}
}
