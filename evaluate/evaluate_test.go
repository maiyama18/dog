package evaluate

import (
	"testing"

	"github.com/maiyama18/dog/lex"
	"github.com/maiyama18/dog/object"
	"github.com/maiyama18/dog/parse"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input string
		want  int64
	}{
		{
			input: "5",
			want:  5,
		},
		{
			input: "42",
			want:  42,
		},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			got := eval(test.input)
			testInteger(t, got, test.want)
		})
	}
}

func testInteger(t *testing.T, got object.Object, want int64) {
	t.Helper()

	integer, ok := got.(*object.Integer)
	if !ok {
		t.Fatalf("not Integer: %+v", got)
	}

	if integer.Value != want {
		t.Fatalf("integer value wrong. want=%d, got=%d", want, integer.Value)
	}
}

func eval(input string) object.Object {
	l := lex.NewLexer(input)
	p := parse.NewParser(l)

	program := p.ParseProgram()

	return Eval(program)
}
