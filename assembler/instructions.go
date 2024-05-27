package assembler

type Instruction struct {
	Opcode byte
	Length int
}

// Grouped by instruction set group as per "Table 2. Instruction Set Summary",
// in the Intel 8080A 8-BIT N-CHANNEL MICROPROCESSOR datasheet.
// Opcode = the 1 byte hexadecimal representation of the instruction for the 8080
// Length = how many bytes in total the instruction has (more succintly, how many n-1 bytes _after_
// the opcode to interpret as the operand)
var instructionSet = map[string]Instruction{
	// MOVE, LOAD AND STORE
	"MOV A,B": {Opcode: 0x78, Length: 1},
	"MOV A,C": {Opcode: 0x79, Length: 1},
	"MOV A,D": {Opcode: 0x7A, Length: 1},
	"MOV A,E": {Opcode: 0x7B, Length: 1},
	"MOV A,H": {Opcode: 0x7C, Length: 1},
	"MOV A,L": {Opcode: 0x7D, Length: 1},
	"MOV A,M": {Opcode: 0x7E, Length: 1},
	"MOV A,A": {Opcode: 0x7F, Length: 1},

	"MOV B,B": {Opcode: 0x40, Length: 1},
	"MOV B,C": {Opcode: 0x41, Length: 1},
	"MOV B,D": {Opcode: 0x42, Length: 1},
	"MOV B,E": {Opcode: 0x43, Length: 1},
	"MOV B,H": {Opcode: 0x44, Length: 1},
	"MOV B,L": {Opcode: 0x45, Length: 1},
	"MOV B,M": {Opcode: 0x46, Length: 1},
	"MOV B,A": {Opcode: 0x47, Length: 1},

	"MOV C,B": {Opcode: 0x48, Length: 1},
	"MOV C,C": {Opcode: 0x49, Length: 1},
	"MOV C,D": {Opcode: 0x4A, Length: 1},
	"MOV C,E": {Opcode: 0x4B, Length: 1},
	"MOV C,H": {Opcode: 0x4C, Length: 1},
	"MOV C,L": {Opcode: 0x4D, Length: 1},
	"MOV C,M": {Opcode: 0x4E, Length: 1},
	"MOV C,A": {Opcode: 0x4F, Length: 1},

	"MOV D,B": {Opcode: 0x50, Length: 1},
	"MOV D,C": {Opcode: 0x51, Length: 1},
	"MOV D,D": {Opcode: 0x52, Length: 1},
	"MOV D,E": {Opcode: 0x53, Length: 1},
	"MOV D,H": {Opcode: 0x54, Length: 1},
	"MOV D,L": {Opcode: 0x55, Length: 1},
	"MOV D,M": {Opcode: 0x56, Length: 1},
	"MOV D,A": {Opcode: 0x57, Length: 1},

	"MOV E,B": {Opcode: 0x58, Length: 1},
	"MOV E,C": {Opcode: 0x59, Length: 1},
	"MOV E,D": {Opcode: 0x5A, Length: 1},
	"MOV E,E": {Opcode: 0x5B, Length: 1},
	"MOV E,H": {Opcode: 0x5C, Length: 1},
	"MOV E,L": {Opcode: 0x5D, Length: 1},
	"MOV E,M": {Opcode: 0x5E, Length: 1},
	"MOV E,A": {Opcode: 0x5F, Length: 1},

	"MOV H,B": {Opcode: 0x60, Length: 1},
	"MOV H,C": {Opcode: 0x61, Length: 1},
	"MOV H,D": {Opcode: 0x62, Length: 1},
	"MOV H,E": {Opcode: 0x63, Length: 1},
	"MOV H,H": {Opcode: 0x64, Length: 1},
	"MOV H,L": {Opcode: 0x65, Length: 1},
	"MOV H,M": {Opcode: 0x66, Length: 1},
	"MOV H,A": {Opcode: 0x67, Length: 1},

	"MOV L,B": {Opcode: 0x68, Length: 1},
	"MOV L,C": {Opcode: 0x69, Length: 1},
	"MOV L,D": {Opcode: 0x6A, Length: 1},
	"MOV L,E": {Opcode: 0x6B, Length: 1},
	"MOV L,H": {Opcode: 0x6C, Length: 1},
	"MOV L,L": {Opcode: 0x6D, Length: 1},
	"MOV L,M": {Opcode: 0x6E, Length: 1},
	"MOV L,A": {Opcode: 0x6F, Length: 1},

	"MOV M,B": {Opcode: 0x70, Length: 1},
	"MOV M,C": {Opcode: 0x71, Length: 1},
	"MOV M,D": {Opcode: 0x72, Length: 1},
	"MOV M,E": {Opcode: 0x73, Length: 1},
	"MOV M,H": {Opcode: 0x74, Length: 1},
	"MOV M,L": {Opcode: 0x75, Length: 1},
	"MOV M,A": {Opcode: 0x77, Length: 1},

	"MVI B": {Opcode: 0x06, Length: 2},
	"MVI C": {Opcode: 0x0E, Length: 2},
	"MVI D": {Opcode: 0x16, Length: 2},
	"MVI E": {Opcode: 0x1E, Length: 2},
	"MVI H": {Opcode: 0x26, Length: 2},
	"MVI L": {Opcode: 0x2E, Length: 2},
	"MVI M": {Opcode: 0x36, Length: 2},
	"MVI A": {Opcode: 0x3E, Length: 2},

	"LXI B":  {Opcode: 0x01, Length: 3},
	"LXI D":  {Opcode: 0x11, Length: 3},
	"LXI H":  {Opcode: 0x21, Length: 3},
	"STAX B": {Opcode: 0x02, Length: 1},
	"STAX D": {Opcode: 0x12, Length: 1},
	"LDAX B": {Opcode: 0x0A, Length: 1},
	"LDAX D": {Opcode: 0x1A, Length: 1},
	"STA":    {Opcode: 0x32, Length: 3},
	"LDA":    {Opcode: 0x3A, Length: 3},
	"SHLD":   {Opcode: 0x22, Length: 3},
	"LHLD":   {Opcode: 0x2A, Length: 3},
	"XCHG":   {Opcode: 0xEB, Length: 1},

	// STACK OPERATIONS
	"PUSH B":   {Opcode: 0xC5, Length: 1},
	"PUSH D":   {Opcode: 0xD5, Length: 1},
	"PUSH H":   {Opcode: 0xE5, Length: 1},
	"PUSH PSW": {Opcode: 0xF5, Length: 1},
	"POP B":    {Opcode: 0xC1, Length: 1},
	"POP D":    {Opcode: 0xD1, Length: 1},
	"POP H":    {Opcode: 0xE1, Length: 1},
	"POP PSW":  {Opcode: 0xF1, Length: 1},
	"XTHL":     {Opcode: 0xE3, Length: 1},
	"SPHL":     {Opcode: 0xF9, Length: 1},
	"LXI SP":   {Opcode: 0x31, Length: 3},
	"INX SP":   {Opcode: 0x33, Length: 1},
	"DCX SP":   {Opcode: 0x3B, Length: 1},
	"DAD SP":   {Opcode: 0x39, Length: 1},

	// JUMP
	"JMP":  {Opcode: 0xC3, Length: 3},
	"JC":   {Opcode: 0xDA, Length: 3},
	"JNC":  {Opcode: 0xD2, Length: 3},
	"JZ":   {Opcode: 0xCA, Length: 3},
	"JNZ":  {Opcode: 0xC2, Length: 3},
	"JP":   {Opcode: 0xF2, Length: 3},
	"JM":   {Opcode: 0xFA, Length: 3},
	"JPE":  {Opcode: 0xEA, Length: 3},
	"JPO":  {Opcode: 0xE2, Length: 3},
	"PCHL": {Opcode: 0xE9, Length: 1},

	// CALL
	"CALL": {Opcode: 0xCD, Length: 3},
	"CC":   {Opcode: 0xDC, Length: 3},
	"CNC":  {Opcode: 0xD4, Length: 3},
	"CZ":   {Opcode: 0xCC, Length: 3},
	"CNZ":  {Opcode: 0xC4, Length: 3},
	"CP":   {Opcode: 0xF4, Length: 3},
	"CM":   {Opcode: 0xFC, Length: 3},
	"CPE":  {Opcode: 0xEC, Length: 3},
	"CPO":  {Opcode: 0xE4, Length: 3},

	// RETURN
	"RET": {Opcode: 0xC9, Length: 1},
	"RC":  {Opcode: 0xD8, Length: 1},
	"RNC": {Opcode: 0xD0, Length: 1},
	"RZ":  {Opcode: 0xC8, Length: 1},
	"RNZ": {Opcode: 0xC0, Length: 1},
	"RP":  {Opcode: 0xF0, Length: 1},
	"RM":  {Opcode: 0xF8, Length: 1},
	"RPE": {Opcode: 0xE8, Length: 1},
	"RPO": {Opcode: 0xE0, Length: 1},

	// RESTART
	"RST 0": {Opcode: 0xC7, Length: 1},
	"RST 1": {Opcode: 0xCF, Length: 1},
	"RST 2": {Opcode: 0xD7, Length: 1},
	"RST 3": {Opcode: 0xDF, Length: 1},
	"RST 4": {Opcode: 0xE7, Length: 1},
	"RST 5": {Opcode: 0xEF, Length: 1},
	"RST 6": {Opcode: 0xF7, Length: 1},
	"RST 7": {Opcode: 0xFF, Length: 1},

	// INCREMENT AND DECREMENT
	"INR B": {Opcode: 0x04, Length: 1},
	"INR C": {Opcode: 0x0C, Length: 1},
	"INR D": {Opcode: 0x14, Length: 1},
	"INR E": {Opcode: 0x1C, Length: 1},
	"INR H": {Opcode: 0x24, Length: 1},
	"INR L": {Opcode: 0x2C, Length: 1},
	"INR M": {Opcode: 0x34, Length: 1},
	"INR A": {Opcode: 0x3C, Length: 1},

	"DCR B": {Opcode: 0x05, Length: 1},
	"DCR C": {Opcode: 0x0D, Length: 1},
	"DCR D": {Opcode: 0x15, Length: 1},
	"DCR E": {Opcode: 0x1D, Length: 1},
	"DCR H": {Opcode: 0x25, Length: 1},
	"DCR L": {Opcode: 0x2D, Length: 1},
	"DCR M": {Opcode: 0x35, Length: 1},
	"DCR A": {Opcode: 0x3D, Length: 1},

	"INX B": {Opcode: 0x03, Length: 1},
	"INX D": {Opcode: 0x13, Length: 1},
	"INX H": {Opcode: 0x23, Length: 1},

	"DCX B": {Opcode: 0x0B, Length: 1},
	"DCX D": {Opcode: 0x1B, Length: 1},
	"DCX H": {Opcode: 0x2B, Length: 1},

	// ADD
	"ADD B": {Opcode: 0x80, Length: 1},
	"ADD C": {Opcode: 0x81, Length: 1},
	"ADD D": {Opcode: 0x82, Length: 1},
	"ADD E": {Opcode: 0x83, Length: 1},
	"ADD H": {Opcode: 0x84, Length: 1},
	"ADD L": {Opcode: 0x85, Length: 1},
	"ADD M": {Opcode: 0x86, Length: 1},
	"ADD A": {Opcode: 0x87, Length: 1},

	"ADC B": {Opcode: 0x88, Length: 1},
	"ADC C": {Opcode: 0x89, Length: 1},
	"ADC D": {Opcode: 0x8A, Length: 1},
	"ADC E": {Opcode: 0x8B, Length: 1},
	"ADC H": {Opcode: 0x8C, Length: 1},
	"ADC L": {Opcode: 0x8D, Length: 1},
	"ADC M": {Opcode: 0x8E, Length: 1},
	"ADC A": {Opcode: 0x8F, Length: 1},

	"ADI": {Opcode: 0xC6, Length: 2},
	"ACI": {Opcode: 0xCE, Length: 2},

	"DAD B": {Opcode: 0x09, Length: 1},
	"DAD D": {Opcode: 0x19, Length: 1},
	"DAD H": {Opcode: 0x29, Length: 1},

	// SUBTRACT
	"SUB B": {Opcode: 0x90, Length: 1},
	"SUB C": {Opcode: 0x91, Length: 1},
	"SUB D": {Opcode: 0x92, Length: 1},
	"SUB E": {Opcode: 0x93, Length: 1},
	"SUB H": {Opcode: 0x94, Length: 1},
	"SUB L": {Opcode: 0x95, Length: 1},
	"SUB M": {Opcode: 0x96, Length: 1},
	"SUB A": {Opcode: 0x97, Length: 1},

	"SBB B": {Opcode: 0x98, Length: 1},
	"SBB C": {Opcode: 0x99, Length: 1},
	"SBB D": {Opcode: 0x9A, Length: 1},
	"SBB E": {Opcode: 0x9B, Length: 1},
	"SBB H": {Opcode: 0x9C, Length: 1},
	"SBB L": {Opcode: 0x9D, Length: 1},
	"SBB M": {Opcode: 0x9E, Length: 1},
	"SBB A": {Opcode: 0x9F, Length: 1},

	"SUI": {Opcode: 0xD6, Length: 2},
	"SBI": {Opcode: 0xDE, Length: 2},

	// LOGICAL
	"ANA B": {Opcode: 0xA0, Length: 1},
	"ANA C": {Opcode: 0xA1, Length: 1},
	"ANA D": {Opcode: 0xA2, Length: 1},
	"ANA E": {Opcode: 0xA3, Length: 1},
	"ANA H": {Opcode: 0xA4, Length: 1},
	"ANA L": {Opcode: 0xA5, Length: 1},
	"ANA M": {Opcode: 0xA6, Length: 1},
	"ANA A": {Opcode: 0xA7, Length: 1},

	"XRA B": {Opcode: 0xA8, Length: 1},
	"XRA C": {Opcode: 0xA9, Length: 1},
	"XRA D": {Opcode: 0xAA, Length: 1},
	"XRA E": {Opcode: 0xAB, Length: 1},
	"XRA H": {Opcode: 0xAC, Length: 1},
	"XRA L": {Opcode: 0xAD, Length: 1},
	"XRA M": {Opcode: 0xAE, Length: 1},
	"XRA A": {Opcode: 0xAF, Length: 1},

	"ORA B": {Opcode: 0xB0, Length: 1},
	"ORA C": {Opcode: 0xB1, Length: 1},
	"ORA D": {Opcode: 0xB2, Length: 1},
	"ORA E": {Opcode: 0xB3, Length: 1},
	"ORA H": {Opcode: 0xB4, Length: 1},
	"ORA L": {Opcode: 0xB5, Length: 1},
	"ORA M": {Opcode: 0xB6, Length: 1},
	"ORA A": {Opcode: 0xB7, Length: 1},

	"CMP B": {Opcode: 0xB8, Length: 1},
	"CMP C": {Opcode: 0xB9, Length: 1},
	"CMP D": {Opcode: 0xBA, Length: 1},
	"CMP E": {Opcode: 0xBB, Length: 1},
	"CMP H": {Opcode: 0xBC, Length: 1},
	"CMP L": {Opcode: 0xBD, Length: 1},
	"CMP M": {Opcode: 0xBE, Length: 1},
	"CMP A": {Opcode: 0xBF, Length: 1},

	"ANI": {Opcode: 0xE6, Length: 2},
	"XRI": {Opcode: 0xEE, Length: 2},
	"ORI": {Opcode: 0xF6, Length: 2},
	"CPI": {Opcode: 0xFE, Length: 2},

	// ROTATE
	"RLC": {Opcode: 0x07, Length: 1},
	"RRC": {Opcode: 0x0F, Length: 1},
	"RAL": {Opcode: 0x17, Length: 1},
	"RAR": {Opcode: 0x1F, Length: 1},

	// SPECIALS
	"CMA": {Opcode: 0x2F, Length: 1},
	"STC": {Opcode: 0x37, Length: 1},
	"CMC": {Opcode: 0x3F, Length: 1},
	"DAA": {Opcode: 0x27, Length: 1},

	// INPUT/OUTPUT
	"IN":  {Opcode: 0xDB, Length: 2},
	"OUT": {Opcode: 0xD3, Length: 2},

	// CONTROL
	"EI":  {Opcode: 0xFB, Length: 1},
	"DI":  {Opcode: 0xF3, Length: 1},
	"NOP": {Opcode: 0x00, Length: 1},
	"HLT": {Opcode: 0x76, Length: 1},
}
