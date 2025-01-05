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
	; Simple 8080 Assembly Program to Test Lexer
	START:      LXI H, MSG          ; Load the address of MSG into HL
	            MVI C, 9            ; Load 9 (print string call) into register C
	            CALL PRINT          ; Call the PRINT subroutine

	            HLT                 ; Halt the program

	PRINT:      MOV A, M            ; Load the character at HL into A
	            ORA A               ; Check if the character is null (A OR A sets Z flag if A is 0)
	            RZ                  ; Return if zero (end of string)
	            OUT 1               ; Output the character (assumes device 1 is stdout)
	            INX H               ; Increment HL to point to the next character
	            JMP PRINT           ; Repeat the process

	MSG:        DB 'Hello, 8080!', 0 ; Message string (null-terminated)
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
