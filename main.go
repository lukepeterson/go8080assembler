package main

import (
	"fmt"
	"log"

	"github.com/lukepeterson/go8080assembler/assembler"
)

func main() {
	code := `
		MVI A, 0x33
TESTX:	MOV B, C
		LDA, 0x1234
		JMP TESTX
		HLT
	`

	asm := assembler.New()

	err := asm.Assemble(code)
	if err != nil {
		log.Fatal(err)
	}

	for _, instruction := range asm.ByteCode {
		fmt.Printf("%02X ", instruction)
	}
}
