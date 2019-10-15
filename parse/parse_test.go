package parse

import (
	"strings"
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
			testPrefixExpression(t, expStmt.Expression, test.want.operator, test.want.value)
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
			testInfixExpression(t, expStmt.Expression, test.want.operator, test.want.left, test.want.right)
		})
	}
}

func TestIfExpressions(t *testing.T) {
	input := `if (x > y) { x }; if (true) { x } else { -x };`

	program := parseProgram(t, input)

	if len(program.Statements) != 2 {
		t.Fatalf("program statements length wrong. want=%d, got=%d", 2, len(program.Statements))
	}

	expStmt1, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("not ExpressionStatement: %+v", program.Statements[0])
	}
	ifExp1, ok := expStmt1.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("not IfExpression: %+v", expStmt1.Expression)
	}
	testInfixExpression(t, ifExp1.Condition, ">", "x", "y")
	if len(ifExp1.Consequence.Statements) != 1 {
		t.Fatalf("consequence statements length wrong. want=%d, got=%d", 1, len(ifExp1.Consequence.Statements))
	}
	consqExpStmt1, ok := ifExp1.Consequence.Statements[0].(*ast.ExpressionStatement)
	testLiteralExpression(t, consqExpStmt1.Expression, "x")
	if ifExp1.Alternative != nil {
		t.Fatalf("alternative not nil")
	}

	expStmt2, ok := program.Statements[1].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("not ExpressionStatement: %+v", program.Statements[1])
	}
	ifExp2, ok := expStmt2.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("not IfExpression: %+v", expStmt1.Expression)
	}
	testLiteralExpression(t, ifExp2.Condition, true)
	if len(ifExp1.Consequence.Statements) != 1 {
		t.Fatalf("consequence statements length wrong. want=%d, got=%d", 1, len(ifExp1.Consequence.Statements))
	}
	consqExpStmt2, ok := ifExp2.Consequence.Statements[0].(*ast.ExpressionStatement)
	testLiteralExpression(t, consqExpStmt2.Expression, "x")
	if ifExp2.Alternative == nil {
		t.Fatalf("alternative nil")
	}
	alterExpStmt, ok := ifExp2.Alternative.Statements[0].(*ast.ExpressionStatement)
	testPrefixExpression(t, alterExpStmt.Expression, "-", "x")
}

func TestFunctionLiterals(t *testing.T) {
	input := `fn (x, y) { x + y }`

	program := parseProgram(t, input)

	if len(program.Statements) != 1 {
		t.Fatalf("program statements length wrong. want=%d, got=%d", 1, len(program.Statements))
	}

	expStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("not ExpressionStatement: %+v", program.Statements[0])
	}
	funcLiteral, ok := expStmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("not FunctionLiteral: %+v", expStmt.Expression)
	}

	params := funcLiteral.Parameters
	if len(params) != 2 {
		t.Fatalf("parameters size wrong. want=%d, got=%d", 2, len(params))
	}
	testLiteralExpression(t, &params[0], "x")
	testLiteralExpression(t, &params[1], "y")

	stmts := funcLiteral.Body.Statements
	if len(stmts) != 1 {
		t.Fatalf("function literal body statments length wrong. want=%d, got=%d", 1, len(stmts))
	}
	bodyExpStmt, ok := stmts[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("not ExpressionStatement: %+v", stmts[0])
	}
	testInfixExpression(t, bodyExpStmt.Expression, "+", "x", "y")
}

func TestFunctionParameters(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{
			input: `fn () {}`,
			want:  []string{},
		},
		{
			input: `fn (x) {}`,
			want:  []string{"x"},
		},
		{
			input: `fn (x, y) {}`,
			want:  []string{"x", "y"},
		},
		{
			input: `fn (x, y, z) {}`,
			want:  []string{"x", "y", "z"},
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
			funcLiteral, ok := expStmt.Expression.(*ast.FunctionLiteral)
			if !ok {
				t.Fatalf("not FunctionLiteral: %+v", expStmt.Expression)
			}

			if len(funcLiteral.Parameters) != len(test.want) {
				t.Fatalf("parameters size wrong. want=%d, got=%d", len(test.want), len(funcLiteral.Parameters))
			}
			for i, w := range test.want {
				testLiteralExpression(t, &funcLiteral.Parameters[i], w)
			}
		})
	}
}

func TestCallExpression(t *testing.T) {
	input := `add(1, x * y);`

	program := parseProgram(t, input)

	if len(program.Statements) != 1 {
		t.Fatalf("program statements length wrong. want=%d, got=%d", 1, len(program.Statements))
	}

	expStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("not ExpressionStatement: %+v", program.Statements[0])
	}
	callExp, ok := expStmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("not CallExpression: %+v", expStmt.Expression)
	}

	testLiteralExpression(t, callExp.Function, "add")

	if len(callExp.Arguments) != 2 {
		t.Fatalf("function literal body statments length wrong. want=%d, got=%d", 2, len(callExp.Arguments))
	}
	testLiteralExpression(t, callExp.Arguments[0], 1)
	testInfixExpression(t, callExp.Arguments[1], "*", "x", "y")
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
		{
			input: "3 + (5 + 7);",
			want:  "(3 + (5 + 7));",
		},
		{
			input: "3 * (5 + 7);",
			want:  "(3 * (5 + 7));",
		},
		{
			input: "!(true == true);",
			want:  "(!(true == true));",
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
	t.Helper()

	lexer := lex.NewLexer(input)
	parser := NewParser(lexer)

	program := parser.ParseProgram()

	if len(parser.Errors()) > 0 {
		var errMsgs []string
		for _, err := range parser.Errors() {
			errMsgs = append(errMsgs, err.Error())
		}
		t.Fatalf("got parser errors: \n%s", strings.Join(errMsgs, "\n"))
	}
	if program == nil {
		t.Fatalf("program is nil")
	}

	return program
}

func testPrefixExpression(t *testing.T, exp ast.Expression, wantedOperator string, wantedRight interface{}) {
	prefixExp, ok := exp.(*ast.PrefixExpression)
	if !ok {
		t.Fatalf("not PrefixExp: %+v", exp)
	}
	if prefixExp.Operator != wantedOperator {
		t.Fatalf("operator wrong. want=%q, got=%q", wantedOperator, prefixExp.Operator)
	}
	testLiteralExpression(t, prefixExp.Right, wantedRight)
}

func testInfixExpression(t *testing.T, exp ast.Expression, wantedOperator string, wantedLeft, wantedRight interface{}) {
	infixExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("not PrefixExp: %+v", exp)
	}
	if infixExp.Operator != wantedOperator {
		t.Fatalf("operator wrong. want=%q, got=%q", wantedOperator, infixExp.Operator)
	}
	testLiteralExpression(t, infixExp.Left, wantedLeft)
	testLiteralExpression(t, infixExp.Right, wantedRight)
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
