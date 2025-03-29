package assembler

import (
	"github.com/lukepeterson/go8080assembler/pkg/lexer"
	"github.com/lukepeterson/go8080assembler/pkg/parser"
)

type Assembler struct {
	input    string
	bytecode []byte
}

func New(input string) *Assembler {
	return &Assembler{input: input}
}

func (a *Assembler) Assemble() ([]byte, error) {
	l := lexer.New(a.input)
	tokens, err := l.Lex()
	if err != nil {
		return nil, err
	}

	p := parser.New(tokens)
	a.bytecode, err = p.Parse()
	if err != nil {
		return nil, err
	}

	return a.bytecode, nil
}
