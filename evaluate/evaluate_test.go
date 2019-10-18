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
		{input: "5", want: 5},
		{input: "42", want: 42},
		{input: "-42", want: -42},
		{input: "8 + 2", want: 10},
		{input: "8 - 2", want: 6},
		{input: "2 - 8", want: -6},
		{input: "8 * 2", want: 16},
		{input: "8 / 2", want: 4},
		{input: "-8 / 2", want: -4},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			got := eval(test.input)
			testInteger(t, got, test.want)
		})
	}
}

func TestEvalBoolean(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{input: "true", want: true},
		{input: "false", want: false},
		{input: "!true", want: false},
		{input: "!false", want: true},
		{input: "!!true", want: true},
		{input: "!!false", want: false},
		{input: "!5", want: false},
		{input: "!!5", want: true},
		{input: "1 == 1", want: true},
		{input: "2 == 1", want: false},
		{input: "1 != 1", want: false},
		{input: "2 != 1", want: true},
		{input: "2 > 1", want: true},
		{input: "1 > 2", want: false},
		{input: "2 < 1", want: false},
		{input: "1 < 2", want: true},
		{input: "true == true", want: true},
		{input: "true == false", want: false},
		{input: "true != true", want: false},
		{input: "true != false", want: true},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			got := eval(test.input)
			testBoolean(t, got, test.want)
		})
	}
}

func TestEvalIfElseExpression(t *testing.T) {
	tests := []struct {
		input string
		want  interface{}
	}{
		{input: "if (true) { 10 }", want: 10},
		{input: "if (false) { 10 }", want: nil},
		{input: "if (1) { 10 }", want: 10},
		{input: "if (1 < 2) { 10 }", want: 10},
		{input: "if (1 > 2) { 10 }", want: nil},
		{input: "if (1 < 2) { 10 } else { 20 }", want: 10},
		{input: "if (1 > 2) { 10 } else { 20 }", want: 20},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			got := eval(test.input)
			wantedInt, ok := test.want.(int)
			if ok {
				testInteger(t, got, int64(wantedInt))
			} else {
				testNull(t, got)
			}
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

func testNull(t *testing.T, got object.Object) {
	t.Helper()

	_, ok := got.(*object.Null)
	if !ok {
		t.Fatalf("not Integer: %+v", got)
	}
}

func eval(input string) object.Object {
	l := lex.NewLexer(input)
	p := parse.NewParser(l)

	program := p.ParseProgram()

	return Eval(program)
}
