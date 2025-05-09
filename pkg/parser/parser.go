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
	labelDefinitions map[string]uint16   // Stores resolved label addresses
	labelReferences  map[string][]uint16 // Tracks unresolved label usages
}

func New(tokens []lexer.Token) *Parser {
	return &Parser{
		tokens:           tokens,
		position:         0,
		labelDefinitions: make(map[string]uint16),
		labelReferences:  make(map[string][]uint16),
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
			// comments aren't assembled, so we simply skip the token

		case lexer.LABEL:
			label := p.currentToken().Literal
			_, labelExists := p.labelDefinitions[label]
			if labelExists {
				return nil, fmt.Errorf("duplicate label found: %s", label)
			}
			p.labelDefinitions[label] = uint16(len(p.bytecode))
			p.advanceToken()

		default:
			return nil, fmt.Errorf("unexpected token type \"%s\", literal: \"%s\"", p.currentToken().Type, p.currentToken().Literal)
		}

		p.advanceToken()
	}

	// TODO: Split this out into a second pass
	for label, positions := range p.labelReferences {
		targetAddr, targetExists := p.labelDefinitions[label]
		if targetExists {
			highByte := uint8(targetAddr >> 8)
			lowByte := uint8(targetAddr & 0x00FF)

			// Update all instances of the label reference
			for _, position := range positions {
				p.bytecode[position] = lowByte
				p.bytecode[position+1] = highByte
			}
		} else {
			return nil, fmt.Errorf("label definition not found: %s", label)
		}
	}

	return p.bytecode, nil
}

type parseFunc func(*Parser) ([]byte, error)

var instructionMap = map[string]parseFunc{

	// MOVE, LOAD AND STORE
	"MOV":  (*Parser).parseMOV,
	"MVI":  (*Parser).parseMVI,
	"LXI":  (*Parser).parseLXI,
	"STAX": (*Parser).parseSTAXandLDAX,
	"LDAX": (*Parser).parseSTAXandLDAX,
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
	"JMP":  (*Parser).parseDirectAddressInstruction,
	"JC":   (*Parser).parseDirectAddressInstruction,
	"JNC":  (*Parser).parseDirectAddressInstruction,
	"JZ":   (*Parser).parseDirectAddressInstruction,
	"JNZ":  (*Parser).parseDirectAddressInstruction,
	"JP":   (*Parser).parseDirectAddressInstruction,
	"JM":   (*Parser).parseDirectAddressInstruction,
	"JPE":  (*Parser).parseDirectAddressInstruction,
	"JPO":  (*Parser).parseDirectAddressInstruction,
	"PCHL": (*Parser).parseSingleByteInstruction,
	"CALL": (*Parser).parseDirectAddressInstruction,
	"CC":   (*Parser).parseDirectAddressInstruction,
	"CNC":  (*Parser).parseDirectAddressInstruction,
	"CZ":   (*Parser).parseDirectAddressInstruction,
	"CNZ":  (*Parser).parseDirectAddressInstruction,
	"CP":   (*Parser).parseDirectAddressInstruction,
	"CM":   (*Parser).parseDirectAddressInstruction,
	"CPE":  (*Parser).parseDirectAddressInstruction,
	"CPO":  (*Parser).parseDirectAddressInstruction,

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

	// INCREMENT AND DECREMENT
	"INR": (*Parser).parseRegister8Instruction,
	"DCR": (*Parser).parseRegister8Instruction,

	// ADD AND SUBTRACT
	"ADD": (*Parser).parseRegister8Instruction,
	"ADC": (*Parser).parseRegister8Instruction,
	"ADI": (*Parser).parseImmediateInstruction,
	"ACI": (*Parser).parseImmediateInstruction,
	"SUB": (*Parser).parseRegister8Instruction,
	"SBB": (*Parser).parseRegister8Instruction,
	"SUI": (*Parser).parseImmediateInstruction,
	"SBI": (*Parser).parseImmediateInstruction,
	"ANI": (*Parser).parseImmediateInstruction,
	"XRI": (*Parser).parseImmediateInstruction,
	"ORI": (*Parser).parseImmediateInstruction,
	"CPI": (*Parser).parseImmediateInstruction,

	// LOGICAL
	"ANA": (*Parser).parseRegister8Instruction,
	"XRA": (*Parser).parseRegister8Instruction,
	"ORA": (*Parser).parseRegister8Instruction,
	"CMP": (*Parser).parseRegister8Instruction,

	// ROTATE
	"RLC": (*Parser).parseSingleByteInstruction,
	"RRC": (*Parser).parseSingleByteInstruction,
	"RAL": (*Parser).parseSingleByteInstruction,
	"RAR": (*Parser).parseSingleByteInstruction,

	// SPECIALS
	"CMA": (*Parser).parseSingleByteInstruction,
	"STC": (*Parser).parseSingleByteInstruction,
	"CMC": (*Parser).parseSingleByteInstruction,
	"DAA": (*Parser).parseSingleByteInstruction,

	// INPUT/OUTPUT
	"IN":  (*Parser).parseImmediateInstruction,
	"OUT": (*Parser).parseImmediateInstruction,

	// CONTROL
	"EI":  (*Parser).parseSingleByteInstruction,
	"DI":  (*Parser).parseSingleByteInstruction,
	"NOP": (*Parser).parseSingleByteInstruction,
	"HLT": (*Parser).parseSingleByteInstruction,

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

	destRegister, exists := registerMap8[dest]
	if !exists {
		return nil, fmt.Errorf("invalid destination register for MOV: %s", dest)
	}

	_, lowByte, err := parseHex(byteToStore)
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
		label := p.currentToken().Literal
		_, labelExists := p.labelReferences[label]
		if !labelExists {
			p.labelReferences[label] = []uint16{}
		}
		p.labelReferences[label] = append(p.labelReferences[label], uint16(len(p.bytecode)+1))

		return []byte{opcode, 0x00, 0x00}, nil
	}

	return nil, fmt.Errorf("expected address or label, got: %s", p.currentToken().Type)
}

func (p *Parser) parseSTAXandLDAX() ([]byte, error) {
	opcodes := map[string]byte{
		"STAX": 0x02,
		"LDAX": 0x0A,
	}

	opcode, valid := opcodes[p.currentToken().Literal]
	if !valid {
		return nil, fmt.Errorf("invalid STAX/LDAX instruction: %s", p.currentToken().Literal)
	}
	p.advanceToken()

	if p.currentToken().Type != lexer.REGISTER {
		return nil, fmt.Errorf("expected register, got: %s", p.currentToken().Literal)
	}
	dest := p.currentToken().Literal

	registerMap := map[string]byte{
		"B": 0x00, "D": 0x01,
	}

	destRegister, exists := registerMap[dest]
	if !exists {
		return nil, fmt.Errorf("invalid destination register for STAX: %s", dest)
	}

	opcode = opcode | (destRegister << 4)
	return []byte{opcode}, nil
}

func (p *Parser) parseDirectAddressInstruction() ([]byte, error) {
	opcodes := map[string]byte{
		"STA":  0x32,
		"LDA":  0x3A,
		"SHLD": 0x22,
		"LHLD": 0x2A,
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
		return nil, fmt.Errorf("invalid direct address instruction: %s", p.currentToken().Literal)
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
		label := p.currentToken().Literal
		_, labelExists := p.labelReferences[label]
		if !labelExists {
			p.labelReferences[label] = []uint16{}
		}
		p.labelReferences[label] = append(p.labelReferences[label], uint16(len(p.bytecode)+1))
		return []byte{opcode, 0x00, 0x00}, nil
	}

	return nil, fmt.Errorf("expected address or label, got: %s", p.currentToken().Type)
}

func (p *Parser) parseImmediateInstruction() ([]byte, error) {
	opcodes := map[string]byte{
		"ADI": 0xC6,
		"ACI": 0xCE,
		"SUI": 0xD6,
		"SBI": 0xDE,
		"ANI": 0xE6,
		"XRI": 0xEE,
		"ORI": 0xF6,
		"CPI": 0xFE,
		"IN":  0xDB,
		"OUT": 0xD3,
	}

	opcode, valid := opcodes[p.currentToken().Literal]
	if !valid {
		return nil, fmt.Errorf("invalid immediate instruction: %s, ", p.currentToken().Literal)
	}
	p.advanceToken()

	if p.currentToken().Type == lexer.NUMBER {
		highByte, lowByte, err := parseHex(p.currentToken().Literal)
		if err != nil {
			return nil, err
		}

		if highByte != 0x00 {
			return nil, fmt.Errorf("expected single byte of data, got: %s", p.currentToken().Literal)
		}

		return []byte{opcode, lowByte}, nil
	}

	return nil, fmt.Errorf("expected number, got: %s", p.currentToken().Type)
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
		"RLC":  0x07,
		"RRC":  0x0F,
		"RAL":  0x17,
		"RAR":  0x1F,
		"CMA":  0x2F,
		"STC":  0x37,
		"CMC":  0x3F,
		"DAA":  0x27,
		"EI":   0xFB,
		"DI":   0xF3,
		"NOP":  0x00,
		"HLT":  0x76,
	}

	opcode, valid := opcodes[p.currentToken().Literal]
	if !valid {
		return nil, fmt.Errorf("invalid instruction: %s, ", p.currentToken().Literal)
	}

	return []byte{opcode}, nil
}

// TODO: make this parseRegister16Instruction?
// Need to work out how to deal with PSW.
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

	destRegister, exists := registerMap[dest]
	if !exists {
		return nil, fmt.Errorf("invalid destination register for %s: %s", p.currentToken().Literal, dest)
	}

	opcode = opcode | (destRegister << 4)
	return []byte{opcode}, nil
}

func (p *Parser) parseRestartInstruction() ([]byte, error) {
	p.advanceToken()

	if p.currentToken().Type != lexer.NUMBER {
		return nil, fmt.Errorf("expected number, got: %s", p.currentToken().Literal)
	}

	routine, err := strconv.ParseUint(p.currentToken().Literal, 16, 8)
	if err != nil || routine > 7 {
		return nil, fmt.Errorf("expected routine value between 0 and 7, got: %s", p.currentToken().Literal)
	}

	opcode := byte(0xC7 + routine<<3)
	return []byte{opcode}, nil
}

func (p *Parser) parseRegister8Instruction() ([]byte, error) {
	opcodes := map[string]byte{
		"INR": 0x04,
		"DCR": 0x05,
		"ADD": 0x80,
		"ADC": 0x88,
		"SUB": 0x90,
		"SBB": 0x98,
		"ANA": 0xA0,
		"XRA": 0xA8,
		"ORA": 0xB0,
		"CMP": 0xB8,
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

	destRegister, exists := registerMap8[dest]
	if !exists {
		return nil, fmt.Errorf("invalid destination register for %s: %s", p.currentToken().Literal, dest)
	}

	// TODO: Find a better way to do this
	// If INR or DCR, set shift to 3
	shiftAmount := 0
	if opcode == 0x04 || opcode == 0x05 {
		shiftAmount = 3
	}

	opcode = opcode | (destRegister << shiftAmount)
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
