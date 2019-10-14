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
		t.Fatalf("program statments length wrong. want=%d, got=%d", len(expectedIdentNames), len(program.Statements))
	}

	for i, s := range program.Statements {
		testLetStatement(t, s, expectedIdentNames[i])
	}
}

func testLetStatement(t *testing.T, statement ast.Statement, expectedIdentName string) {
	t.Helper()

	letStatement, ok := statement.(*ast.LetStatement)
	if !ok {
		t.Fatalf("not LetStatement: %+v", statement)
	}

	if letStatement.Identifier.Name != expectedIdentName {
		t.Fatalf("identifier name wrong. want=%q, got=%q", expectedIdentName, letStatement.Identifier.Name)
	}
}

func TestReturnStatements(t *testing.T) {
	input := `
return 5;
return 42;
`

	program := parseProgram(t, input)

	if len(program.Statements) != 2 {
		t.Fatalf("program statments length wrong. want=%d, got=%d", 2, len(program.Statements))
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
	input := `
foo;
bar;
`

	program := parseProgram(t, input)

	expectedIdentNames := []string{"foo", "bar"}

	if len(program.Statements) != len(expectedIdentNames) {
		t.Fatalf("program statments length wrong. want=%d, got=%d", len(expectedIdentNames), len(program.Statements))
	}

	for i, s := range program.Statements {
		expStatement, ok := s.(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("not ExpressionStatement: %+v", s)
		}
		ident, ok := expStatement.Expression.(*ast.Identifier)
		if !ok {
			t.Fatalf("not Identifier: %+v", expStatement.Expression)
		}
		if ident.Name != expectedIdentNames[i] {
			t.Fatalf("identifier name wrong. want=%q, got=%q", expectedIdentNames[i], ident.Name)
		}
	}
}

func TestIntegerLiterals(t *testing.T) {
	input := `
42;
3;
`

	program := parseProgram(t, input)

	expectedInts := []int64{42, 3}

	if len(program.Statements) != len(expectedInts) {
		t.Fatalf("program statments length wrong. want=%d, got=%d", len(expectedInts), len(program.Statements))
	}

	for i, s := range program.Statements {
		expStatement, ok := s.(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("not ExpressionStatement: %+v", s)
		}
		il, ok := expStatement.Expression.(*ast.IntegerLiteral)
		if !ok {
			t.Fatalf("not IntegerLiteral: %+v", expStatement.Expression)
		}
		if il.Value != expectedInts[i] {
			t.Fatalf("integer value wrong. want=%q, got=%q", expectedInts[i], il.Value)
		}
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
