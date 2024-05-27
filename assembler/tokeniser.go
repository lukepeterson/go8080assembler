package assembler

import (
	"fmt"
	"strings"
)

// tokenise takes a line of code as an input, removes extra spacing/commas,
// and then returns the opcode and operands as separate strings.  The most
// important feature of this function is the ability to distinguish between
// single and multi-byte instructions and treat them as one.  For example:
// "MOV B, B" (one byte) and "MVI B" (two bytes).
func tokenise(line string) ([]string, error) {
	line = normalise(line)

	// Try and find a match for the opcode in our map.  We're doing an O(n)
	// search of the map as we don't yet know how long the instruction is.
	opcode, operand, err := extractOpcodeAndOperand(line)
	if err != nil {
		return nil, err
	}

	tokens := []string{opcode}
	if operand != "" { // we have an instruction that's more than one byte ..
		tokens = append(tokens, operand) // .. so append, ready for the parser
	}

	return tokens, nil
}

// normalise takes an input string, converts it to upper case, strips out extra spaces,
// and makes all comma formatting consistent, helping us match against out opcode map.
func normalise(line string) string {
	line = strings.ToUpper(line)

	// Remove all extra spaces to make matching opcode possible
	normalised := strings.Join(strings.Fields(line), " ")

	// Normalise spaces between and after the opcode/operand to match the opcodes
	// opcodes in our instruction set.  For example: "MOV A , B" -> "MOV A,B"
	normalised = strings.ReplaceAll(normalised, " ,", ",")
	normalised = strings.ReplaceAll(normalised, ", ", ",")

	return normalised
}

// extractOpcodeAndOperand extracts the opcode and operand from a normalised line.
func extractOpcodeAndOperand(line string) (string, string, error) {
	for opcode := range instructionSet {
		if strings.HasPrefix(line, opcode) {
			remaining := strings.TrimSpace(strings.TrimPrefix(line, opcode))
			operand := strings.TrimPrefix(remaining, ",")
			return opcode, operand, nil
		}
	}

	// If foundOpcode is still empty, we've found no match for the opcode in the instruction set
	return "", "", fmt.Errorf("no valid opcode found in line: %s", line)
}
