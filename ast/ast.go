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

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (b *BlockStatement) statement()           {}
func (b *BlockStatement) TokenLiteral() string { return b.Token.Literal }
func (b *BlockStatement) String() string {
	var buff strings.Builder
	for _, s := range b.Statements {
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

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (i *IfExpression) expression()          {}
func (i *IfExpression) TokenLiteral() string { return i.Token.Literal }
func (i *IfExpression) String() string {
	var buff strings.Builder
	buff.WriteString(fmt.Sprintf("if (%s) { %s }", i.Condition.String(), i.Consequence.String()))
	if i.Alternative != nil {
		buff.WriteString(fmt.Sprintf(" else { %s }", i.Alternative.String()))
	}
	return buff.String()
}

type FunctionLiteral struct {
	Token      token.Token
	Parameters []Identifier
	Body       *BlockStatement
}

func (f *FunctionLiteral) expression()          {}
func (f *FunctionLiteral) TokenLiteral() string { return f.Token.Literal }
func (f *FunctionLiteral) String() string {
	var paramNames []string
	for _, p := range f.Parameters {
		paramNames = append(paramNames, p.Name)
	}
	return fmt.Sprintf("fn (%s) { %s }", strings.Join(paramNames, ", "), f.Body.String())
}

type CallExpression struct {
	Token     token.Token
	Function  Expression // Identifier or FunctionLiteral
	Arguments []Expression
}

func (c *CallExpression) expression()          {}
func (c *CallExpression) TokenLiteral() string { return c.Token.Literal }
func (c *CallExpression) String() string {
	var argStrs []string
	for _, a := range c.Arguments {
		argStrs = append(argStrs, a.String())
	}
	return fmt.Sprintf("%s(%s)", c.Function.String(), strings.Join(argStrs, ", "))
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (p *PrefixExpression) expression()          {}
func (p *PrefixExpression) TokenLiteral() string { return p.Token.Literal }
func (p *PrefixExpression) String() string       { return fmt.Sprintf("(%s%s)", p.Operator, p.Right.String()) }

type InfixExpression struct {
	Token    token.Token
	Operator string
	Left     Expression
	Right    Expression
}

func (i *InfixExpression) expression()          {}
func (i *InfixExpression) TokenLiteral() string { return i.Token.Literal }
func (i *InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", i.Left.String(), i.Operator, i.Right.String())
}

type Identifier struct {
	Token token.Token
	Name  string
}

func (i *Identifier) expression()          {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Name }

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (i *IntegerLiteral) expression()          {}
func (i *IntegerLiteral) TokenLiteral() string { return i.Token.Literal }
func (i *IntegerLiteral) String() string       { return i.Token.Literal }

type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (b *BooleanLiteral) expression()          {}
func (b *BooleanLiteral) TokenLiteral() string { return b.Token.Literal }
func (b *BooleanLiteral) String() string       { return b.Token.Literal }
