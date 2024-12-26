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

				// MOV A, B
				// START:
				// MOV C, D
				// JMP START

			},
			wantBytecode: []byte{0x78, 0x4A, 0xC3, 0x01, 0x00},
			wantErr:      false,
		},
		{
			name: "Multiple labels",
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

				// MOV A, B
				// END2:
				// JMP END
				// MOV C, D
				// JMP END2
				// END:
			},
			wantBytecode: []byte{0x78, 0xC3, 0x08, 0x00, 0x4A, 0xC3, 0x01, 0x00},
			wantErr:      false,
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

				// MOV A, B
				// JMP END
				// MOV C, D
				// END:
			},
			wantBytecode: []byte{0x78, 0xC3, 0x05, 0x00, 0x4A},
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			p := New(tt.tokens)

			err := p.Parse()
			got := p.Bytecode

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
