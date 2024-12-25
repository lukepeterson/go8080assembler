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
	START:
	MOV C, D
	JMP START

	`

	// Lex
	l := lexer.New(input)
	var tokens []lexer.Token
	for token := l.NextToken(); token.Type != "EOF"; token = l.NextToken() {
		tokens = append(tokens, token)
	}
	tokens = append(tokens, lexer.Token{Type: lexer.EOF})
	fmt.Printf("tokens: %v\n", tokens)

	// Parse
	p := parser.New(tokens)
	program, err := p.Parse()
	if err != nil {
		fmt.Println(err)
	}

	for _, b := range program {
		fmt.Printf("%02X ", b)
	}
	fmt.Println()

}
