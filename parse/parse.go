package parse

import (
	"errors"
	"fmt"
	"strconv"

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

func getPrecedence(tokenType token.Type) Precedence {
	switch tokenType {
	case token.EQ, token.NOTEQ:
		return EQUALS
	case token.LT, token.GT:
		return LESSGREATER
	case token.PLUS, token.MINUS:
		return SUM
	case token.ASTERISK, token.SLASH:
		return PRODUCT
	default:
		return LOWEST
	}
}

type (
	parsePrefixFunc func() ast.Expression
	parseInfixFunc  func(ast.Expression) ast.Expression
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

func (p *Parser) nextPrecedence() Precedence {
	return getPrecedence(p.nextToken.Type)
}

func (p *Parser) expectNextTokenType(tokenType token.Type) error {
	if !p.isNextTokenType(tokenType) {
		return fmt.Errorf("expect token type %q, bug got %q", tokenType, p.nextToken.Type)
	}
	p.consumeToken()
	return nil
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

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	tok := p.currentToken

	var statements []ast.Statement
	for !p.isNextTokenType(token.RBRACE) && !p.isNextTokenType(token.EOF) {
		p.consumeToken()
		s := p.parseStatement()
		if s != nil {
			statements = append(statements, s)
		}
	}
	if err := p.expectNextTokenType(token.RBRACE); err != nil {
		p.addError(err)
		return nil
	}

	return &ast.BlockStatement{Token: tok, Statements: statements}
}

func (p *Parser) parseExpression(precedence Precedence) ast.Expression {
	parsePrefix, err := p.getParsePrefixFunc()
	if err != nil {
		p.addError(err)
		return nil
	}

	left := parsePrefix()
	for precedence < p.nextPrecedence() {
		parseInfix, err := p.getParseInfixFunc()
		if err != nil {
			p.addError(err)
			return nil
		}

		p.consumeToken()
		left = parseInfix(left)
	}

	return left
}

func (p *Parser) getParsePrefixFunc() (parsePrefixFunc, error) {
	switch p.currentToken.Type {
	case token.IDENT:
		return p.parseIdentifier, nil
	case token.INT:
		return p.parseIntegerLiteral, nil
	case token.TRUE, token.FALSE:
		return p.parseBooleanLiteral, nil
	case token.BANG, token.MINUS:
		return p.parsePrefixExpression, nil
	case token.LPAREN:
		return p.parseGroupedExpression, nil
	case token.IF:
		return p.parseIfExpression, nil
	default:
		return nil, fmt.Errorf("could not find to parse prefix function for token type %+v", p.currentToken)
	}
}

func (p *Parser) getParseInfixFunc() (parseInfixFunc, error) {
	switch p.nextToken.Type {
	case token.PLUS, token.MINUS, token.ASTERISK, token.SLASH, token.EQ, token.NOTEQ, token.GT, token.LT:
		return p.parseInfixExpression, nil
	default:
		return nil, fmt.Errorf("could not find to parse infix function for token type %+v", p.currentToken)
	}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	tok := p.currentToken
	p.consumeToken()
	right := p.parseExpression(PREFIX)
	return &ast.PrefixExpression{Token: tok, Operator: tok.Literal, Right: right}
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	opToken := p.currentToken
	p.consumeToken()
	precedence := getPrecedence(opToken.Type)
	right := p.parseExpression(precedence)
	return &ast.InfixExpression{Token: opToken, Operator: opToken.Literal, Left: left, Right: right}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.consumeToken()
	exp := p.parseExpression(LOWEST)
	if err := p.expectNextTokenType(token.RPAREN); err != nil {
		p.addError(err)
		return nil
	}
	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	tok := p.currentToken
	if err := p.expectNextTokenType(token.LPAREN); err != nil {
		p.addError(err)
		return nil
	}
	p.consumeToken()

	condition := p.parseExpression(LOWEST)
	if condition == nil {
		p.addError(errors.New("failed to parse condition of if expression"))
	}

	if err := p.expectNextTokenType(token.RPAREN); err != nil {
		p.addError(err)
		return nil
	}
	if err := p.expectNextTokenType(token.LBRACE); err != nil {
		p.addError(err)
		return nil
	}

	consequence := p.parseBlockStatement()

	var alternative *ast.BlockStatement
	if p.isNextTokenType(token.ELSE) {
		p.consumeToken()
		if err := p.expectNextTokenType(token.LBRACE); err != nil {
			p.addError(err)
			return nil
		}
		alternative = p.parseBlockStatement()
	}

	return &ast.IfExpression{Token: tok, Condition: condition, Consequence: consequence, Alternative: alternative}
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currentToken, Name: p.currentToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	i, err := strconv.ParseInt(p.currentToken.Literal, 10, 64)
	if err != nil {
		p.addError(fmt.Errorf("failed to parse %q as integer literal: %v", p.currentToken.Literal, err))
		return nil
	}
	return &ast.IntegerLiteral{Token: p.currentToken, Value: i}
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{Token: p.currentToken, Value: p.isCurrentTokenType(token.TRUE)}
}
