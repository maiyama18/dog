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

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	case *ast.IntegerLiteral:
		return object.NewInteger(node.Value)
	case *ast.BooleanLiteral:
		if node.Value {
			return TRUE
		} else {
			return FALSE
		}
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
