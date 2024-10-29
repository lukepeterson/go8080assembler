# Intel 8080 CPU assembler in Go

This is my Intel 8080 CPU assembler, which I wrote to support [my Intel 8080 CPU emulator](https://github.com/lukepeterson/go8080cpu).  

It takes a newline separated string of Intel 8080 instructions, parses and validates the tokens, and then returns the assembled byte code.

[![Tests](https://github.com/lukepeterson/go8080assembler/actions/workflows/go.yml/badge.svg)](https://github.com/lukepeterson/go8080assembler/actions/workflows/go.yml)
![Go Report Card](https://goreportcard.com/badge/github.com/lukepeterson/go8080assembler)
![GitHub release](https://img.shields.io/github/v/release/lukepeterson/go8080assembler)

# Features

- :white_check_mark: Tokeniser
- :white_check_mark: Parser
- :white_check_mark: Comment support
- :white_check_mark: Supports full (244) 8080 CPU instructions

# TODO

- Label support
- Data support (eg, `DB`, `DW`)
- Input from `STDIN`

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
