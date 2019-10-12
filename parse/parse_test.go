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

	lexer := lex.NewLexer(input)
	parser := NewParser(lexer)

	program := parser.ParseProgram()

	if len(parser.Errors()) > 0 {
		t.Fatalf("got parser errors: %+v", parser.Errors())
	}
	if program == nil {
		t.Fatalf("program is nil")
	}
	if len(program.Statements) != 2 {
		t.Fatalf("program statments length wrong. want=%d, got=%d", 2, len(program.Statements))
	}

	expectedIdentNames := []string{"x", "foo"}

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

	lexer := lex.NewLexer(input)
	parser := NewParser(lexer)

	program := parser.ParseProgram()

	if len(parser.Errors()) > 0 {
		t.Fatalf("got parser errors: %+v", parser.Errors())
	}
	if program == nil {
		t.Fatalf("program is nil")
	}
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
