package main

import (
	"fmt"

	"github.com/lukepeterson/go8080assembler/assembler"
)

func main() {
	code := `
		MVI A, 34h
		MOV B, C
		LDA, 1234h
		HLT
	`

	assembler := &assembler.Assembler{}
	assembler.Assemble(code)

	for _, instruction := range assembler.ByteCode {
		fmt.Printf("%02X ", instruction)
	}
}
