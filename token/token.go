package token

type Type string

type Token struct {
	Type    Type
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT = "IDENT"
	INT   = "INT"

	// operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"
	BANG     = "!"

	EQ    = "=="
	NOTEQ = "!="
	LT    = "<"
	GT    = ">"

	// delimiters
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
)

func TypeFromLiteral(literal string) Type {
	switch literal {
	case "fn":
		return FUNCTION
	case "let":
		return LET
	case "if":
		return IF
	case "else":
		return ELSE
	case "return":
		return RETURN
	case "true":
		return TRUE
	case "false":
		return FALSE
	default:
		return IDENT
	}
}
