package evaluate

import (
	"testing"

	"github.com/maiyama18/dog/lex"
	"github.com/maiyama18/dog/object"
	"github.com/maiyama18/dog/parse"
)

func TestEvalInteger(t *testing.T) {
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
		{
			input: "-42",
			want:  -42,
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

func TestEvalBoolean(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{
			input: "true",
			want:  true,
		},
		{
			input: "false",
			want:  false,
		},
		{
			input: "!true",
			want:  false,
		},
		{
			input: "!false",
			want:  true,
		},
		{
			input: "!!true",
			want:  true,
		},
		{
			input: "!!false",
			want:  false,
		},
		{
			input: "!5",
			want:  false,
		},
		{
			input: "!!5",
			want:  true,
		},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			got := eval(test.input)
			testBoolean(t, got, test.want)
		})
	}
}

func testBoolean(t *testing.T, got object.Object, want bool) {
	t.Helper()

	boolean, ok := got.(*object.Boolean)
	if !ok {
		t.Fatalf("not Integer: %+v", got)
	}

	if boolean.Value != want {
		t.Fatalf("integer value wrong. want=%t, got=%t", want, boolean.Value)
	}
}

func eval(input string) object.Object {
	l := lex.NewLexer(input)
	p := parse.NewParser(l)

	program := p.ParseProgram()

	return Eval(program)
}
