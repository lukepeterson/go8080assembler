package assembler

import (
	"fmt"
	"strings"
)

func tokenise(line string) ([]string, error) {
	line = strings.ToUpper(line)
	// Remove all extra spaces that we don't need
	normalised := strings.Join(strings.Fields(line), " ")

	// Normalise spaces between and after the opcode/operand to match our
	// instruction set.  For example: "MOV A , B" -> "MOV A,B"
	normalised = strings.ReplaceAll(normalised, " ,", ",")
	normalised = strings.ReplaceAll(normalised, ", ", ",")

	tokens := []string{}

	// Try and find a match for the opcode in our map.  We're doing an O(n)
	// search of the map as we don't yet know how long the instruction is.
	var foundOpcode string
	for opcode := range instructionSet {
		if strings.HasPrefix(normalised, opcode) { // Thankfully, there are no subsets in our instruction set
			foundOpcode = opcode
			tokens = append(tokens, foundOpcode)
			// Remove the opcode from the remainder of the instruction, so we can extract the operand
			normalised = strings.TrimSpace(strings.TrimPrefix(normalised, foundOpcode))
			break
		}
	}

	if foundOpcode == "" {
		return nil, fmt.Errorf("no opcode found")
	}

	// If we have anything left, we have a multi-byte instruction ..
	// fmt.Printf("normalised: %v\n", normalised)
	if normalised != "" {
		// .. so treat anything after the "," as the operand
		operand := strings.TrimPrefix(normalised, ",")
		tokens = append(tokens, operand)
	}

	return tokens, nil
}
