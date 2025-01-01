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
			name: "MOV with valid source and destination registers",
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
			name: "MOV with invalid source register",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "X"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "A"},
			},
			wantErr: true,
		},
		{
			name: "MOV with invalid destination register",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "A"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "X"},
			},
			wantErr: true,
		},
		{
			name: "MOV with missing comma",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "X"},
				{Type: lexer.REGISTER, Literal: "A"},
			},
			wantErr: true,
		},
		{
			name: "MOV with missing register operand one",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.NUMBER, Literal: "1"},
			},
			wantErr: true,
		},
		{
			name: "MOV with missing register operand two",
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
