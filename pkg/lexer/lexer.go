package lexer

import "fmt"

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
	input        string
	position     int
	readPosition int
	currentChar  byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.currentChar = 0
	} else {
		l.currentChar = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++

	fmt.Printf("l.readPosition: %v\n", l.readPosition)
}
