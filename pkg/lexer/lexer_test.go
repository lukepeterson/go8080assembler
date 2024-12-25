package lexer

import (
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `
		START: MOV A, B ; 1st comment
		INR A
				; 2nd comment
		`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{LABEL, "START"},
		{COLON, ":"},
		{MNEMONIC, "MOV"},
		{REGISTER, "A"},
		{COMMA, ","},
		{REGISTER, "B"},
		{COMMENT, "; 1st comment"},

		{MNEMONIC, "INR"},
		{REGISTER, "A"},
		{COMMENT, "; 2nd comment"},

		{EOF, ""},
	}

	l := New(input)
	for _, tt := range tests {
		token := l.NextToken()
		if token.Type != tt.expectedType {
			t.Fatalf("NextToken()\ngot type = %+q, \nwant type = %+q", tt.expectedType, token.Type)

		}
		if token.Literal != tt.expectedLiteral {
			t.Fatalf("NextToken()\ngot literal = %+q, \nwant literal = %+q", tt.expectedLiteral, token.Literal)
		}
	}
}
