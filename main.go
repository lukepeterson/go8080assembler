package main

import (
	"fmt"

	"github.com/lukepeterson/go8080assembler/pkg/lexer"
	"github.com/lukepeterson/go8080assembler/pkg/parser"
)

func main() {

	input := `
			MVI A, 0x33
	START:	MOV B, C	; First comment
			LDA 0x1234	; Second comment
			JMP START
	`

	// Lex
	l := lexer.New(input)
	tokens, err := l.Lex()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("tokens: %v\n", tokens)

	// Parse
	p := parser.New(tokens)
	bytecode, err := p.Parse()
	if err != nil {
		fmt.Println(err)
	}

	// Print bytecode
	for _, b := range bytecode {
		fmt.Printf("%02X ", b)
	}
	fmt.Println()
}
