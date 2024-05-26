package assembler

import (
	"reflect"
	"testing"
)

func TestTokenise(t *testing.T) {
	tests := []struct {
		name string
		code string
		want []string
	}{
		{
			name: "case sensitivity",
			code: "hlt",
			want: []string{"HLT"},
		},
		{
			name: "single byte instruction without space",
			code: "HLT",
			want: []string{"HLT"},
		},
		{
			name: "single byte instruction with spaces",
			code: "  HLT  ",
			want: []string{"HLT"},
		},
		{
			name: "single byte instruction with comma and space after comma",
			code: "MOV B, H",
			want: []string{"MOV B,H"},
		},
		{
			name: "single byte instruction with comma and space before comma",
			code: "MOV B ,H",
			want: []string{"MOV B,H"},
		},
		{
			name: "two byte instruction",
			code: "MVI B, 34h",
			want: []string{"MVI B", "34H"},
		},
		{
			name: "three byte instruction",
			code: "LDA, 3412h",
			want: []string{"LDA", "3412H"},
		},
		{
			name: "lots of extra space",
			code: "    mov B ,      C   ",
			want: []string{"MOV B,C"},
		},
		{
			name: "invalid instruction",
			code: "FOO",
			want: nil,
		},

		// TODO
		// Test for length of results based on opcodes
		// Proper test for errors being returned
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := tokenise(tt.code)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tokenise() got %q, want %q", got, tt.want)
			}
		})
	}
}
