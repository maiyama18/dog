package parse

import (
	"fmt"

	"github.com/maiyama18/dog/ast"
	"github.com/maiyama18/dog/lex"
	"github.com/maiyama18/dog/token"
)

type Precedence int

const (
	LOWEST Precedence = iota + 1
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

type (
	parsePrefix func() ast.Expression
	parseInfix  func(ast.Expression) ast.Expression
)

type Parser struct {
	lexer        *lex.Lexer
	currentToken token.Token
	nextToken    token.Token

	parsePrefixFuncs map[token.Type]parsePrefix
	parseInfixFuncs  map[token.Type]parseInfix

	errors []error
}

func NewParser(lexer *lex.Lexer) *Parser {
	p := &Parser{lexer: lexer}

	p.consumeToken()
	p.consumeToken()

	p.parsePrefixFuncs = map[token.Type]parsePrefix{
		token.IDENT: p.parseIdentifier,
	}

	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	var statements []ast.Statement

	for !p.isCurrentTokenType(token.EOF) {
		s := p.parseStatement()
		if s != nil {
			statements = append(statements, s)
		}
		p.consumeToken()
	}

	return &ast.Program{Statements: statements}
}

func (p *Parser) Errors() []error {
	return p.errors
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	tok := p.currentToken

	if err := p.expectNextTokenType(token.IDENT); err != nil {
		p.addError(err)
		return nil
	}
	ident := &ast.Identifier{Token: p.currentToken, Name: p.currentToken.Literal}

	if err := p.expectNextTokenType(token.ASSIGN); err != nil {
		p.addError(err)
		return nil
	}

	// TODO: parse expression
	for !p.isCurrentTokenType(token.SEMICOLON) {
		p.consumeToken()
	}

	return &ast.LetStatement{Token: tok, Identifier: ident}
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	tok := p.currentToken

	// TODO: parse expression
	for !p.isCurrentTokenType(token.SEMICOLON) {
		p.consumeToken()
	}

	return &ast.ReturnStatement{Token: tok}
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	tok := p.currentToken

	// TODO: parse expression
	expression := p.parseExpression(LOWEST)

	if p.isNextTokenType(token.SEMICOLON) {
		p.consumeToken()
	}

	return &ast.ExpressionStatement{Token: tok, Expression: expression}
}

func (p *Parser) parseExpression(precedence Precedence) ast.Expression {
	f := p.getParsePrefixFunc()
	if f == nil {
		p.addError(fmt.Errorf("could not find to find parse function for token type %s", p.currentToken))
		return nil
	}

	left := f()

	return left
}

func (p *Parser) consumeToken() {
	p.currentToken = p.nextToken
	p.nextToken = p.lexer.NextToken()
}

func (p *Parser) addError(err error) {
	p.errors = append(p.errors, err)
}

func (p *Parser) isCurrentTokenType(tokenType token.Type) bool {
	return p.currentToken.Type == tokenType
}

func (p *Parser) isNextTokenType(tokenType token.Type) bool {
	return p.nextToken.Type == tokenType
}

func (p *Parser) expectNextTokenType(tokenType token.Type) error {
	if !p.isNextTokenType(tokenType) {
		return fmt.Errorf("expect token type %v, bug got %v", tokenType, p.currentToken)
	}
	p.consumeToken()
	return nil
}

func (p *Parser) getParsePrefixFunc() parsePrefix {
	return p.parsePrefixFuncs[p.currentToken.Type]
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currentToken, Name: p.currentToken.Literal}
}
