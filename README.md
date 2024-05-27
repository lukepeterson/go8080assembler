# Intel 8080 CPU assembler in Go

What good is [an 8080 CPU emulator](https://github.com/lukepeterson/go8080cpu) without an assembler?  

This project takes an input string, tokenises the string into a slice of tokens, then parses those tokens, converting each to a valid 8080 opcode.

[![Tests](https://github.com/lukepeterson/go8080assembler/actions/workflows/go.yml/badge.svg)](https://github.com/lukepeterson/go8080assembler/actions/workflows/go.yml)

# Features

- :white_check_mark: Tokeniser
- :white_check_mark: Parser
- :x: Support for the full 8080 instruction set
- :x: Labels
- :x: Comments
- :x: Data (define byte, word, storage)
- :x: Input from STDIN

# Usage

```
code := `
	MVI A, 34h
	MOV B, C
	LDA 1234h
	HLT
`

assembler := &assembler.Assembler{}
assembler.Assemble(code)
for _, instruction := range assembler.ByteCode {
	fmt.Printf("%02X ", instruction)
}

// Prints "3E 34 41 3A 34 12 76"
```

# Running tests

Run `go test ./...`.