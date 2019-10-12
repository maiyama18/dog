package parse

import (
	"fmt"

	"github.com/maiyama18/dog/ast"
	"github.com/maiyama18/dog/lex"
	"github.com/maiyama18/dog/token"
)

type Parser struct {
	lexer        *lex.Lexer
	currentToken token.Token
	nextToken    token.Token

	errors []error
}

func NewParser(lexer *lex.Lexer) *Parser {
	p := &Parser{lexer: lexer}

	p.consumeToken()
	p.consumeToken()

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
		return nil
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
