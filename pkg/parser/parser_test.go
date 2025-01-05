package parser

import (
	"reflect"
	"testing"

	"github.com/lukepeterson/go8080assembler/pkg/lexer"
)

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name         string
		tokens       []lexer.Token
		wantBytecode []byte
		wantErr      bool
	}{
		{
			name: "label defined before JMP",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "A"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "B"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.COLON, Literal: ":"},
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "C"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "D"},
				{Type: lexer.MNEMONIC, Literal: "JMP"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x78, 0x4A, 0xC3, 0x01, 0x00},
		},
		{
			name: "multiple labels",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "A"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "B"},
				{Type: lexer.LABEL, Literal: "END2"},
				{Type: lexer.COLON, Literal: ":"},
				{Type: lexer.MNEMONIC, Literal: "JMP"},
				{Type: lexer.LABEL, Literal: "END"},
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "C"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "D"},
				{Type: lexer.MNEMONIC, Literal: "JMP"},
				{Type: lexer.LABEL, Literal: "END2"},
				{Type: lexer.LABEL, Literal: "END"},
				{Type: lexer.COLON, Literal: ":"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x78, 0xC3, 0x08, 0x00, 0x4A, 0xC3, 0x01, 0x00},
		},
		{
			name: "label defined after JMP",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "A"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "B"},
				{Type: lexer.MNEMONIC, Literal: "JMP"},
				{Type: lexer.LABEL, Literal: "END"},
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "C"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "D"},
				{Type: lexer.LABEL, Literal: "END"},
				{Type: lexer.COLON, Literal: ":"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x78, 0xC3, 0x05, 0x00, 0x4A},
		},
		{
			name: "MOV A, B",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "A"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "B"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x78},
		},
		// TODO: Add all the other MOV instructions here to confirm all their bytecodes (64 in total)
		{
			name: "MOV X, A (invalid source register)",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "X"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "A"},
			},
			wantErr: true,
		},
		{
			name: "MOV A, X (invalid destination register)",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "A"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "X"},
			},
			wantErr: true,
		},
		{
			name: "MOV X A (missing comma)",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "X"},
				{Type: lexer.REGISTER, Literal: "A"},
			},
			wantErr: true,
		},
		{
			name: "MOV 1 (missing register as first operand)",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.NUMBER, Literal: "1"},
			},
			wantErr: true,
		},
		{
			name: "MOV A, 1 (missing register as second operand)",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "A"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.NUMBER, Literal: "1"},
			},
			wantErr: true,
		},

		{
			name: "MVI B, 0x55",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MVI"},
				{Type: lexer.REGISTER, Literal: "B"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.NUMBER, Literal: "0x55"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x06, 0x55},
		},
		{
			name: "MVI C, 0x55",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MVI"},
				{Type: lexer.REGISTER, Literal: "C"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.NUMBER, Literal: "0x55"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x0E, 0x55},
		},
		{
			name: "MVI D, 0x55",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MVI"},
				{Type: lexer.REGISTER, Literal: "D"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.NUMBER, Literal: "0x55"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x16, 0x55},
		},
		{
			name: "MVI E, 0x55",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MVI"},
				{Type: lexer.REGISTER, Literal: "E"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.NUMBER, Literal: "0x55"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x1E, 0x55},
		},
		{
			name: "MVI H, 0x55",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MVI"},
				{Type: lexer.REGISTER, Literal: "H"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.NUMBER, Literal: "0x55"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x26, 0x55},
		},
		{
			name: "MVI L, 0x55",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MVI"},
				{Type: lexer.REGISTER, Literal: "L"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.NUMBER, Literal: "0x55"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x2E, 0x55},
		},
		{
			name: "MVI M, 0x55",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MVI"},
				{Type: lexer.REGISTER, Literal: "M"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.NUMBER, Literal: "0x55"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x36, 0x55},
		},
		{
			name: "MVI A, 0x55",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MVI"},
				{Type: lexer.REGISTER, Literal: "A"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.NUMBER, Literal: "0x55"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x3E, 0x55},
		},
		{
			name: "MVI A, B (invalid immediate data)",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MVI"},
				{Type: lexer.REGISTER, Literal: "A"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "B"},
				{Type: lexer.EOF},
			},
			wantErr: true,
		},
		{
			name: "LXI B, 0x4455",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "LXI"},
				{Type: lexer.REGISTER, Literal: "B"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.NUMBER, Literal: "0x4455"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x01, 0x55, 0x44},
		},
		{
			name: "LXI C, 0x4455 (invalid destination register)",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "LXI"},
				{Type: lexer.REGISTER, Literal: "C"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.NUMBER, Literal: "0x4455"},
				{Type: lexer.EOF},
			},
			wantErr: true,
		},
		{
			name: "LXI D, 0x4455",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "LXI"},
				{Type: lexer.REGISTER, Literal: "D"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.NUMBER, Literal: "0x4455"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x11, 0x55, 0x44},
		},
		{
			name: "LXI E, 0x4455 (invalid destination register)",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "LXI"},
				{Type: lexer.REGISTER, Literal: "E"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.NUMBER, Literal: "0x4455"},
				{Type: lexer.EOF},
			},
			wantErr: true,
		},
		{
			name: "LXI H, 0x4455",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "LXI"},
				{Type: lexer.REGISTER, Literal: "H"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.NUMBER, Literal: "0x4455"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x21, 0x55, 0x44},
		},
		{
			name: "LXI L, 0x4455 (invalid destination register)",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "LXI"},
				{Type: lexer.REGISTER, Literal: "L"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.NUMBER, Literal: "0x4455"},
				{Type: lexer.EOF},
			},
			wantErr: true,
		},
		{
			name: "LXI SP, 0x4455",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "LXI"},
				{Type: lexer.REGISTER, Literal: "SP"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.NUMBER, Literal: "0x4455"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x31, 0x55, 0x44},
		},
		{
			name: "STAX B",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "STAX"},
				{Type: lexer.REGISTER, Literal: "B"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x02},
		},
		{
			name: "STAX D",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "STAX"},
				{Type: lexer.REGISTER, Literal: "D"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x12},
		},
		{
			name: "STAX H",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "STAX"},
				{Type: lexer.REGISTER, Literal: "H"},
				{Type: lexer.EOF},
			},
			wantErr: true,
		},
		{
			name: "LDAX B",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "LDAX"},
				{Type: lexer.REGISTER, Literal: "B"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x0A},
		},
		{
			name: "LDAX D",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "LDAX"},
				{Type: lexer.REGISTER, Literal: "D"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x1A},
		},
		{
			name: "LDAX H",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "LDAX"},
				{Type: lexer.REGISTER, Literal: "H"},
				{Type: lexer.EOF},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			p := New(tt.tokens)
			got, err := p.Parse()

			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.wantBytecode) {
				t.Errorf("Parser.Parse() = %X, want %X", got, tt.wantBytecode)
			}
		})
	}
}
