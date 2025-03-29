# Intel 8080 CPU assembler in Go

This is my Intel 8080 CPU assembler, which I wrote to support [my Intel 8080 CPU emulator](https://github.com/lukepeterson/go8080cpu).

It takes a newline separated string of Intel 8080 instructions, performs a lexical analysis on the input, parses and validates the tokens, and then converts the tokens to an assembled byte code.

[![Tests](https://github.com/lukepeterson/go8080assembler/actions/workflows/go.yml/badge.svg)](https://github.com/lukepeterson/go8080assembler/actions/workflows/go.yml)
![Go Report Card](https://goreportcard.com/badge/github.com/lukepeterson/go8080assembler)
![GitHub release](https://img.shields.io/github/v/release/lukepeterson/go8080assembler)

# Features

- :white_check_mark: Lexer
- :white_check_mark: Parser
- :white_check_mark: Comment support
- :white_check_mark: Label support
- :white_check_mark: Supports all 244 8080 CPU instructions

# TODO

- Data support (eg, `DB`, `DW`)
- Input from `STDIN`

# Usage

See `main.go` for examples on how to use both the lexer and the parser.

# Running tests

Run `go test ./...`.
