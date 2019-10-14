package ast

import (
	"fmt"
	"strings"

	"github.com/maiyama18/dog/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statement()
}

type Expression interface {
	Node
	expression()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) == 0 {
		return ""
	}
	return p.Statements[0].TokenLiteral()
}
func (p *Program) String() string {
	var buff strings.Builder
	for _, s := range p.Statements {
		buff.WriteString(s.String())
	}
	return buff.String()
}

type LetStatement struct {
	Token      token.Token
	Identifier *Identifier
	Expression Expression
}

func (l *LetStatement) statement()           {}
func (l *LetStatement) TokenLiteral() string { return l.Token.Literal }
func (l *LetStatement) String() string {
	var buff strings.Builder
	buff.WriteString(fmt.Sprintf("let %s = ", l.Identifier.Name))
	if l.Expression != nil {
		buff.WriteString(l.Expression.String())
	}
	buff.WriteString(";")
	return buff.String()
}

type ReturnStatement struct {
	Token      token.Token
	Expression Expression
}

func (r *ReturnStatement) statement()           {}
func (r *ReturnStatement) TokenLiteral() string { return r.Token.Literal }
func (r *ReturnStatement) String() string {
	var buff strings.Builder
	buff.WriteString("return ")
	if r.Expression != nil {
		buff.WriteString(r.Expression.String())
	}
	buff.WriteString(";")
	return buff.String()
}

type ExpressionStatement struct {
	Token      token.Token // first token of the expression
	Expression Expression
}

func (e *ExpressionStatement) statement()           {}
func (e *ExpressionStatement) TokenLiteral() string { return e.Token.Literal }
func (e *ExpressionStatement) String() string {
	var buff strings.Builder
	if e.Expression != nil {
		buff.WriteString(e.Expression.String())
	}
	buff.WriteString(";")
	return buff.String()
}

type Identifier struct {
	Token token.Token // token.IDENT
	Name  string
}

func (i *Identifier) expression()          {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Name }

type IntegerLiteral struct {
	Token token.Token // token.INT
	Value int64
}

func (i *IntegerLiteral) expression()          {}
func (i *IntegerLiteral) TokenLiteral() string { return i.Token.Literal }
func (i *IntegerLiteral) String() string       { return i.Token.Literal }
