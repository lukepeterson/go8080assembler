package assembler

import (
	"fmt"
	"strings"
)

// tokenise takes a line of code as an input, removes extra spacing/commas, and then returns the opcode
// and operands as separate strings.  Does no parsing.  Most important feature of this function is the
// ability to distinguish between single and multi-byte instructions and treat them as one.  For example:
// "MOV B, B" (one byte) and "MVI B" (two bytes).
func tokenise(line string) ([]string, error) {

	line = strings.ToUpper(line)

	// Remove all extra spaces that we don't need
	normalised := strings.Join(strings.Fields(line), " ")

	// Normalise spaces between and after the opcode/operand to match the opcodes
	// opcodes in our instruction set.  For example: "MOV A , B" -> "MOV A,B"
	normalised = strings.ReplaceAll(normalised, " ,", ",")
	normalised = strings.ReplaceAll(normalised, ", ", ",")

	tokens := []string{}

	// Try and find a match for the opcode in our map.  We're doing an O(n)
	// search of the map as we don't yet know how long the instruction is.
	var foundOpcode string
	for opcode := range instructionSet {
		if strings.HasPrefix(normalised, opcode) { // Check the start of the instruction matches an opcode
			foundOpcode = opcode
			tokens = append(tokens, foundOpcode)
			// Remove the opcode from the remainder of the instruction, so we can extract the operand
			normalised = strings.TrimSpace(strings.TrimPrefix(normalised, foundOpcode))
			break
		}
	}

	// If foundOpcode is still empty, we've found no match for the opcode in the instruction set
	if foundOpcode == "" {
		return nil, fmt.Errorf("no valid opcode found in line: %s", line)
	}

	// If we have anything left, we have a multi-byte instruction ..
	if normalised != "" {
		// .. so treat anything after the "," as the operand, ready for the parser to interpret
		operand := strings.TrimPrefix(normalised, ",")
		tokens = append(tokens, operand)
	}

	return tokens, nil
}
