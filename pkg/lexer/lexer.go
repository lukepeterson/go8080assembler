package lexer

import (
	"strings"
	"unicode"
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	COMMA = "COMMA"
	COLON = "COLON"

	IDENT    = "IDENT"
	NUMBER   = "NUMBER"
	REGISTER = "REGISTER"
	MNEMONIC = "MNEMONIC"

	COMMENT = "COMMENT"
	EOF     = "EOF"

	UNKNOWN = "UNKNOWN"
)

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

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.currentChar = 0
	} else {
		l.currentChar = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++

}

func (l *Lexer) NextToken() Token {
	var token Token

	l.skipWhitespace()

	switch l.currentChar {
	case ',':
		token = newToken(COMMA, l.currentChar)
	case ':':
		token = newToken(COLON, l.currentChar)
	case ';':
		token.Type = COMMENT
		token.Literal = l.readComment()
	case 0:
		token.Literal = ""
		token.Type = EOF
	default:
		if isLetter(l.currentChar) || l.currentChar == '_' {
			literal := l.readIdentifier()
			token.Literal = literal
			token.Type = l.lookupIdentifier(strings.ToUpper(literal))
			return token
		} else if isDigit(l.currentChar) {
			token.Type = NUMBER
			token.Literal = l.readNumber()
			return token
		} else {
			token = newToken(UNKNOWN, l.currentChar)
		}
	}

	l.readChar()
	return token
}

func newToken(tokenType TokenType, currentChar byte) Token {
	return Token{Type: tokenType, Literal: string(currentChar)}
}

func (l *Lexer) skipWhitespace() {
	for l.currentChar == ' ' || l.currentChar == '\t' || l.currentChar == '\n' {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.currentChar) || isDigit(l.currentChar) || l.currentChar == '_' {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position

	// Check for '0x' prefix
	if l.currentChar == '0' && (l.peekChar() == 'x' || l.peekChar() == 'X') {
		l.readChar() // consume '0'
		l.readChar() // consume 'x' or 'X'
	}

	for isDigit(l.currentChar) || isHex(l.currentChar) {
		l.readChar()
	}

	// Check for 'H' suffix
	if l.currentChar == 'H' || l.currentChar == 'h' {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) readComment() string {
	position := l.position
	for l.currentChar != '\n' && l.currentChar != 0 {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) lookupIdentifier(ident string) TokenType {

	ident = strings.ToUpper(ident)

	if tokenType, exists := mnemonics[ident]; exists {
		return tokenType
	}
	if tokenType, exists := registers[ident]; exists {
		return tokenType
	}
	return IDENT
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
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
