package assembler

import (
	"fmt"
	"strings"
)

type Assembler struct {
}

func parseLine(line string) []byte {

	line = strings.TrimSpace(line)
	if line == "" {
		return nil
	}

	tokens, _ := tokenise(line)
	for i, token := range tokens {
		fmt.Printf("token %d: %#v\n", i, token)
	}

	return []byte{}
}

func (a *Assembler) Assemble(code string) []byte {
	lines := strings.Split(code, "\n")
	for _, line := range lines {
		parseLine(line)
	}

	return []byte{0xFF}
}
