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

	// MOVE, LOAD AND STORE
	"MOV":  (*Parser).parseMOV,
	"MVI":  (*Parser).parseMVI,
	"LXI":  (*Parser).parseLXI,
	"STAX": (*Parser).parseSTAX,
	"LDAX": (*Parser).parseLDAX,
	"STA":  (*Parser).parseDirectAddressInstruction,
	"LDA":  (*Parser).parseDirectAddressInstruction,
	"SHLD": (*Parser).parseDirectAddressInstruction,
	"LHLD": (*Parser).parseDirectAddressInstruction,
	"XCHG": (*Parser).parseSingleByteInstruction,

	// STACK OPERATIONS
	"XTHL": (*Parser).parseSingleByteInstruction,
	"SPHL": (*Parser).parseSingleByteInstruction,

	"PUSH": (*Parser).parseRegisterPairInstruction,
	"POP":  (*Parser).parseRegisterPairInstruction,
	"INX":  (*Parser).parseRegisterPairInstruction,
	"DCX":  (*Parser).parseRegisterPairInstruction,
	"DAD":  (*Parser).parseRegisterPairInstruction,

	// JUMP AND CALL
	"JMP":  (*Parser).parseJumpAndCallInstruction,
	"JC":   (*Parser).parseJumpAndCallInstruction,
	"JNC":  (*Parser).parseJumpAndCallInstruction,
	"JZ":   (*Parser).parseJumpAndCallInstruction,
	"JNZ":  (*Parser).parseJumpAndCallInstruction,
	"JP":   (*Parser).parseJumpAndCallInstruction,
	"JM":   (*Parser).parseJumpAndCallInstruction,
	"JPE":  (*Parser).parseJumpAndCallInstruction,
	"JPO":  (*Parser).parseJumpAndCallInstruction,
	"PCHL": (*Parser).parseSingleByteInstruction,
	"CALL": (*Parser).parseJumpAndCallInstruction,
	"CC":   (*Parser).parseJumpAndCallInstruction,
	"CNC":  (*Parser).parseJumpAndCallInstruction,
	"CZ":   (*Parser).parseJumpAndCallInstruction,
	"CNZ":  (*Parser).parseJumpAndCallInstruction,
	"CP":   (*Parser).parseJumpAndCallInstruction,
	"CM":   (*Parser).parseJumpAndCallInstruction,
	"CPE":  (*Parser).parseJumpAndCallInstruction,
	"CPO":  (*Parser).parseJumpAndCallInstruction,

	// RETURN
	"RET": (*Parser).parseSingleByteInstruction,
	"RC":  (*Parser).parseSingleByteInstruction,
	"RNC": (*Parser).parseSingleByteInstruction,
	"RZ":  (*Parser).parseSingleByteInstruction,
	"RNZ": (*Parser).parseSingleByteInstruction,
	"RP":  (*Parser).parseSingleByteInstruction,
	"RM":  (*Parser).parseSingleByteInstruction,
	"RPE": (*Parser).parseSingleByteInstruction,
	"RPO": (*Parser).parseSingleByteInstruction,

	// RESTART
	"RST": (*Parser).parseRestartInstruction,

	"DB": (*Parser).parseDB,
}

var registerMap8 = map[string]byte{
	"B": 0x00, "C": 0x01, "D": 0x02, "E": 0x03, "H": 0x04, "L": 0x05, "M": 0x06, "A": 0x07,
}

var registerMap16 = map[string]byte{
	"B": 0x00, "D": 0x01, "H": 0x02, "SP": 0x03,
}

func (p *Parser) parseInstruction() ([]byte, error) {
	instruction := p.currentToken().Literal
	parseFunc, instructionExists := instructionMap[instruction]
	if !instructionExists {
		return nil, fmt.Errorf("unknown instruction: %s", instruction)
	}
	return parseFunc(p)
}

func (p *Parser) parseMOV() ([]byte, error) {
	p.advanceToken()

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
	srcRegister, exists := registerMap8[src]
	if !exists {
		return nil, fmt.Errorf("invalid source register for MOV: %s", src)
	}

	destRegister, exists := registerMap8[dest]
	if !exists {
		return nil, fmt.Errorf("invalid destination register for MOV: %s", dest)
	}

	opcode := byte(0x40) | (destRegister << 3) | srcRegister
	return []byte{opcode}, nil
}

func (p *Parser) parseMVI() ([]byte, error) {
	p.advanceToken()

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
	destRegister, exists := registerMap8[dest]
	if !exists {
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
	p.advanceToken()

	if p.currentToken().Type != lexer.REGISTER {
		return nil, fmt.Errorf("expected register, got: %s", p.currentToken().Literal)
	}
	dest := p.currentToken().Literal
	p.advanceToken()

	if p.currentToken().Type != lexer.COMMA {
		return nil, fmt.Errorf("expected comma, got: %s", p.currentToken().Literal)
	}
	p.advanceToken()

	destRegister, exists := registerMap16[dest]
	if !exists {
		return nil, fmt.Errorf("invalid destination register for LXI: %s", dest)
	}
	opcode := byte(0x01) | (destRegister << 4)

	if p.currentToken().Type == lexer.NUMBER {
		highByte, lowByte, err := parseHex(p.currentToken().Literal)
		if err != nil {
			return nil, err
		}
		return []byte{opcode, lowByte, highByte}, nil
	}

	if p.currentToken().Type == lexer.LABEL {
		p.labelReferences[p.currentToken().Literal] = uint16(len(p.bytecode) + 1)
		return []byte{opcode, 0x00, 0x00}, nil
	}

	return nil, fmt.Errorf("expected address or label, got: %s", p.currentToken().Type)
}

// TODO: Combine STAX and LDAX instructions
func (p *Parser) parseSTAX() ([]byte, error) {
	p.advanceToken()

	if p.currentToken().Type != lexer.REGISTER {
		return nil, fmt.Errorf("expected register, got: %s", p.currentToken().Literal)
	}
	dest := p.currentToken().Literal
	p.advanceToken()

	return generateSTAXHex(dest)
}

func generateSTAXHex(dest string) ([]byte, error) {
	registerMap := map[string]byte{
		"B": 0x00, "D": 0x01,
	}

	destRegister, exists := registerMap[dest]
	if !exists {
		return nil, fmt.Errorf("invalid destination register for STAX: %s", dest)
	}

	opcode := byte(0x02) | (destRegister << 4)
	return []byte{opcode}, nil
}

func (p *Parser) parseLDAX() ([]byte, error) {
	p.advanceToken()

	if p.currentToken().Type != lexer.REGISTER {
		return nil, fmt.Errorf("expected register, got: %s", p.currentToken().Literal)
	}
	dest := p.currentToken().Literal
	p.advanceToken()

	return generateLDAXHex(dest)
}

func generateLDAXHex(dest string) ([]byte, error) {
	registerMap := map[string]byte{
		"B": 0x00, "D": 0x01,
	}

	destRegister, exists := registerMap[dest]
	if !exists {
		return nil, fmt.Errorf("invalid destination register for LDAX: %s", dest)
	}

	opcode := byte(0x0A) | (destRegister << 4)
	return []byte{opcode}, nil
}

func (p *Parser) parseDirectAddressInstruction() ([]byte, error) {
	opcodes := map[string]byte{
		"STA":  0x32,
		"LDA":  0x3A,
		"SHLD": 0x22,
		"LHLD": 0x2A,
	}

	opcode, valid := opcodes[p.currentToken().Literal]
	if !valid {
		return nil, fmt.Errorf("invalid direct address instruction: %s, ", p.currentToken().Literal)
	}
	p.advanceToken()

	if p.currentToken().Type == lexer.NUMBER {
		highByte, lowByte, err := parseHex(p.currentToken().Literal)
		if err != nil {
			return nil, err
		}
		return []byte{opcode, lowByte, highByte}, nil
	}

	if p.currentToken().Type == lexer.LABEL {
		p.labelReferences[p.currentToken().Literal] = uint16(len(p.bytecode) + 1)
		return []byte{opcode, 0x00, 0x00}, nil
	}

	return nil, fmt.Errorf("expected address or label, got: %s", p.currentToken().Type)
}

func (p *Parser) parseSingleByteInstruction() ([]byte, error) {
	opcodes := map[string]byte{
		"XCHG": 0xEB,
		"XTHL": 0xE3,
		"SPHL": 0xF9,
		"PCHL": 0xE9,
		"RET":  0xC9,
		"RC":   0xD8,
		"RNC":  0xD0,
		"RZ":   0xC8,
		"RNZ":  0xC0,
		"RP":   0xF0,
		"RM":   0xF8,
		"RPE":  0xE8,
		"RPO":  0xE0,
	}

	opcode, valid := opcodes[p.currentToken().Literal]
	if !valid {
		return nil, fmt.Errorf("invalid instruction: %s, ", p.currentToken().Literal)
	}
	p.advanceToken()

	return []byte{opcode}, nil
}

func (p *Parser) parseRegisterPairInstruction() ([]byte, error) {
	opcodes := map[string]byte{
		"PUSH": 0xC5,
		"POP":  0xC1,
		"INX":  0x03,
		"DCX":  0x0B,
		"DAD":  0x09,
	}

	registerMap := map[string]byte{
		"B": 0x00, "D": 0x01, "H": 0x02, "PSW": 0x03, "SP": 0x03,
	}
	opcode, valid := opcodes[p.currentToken().Literal]
	if !valid {
		return nil, fmt.Errorf("invalid instruction: %s, ", p.currentToken().Literal)
	}
	p.advanceToken()

	if p.currentToken().Type != lexer.REGISTER {
		return nil, fmt.Errorf("expected register, got: %s", p.currentToken().Literal)
	}
	dest := p.currentToken().Literal
	p.advanceToken()

	destRegister, exists := registerMap[dest]
	if !exists {
		return nil, fmt.Errorf("invalid destination register for %s: %s", p.currentToken().Literal, dest)
	}

	opcode = opcode | (destRegister << 4)
	return []byte{opcode}, nil
}

func (p *Parser) parseJumpAndCallInstruction() ([]byte, error) {
	opcodes := map[string]byte{
		"JMP":  0xC3,
		"JC":   0xDA,
		"JNC":  0xD2,
		"JZ":   0xCA,
		"JNZ":  0xC2,
		"JP":   0xF2,
		"JM":   0xFA,
		"JPE":  0xEA,
		"JPO":  0xE2,
		"CALL": 0xCD,
		"CC":   0xDC,
		"CNC":  0xD4,
		"CZ":   0xCC,
		"CNZ":  0xC4,
		"CP":   0xF4,
		"CM":   0xFC,
		"CPE":  0xEC,
		"CPO":  0xE4,
	}

	opcode, valid := opcodes[p.currentToken().Literal]
	if !valid {
		return nil, fmt.Errorf("invalid instruction: %s, ", p.currentToken().Literal)
	}
	p.advanceToken()

	if p.currentToken().Type == lexer.NUMBER {
		highByte, lowByte, err := parseHex(p.currentToken().Literal)
		if err != nil {
			return nil, err
		}
		return []byte{opcode, lowByte, highByte}, nil
	}

	if p.currentToken().Type == lexer.LABEL {
		p.labelReferences[p.currentToken().Literal] = uint16(len(p.bytecode) + 1)
		return []byte{opcode, 0x00, 0x00}, nil
	}

	return nil, fmt.Errorf("expected address or label, got: %s", p.currentToken().Type)
}

func (p *Parser) parseRestartInstruction() ([]byte, error) {
	p.advanceToken()

	if p.currentToken().Type != lexer.NUMBER {
		return nil, fmt.Errorf("expected number, got: %s", p.currentToken().Literal)
	}

	routine, err := strconv.ParseUint(p.currentToken().Literal, 16, 8)
	if err != nil || routine > 7 {
		return nil, fmt.Errorf("invalid routine value: %s", p.currentToken().Literal)
	}

	opcode := byte(0xC7 + routine<<3)
	return []byte{opcode}, nil
}

func (p *Parser) parseDB() ([]byte, error) {
	p.advanceToken()

	data := []byte{}

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
