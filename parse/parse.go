package parse

import (
	"github.com/maiyama18/dog/ast"
	"github.com/maiyama18/dog/lex"
	"github.com/maiyama18/dog/token"
)

type Parser struct {
	lexer        *lex.Lexer
	currentToken token.Token
	nextToken    token.Token
}

func New(lexer *lex.Lexer) *Parser {
	p := &Parser{lexer: lexer}

	p.consumeToken()
	p.consumeToken()

	return p
}

func (p *Parser) consumeToken() {
	p.currentToken = p.nextToken
	p.nextToken = p.lexer.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	return nil
}
