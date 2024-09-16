package lexer

type TokenType string

const (
	COMMA TokenType = "COMMA"
	COLON TokenType = "COLON"

	IDENT    TokenType = "IDENT"
	NUMBER   TokenType = "NUMBER"
	REGISTER TokenType = "REGISTER"
	MNEMONIC TokenType = "MNEMONIC"

	COMMENT TokenType = "COMMENT"
	EOF     TokenType = "EOF"

	UNKNOWN TokenType = "UNKNOWN"
)

type Token struct {
	Type    TokenType
	Literal string
}

type Lexer struct {
	input string
	// position int
}

func New(input string) *Lexer {
	lexer := &Lexer{input: input}
	// lexer.readChar()
	return lexer
}
