package parse

import (
	"testing"

	"github.com/maiyama18/dog/ast"

	"github.com/maiyama18/dog/lex"
)

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let foo = 42;
`

	program := parseProgram(t, input)

	expectedIdentNames := []string{"x", "foo"}

	if len(program.Statements) != len(expectedIdentNames) {
		t.Fatalf("program statements length wrong. want=%d, got=%d", len(expectedIdentNames), len(program.Statements))
	}

	for i, s := range program.Statements {
		testLetStatement(t, s, expectedIdentNames[i])
	}
}

func testLetStatement(t *testing.T, statement ast.Statement, expectedIdentName string) {
	t.Helper()

	letStmt, ok := statement.(*ast.LetStatement)
	if !ok {
		t.Fatalf("not LetStatement: %+v", statement)
	}

	if letStmt.Identifier.Name != expectedIdentName {
		t.Fatalf("identifier name wrong. want=%q, got=%q", expectedIdentName, letStmt.Identifier.Name)
	}
}

func TestReturnStatements(t *testing.T) {
	input := `
return 5;
return 42;
`

	program := parseProgram(t, input)

	if len(program.Statements) != 2 {
		t.Fatalf("program statements length wrong. want=%d, got=%d", 2, len(program.Statements))
	}

	for _, s := range program.Statements {
		testReturnStatement(t, s)
	}
}

func testReturnStatement(t *testing.T, statement ast.Statement) {
	t.Helper()

	_, ok := statement.(*ast.ReturnStatement)
	if !ok {
		t.Fatalf("not ReturnStatement: %+v", statement)
	}
}

func TestIdentifiers(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			input: "a;",
			want:  "a",
		},
		{
			input: "foobar;",
			want:  "foobar",
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			program := parseProgram(t, test.input)

			if len(program.Statements) != 1 {
				t.Fatalf("program statements length wrong. want=%d, got=%d", 1, len(program.Statements))
			}

			expStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("not ExpressionStatement: %+v", expStmt)
			}
			testIdentifier(t, expStmt.Expression, test.want)
		})
	}
}

func TestIntegerLiterals(t *testing.T) {
	tests := []struct {
		input string
		want  int64
	}{
		{
			input: "3;",
			want:  3,
		},
		{
			input: "42;",
			want:  42,
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			program := parseProgram(t, test.input)

			if len(program.Statements) != 1 {
				t.Fatalf("program statements length wrong. want=%d, got=%d", 1, len(program.Statements))
			}

			expStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("not ExpressionStatement: %+v", expStmt)
			}
			testIntegerLiteral(t, expStmt.Expression, test.want)
		})
	}
}

func TestBooleanLiterals(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{
			input: "true;",
			want:  true,
		},
		{
			input: "false;",
			want:  false,
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			program := parseProgram(t, test.input)

			if len(program.Statements) != 1 {
				t.Fatalf("program statements length wrong. want=%d, got=%d", 1, len(program.Statements))
			}

			expStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("not ExpressionStatement: %+v", expStmt)
			}
			testBooleanLiteral(t, expStmt.Expression, test.want)
		})
	}
}

func TestPrefixExpressions(t *testing.T) {
	type want struct {
		operator string
		value    interface{}
	}
	tests := []struct {
		input string
		want  want
	}{
		{
			input: "!5;",
			want:  want{operator: "!", value: 5},
		},
		{
			input: "-42;",
			want:  want{operator: "-", value: 42},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			program := parseProgram(t, test.input)

			if len(program.Statements) != 1 {
				t.Fatalf("program statements length wrong. want=%d, got=%d", 1, len(program.Statements))
			}

			expStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("not ExpressionStatement: %+v", program.Statements[0])
			}
			prefixExp, ok := expStmt.Expression.(*ast.PrefixExpression)
			if !ok {
				t.Fatalf("not PrefixExp: %+v", expStmt.Expression)
			}
			if prefixExp.Operator != test.want.operator {
				t.Fatalf("operator wrong. want=%q, got=%q", test.want.operator, prefixExp.Operator)
			}
			testLiteralExpression(t, prefixExp.Right, test.want.value)
		})
	}
}

func TestInfixExpressions(t *testing.T) {
	type want struct {
		operator string
		left     interface{}
		right    interface{}
	}
	tests := []struct {
		input string
		want  want
	}{
		{
			input: "5 + 6;",
			want:  want{operator: "+", left: 5, right: 6},
		},
		{
			input: "5 - 6;",
			want:  want{operator: "-", left: 5, right: 6},
		},
		{
			input: "5 * 6;",
			want:  want{operator: "*", left: 5, right: 6},
		},
		{
			input: "5 / 6;",
			want:  want{operator: "/", left: 5, right: 6},
		},
		{
			input: "5 == 6;",
			want:  want{operator: "==", left: 5, right: 6},
		},
		{
			input: "5 != 6;",
			want:  want{operator: "!=", left: 5, right: 6},
		},
		{
			input: "5 < 6;",
			want:  want{operator: "<", left: 5, right: 6},
		},
		{
			input: "5 > 6;",
			want:  want{operator: ">", left: 5, right: 6},
		},
		{
			input: "alice == bob",
			want:  want{operator: "==", left: "alice", right: "bob"},
		},
		{
			input: "alice != bob",
			want:  want{operator: "!=", left: "alice", right: "bob"},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			program := parseProgram(t, test.input)

			if len(program.Statements) != 1 {
				t.Fatalf("program statements length wrong. want=%d, got=%d", 1, len(program.Statements))
			}

			expStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("not ExpressionStatement: %+v", program.Statements[0])
			}
			infixExp, ok := expStmt.Expression.(*ast.InfixExpression)
			if !ok {
				t.Fatalf("not PrefixExp: %+v", expStmt.Expression)
			}
			if infixExp.Operator != test.want.operator {
				t.Fatalf("operator wrong. want=%q, got=%q", test.want.operator, infixExp.Operator)
			}
			testLiteralExpression(t, infixExp.Left, test.want.left)
			testLiteralExpression(t, infixExp.Right, test.want.right)
		})
	}
}

func TestOperatorPrecedences(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			input: "!-5;",
			want:  "(!(-5));",
		},
		{
			input: "3 + -5;",
			want:  "(3 + (-5));",
		},
		{
			input: "3 + 5 + 7;",
			want:  "((3 + 5) + 7);",
		},
		{
			input: "3 + 5 * 7;",
			want:  "(3 + (5 * 7));",
		},
		{
			input: "3 + 5 * 7 + 9;",
			want:  "((3 + (5 * 7)) + 9);",
		},
		{
			input: "3 + 5 > 7;",
			want:  "((3 + 5) > 7);",
		},
		{
			input: "3 + 5 != 7 * 9;",
			want:  "((3 + 5) != (7 * 9));",
		},
		{
			input: "true == !false;",
			want:  "(true == (!false));",
		},
		{
			input: "true == false != false;",
			want:  "((true == false) != false);",
		},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			program := parseProgram(t, test.input)

			if program.String() != test.want {
				t.Fatalf("program string wrong. want=%q, got=%q", test.want, program.String())
			}
		})
	}
}

func parseProgram(t *testing.T, input string) *ast.Program {
	lexer := lex.NewLexer(input)
	parser := NewParser(lexer)

	program := parser.ParseProgram()

	if len(parser.Errors()) > 0 {
		t.Fatalf("got parser errors: %+v", parser.Errors())
	}
	if program == nil {
		t.Fatalf("program is nil")
	}

	return program
}

func testLiteralExpression(t *testing.T, exp ast.Expression, want interface{}) {
	switch want := want.(type) {
	case int:
		testIntegerLiteral(t, exp, int64(want))
	case int64:
		testIntegerLiteral(t, exp, want)
	case bool:
		testBooleanLiteral(t, exp, want)
	case string:
		testIdentifier(t, exp, want)
	default:
		t.Fatalf("literal %+v not handled", exp)
	}
}

func testIntegerLiteral(t *testing.T, exp ast.Expression, want int64) {
	t.Helper()

	intLiteral, ok := exp.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("not IntegerLiteral: %+v", exp)
	}
	if intLiteral.Value != want {
		t.Fatalf("integer value wrong. want=%q, got=%q", want, intLiteral.Value)
	}
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, want bool) {
	t.Helper()

	boolLiteral, ok := exp.(*ast.BooleanLiteral)
	if !ok {
		t.Fatalf("not BooleanLiteral: %+v", exp)
	}
	if boolLiteral.Value != want {
		t.Fatalf("integer value wrong. want=%t, got=%t", want, boolLiteral.Value)
	}
}

func testIdentifier(t *testing.T, exp ast.Expression, want string) {
	t.Helper()

	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Fatalf("not Identifier: %+v", exp)
	}
	if ident.Name != want {
		t.Fatalf("identifier name wrong. want=%q, got=%q", want, ident.Name)
	}
}
