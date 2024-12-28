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
				t.Errorf("Lexer.Lex() = %v, want %v", got, tt.want)
			}
		})
	}
}
