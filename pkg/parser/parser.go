package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/lukepeterson/go8080assembler/pkg/lexer"
)

type Parser struct {
	tokens           []lexer.Token
	position         int
	bytecode         []byte
	labelDefinitions map[string]uint16 // Stores resolved label addresses
	labelReferences  map[string]uint16 // Tracks unresolved label usages
}

func New(tokens []lexer.Token) *Parser {
	return &Parser{
		tokens:           tokens,
		position:         0,
		labelDefinitions: make(map[string]uint16),
		labelReferences:  make(map[string]uint16),
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
	for p.currentToken().Type != lexer.EOF {

		switch p.currentToken().Type {
		case lexer.MNEMONIC:
			hexCode, err := p.parseInstruction()
			if err != nil {
				return nil, err
			}
			p.bytecode = append(p.bytecode, hexCode...)

		case lexer.COMMENT:
			// do nothing

		case lexer.LABEL:

			// TODO: Check for duplicate labels and error if found
			p.labelDefinitions[p.currentToken().Literal] = uint16(len(p.bytecode))

			p.advanceToken()

		default:
			return nil, fmt.Errorf("unexpected token type \"%s\", literal: \"%s\"", p.currentToken().Type, p.currentToken().Literal)
		}

		p.advanceToken()
	}

	// TODO: Split this out into a second pass
	for labelReference, location := range p.labelReferences {
		hex := p.labelDefinitions[labelReference]

		highByte := uint8(hex >> 8)
		lowByte := uint8(hex & 0x00FF)

		p.bytecode[location] = lowByte
		p.bytecode[location+1] = highByte
	}

	return p.bytecode, nil
}

func (p *Parser) parseInstruction() ([]byte, error) {
	instruction := p.currentToken().Literal
	p.advanceToken()

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
	dest := p.currentToken().Literal
	p.advanceToken()

	if p.currentToken().Type != lexer.COMMA {
		return nil, fmt.Errorf("expected comma, got: %s", p.currentToken().Literal)
	}
	p.advanceToken()

	if p.currentToken().Type != lexer.REGISTER {
		return nil, fmt.Errorf("expected register, got: %s", p.currentToken().Literal)
	}
	src := p.currentToken().Literal

	return generateMOVHex(src, dest)
}

func generateMOVHex(src, dest string) ([]byte, error) {
	registerMap := map[string]byte{
		"A": 0x07, "B": 0x00, "C": 0x01, "D": 0x02, "E": 0x03, "H": 0x04, "L": 0x05, "M": 0x06,
	}

	srcRegister, valid := registerMap[src]
	if !valid {
		return nil, fmt.Errorf("invalid source register in MOV: %s", src)
	}

	destRegister, valid := registerMap[dest]
	if !valid {
		return nil, fmt.Errorf("invalid destination register in MOV: %s", dest)
	}

	opcode := byte(0x40) | (destRegister << 3) | srcRegister
	return []byte{opcode}, nil
}

func (p *Parser) parseJMP() ([]byte, error) {

	if p.currentToken().Type == lexer.NUMBER {
		highByte, lowByte, err := parseHex(p.currentToken().Literal)
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}

		return []byte{0xC3, lowByte, highByte}, nil
	}

	// If the destination isn't a number, assume it's a label
	if p.currentToken().Type == lexer.LABEL {

		p.labelReferences[p.currentToken().Literal] = uint16(len(p.bytecode) + 1)

		// Labels can be used before they exist, so we use 0x0000 as a placeholder
		// until we know all the label locations.
		return []byte{0xC3, 0x00, 0x00}, nil
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
