package lexer

import "github.com/maiyama18/dog/token"

type Lexer struct {
	input        []rune
	position     int
	nextPosition int
	currentRune  rune
}

func New(input string) *Lexer {
	l := &Lexer{input: []rune(input)}
	l.consumeRune()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var t token.Token

	switch l.currentRune {
	case '=':
		t = newToken(token.ASSIGN, l.currentRune)
	case '+':
		t = newToken(token.PLUS, l.currentRune)
	case '(':
		t = newToken(token.LPAREN, l.currentRune)
	case ')':
		t = newToken(token.RPAREN, l.currentRune)
	case '{':
		t = newToken(token.LBRACE, l.currentRune)
	case '}':
		t = newToken(token.RBRACE, l.currentRune)
	case ',':
		t = newToken(token.COMMA, l.currentRune)
	case ';':
		t = newToken(token.SEMICOLON, l.currentRune)
	case 0:
		t = newToken(token.EOF, ' ')
	}

	l.consumeRune()
	return t
}

func (l *Lexer) consumeRune() {
	if l.nextPosition >= len(l.input) {
		l.currentRune = 0
	} else {
		l.currentRune = l.input[l.nextPosition]
	}
	l.position = l.nextPosition
	l.nextPosition++
}

func newToken(tokenType token.Type, literal rune) token.Token {
	return token.Token{Type: tokenType, Literal: string(literal)}
}
