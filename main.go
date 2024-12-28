package main

import (
	// import "github.com/lukepeterson/go8080assembler/assembler"

	"fmt"

	"github.com/lukepeterson/go8080assembler/pkg/lexer"
	"github.com/lukepeterson/go8080assembler/pkg/parser"
)

func main() {
	// code := `
	// 	MVI A, 34h
	// 	MOV B, C
	// 	LDA, 1234h
	// 	HLT
	// `

	// assembler := assembler.New()

	// err := assembler.Assemble(code)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// for _, instruction := range assembler.ByteCode {
	// 	fmt.Printf("%02X ", instruction)
	// }

	// input := `MVI A, 0x33
	// TEST:	MOV B, C   ; The first comment
	// 		LDA, 1234h ; The second comment
	// 		JMP TEST

	// 	; The third comment with another colon ; here
	// 		`

	input := `
				MOV A, B
				END2:
				JMP END
				MOV C, D
				JMP END2
				END:
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

	for _, b := range bytecode {
		fmt.Printf("%02X ", b)
	}
	fmt.Println()

}
