package lexer

import (
	"testing"

	"github.com/maiyama18/dog/token"
)

func TestNextToken(t *testing.T) {
	input := `=+(){},;`

	expectedTokens := []token.Token{
		{Type: token.ASSIGN, Literal: "="},
		{Type: token.PLUS, Literal: "+"},
		{Type: token.LPAREN, Literal: "("},
		{Type: token.RPAREN, Literal: ")"},
		{Type: token.LBRACE, Literal: "{"},
		{Type: token.RBRACE, Literal: "}"},
		{Type: token.COMMA, Literal: ","},
		{Type: token.SEMICOLON, Literal: ";"},
	}

	l := New(input)

	for i, expected := range expectedTokens {
		actual := l.NextToken()

		if actual.Type != expected.Type || actual.Literal != actual.Literal {
			t.Fatalf("[%d] token wrong. want=%+v, got=%+v", i, expected, actual)
		}
	}
}
