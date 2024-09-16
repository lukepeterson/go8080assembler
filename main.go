package main

// import "github.com/lukepeterson/go8080assembler/assembler"
import (
	"fmt"

	"github.com/lukepeterson/go8080assembler/pkg/lexer"
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

	input := `
			MVI A, 0x33
	TEST:	MOV B, C   ; The first comment
			LDA, 1234h ; The second comment
			JMP TEST

		; The Third comment with a colon ; here
			`

	myLexer := lexer.New(input)
	fmt.Printf("myLexer: %v\n", myLexer)

}
