package object

import "fmt"

type Type string

const (
	IntegerType = "INTEGER"
	BooleanType = "BOOLEAN"
	NullType    = "NULL"
)

type Object interface {
	Type() Type
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() Type      { return IntegerType }
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() Type      { return BooleanType }
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }

type Null struct{}

func (n *Null) Type() Type      { return NullType }
func (n *Null) Inspect() string { return "null" }
