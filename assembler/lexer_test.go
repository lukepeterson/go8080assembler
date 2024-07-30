package assembler

import (
	"reflect"
	"testing"
)

func TestTokeniseLine(t *testing.T) {
	tests := []struct {
		name string
		line string
		want []Token
	}{
		{
			name: "opcode with no operand",
			line: "STC",
			want: []Token{
				{
					Type:  OPCODE,
					Value: "STC",
				},
			},
		},
		{
			name: "opcode with one operand",
			line: "MVI C",
			want: []Token{
				{
					Type:  OPCODE,
					Value: "MVI",
				},
				{
					Type:  OPERAND1,
					Value: "C",
				},
			},
		},
		{
			name: "opcode with one operand and comment",
			line: "MVI C ;comment",
			want: []Token{
				{
					Type:  OPCODE,
					Value: "MVI",
				},
				{
					Type:  OPERAND1,
					Value: "C",
				},
			},
		},
		{
			name: "opcode with two operands",
			line: "MOV A, B",
			want: []Token{
				{
					Type:  OPCODE,
					Value: "MOV",
				},
				{
					Type:  OPERAND1,
					Value: "A",
				},
				{
					Type:  OPERAND2,
					Value: "B",
				},
			},
		},
		{
			name: "opcode with two operands and comment",
			line: "MOV A, B ;comment here",
			want: []Token{
				{
					Type:  OPCODE,
					Value: "MOV",
				},
				{
					Type:  OPERAND1,
					Value: "A",
				},
				{
					Type:  OPERAND2,
					Value: "B",
				},
			},
		},
		{
			name: "label with opcode, two operands and comment",
			line: "START: MOV A, B ;comment here",
			want: []Token{
				{
					Type:  LABEL,
					Value: "START",
				},
				{
					Type:  OPCODE,
					Value: "MOV",
				},
				{
					Type:  OPERAND1,
					Value: "A",
				},
				{
					Type:  OPERAND2,
					Value: "B",
				},
			},
		},
		{
			name: "label with opcode, two operands (one as a label) and comment",
			line: "START: MOV A, END ;comment here",
			want: []Token{
				{
					Type:  LABEL,
					Value: "START",
				},
				{
					Type:  OPCODE,
					Value: "MOV",
				},
				{
					Type:  OPERAND1,
					Value: "A",
				},
				{
					Type:  OPERAND2,
					Value: "END",
				},
			},
		},
		{
			name: "label only",
			line: "START:",
			want: []Token{
				{
					Type:  LABEL,
					Value: "START",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tokeniseLine(tt.line); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tokeniseLine()\ngot  = %+v, \nwant = %+v", got, tt.want)
			}
		})
	}
}

func TestStripComment(t *testing.T) {
	tests := []struct {
		name string
		line string
		want string
	}{
		{
			name: "single comment at start of line",
			line: ";comment",
			want: "",
		},
		{
			name: "multiple comments over multiple lines, with empty new lines",
			line: `

				; comment 1
				; comment 2 ;;;

			`,
			want: "",
		},
		{
			name: "comment after instruction",
			line: "MOV A, B ;comment",
			want: "MOV A, B",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stripComment(tt.line); got != tt.want {
				t.Errorf("stripComment()\ngot  = %q,\nwant = %q", got, tt.want)
			}
		})
	}
}

func TestLexicalAnalysis(t *testing.T) {
	tests := []struct {
		name string
		code string
		want []Token
	}{
		{
			name: "multiple lines",
			code: `
				ORG 100H
				; Comments here
		START:	MOV A, B
				MVI C, 1
				JMP START ; jump to start
		`,
			want: []Token{
				{
					Type:  OPCODE,
					Value: "ORG",
				},
				{
					Type:  OPERAND1,
					Value: "100H",
				},
				{
					Type:  LABEL,
					Value: "START",
				},
				{
					Type:  OPCODE,
					Value: "MOV",
				},
				{
					Type:  OPERAND1,
					Value: "A",
				},
				{
					Type:  OPERAND2,
					Value: "B",
				},
				{
					Type:  OPCODE,
					Value: "MVI",
				},
				{
					Type:  OPERAND1,
					Value: "C",
				},
				{
					Type:  OPERAND2,
					Value: "1",
				},
				{
					Type:  OPCODE,
					Value: "JMP",
				},
				{
					Type:  OPERAND1,
					Value: "START",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lexicalAnalysis(tt.code); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("lexicalAnalysis(): \ngot:  %v, \nwant: %v", got, tt.want)
			}
		})
	}
}
