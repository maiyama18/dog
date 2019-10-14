package ast

import "github.com/maiyama18/dog/token"

type Node interface {
	TokenLiteral() string
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

type LetStatement struct {
	Token      token.Token
	Identifier *Identifier
	Expression Expression
}

func (l *LetStatement) statement()           {}
func (l *LetStatement) TokenLiteral() string { return l.Token.Literal }

type ReturnStatement struct {
	Token      token.Token
	Expression Expression
}

func (r *ReturnStatement) statement()           {}
func (r *ReturnStatement) TokenLiteral() string { return r.Token.Literal }

type ExpressionStatement struct {
	Token      token.Token // first token of the expression
	Expression Expression
}

func (e *ExpressionStatement) statement()           {}
func (e *ExpressionStatement) TokenLiteral() string { return e.Token.Literal }

type Identifier struct {
	Token token.Token // token.IDENT
	Name  string
}

func (i *Identifier) expression()          {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
