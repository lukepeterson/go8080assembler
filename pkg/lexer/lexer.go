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
	MNEMONIC = "MNEMONIC"
	REGISTER = "REGISTER"
	NUMBER   = "NUMBER"
	COMMA    = "COMMA"
	COLON    = "COLON"
	LABEL    = "LABEL"
	COMMENT  = "COMMENT"
	EOF      = "EOF"
)

var mnemonics = map[string]TokenType{
	// MOVE, LOAD AND STORE
	"MOV":  MNEMONIC,
	"MVI":  MNEMONIC,
	"LXI":  MNEMONIC,
	"STAX": MNEMONIC,
	"LDAX": MNEMONIC,
	"STA":  MNEMONIC,
	"LDA":  MNEMONIC,
	"SHLD": MNEMONIC,
	"LHLD": MNEMONIC,
	"XCHG": MNEMONIC,

	// STACK OPERATIONS
	"PUSH": MNEMONIC,
	"POP":  MNEMONIC,
	"XTHL": MNEMONIC,
	"SPHL": MNEMONIC,
	"INX":  MNEMONIC,
	"DCX":  MNEMONIC,
	"DAD":  MNEMONIC,

	// JUMP
	"JMP":  MNEMONIC,
	"JC":   MNEMONIC,
	"JNC":  MNEMONIC,
	"JZ":   MNEMONIC,
	"JNZ":  MNEMONIC,
	"JP":   MNEMONIC,
	"JM":   MNEMONIC,
	"JPE":  MNEMONIC,
	"JPO":  MNEMONIC,
	"PCHL": MNEMONIC,

	// CALL
	"CALL": MNEMONIC,
	"CC":   MNEMONIC,
	"CNC":  MNEMONIC,
	"CZ":   MNEMONIC,
	"CNZ":  MNEMONIC,
	"CP":   MNEMONIC,
	"CM":   MNEMONIC,
	"CPE":  MNEMONIC,
	"CPO":  MNEMONIC,

	// RETURN
	"RET": MNEMONIC,
	"RC":  MNEMONIC,
	"RNC": MNEMONIC,
	"RZ":  MNEMONIC,
	"RNZ": MNEMONIC,
	"RP":  MNEMONIC,
	"RM":  MNEMONIC,
	"RPE": MNEMONIC,
	"RPO": MNEMONIC,

	// RESTART
	"RST": MNEMONIC,

	// INCREMENT AND DECREMENT
	"INR": MNEMONIC,
	"DCR": MNEMONIC,

	// ADD
	"ADD": MNEMONIC,
	"ADC": MNEMONIC,
	"ADI": MNEMONIC,
	"ACI": MNEMONIC,

	// SUBTRACT
	"SUB": MNEMONIC,
	"SBB": MNEMONIC,
	"SUI": MNEMONIC,
	"SBI": MNEMONIC,

	// LOGICAL
	"ANA": MNEMONIC,
	"XRA": MNEMONIC,
	"ORA": MNEMONIC,
	"CMP": MNEMONIC,
	"ANI": MNEMONIC,
	"XRI": MNEMONIC,
	"ORI": MNEMONIC,
	"CPI": MNEMONIC,

	// ROTATE
	"RLC": MNEMONIC,
	"RRC": MNEMONIC,
	"RAL": MNEMONIC,
	"RAR": MNEMONIC,

	// SPECIALS
	"CMA": MNEMONIC,
	"STC": MNEMONIC,
	"CMC": MNEMONIC,
	"DAA": MNEMONIC,

	// INPUT/OUTPUT
	"IN":  MNEMONIC,
	"OUT": MNEMONIC,

	// CONTROL
	"EI":  MNEMONIC,
	"DI":  MNEMONIC,
	"NOP": MNEMONIC,
	"HLT": MNEMONIC,
}

// TODO: Split this into 8 and 16-bit registers.  We'll need to differentiate between the
// two to parse registers that have 16-bit operands, eg, PUSH B.
var registers = map[string]TokenType{
	"A": REGISTER,
	"B": REGISTER,
	"C": REGISTER,
	"D": REGISTER,
	"E": REGISTER,
	"H": REGISTER,
	"L": REGISTER,
	"M": REGISTER,
}

type Lexer struct {
	input        string
	position     int
	readPosition int
	currentChar  byte
	Tokens       []Token
}

func New(input string) *Lexer {
	lexer := &Lexer{input: input}
	lexer.readChar()
	return lexer
}

func (l *Lexer) Lex() ([]Token, error) {
	for token := l.NextToken(); token.Type != "EOF"; token = l.NextToken() {
		l.Tokens = append(l.Tokens, token)
	}

	// TODO: Do we need this?
	l.Tokens = append(l.Tokens, Token{Type: EOF})

	// TODO: Make Lex() function actually detect errors
	return l.Tokens, nil
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.currentChar = 0x00
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
		token.Type = COMMA
		token.Literal = ","
	case ':':
		token.Type = COLON
		token.Literal = ":"
	case ';':
		token.Type = COMMENT
		token.Literal = l.readComment()
	case 0x00:
		token.Type = EOF
	default:
		if isLetter(l.currentChar) {
			literal := strings.ToUpper(l.readToken())
			token.Type = l.lookupToken(literal)
			token.Literal = literal
			return token
		}
		if isDigit(l.currentChar) {
			literal := strings.ToUpper(l.readNumber())
			token.Type = NUMBER
			token.Literal = literal
			return token
		}
	}

	l.readChar()
	return token
}

func (l *Lexer) skipWhitespace() {
	for l.currentChar == ' ' || l.currentChar == '\t' || l.currentChar == '\n' {
		l.readChar()
	}
}

func (l *Lexer) readToken() string {
	position := l.position
	for isLetter(l.currentChar) || isDigit(l.currentChar) {
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

	// for isDigit(l.currentChar) || isHex(l.currentChar) {
	for isHex(l.currentChar) {
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
	for l.currentChar != '\n' && l.currentChar != 0x00 {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) lookupToken(token string) TokenType {
	token = strings.ToUpper(token)

	if tokenType, exists := mnemonics[token]; exists {
		return tokenType
	}
	if tokenType, exists := registers[token]; exists {
		return tokenType
	}

	return LABEL
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		// TODO: Find a better way to signify EOF without returning 0x00
		return 0x00
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
	return unicode.Is(unicode.Hex_Digit, rune(char))
}
