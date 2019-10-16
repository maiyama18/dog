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
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
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
