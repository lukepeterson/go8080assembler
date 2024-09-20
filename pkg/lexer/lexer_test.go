package lexer

import (
	"testing"
)

func TestNextToken(t *testing.T) {
	input := "START: MOV A, B ; Comment"

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{IDENT, "START"},
		{COLON, ":"},
		{MNEMONIC, "MOV"},
		{REGISTER, "A"},
		{COMMA, ","},
		{REGISTER, "B"},
		{COMMENT, "; Comment"},
		{EOF, ""},
	}

	l := New(input)
	for _, tt := range tests {
		token := l.NextToken()
		if token.Type != tt.expectedType {
			t.Fatalf("NextToken()\ngot type = %+v, \nwant type = %+v", tt.expectedType, token.Type)

		}
		if token.Literal != tt.expectedLiteral {
			t.Fatalf("NextToken()\ngot literal = %+q, \nwant literal = %+q", tt.expectedLiteral, token.Literal)
		}
	}
}
