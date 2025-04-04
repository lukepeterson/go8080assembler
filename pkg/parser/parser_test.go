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
			name: "two references to the same label",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "JMP"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.COLON, Literal: ":"},
				{Type: lexer.MNEMONIC, Literal: "JMP"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xC3, 0x03, 0x00, 0xC3, 0x03, 0x00},
		},
		{
			name: "three references to the same label",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "JMP"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.COLON, Literal: ":"},
				{Type: lexer.MNEMONIC, Literal: "JMP"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.MNEMONIC, Literal: "JMP"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xC3, 0x03, 0x00, 0xC3, 0x03, 0x00, 0xC3, 0x03, 0x00},
		},
		{
			name: "two definitions of the same label",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "JMP"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.COLON, Literal: ":"},
				{Type: lexer.MNEMONIC, Literal: "JMP"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.COLON, Literal: ":"},
				{Type: lexer.EOF},
			},
			wantErr: true,
		},
		{
			name: "three definitions of the same label",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "JMP"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.COLON, Literal: ":"},
				{Type: lexer.MNEMONIC, Literal: "JMP"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.COLON, Literal: ":"},
				{Type: lexer.MNEMONIC, Literal: "JMP"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.COLON, Literal: ":"},
				{Type: lexer.EOF},
			},
			wantErr: true,
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
		{
			name: "MOV A, C",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "A"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "C"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x79},
		},
		{
			name: "MOV A, D",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "A"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "D"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x7A},
		},
		{
			name: "MOV A, E",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "A"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "E"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x7B},
		},
		{
			name: "MOV A, H",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "A"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "H"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x7C},
		},
		{
			name: "MOV A, L",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "A"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "L"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x7D},
		},
		{
			name: "MOV A, M",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "A"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "M"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x7E},
		},
		{
			name: "MOV A, A",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "A"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "A"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x7F},
		},
		{
			name: "MOV B, B",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "B"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "B"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x40},
		},
		{
			name: "MOV B, C",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "B"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "C"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x41},
		},
		{
			name: "MOV B, D",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "B"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "D"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x42},
		},
		{
			name: "MOV B, E",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "B"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "E"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x43},
		},
		{
			name: "MOV B, H",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "B"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "H"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x44},
		},
		{
			name: "MOV B, L",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "B"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "L"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x45},
		},
		{
			name: "MOV B, M",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "B"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "M"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x46},
		},
		{
			name: "MOV B, A",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "B"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "A"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x47},
		},
		{
			name: "MOV C, B",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "C"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "B"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x48},
		},
		{
			name: "MOV C, C",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "C"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "C"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x49},
		},
		{
			name: "MOV C, D",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "C"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "D"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x4A},
		},
		{
			name: "MOV C, E",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "C"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "E"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x4B},
		},
		{
			name: "MOV C, H",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "C"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "H"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x4C},
		},
		{
			name: "MOV C, L",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "C"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "L"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x4D},
		},
		{
			name: "MOV C, M",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "C"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "M"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x4E},
		},
		{
			name: "MOV C, A",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "C"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "A"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x4F},
		},
		{
			name: "MOV D, B",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "D"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "B"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x50},
		},
		{
			name: "MOV D, C",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "D"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "C"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x51},
		},
		{
			name: "MOV D, D",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "D"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "D"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x52},
		},
		{
			name: "MOV D, E",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "D"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "E"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x53},
		},
		{
			name: "MOV D, H",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "D"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "H"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x54},
		},
		{
			name: "MOV D, L",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "D"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "L"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x55},
		},
		{
			name: "MOV D, M",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "D"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "M"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x56},
		},
		{
			name: "MOV D, A",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "D"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "A"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x57},
		},
		{
			name: "MOV E, B",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "E"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "B"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x58},
		},
		{
			name: "MOV E, C",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "E"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "C"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x59},
		},
		{
			name: "MOV E, D",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "E"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "D"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x5A},
		},
		{
			name: "MOV E, E",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "E"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "E"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x5B},
		},
		{
			name: "MOV E, H",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "E"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "H"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x5C},
		},
		{
			name: "MOV E, L",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "E"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "L"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x5D},
		},
		{
			name: "MOV E, M",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "E"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "M"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x5E},
		},
		{
			name: "MOV E, A",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "E"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "A"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x5F},
		},
		{
			name: "MOV H, B",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "H"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "B"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x60},
		},
		{
			name: "MOV H, C",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "H"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "C"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x61},
		},
		{
			name: "MOV H, D",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "H"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "D"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x62},
		},
		{
			name: "MOV H, E",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "H"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "E"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x63},
		},
		{
			name: "MOV H, H",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "H"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "H"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x64},
		},
		{
			name: "MOV H, L",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "H"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "L"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x65},
		},
		{
			name: "MOV H, M",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "H"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "M"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x66},
		},
		{
			name: "MOV H, A",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "H"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "A"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x67},
		},
		{
			name: "MOV L, B",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "L"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "B"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x68},
		},
		{
			name: "MOV L, C",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "L"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "C"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x69},
		},
		{
			name: "MOV L, D",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "L"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "D"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x6A},
		},
		{
			name: "MOV L, E",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "L"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "E"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x6B},
		},
		{
			name: "MOV L, H",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "L"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "H"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x6C},
		},
		{
			name: "MOV L, L",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "L"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "L"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x6D},
		},
		{
			name: "MOV L, M",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "L"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "M"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x6E},
		},
		{
			name: "MOV L, A",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "L"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "A"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x6F},
		},
		{
			name: "MOV M, B",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "M"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "B"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x70},
		},
		{
			name: "MOV M, C",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "M"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "C"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x71},
		},
		{
			name: "MOV M, D",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "M"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "D"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x72},
		},
		{
			name: "MOV M, E",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "M"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "E"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x73},
		},
		{
			name: "MOV M, H",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "M"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "H"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x74},
		},
		{
			name: "MOV M, L",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "M"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "L"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x75},
		},
		{
			name: "MOV M, M",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "M"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "M"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x76},
		},
		{
			name: "MOV M, A",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "MOV"},
				{Type: lexer.REGISTER, Literal: "M"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.REGISTER, Literal: "A"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x77},
		},
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
			name: "LXI H, MSG (load address from label)",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "LXI"},
				{Type: lexer.REGISTER, Literal: "H"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.LABEL, Literal: "MSG"},
				{Type: lexer.LABEL, Literal: "MSG"},
				{Type: lexer.COLON, Literal: ":"},
				{Type: lexer.MNEMONIC, Literal: "DB"},
				{Type: lexer.STRING, Literal: "Test"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x21, 0x03, 0x00, 0x54, 0x65, 0x73, 0x74},
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
			name: "INX B",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "INX"},
				{Type: lexer.REGISTER, Literal: "B"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x03},
		},
		{
			name: "INX C (invalid destination register)",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "INX"},
				{Type: lexer.REGISTER, Literal: "C"},
				{Type: lexer.EOF},
			},
			wantErr: true,
		},
		{
			name: "INX D",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "INX"},
				{Type: lexer.REGISTER, Literal: "D"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x13},
		},
		{
			name: "INX E (invalid destination register)",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "INX"},
				{Type: lexer.REGISTER, Literal: "E"},
				{Type: lexer.EOF},
			},
			wantErr: true,
		},
		{
			name: "INX H",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "INX"},
				{Type: lexer.REGISTER, Literal: "H"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x23},
		},
		{
			name: "INX L (invalid destination register)",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "INX"},
				{Type: lexer.REGISTER, Literal: "L"},
				{Type: lexer.EOF},
			},
			wantErr: true,
		},
		{
			name: "INX SP",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "INX"},
				{Type: lexer.REGISTER, Literal: "SP"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x33},
		},
		{
			name: "DCX B",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "DCX"},
				{Type: lexer.REGISTER, Literal: "B"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x0B},
		},
		{
			name: "DCX C (invalid destination register)",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "DCX"},
				{Type: lexer.REGISTER, Literal: "C"},
				{Type: lexer.EOF},
			},
			wantErr: true,
		},
		{
			name: "DCX D",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "DCX"},
				{Type: lexer.REGISTER, Literal: "D"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x1B},
		},
		{
			name: "DCX E (invalid destination register)",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "DCX"},
				{Type: lexer.REGISTER, Literal: "E"},
				{Type: lexer.EOF},
			},
			wantErr: true,
		},
		{
			name: "DCX H",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "DCX"},
				{Type: lexer.REGISTER, Literal: "H"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x2B},
		},
		{
			name: "DCX L (invalid destination register)",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "DCX"},
				{Type: lexer.REGISTER, Literal: "L"},
				{Type: lexer.EOF},
			},
			wantErr: true,
		},
		{
			name: "DCX SP",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "DCX"},
				{Type: lexer.REGISTER, Literal: "SP"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x3B},
		},
		{
			name: "DAD B",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "DAD"},
				{Type: lexer.REGISTER, Literal: "B"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x09},
		},
		{
			name: "DAD C (invalid destination register)",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "DAD"},
				{Type: lexer.REGISTER, Literal: "C"},
				{Type: lexer.EOF},
			},
			wantErr: true,
		},
		{
			name: "DAD D",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "DAD"},
				{Type: lexer.REGISTER, Literal: "D"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x19},
		},
		{
			name: "DAD E (invalid destination register)",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "DAD"},
				{Type: lexer.REGISTER, Literal: "E"},
				{Type: lexer.EOF},
			},
			wantErr: true,
		},
		{
			name: "DAD H",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "DAD"},
				{Type: lexer.REGISTER, Literal: "H"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x29},
		},
		{
			name: "DAD L (invalid destination register)",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "DAD"},
				{Type: lexer.REGISTER, Literal: "L"},
				{Type: lexer.EOF},
			},
			wantErr: true,
		},
		{
			name: "DAD SP",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "DAD"},
				{Type: lexer.REGISTER, Literal: "SP"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x39},
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
			name: "STAX B, STAX B",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "STAX"},
				{Type: lexer.REGISTER, Literal: "B"},
				{Type: lexer.MNEMONIC, Literal: "STAX"},
				{Type: lexer.REGISTER, Literal: "B"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x02, 0x02},
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
		{
			name: "DB Hello",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "DB"},
				{Type: lexer.STRING, Literal: "Hello"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x48, 0x65, 0x6C, 0x6C, 0x6F},
		},
		{
			name: "STA 0x4455",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "STA"},
				{Type: lexer.NUMBER, Literal: "0x4455"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x32, 0x55, 0x44},
		},
		{
			name: "STA MSG",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "STA"},
				{Type: lexer.LABEL, Literal: "MSG"},
				{Type: lexer.LABEL, Literal: "MSG"},
				{Type: lexer.COLON, Literal: ":"},
				{Type: lexer.MNEMONIC, Literal: "DB"},
				{Type: lexer.STRING, Literal: "Test"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x32, 0x03, 0x00, 0x54, 0x65, 0x73, 0x74},
		},
		{
			name: "LDA 0x4455",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "LDA"},
				{Type: lexer.NUMBER, Literal: "0x4455"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x3A, 0x55, 0x44},
		},
		{
			name: "LDA MSG",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "LDA"},
				{Type: lexer.LABEL, Literal: "MSG"},
				{Type: lexer.LABEL, Literal: "MSG"},
				{Type: lexer.COLON, Literal: ":"},
				{Type: lexer.MNEMONIC, Literal: "DB"},
				{Type: lexer.STRING, Literal: "Test"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x3A, 0x03, 0x00, 0x54, 0x65, 0x73, 0x74},
		},
		{
			name: "XCHG",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "XCHG"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xEB},
		},
		{
			name: "PUSH B",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "PUSH"},
				{Type: lexer.REGISTER, Literal: "B"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xC5},
		},
		{
			name: "PUSH C (invalid destination register)",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "PUSH"},
				{Type: lexer.REGISTER, Literal: "C"},
				{Type: lexer.EOF},
			},
			wantErr: true,
		},
		{
			name: "PUSH D",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "PUSH"},
				{Type: lexer.REGISTER, Literal: "D"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xD5},
		},
		{
			name: "PUSH E (invalid destination register)",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "PUSH"},
				{Type: lexer.REGISTER, Literal: "E"},
				{Type: lexer.EOF},
			},
			wantErr: true,
		},
		{
			name: "PUSH H",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "PUSH"},
				{Type: lexer.REGISTER, Literal: "H"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xE5},
		},
		{
			name: "PUSH L (invalid destination register)",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "PUSH"},
				{Type: lexer.REGISTER, Literal: "L"},
				{Type: lexer.EOF},
			},
			wantErr: true,
		},
		{
			name: "PUSH M (invalid destination register)",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "PUSH"},
				{Type: lexer.REGISTER, Literal: "M"},
				{Type: lexer.EOF},
			},
			wantErr: true,
		},
		{
			name: "PUSH A (invalid destination register)",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "PUSH"},
				{Type: lexer.REGISTER, Literal: "A"},
				{Type: lexer.EOF},
			},
			wantErr: true,
		},
		{
			name: "PUSH PSW",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "PUSH"},
				{Type: lexer.REGISTER, Literal: "PSW"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xF5},
		}, {
			name: "POP B",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "POP"},
				{Type: lexer.REGISTER, Literal: "B"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xC1},
		},
		{
			name: "POP C (invalid destination register)",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "POP"},
				{Type: lexer.REGISTER, Literal: "C"},
				{Type: lexer.EOF},
			},
			wantErr: true,
		},
		{
			name: "POP D",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "POP"},
				{Type: lexer.REGISTER, Literal: "D"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xD1},
		},
		{
			name: "POP E (invalid destination register)",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "POP"},
				{Type: lexer.REGISTER, Literal: "E"},
				{Type: lexer.EOF},
			},
			wantErr: true,
		},
		{
			name: "POP H",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "POP"},
				{Type: lexer.REGISTER, Literal: "H"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xE1},
		},
		{
			name: "POP L (invalid destination register)",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "POP"},
				{Type: lexer.REGISTER, Literal: "L"},
				{Type: lexer.EOF},
			},
			wantErr: true,
		},
		{
			name: "POP M (invalid destination register)",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "POP"},
				{Type: lexer.REGISTER, Literal: "M"},
				{Type: lexer.EOF},
			},
			wantErr: true,
		},
		{
			name: "POP A (invalid destination register)",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "POP"},
				{Type: lexer.REGISTER, Literal: "A"},
				{Type: lexer.EOF},
			},
			wantErr: true,
		},
		{
			name: "POP PSW",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "POP"},
				{Type: lexer.REGISTER, Literal: "PSW"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xF1},
		},
		{
			name: "XTHL",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "XTHL"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xE3},
		},
		{
			name: "SPHL",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "SPHL"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xF9},
		},
		{
			name: "JMP (to address)",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "JMP"},
				{Type: lexer.NUMBER, Literal: "0x04"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xC3, 0x04, 0x00},
		},
		{
			name: "JMP (to label)",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "JMP"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.COLON, Literal: ":"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xC3, 0x03, 0x00},
		},
		{
			name: "JC",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "JC"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.COLON, Literal: ":"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xDA, 0x03, 0x00},
		},
		{
			name: "JNC",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "JNC"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.COLON, Literal: ":"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xD2, 0x03, 0x00},
		},
		{
			name: "JZ",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "JZ"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.COLON, Literal: ":"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xCA, 0x03, 0x00},
		},
		{
			name: "JNZ",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "JNZ"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.COLON, Literal: ":"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xC2, 0x03, 0x00},
		},
		{
			name: "JP",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "JP"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.COLON, Literal: ":"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xF2, 0x03, 0x00},
		},
		{
			name: "JM",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "JM"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.COLON, Literal: ":"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xFA, 0x03, 0x00},
		},
		{
			name: "JPE",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "JPE"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.COLON, Literal: ":"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xEA, 0x03, 0x00},
		},
		{
			name: "JPO",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "JPO"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.COLON, Literal: ":"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xE2, 0x03, 0x00},
		},
		{
			name: "PCHL",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "PCHL"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xE9},
		},
		{
			name: "CALL (to address)",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "CALL"},
				{Type: lexer.NUMBER, Literal: "0x04"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xCD, 0x04, 0x00},
		},
		{
			name: "CALL (to label)",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "CALL"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.COLON, Literal: ":"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xCD, 0x03, 0x00},
		},
		{
			name: "CC",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "CC"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.COLON, Literal: ":"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xDC, 0x03, 0x00},
		},
		{
			name: "CNC",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "CNC"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.COLON, Literal: ":"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xD4, 0x03, 0x00},
		},
		{
			name: "CZ",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "CZ"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.COLON, Literal: ":"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xCC, 0x03, 0x00},
		},
		{
			name: "CNZ",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "CNZ"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.COLON, Literal: ":"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xC4, 0x03, 0x00},
		},
		{
			name: "CP",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "CP"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.COLON, Literal: ":"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xF4, 0x03, 0x00},
		},
		{
			name: "CM",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "CM"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.COLON, Literal: ":"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xFC, 0x03, 0x00},
		},
		{
			name: "CPE",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "CPE"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.COLON, Literal: ":"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xEC, 0x03, 0x00},
		},
		{
			name: "CPO",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "CPO"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.LABEL, Literal: "START"},
				{Type: lexer.COLON, Literal: ":"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xE4, 0x03, 0x00},
		},
		{
			name: "RET",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "RET"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xC9},
		},
		{
			name: "RC",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "RC"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xD8},
		},
		{
			name: "RNC",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "RNC"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xD0},
		},
		{
			name: "RZ",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "RZ"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xC8},
		},
		{
			name: "RNZ",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "RNZ"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xC0},
		},
		{
			name: "RP",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "RP"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xF0},
		},
		{
			name: "RM",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "RM"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xF8},
		},
		{
			name: "RPE",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "RPE"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xE8},
		},
		{
			name: "RPO",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "RPO"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xE0},
		},
		{
			name: "RST 0",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "RST"},
				{Type: lexer.NUMBER, Literal: "0"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xC7},
		},
		{
			name: "RST 1",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "RST"},
				{Type: lexer.NUMBER, Literal: "1"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xCF},
		},
		{
			name: "RST 2",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "RST"},
				{Type: lexer.NUMBER, Literal: "2"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xD7},
		},
		{
			name: "RST 3",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "RST"},
				{Type: lexer.NUMBER, Literal: "3"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xDF},
		},
		{
			name: "RST 4",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "RST"},
				{Type: lexer.NUMBER, Literal: "4"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xE7},
		},
		{
			name: "RST 5",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "RST"},
				{Type: lexer.NUMBER, Literal: "5"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xEF},
		},
		{
			name: "RST 6",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "RST"},
				{Type: lexer.NUMBER, Literal: "6"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xF7},
		},
		{
			name: "RST 7",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "RST"},
				{Type: lexer.NUMBER, Literal: "7"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xFF},
		},
		{
			name: "RST 8 (routine out of range)",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "RST"},
				{Type: lexer.NUMBER, Literal: "8"},
				{Type: lexer.EOF},
			},
			wantErr: true,
		},
		{
			name: "RST -1 (negative number for routine)",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "RST"},
				{Type: lexer.NUMBER, Literal: "-1"},
				{Type: lexer.EOF},
			},
			wantErr: true,
		},
		{
			name: "RST A (non-number for routine)",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "RST"},
				{Type: lexer.STRING, Literal: "A"},
				{Type: lexer.EOF},
			},
			wantErr: true,
		},
		{
			name: "INR B",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "INR"},
				{Type: lexer.REGISTER, Literal: "B"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x04},
		},
		{
			name: "INR C",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "INR"},
				{Type: lexer.REGISTER, Literal: "C"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x0C},
		},
		{
			name: "INR D",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "INR"},
				{Type: lexer.REGISTER, Literal: "D"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x14},
		},
		{
			name: "INR E",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "INR"},
				{Type: lexer.REGISTER, Literal: "E"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x1C},
		},
		{
			name: "INR H",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "INR"},
				{Type: lexer.REGISTER, Literal: "H"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x24},
		},
		{
			name: "INR L",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "INR"},
				{Type: lexer.REGISTER, Literal: "L"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x2C},
		},
		{
			name: "INR M",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "INR"},
				{Type: lexer.REGISTER, Literal: "M"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x34},
		},
		{
			name: "INR A",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "INR"},
				{Type: lexer.REGISTER, Literal: "A"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x3C},
		},
		{
			name: "INR A, INR B",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "INR"},
				{Type: lexer.REGISTER, Literal: "A"},
				{Type: lexer.MNEMONIC, Literal: "INR"},
				{Type: lexer.REGISTER, Literal: "B"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x3C, 0x04},
		},
		{
			name: "DCR B",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "DCR"},
				{Type: lexer.REGISTER, Literal: "B"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x05},
		},
		{
			name: "DCR C",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "DCR"},
				{Type: lexer.REGISTER, Literal: "C"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x0D},
		},
		{
			name: "DCR D",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "DCR"},
				{Type: lexer.REGISTER, Literal: "D"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x15},
		},
		{
			name: "DCR E",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "DCR"},
				{Type: lexer.REGISTER, Literal: "E"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x1D},
		},
		{
			name: "DCR H",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "DCR"},
				{Type: lexer.REGISTER, Literal: "H"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x25},
		},
		{
			name: "DCR L",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "DCR"},
				{Type: lexer.REGISTER, Literal: "L"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x2D},
		},
		{
			name: "DCR M",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "DCR"},
				{Type: lexer.REGISTER, Literal: "M"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x35},
		},
		{
			name: "DCR A",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "DCR"},
				{Type: lexer.REGISTER, Literal: "A"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x3D},
		},
		{
			name: "ADD B",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "ADD"},
				{Type: lexer.REGISTER, Literal: "B"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x80},
		},
		{
			name: "ADD C",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "ADD"},
				{Type: lexer.REGISTER, Literal: "C"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x81},
		},
		{
			name: "ADD D",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "ADD"},
				{Type: lexer.REGISTER, Literal: "D"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x82},
		},
		{
			name: "ADD E",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "ADD"},
				{Type: lexer.REGISTER, Literal: "E"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x83},
		},
		{
			name: "ADD H",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "ADD"},
				{Type: lexer.REGISTER, Literal: "H"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x84},
		},
		{
			name: "ADD L",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "ADD"},
				{Type: lexer.REGISTER, Literal: "L"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x85},
		},
		{
			name: "ADD M",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "ADD"},
				{Type: lexer.REGISTER, Literal: "M"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x86},
		},
		{
			name: "ADD A",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "ADD"},
				{Type: lexer.REGISTER, Literal: "A"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x87},
		},
		{
			name: "ADC B",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "ADC"},
				{Type: lexer.REGISTER, Literal: "B"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x88},
		},
		{
			name: "ADC C",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "ADC"},
				{Type: lexer.REGISTER, Literal: "C"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x89},
		},
		{
			name: "ADC D",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "ADC"},
				{Type: lexer.REGISTER, Literal: "D"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x8A},
		},
		{
			name: "ADC E",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "ADC"},
				{Type: lexer.REGISTER, Literal: "E"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x8B},
		},
		{
			name: "ADC H",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "ADC"},
				{Type: lexer.REGISTER, Literal: "H"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x8C},
		},
		{
			name: "ADC L",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "ADC"},
				{Type: lexer.REGISTER, Literal: "L"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x8D},
		},
		{
			name: "ADC M",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "ADC"},
				{Type: lexer.REGISTER, Literal: "M"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x8E},
		},
		{
			name: "ADC A",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "ADC"},
				{Type: lexer.REGISTER, Literal: "A"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x8F},
		},
		{
			name: "ADI 0x44",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "ADI"},
				{Type: lexer.NUMBER, Literal: "0x44"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xC6, 0x44},
		},
		{
			name: "ADI 0x4455 (two bytes of data)",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "ADI"},
				{Type: lexer.NUMBER, Literal: "0x4455"},
				{Type: lexer.EOF},
			},
			wantErr: true,
		},
		{
			name: "ACI 0x44",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "ACI"},
				{Type: lexer.NUMBER, Literal: "0x44"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xCE, 0x44},
		},
		{
			name: "ACI 0x4455 (two bytes of data)",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "ACI"},
				{Type: lexer.NUMBER, Literal: "0x4455"},
				{Type: lexer.EOF},
			},
			wantErr: true,
		},
		{
			name: "SUB B",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "SUB"},
				{Type: lexer.REGISTER, Literal: "B"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x90},
		},
		{
			name: "SUB C",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "SUB"},
				{Type: lexer.REGISTER, Literal: "C"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x91},
		},
		{
			name: "SUB D",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "SUB"},
				{Type: lexer.REGISTER, Literal: "D"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x92},
		},
		{
			name: "SUB E",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "SUB"},
				{Type: lexer.REGISTER, Literal: "E"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x93},
		},
		{
			name: "SUB H",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "SUB"},
				{Type: lexer.REGISTER, Literal: "H"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x94},
		},
		{
			name: "SUB L",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "SUB"},
				{Type: lexer.REGISTER, Literal: "L"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x95},
		},
		{
			name: "SUB M",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "SUB"},
				{Type: lexer.REGISTER, Literal: "M"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x96},
		},
		{
			name: "SUB A",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "SUB"},
				{Type: lexer.REGISTER, Literal: "A"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x97},
		},
		{
			name: "SBB B",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "SBB"},
				{Type: lexer.REGISTER, Literal: "B"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x98},
		},
		{
			name: "SBB C",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "SBB"},
				{Type: lexer.REGISTER, Literal: "C"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x99},
		},
		{
			name: "SBB D",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "SBB"},
				{Type: lexer.REGISTER, Literal: "D"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x9A},
		},
		{
			name: "SBB E",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "SBB"},
				{Type: lexer.REGISTER, Literal: "E"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x9B},
		},
		{
			name: "SBB H",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "SBB"},
				{Type: lexer.REGISTER, Literal: "H"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x9C},
		},
		{
			name: "SBB L",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "SBB"},
				{Type: lexer.REGISTER, Literal: "L"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x9D},
		},
		{
			name: "SBB M",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "SBB"},
				{Type: lexer.REGISTER, Literal: "M"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x9E},
		},
		{
			name: "SBB A",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "SBB"},
				{Type: lexer.REGISTER, Literal: "A"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x9F},
		},
		{
			name: "SUI 0x44",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "SUI"},
				{Type: lexer.NUMBER, Literal: "0x44"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xD6, 0x44},
		},
		{
			name: "SUI 0x4455 (two bytes of data)",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "SUI"},
				{Type: lexer.NUMBER, Literal: "0x4455"},
				{Type: lexer.EOF},
			},
			wantErr: true,
		},
		{
			name: "SBI 0x44",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "SBI"},
				{Type: lexer.NUMBER, Literal: "0x44"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xDE, 0x44},
		},
		{
			name: "SBI 0x4455 (two bytes of data)",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "SBI"},
				{Type: lexer.NUMBER, Literal: "0x4455"},
				{Type: lexer.EOF},
			},
			wantErr: true,
		},
		{
			name: "ANI 0x44",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "ANI"},
				{Type: lexer.NUMBER, Literal: "0x44"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xE6, 0x44},
		},
		{
			name: "ANI 0x4455 (two bytes of data)",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "SBI"},
				{Type: lexer.NUMBER, Literal: "0x4455"},
				{Type: lexer.EOF},
			},
			wantErr: true,
		},
		{
			name: "XRI 0x44",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "XRI"},
				{Type: lexer.NUMBER, Literal: "0x44"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xEE, 0x44},
		},
		{
			name: "XRI 0x4455 (two bytes of data)",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "XRI"},
				{Type: lexer.NUMBER, Literal: "0x4455"},
				{Type: lexer.EOF},
			},
			wantErr: true,
		},
		{
			name: "ORI 0x44",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "ORI"},
				{Type: lexer.NUMBER, Literal: "0x44"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xF6, 0x44},
		},
		{
			name: "ORI 0x4455 (two bytes of data)",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "ORI"},
				{Type: lexer.NUMBER, Literal: "0x4455"},
				{Type: lexer.EOF},
			},
			wantErr: true,
		},
		{
			name: "CPI 0x44",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "CPI"},
				{Type: lexer.NUMBER, Literal: "0x44"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xFE, 0x44},
		},
		{
			name: "CPI 0x4455 (two bytes of data)",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "CPI"},
				{Type: lexer.NUMBER, Literal: "0x4455"},
				{Type: lexer.EOF},
			},
			wantErr: true,
		},
		{
			name: "ANA B",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "ANA"},
				{Type: lexer.REGISTER, Literal: "B"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xA0},
		},
		{
			name: "ANA C",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "ANA"},
				{Type: lexer.REGISTER, Literal: "C"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xA1},
		},
		{
			name: "ANA D",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "ANA"},
				{Type: lexer.REGISTER, Literal: "D"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xA2},
		},
		{
			name: "ANA E",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "ANA"},
				{Type: lexer.REGISTER, Literal: "E"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xA3},
		},
		{
			name: "ANA H",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "ANA"},
				{Type: lexer.REGISTER, Literal: "H"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xA4},
		},
		{
			name: "ANA L",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "ANA"},
				{Type: lexer.REGISTER, Literal: "L"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xA5},
		},
		{
			name: "ANA M",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "ANA"},
				{Type: lexer.REGISTER, Literal: "M"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xA6},
		},
		{
			name: "ANA A",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "ANA"},
				{Type: lexer.REGISTER, Literal: "A"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xA7},
		},
		{
			name: "XRA B",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "XRA"},
				{Type: lexer.REGISTER, Literal: "B"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xA8},
		},
		{
			name: "XRA C",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "XRA"},
				{Type: lexer.REGISTER, Literal: "C"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xA9},
		},
		{
			name: "XRA D",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "XRA"},
				{Type: lexer.REGISTER, Literal: "D"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xAA},
		},
		{
			name: "XRA E",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "XRA"},
				{Type: lexer.REGISTER, Literal: "E"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xAB},
		},
		{
			name: "XRA H",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "XRA"},
				{Type: lexer.REGISTER, Literal: "H"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xAC},
		},
		{
			name: "XRA L",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "XRA"},
				{Type: lexer.REGISTER, Literal: "L"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xAD},
		},
		{
			name: "XRA M",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "XRA"},
				{Type: lexer.REGISTER, Literal: "M"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xAE},
		},
		{
			name: "XRA A",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "XRA"},
				{Type: lexer.REGISTER, Literal: "A"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xAF},
		},
		{
			name: "ORA B",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "ORA"},
				{Type: lexer.REGISTER, Literal: "B"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xB0},
		},
		{
			name: "ORA C",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "ORA"},
				{Type: lexer.REGISTER, Literal: "C"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xB1},
		},
		{
			name: "ORA D",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "ORA"},
				{Type: lexer.REGISTER, Literal: "D"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xB2},
		},
		{
			name: "ORA E",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "ORA"},
				{Type: lexer.REGISTER, Literal: "E"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xB3},
		},
		{
			name: "ORA H",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "ORA"},
				{Type: lexer.REGISTER, Literal: "H"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xB4},
		},
		{
			name: "ORA L",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "ORA"},
				{Type: lexer.REGISTER, Literal: "L"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xB5},
		},
		{
			name: "ORA M",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "ORA"},
				{Type: lexer.REGISTER, Literal: "M"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xB6},
		},
		{
			name: "ORA A",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "ORA"},
				{Type: lexer.REGISTER, Literal: "A"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xB7},
		},
		{
			name: "CMP B",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "CMP"},
				{Type: lexer.REGISTER, Literal: "B"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xB8},
		},
		{
			name: "CMP C",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "CMP"},
				{Type: lexer.REGISTER, Literal: "C"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xB9},
		},
		{
			name: "CMP D",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "CMP"},
				{Type: lexer.REGISTER, Literal: "D"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xBA},
		},
		{
			name: "CMP E",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "CMP"},
				{Type: lexer.REGISTER, Literal: "E"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xBB},
		},
		{
			name: "CMP H",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "CMP"},
				{Type: lexer.REGISTER, Literal: "H"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xBC},
		},
		{
			name: "CMP L",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "CMP"},
				{Type: lexer.REGISTER, Literal: "L"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xBD},
		},
		{
			name: "CMP M",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "CMP"},
				{Type: lexer.REGISTER, Literal: "M"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xBE},
		},
		{
			name: "CMP A",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "CMP"},
				{Type: lexer.REGISTER, Literal: "A"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xBF},
		},
		{
			name: "RLC",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "RLC"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x07},
		},
		{
			name: "RRC",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "RRC"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x0F},
		},
		{
			name: "RAL",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "RAL"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x17},
		},
		{
			name: "RAR",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "RAR"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x1F},
		},
		{
			name: "CMA",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "CMA"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x2F},
		},
		{
			name: "STC",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "STC"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x37},
		},
		{
			name: "CMC",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "CMC"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x3F},
		},
		{
			name: "DAA",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "DAA"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x27},
		},
		{
			name: "IN 0x44",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "IN"},
				{Type: lexer.NUMBER, Literal: "0x44"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xDB, 0x44},
		},
		{
			name: "IN 0x4455 (two bytes of data)",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "IN"},
				{Type: lexer.NUMBER, Literal: "0x4455"},
				{Type: lexer.EOF},
			},
			wantErr: true,
		},
		{
			name: "OUT 0x44",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "OUT"},
				{Type: lexer.NUMBER, Literal: "0x44"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xD3, 0x44},
		},
		{
			name: "OUT 0x4455 (two bytes of data)",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "OUT"},
				{Type: lexer.NUMBER, Literal: "0x4455"},
				{Type: lexer.EOF},
			},
			wantErr: true,
		},
		{
			name: "EI",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "EI"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xFB},
		},
		{
			name: "DI",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "DI"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0xF3},
		},
		{
			name: "NOP",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "NOP"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x00},
		},
		{
			name: "HLT",
			tokens: []lexer.Token{
				{Type: lexer.MNEMONIC, Literal: "HLT"},
				{Type: lexer.EOF},
			},
			wantBytecode: []byte{0x76},
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
