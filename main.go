package main

import (
	"fmt"

	"github.com/lukepeterson/go8080assembler/pkg/assembler"
)

func main() {

	input := `
			MVI A, 0x33
	START:	MOV B, C	; First comment
			LDA 0x1234	; Second comment
			JMP START
	`

	asm := assembler.New(input)
	bytecode, err := asm.Assemble()
	if err != nil {
		fmt.Println(err)
	}

	for _, b := range bytecode {
		fmt.Printf("%02X ", b)
	}
	fmt.Println()
}
