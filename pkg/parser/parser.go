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

type parseFunc func(*Parser) ([]byte, error)

var instructionMap = map[string]parseFunc{
	"MOV":  (*Parser).parseMOV,
	"MVI":  (*Parser).parseMVI,
	"LXI":  (*Parser).parseLXI,
	"STAX": (*Parser).parseSTAX,
	"LDAX": (*Parser).parseLDAX,
	"JMP":  (*Parser).parseJMP,
	"DB":   (*Parser).parseDB,
}

var registerMap8 = map[string]byte{
	"B": 0x00, "C": 0x01, "D": 0x02, "E": 0x03, "H": 0x04, "L": 0x05, "M": 0x06, "A": 0x07,
}

var registerMap16 = map[string]byte{
	"B": 0x00, "D": 0x01, "H": 0x02, "SP": 0x03,
}

func (p *Parser) parseInstruction() ([]byte, error) {
	instruction := p.currentToken().Literal
	p.advanceToken()

	parseFunc, exists := instructionMap[instruction]
	if !exists {
		return nil, fmt.Errorf("unknown instruction: %s", instruction)
	}
	return parseFunc(p)
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
	srcRegister, valid := registerMap8[src]
	if !valid {
		return nil, fmt.Errorf("invalid source register for MOV: %s", src)
	}

	destRegister, valid := registerMap8[dest]
	if !valid {
		return nil, fmt.Errorf("invalid destination register for MOV: %s", dest)
	}

	// TODO: This is likely common to other instructions - split this apart when you have more examples
	opcode := byte(0x40) | (destRegister << 3) | srcRegister
	return []byte{opcode}, nil
}

func (p *Parser) parseMVI() ([]byte, error) {
	if p.currentToken().Type != lexer.REGISTER {
		return nil, fmt.Errorf("expected register, got: %s", p.currentToken().Literal)
	}
	dest := p.currentToken().Literal
	p.advanceToken()

	if p.currentToken().Type != lexer.COMMA {
		return nil, fmt.Errorf("expected comma, got: %s", p.currentToken().Literal)
	}
	p.advanceToken()

	if p.currentToken().Type != lexer.NUMBER {
		return nil, fmt.Errorf("expected number, got: %s", p.currentToken().Literal)
	}
	byteToStore := p.currentToken().Literal

	return generateMVIHex(dest, byteToStore)
}

func generateMVIHex(dest string, data string) ([]byte, error) {
	destRegister, valid := registerMap8[dest]
	if !valid {
		return nil, fmt.Errorf("invalid destination register for MOV: %s", dest)
	}

	_, lowByte, err := parseHex(data)
	if err != nil {
		return nil, err
	}

	opcode := byte(0x06) | (destRegister << 3)
	return []byte{opcode, lowByte}, nil
}

func (p *Parser) parseLXI() ([]byte, error) {
	if p.currentToken().Type != lexer.REGISTER {
		return nil, fmt.Errorf("expected register, got: %s", p.currentToken().Literal)
	}
	dest := p.currentToken().Literal
	p.advanceToken()

	if p.currentToken().Type != lexer.COMMA {
		return nil, fmt.Errorf("expected comma, got: %s", p.currentToken().Literal)
	}
	p.advanceToken()

	destRegister, valid := registerMap16[dest]
	if !valid {
		return nil, fmt.Errorf("invalid destination register for LXI: %s", dest)
	}

	if p.currentToken().Type == lexer.NUMBER {
		highByte, lowByte, err := parseHex(p.currentToken().Literal)
		if err != nil {
			return nil, err
		}

		// TODO: Reduce the duplication here
		opcode := byte(0x01) | (destRegister << 4)
		return []byte{opcode, lowByte, highByte}, nil
	}

	if p.currentToken().Type == lexer.LABEL {
		p.labelReferences[p.currentToken().Literal] = uint16(len(p.bytecode) + 1)
	}
	// TODO: Check for errors here
	// TODO: Reduce the duplication with the opcode code
	// TODO: Switch this out for generateLXIHex (if possible?)

	opcode := byte(0x01) | (destRegister << 4)
	return []byte{opcode, 0x00, 0x00}, nil
}

// func generateLXIHex(dest string, data string) ([]byte, error) {
// 	destRegister, valid := registerMap16[dest]
// 	if !valid {
// 		return nil, fmt.Errorf("invalid destination register for LXI: %s", dest)
// 	}

// 	highByte, lowByte, err := parseHex(data)
// 	if err != nil {
// 		return nil, err
// 	}

// 	opcode := byte(0x01) | (destRegister << 4)
// 	return []byte{opcode, lowByte, highByte}, nil
// }

func (p *Parser) parseSTAX() ([]byte, error) {
	if p.currentToken().Type != lexer.REGISTER {
		return nil, fmt.Errorf("expected register, got: %s", p.currentToken().Literal)
	}
	dest := p.currentToken().Literal
	p.advanceToken()

	return generateSTAXHex(dest)
}

func generateSTAXHex(dest string) ([]byte, error) {
	var registerMap = map[string]byte{
		"B": 0x00, "D": 0x01,
	}

	destRegister, valid := registerMap[dest]
	if !valid {
		return nil, fmt.Errorf("invalid destination register for STAX: %s", dest)
	}

	opcode := byte(0x02) | (destRegister << 4)
	return []byte{opcode}, nil
}

func (p *Parser) parseLDAX() ([]byte, error) {
	if p.currentToken().Type != lexer.REGISTER {
		return nil, fmt.Errorf("expected register, got: %s", p.currentToken().Literal)
	}
	dest := p.currentToken().Literal
	p.advanceToken()

	return generateLDAXHex(dest)
}

func generateLDAXHex(dest string) ([]byte, error) {
	var registerMap = map[string]byte{
		"B": 0x00, "D": 0x01,
	}

	destRegister, valid := registerMap[dest]
	if !valid {
		return nil, fmt.Errorf("invalid destination register for LDAX: %s", dest)
	}

	opcode := byte(0x0A) | (destRegister << 4)
	return []byte{opcode}, nil
}

func (p *Parser) parseJMP() ([]byte, error) {

	if p.currentToken().Type == lexer.NUMBER {
		highByte, lowByte, err := parseHex(p.currentToken().Literal)
		if err != nil {
			// TODO: Return here instead of print
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

func (p *Parser) parseDB() ([]byte, error) {
	var data []byte

	for p.currentToken().Type == lexer.NUMBER || p.currentToken().Type == lexer.STRING {
		if p.currentToken().Type == lexer.NUMBER {
			value, err := strconv.ParseUint(p.currentToken().Literal, 16, 8)
			if err != nil {
				return nil, fmt.Errorf("invalid byte value: %s", p.currentToken().Literal)
			}
			data = append(data, byte(value))
		} else if p.currentToken().Type == lexer.STRING {
			// Convert string into bytes
			data = append(data, []byte(p.currentToken().Literal)...)
		}
		p.advanceToken()

		// Handle optional commas???
		if p.currentToken().Type == lexer.COMMA {
			p.advanceToken()
		}
	}

	return data, nil
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
