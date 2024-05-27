package main

import (
	"fmt"
	"log"

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
	err := assembler.Assemble(code)
	if err != nil {
		log.Fatal(err)
	}

	for _, instruction := range assembler.ByteCode {
		fmt.Printf("%02X ", instruction)
	}
}
