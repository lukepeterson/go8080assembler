package assembler

type Instruction struct {
	Opcode byte
	Length int
}

var instructionSet = map[string]Instruction{
	"NOP": {Opcode: 0x00, Length: 1},

	"MVI B": {Opcode: 0x06, Length: 2},
	"MVI C": {Opcode: 0x0E, Length: 2},
	"MVI D": {Opcode: 0x16, Length: 2},
	"MVI E": {Opcode: 0x1E, Length: 2},
	"MVI H": {Opcode: 0x26, Length: 2},
	"MVI L": {Opcode: 0x2E, Length: 2},
	"MVI M": {Opcode: 0x36, Length: 2},
	"MVI A": {Opcode: 0x3E, Length: 2},

	"MOV B,B": {Opcode: 0x40, Length: 1},
	"MOV B,C": {Opcode: 0x41, Length: 1},
	"MOV B,D": {Opcode: 0x42, Length: 1},
	"MOV B,E": {Opcode: 0x43, Length: 1},
	"MOV B,H": {Opcode: 0x44, Length: 1},
	"MOV B,L": {Opcode: 0x45, Length: 1},
	"MOV B,M": {Opcode: 0x46, Length: 1},
	"MOV B,A": {Opcode: 0x47, Length: 1},

	"LDA": {Opcode: 0x3A, Length: 3},

	"HLT": {Opcode: 0x4D, Length: 1},
}
