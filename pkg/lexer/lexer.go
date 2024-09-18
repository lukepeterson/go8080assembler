package lexer

import (
	"strings"
	"unicode"
)

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

var mnemonics = map[string]TokenType{
	"MOV": MNEMONIC,
	"MVI": MNEMONIC,
	"ADD": MNEMONIC,
	"SUB": MNEMONIC,
	"INX": MNEMONIC,
	"INR": MNEMONIC,
	"JMP": MNEMONIC,
	// and the rest
}

var registers = map[string]TokenType{
	"A": REGISTER,
	"B": REGISTER,
	"C": REGISTER,
	"D": REGISTER,
	// and the rest
}

type Lexer struct {
	input        string
	position     int
	readPosition int
	currentChar  byte
}

func New(input string) *Lexer {
	lexer := &Lexer{input: input}
	lexer.readChar()
	return lexer
}

func (lexer *Lexer) readChar() {
	if lexer.readPosition >= len(lexer.input) {
		lexer.currentChar = 0
	} else {
		lexer.currentChar = lexer.input[lexer.readPosition]
	}
	lexer.position = lexer.readPosition
	lexer.readPosition++

}

func (lexer *Lexer) NextToken() Token {
	var token Token

	lexer.skipWhitespace()

	switch lexer.currentChar {
	case ',':
		token = newToken(COMMA, lexer.currentChar)
	case ':':
		token = newToken(COLON, lexer.currentChar)
	case ';':
		token.Type = COMMENT
		token.Literal = lexer.readComment()
	case 0:
		token.Literal = ""
		token.Type = EOF
	default:
		if isLetter(lexer.currentChar) || lexer.currentChar == '_' {
			literal := lexer.readIdentifier()
			token.Literal = literal
			token.Type = lexer.lookupIdentifier(strings.ToUpper(literal))
			return token
		} else if isDigit(lexer.currentChar) {
			token.Type = NUMBER
			token.Literal = lexer.readNumber()
			return token
		} else {
			token = newToken(UNKNOWN, lexer.currentChar)
		}
	}

	lexer.readChar()
	return token
}

func newToken(tokenType TokenType, currentChar byte) Token {
	return Token{Type: tokenType, Literal: string(currentChar)}
}

func (lexer *Lexer) skipWhitespace() {
	for lexer.currentChar == ' ' || lexer.currentChar == '\t' || lexer.currentChar == '\n' {
		lexer.readChar()
	}
}

func (lexer *Lexer) readIdentifier() string {
	position := lexer.position
	for isLetter(lexer.currentChar) || isDigit(lexer.currentChar) || lexer.currentChar == '_' {
		lexer.readChar()
	}
	return lexer.input[position:lexer.position]
}

func (lexer *Lexer) readNumber() string {
	position := lexer.position

	// Check for '0x' prefix
	if lexer.currentChar == '0' && (lexer.peekChar() == 'x' || lexer.peekChar() == 'X') {
		lexer.readChar() // consume '0'
		lexer.readChar() // consume 'x' or 'X'
	}

	for isDigit(lexer.currentChar) || isHex(lexer.currentChar) {
		lexer.readChar()
	}

	// Check for 'H' suffix
	if lexer.currentChar == 'H' || lexer.currentChar == 'h' {
		lexer.readChar()
	}

	return lexer.input[position:lexer.position]
}

func (lexer *Lexer) readComment() string {
	position := lexer.position
	for lexer.currentChar != '\n' && lexer.currentChar != 0 {
		lexer.readChar()
	}
	return lexer.input[position:lexer.position]
}

func (lexer *Lexer) lookupIdentifier(ident string) TokenType {

	ident = strings.ToUpper(ident)

	if tokenType, ok := mnemonics[ident]; ok {
		return tokenType
	}
	if tokenType, ok := registers[ident]; ok {
		return tokenType
	}
	return IDENT
}

func (lexer *Lexer) peekChar() byte {
	if lexer.readPosition >= len(lexer.input) {
		return 0
	} else {
		return lexer.input[lexer.readPosition]
	}
}

func isLetter(char byte) bool {
	return unicode.IsLetter(rune(char))
}

func isDigit(char byte) bool {
	return unicode.IsDigit(rune(char))
}

func isHex(char byte) bool {
	return ('0' <= char && char <= '9') ||
		('a' <= char && char <= 'f') ||
		('A' <= char && char <= 'F')
}
