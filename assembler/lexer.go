package assembler

import (
	"strings"
)

type TokenType int

const (
	LABEL TokenType = iota
	OPCODE
	OPERAND1
	OPERAND2
)

type Token struct {
	Type  TokenType
	Value string
}

func NewToken(Type TokenType, Value string) *Token {
	return &Token{
		Type:  Type,
		Value: Value,
	}
}

func stripComment(line string) string {
	index := strings.Index(line, ";")
	if index != -1 {
		return strings.TrimSpace(line[:index])
	}
	return line
}

func tokeniseLine(line string) []Token {
	tokens := []Token{}

	line = stripComment(line)
	line = strings.TrimSpace(line)
	// fmt.Printf("\nline: %v\n", line)

	parts := strings.Split(line, ":")
	if len(parts) == 2 { // We must have a label
		tokens = append(tokens, Token{Type: LABEL, Value: parts[0]})
		line = strings.TrimSpace(parts[1])
	}

	parts = strings.Split(line, ",")
	operand2 := ""
	if len(parts) == 2 { // we must have a second operand
		operand2 = strings.TrimSpace(parts[1])
		line = strings.TrimSpace(parts[0])
	}

	parts = strings.Split(line, " ")
	operand1 := ""
	if len(parts) == 2 { // we must have a first operand
		operand1 = strings.TrimSpace(parts[1])
		line = strings.TrimSpace(parts[0])
	}

	if len(line) > 0 { // remainder must be the opcode
		tokens = append(tokens, Token{Type: OPCODE, Value: line})
	}

	if operand1 != "" {
		tokens = append(tokens, Token{Type: OPERAND1, Value: operand1})
	}

	if operand2 != "" {
		tokens = append(tokens, Token{Type: OPERAND2, Value: operand2})
	}

	return tokens
}

func lexicalAnalysis(assemblyCode string) []Token {
	tokens := []Token{}
	lines := strings.Split(assemblyCode, "\n")

	for _, line := range lines {
		lineTokens := tokeniseLine(line)
		tokens = append(tokens, lineTokens...)
	}

	return tokens
}
