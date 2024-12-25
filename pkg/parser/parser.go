package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/lukepeterson/go8080assembler/pkg/lexer"
)

type Parser struct {
	tokens   []lexer.Token
	position int
}

func New(tokens []lexer.Token) *Parser {
	return &Parser{
		tokens:   tokens,
		position: 0,
	}
}

func (p *Parser) advanceToken() {
	if p.position < len(p.tokens)-1 {
		p.position++
	}
}

func (p *Parser) currentToken() lexer.Token {
	if p.position < len(p.tokens) {
		return p.tokens[p.position]
	}
	return lexer.Token{Type: lexer.EOF}
}

func (p *Parser) Parse() ([]byte, error) {
	hexCodes := []byte{}

	for p.currentToken().Type != lexer.EOF {
		// fmt.Printf("p.currentToken(): %v\n", p.currentToken())

		switch p.currentToken().Type {
		case lexer.MNEMONIC:
			hexCode, err := p.parseInstruction()
			if err != nil {
				return nil, err
			}
			hexCodes = append(hexCodes, hexCode...)

		case lexer.COMMENT:
			// do nothing

		case lexer.LABEL:
			fmt.Printf("Found label definition (\"%s\") at position 0x%02X\n", p.currentToken().Literal, len(hexCodes))

			// store our label in a lookup table?
			p.advanceToken() // skip past the comma

		default:
			return nil, fmt.Errorf("unexpected token type \"%s\", literal: \"%s\"", p.currentToken().Type, p.currentToken().Literal)
		}

		p.advanceToken()
	}

	return hexCodes, nil
}

func (p *Parser) parseInstruction() ([]byte, error) {
	instruction := p.currentToken().Literal
	p.advanceToken() // Move past the mnemonic

	switch instruction {
	case "MOV":
		return p.parseMOV()
	case "JMP":
		return p.parseJMP()
	default:
		return nil, fmt.Errorf("unknown instruction: %s", instruction)
	}
}

func (p *Parser) parseMOV() ([]byte, error) {
	if p.currentToken().Type != lexer.REGISTER {
		return nil, fmt.Errorf("expected register, got: %s", p.currentToken().Literal)
	}
	dest := p.currentToken().Literal // Destination register
	p.advanceToken()

	if p.currentToken().Type != lexer.COMMA {
		return nil, fmt.Errorf("expected comma, got: %s", p.currentToken().Literal)
	}
	p.advanceToken() // Skip COMMA

	if p.currentToken().Type != lexer.REGISTER {
		return nil, fmt.Errorf("expected register, got: %s", p.currentToken().Literal)
	}
	src := p.currentToken().Literal // Source register

	return generateMOVHex(src, dest)
}

func generateMOVHex(src, dest string) ([]byte, error) {
	registerMap := map[string]byte{
		"A": 0x07, "B": 0x00, "C": 0x01, "D": 0x02, "E": 0x03, "H": 0x04, "L": 0x05, "M": 0x06,
	}

	srcCode, srcOk := registerMap[src]
	destCode, destOk := registerMap[dest]
	if !srcOk || !destOk {
		return nil, fmt.Errorf("invalid register in MOV: src=%s, dest=%s", src, dest)
	}

	opcode := byte(0x40) | (destCode << 3) | srcCode
	return []byte{opcode}, nil
}

func (p *Parser) parseJMP() ([]byte, error) {

	if p.currentToken().Type == lexer.NUMBER {
		// fmt.Printf("p.currentToken().Literal: %T\n", p.currentToken().Literal)
		highByte, lowByte, err := parseHex(p.currentToken().Literal)
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}

		return []byte{0xC3, lowByte, highByte}, nil
	}

	if p.currentToken().Type == lexer.LABEL {

		// Do the label lookup steps here?
		// We might not have a label yet, so perhaps store a placeholder for a second parser round?

		lowByte := byte(0xFA)
		highByte := byte(0xFB)

		return []byte{0xC3, lowByte, highByte}, nil
	}

	return nil, fmt.Errorf("expected address or label, got: %s", p.currentToken().Type)

}

func parseHex(token string) (uint8, uint8, error) {
	token = strings.TrimSuffix(token, "H")
	token = strings.TrimSuffix(token, "h")
	token = strings.TrimPrefix(token, "0X")
	token = strings.TrimPrefix(token, "0x")
	hex, err := strconv.ParseUint(token, 16, 16)
	if err != nil {
		return 0, 0, err
	}

	highByte := uint8(hex >> 8)
	lowByte := uint8(hex & 0x00FF)

	return highByte, lowByte, nil
}
