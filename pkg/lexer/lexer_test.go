package lexer

import (
	"reflect"
	"testing"
)

func TestLexer_Lex(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []Token
		wantErr bool
	}{
		{
			name:  "case sensitivity",
			input: "mov",
			want: []Token{
				{Type: MNEMONIC, Literal: "MOV"},
				{Type: EOF},
			},
		},
		{
			name:  "single byte instruction without space",
			input: "HLT",
			want: []Token{
				{Type: MNEMONIC, Literal: "HLT"},
				{Type: EOF},
			},
		},
		{
			name:  "single byte instruction with space",
			input: "  HLT  ",
			want: []Token{
				{Type: MNEMONIC, Literal: "HLT"},
				{Type: EOF},
			},
		},
		{
			name:  "single byte instruction with comma and space after comma",
			input: "MOV B, H",
			want: []Token{
				{Type: MNEMONIC, Literal: "MOV"},
				{Type: REGISTER, Literal: "B"},
				{Type: COMMA, Literal: ","},
				{Type: REGISTER, Literal: "H"},
				{Type: EOF},
			},
		},
		{
			name:  "single byte instruction with comma and space before comma",
			input: "MOV B ,H",
			want: []Token{
				{Type: MNEMONIC, Literal: "MOV"},
				{Type: REGISTER, Literal: "B"},
				{Type: COMMA, Literal: ","},
				{Type: REGISTER, Literal: "H"},
				{Type: EOF},
			},
		},
		{
			name:  "two byte instruction",
			input: "MVI B, 34h",
			want: []Token{
				{Type: MNEMONIC, Literal: "MVI"},
				{Type: REGISTER, Literal: "B"},
				{Type: COMMA, Literal: ","},
				{Type: NUMBER, Literal: "34H"},
				{Type: EOF},
			},
		},
		{
			name:  "three byte instruction",
			input: "LDA, 3412h",
			want: []Token{
				{Type: MNEMONIC, Literal: "LDA"},
				{Type: COMMA, Literal: ","},
				{Type: NUMBER, Literal: "3412H"},
				{Type: EOF},
			},
		},
		{
			name:  "lots of extra space",
			input: "    mov B ,      C   ",
			want: []Token{
				{Type: MNEMONIC, Literal: "MOV"},
				{Type: REGISTER, Literal: "B"},
				{Type: COMMA, Literal: ","},
				{Type: REGISTER, Literal: "C"},
				{Type: EOF},
			},
		},
		{
			name:  "move, load and store mnemonics",
			input: "MOV MVI LXI STAX LDAX STA LDA SHLD LHLD XCHG",
			want: []Token{
				{Type: MNEMONIC, Literal: "MOV"},
				{Type: MNEMONIC, Literal: "MVI"},
				{Type: MNEMONIC, Literal: "LXI"},
				{Type: MNEMONIC, Literal: "STAX"},
				{Type: MNEMONIC, Literal: "LDAX"},
				{Type: MNEMONIC, Literal: "STA"},
				{Type: MNEMONIC, Literal: "LDA"},
				{Type: MNEMONIC, Literal: "SHLD"},
				{Type: MNEMONIC, Literal: "LHLD"},
				{Type: MNEMONIC, Literal: "XCHG"},
				{Type: EOF},
			},
		},
		{
			name:  "stack operation mnemonics",
			input: "PUSH POP XTHL SPHL INX DCX DAD",
			want: []Token{
				{Type: MNEMONIC, Literal: "PUSH"},
				{Type: MNEMONIC, Literal: "POP"},
				{Type: MNEMONIC, Literal: "XTHL"},
				{Type: MNEMONIC, Literal: "SPHL"},
				{Type: MNEMONIC, Literal: "INX"},
				{Type: MNEMONIC, Literal: "DCX"},
				{Type: MNEMONIC, Literal: "DAD"},
				{Type: EOF},
			},
		},
		{
			name:  "jump mnemonics",
			input: "JMP JC JNC JZ JNZ JP JM JPE JPO PCHL",
			want: []Token{
				{Type: MNEMONIC, Literal: "JMP"},
				{Type: MNEMONIC, Literal: "JC"},
				{Type: MNEMONIC, Literal: "JNC"},
				{Type: MNEMONIC, Literal: "JZ"},
				{Type: MNEMONIC, Literal: "JNZ"},
				{Type: MNEMONIC, Literal: "JP"},
				{Type: MNEMONIC, Literal: "JM"},
				{Type: MNEMONIC, Literal: "JPE"},
				{Type: MNEMONIC, Literal: "JPO"},
				{Type: MNEMONIC, Literal: "PCHL"},
				{Type: EOF},
			},
		},
		{
			name:  "call mnemonics",
			input: "CALL CC CNC CZ CNZ CP CM CPE CPO",
			want: []Token{
				{Type: MNEMONIC, Literal: "CALL"},
				{Type: MNEMONIC, Literal: "CC"},
				{Type: MNEMONIC, Literal: "CNC"},
				{Type: MNEMONIC, Literal: "CZ"},
				{Type: MNEMONIC, Literal: "CNZ"},
				{Type: MNEMONIC, Literal: "CP"},
				{Type: MNEMONIC, Literal: "CM"},
				{Type: MNEMONIC, Literal: "CPE"},
				{Type: MNEMONIC, Literal: "CPO"},
				{Type: EOF},
			},
		},
		{
			name:  "return mnemonics",
			input: "RET RC RNC RZ RNZ RP RM RPE RPO",
			want: []Token{
				{Type: MNEMONIC, Literal: "RET"},
				{Type: MNEMONIC, Literal: "RC"},
				{Type: MNEMONIC, Literal: "RNC"},
				{Type: MNEMONIC, Literal: "RZ"},
				{Type: MNEMONIC, Literal: "RNZ"},
				{Type: MNEMONIC, Literal: "RP"},
				{Type: MNEMONIC, Literal: "RM"},
				{Type: MNEMONIC, Literal: "RPE"},
				{Type: MNEMONIC, Literal: "RPO"},
				{Type: EOF},
			},
		},
		{
			name:  "restart mnemonic",
			input: "RST",
			want: []Token{
				{Type: MNEMONIC, Literal: "RST"},
				{Type: EOF},
			},
		},
		{
			name:  "increment and decrement mnemonics",
			input: "INR DCR",
			want: []Token{
				{Type: MNEMONIC, Literal: "INR"},
				{Type: MNEMONIC, Literal: "DCR"},
				{Type: EOF},
			},
		},
		{
			name:  "add and subtract mnemonics",
			input: "ADD ADC ADI ACI SUB SBB SUI SBI",
			want: []Token{
				{Type: MNEMONIC, Literal: "ADD"},
				{Type: MNEMONIC, Literal: "ADC"},
				{Type: MNEMONIC, Literal: "ADI"},
				{Type: MNEMONIC, Literal: "ACI"},
				{Type: MNEMONIC, Literal: "SUB"},
				{Type: MNEMONIC, Literal: "SBB"},
				{Type: MNEMONIC, Literal: "SUI"},
				{Type: MNEMONIC, Literal: "SBI"},
				{Type: EOF},
			},
		},
		{
			name:  "logical mnemonics",
			input: "ANA XRA ORA CMP ANI XRI ORI CPI",
			want: []Token{
				{Type: MNEMONIC, Literal: "ANA"},
				{Type: MNEMONIC, Literal: "XRA"},
				{Type: MNEMONIC, Literal: "ORA"},
				{Type: MNEMONIC, Literal: "CMP"},
				{Type: MNEMONIC, Literal: "ANI"},
				{Type: MNEMONIC, Literal: "XRI"},
				{Type: MNEMONIC, Literal: "ORI"},
				{Type: MNEMONIC, Literal: "CPI"},
				{Type: EOF},
			},
		},
		{
			name:  "rotate mnemonics",
			input: "RLC RRC RAL RAR",
			want: []Token{
				{Type: MNEMONIC, Literal: "RLC"},
				{Type: MNEMONIC, Literal: "RRC"},
				{Type: MNEMONIC, Literal: "RAL"},
				{Type: MNEMONIC, Literal: "RAR"},
				{Type: EOF},
			},
		},
		{
			name:  "special mnemonics",
			input: "CMA STC CMC DAA",
			want: []Token{
				{Type: MNEMONIC, Literal: "CMA"},
				{Type: MNEMONIC, Literal: "STC"},
				{Type: MNEMONIC, Literal: "CMC"},
				{Type: MNEMONIC, Literal: "DAA"},
				{Type: EOF},
			},
		},
		{
			name:  "input/output mnemonics",
			input: "IN OUT",
			want: []Token{
				{Type: MNEMONIC, Literal: "IN"},
				{Type: MNEMONIC, Literal: "OUT"},
				{Type: EOF},
			},
		},
		{
			name:  "control mnemonics",
			input: "EI DI NOP HLT",
			want: []Token{
				{Type: MNEMONIC, Literal: "EI"},
				{Type: MNEMONIC, Literal: "DI"},
				{Type: MNEMONIC, Literal: "NOP"},
				{Type: MNEMONIC, Literal: "HLT"},
				{Type: EOF},
			},
		},
		{
			name:  "general purpose registers",
			input: "A B C D E H L M",
			want: []Token{
				{Type: REGISTER, Literal: "A"},
				{Type: REGISTER, Literal: "B"},
				{Type: REGISTER, Literal: "C"},
				{Type: REGISTER, Literal: "D"},
				{Type: REGISTER, Literal: "E"},
				{Type: REGISTER, Literal: "H"},
				{Type: REGISTER, Literal: "L"},
				{Type: REGISTER, Literal: "M"},
				{Type: EOF},
			},
		},
		{
			name:  "special purpose registers",
			input: "SP PSW",
			want: []Token{
				{Type: REGISTER, Literal: "SP"},
				{Type: REGISTER, Literal: "PSW"},
				{Type: EOF},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := New(tt.input)
			got, err := lexer.Lex()
			if (err != nil) != tt.wantErr {
				t.Errorf("Lexer.Lex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Lexer.Lex() got = %v, want %v", got, tt.want)
			}
		})
	}
}
