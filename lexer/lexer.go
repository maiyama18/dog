package lexer

import (
	"unicode"

	"github.com/maiyama18/dog/token"
)

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

	l.skipSpaces()

	switch l.currentRune {
	case '=':
		t = newToken(token.ASSIGN, l.currentRune)
	case '+':
		t = newToken(token.PLUS, l.currentRune)
	case '-':
		t = newToken(token.MINUS, l.currentRune)
	case '*':
		t = newToken(token.ASTERISK, l.currentRune)
	case '/':
		t = newToken(token.SLASH, l.currentRune)
	case '!':
		t = newToken(token.BANG, l.currentRune)
	case '<':
		t = newToken(token.LT, l.currentRune)
	case '>':
		t = newToken(token.GT, l.currentRune)
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
	default:
		if isLetter(l.currentRune) {
			literal := l.readIdentifier()
			tokenType := token.TypeFromLiteral(literal)
			t = token.Token{Type: tokenType, Literal: literal}
		} else if isDigit(l.currentRune) {
			literal := l.readNumber()
			t = token.Token{Type: token.INT, Literal: literal}
		}
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

func (l *Lexer) peekRune() rune {
	if l.nextPosition >= len(l.input) {
		return 0
	}
	return l.input[l.nextPosition]
}

func (l *Lexer) readIdentifier() string {
	start := l.position
	for isLetter(l.peekRune()) {
		l.consumeRune()
	}
	return string(l.input[start : l.position+1])
}

func (l *Lexer) readNumber() string {
	start := l.position
	for isDigit(l.peekRune()) {
		l.consumeRune()
	}
	return string(l.input[start : l.position+1])
}

func (l *Lexer) skipSpaces() {
	for unicode.IsSpace(l.currentRune) {
		l.consumeRune()
	}
}

func newToken(tokenType token.Type, literal rune) token.Token {
	return token.Token{Type: tokenType, Literal: string(literal)}
}

func isLetter(r rune) bool {
	return r == '_' || unicode.IsLetter(r)
}

func isDigit(r rune) bool {
	return unicode.IsDigit(r)
}
