// Package assembler provides functions for assembling Intel 8080 CPU assembler commands.
package assembler

import (
	"fmt"
	"strconv"
	"strings"
)

// Assembler stores our assembled bytecode
type Assembler struct {
	ByteCode []byte
	// TODO: Add the types.Word type in place of uint16
	Labels map[string]uint16
}

func New() *Assembler {
	return &Assembler{
		// TODO: Test this:
		// ByteCode: make([]byte, 0),
		Labels: make(map[string]uint16, 100),
	}
}

// parseLine takes an input string, and then performs a number of operations:
// - trims all empty lines
// - passes each string to tokenise()
// - converts each instruction to its corresponding opcode byte
// - converts all two and three byte operands to little endian (low byte first)
// - updates the ByteCode field with converted values
func (a *Assembler) parseLine(line string) error {

	line = removeComment(line)

	// Trim the line to remove leading and trailing whitespace
	// If the line is empty after trimming, return nil to skip processing
	line = strings.TrimSpace(line)
	if line == "" {
		return nil
	}

	// fmt.Printf("\nline: %v\n", line)
	// fmt.Printf("len(a.ByteCode): %v\n", len(a.ByteCode))

	// Check for a label
	index := strings.Index(line, ":")
	var label string
	if index != -1 {
		// found a label
		// fmt.Printf("********************* found a label at index %v\n", index)

		label = line[0:index]

		a.Labels[label] = uint16(len(a.ByteCode))

		line = strings.TrimSpace(line[index+1:])

	}
	// fmt.Printf("line:  \"%v\"\n", line)
	// fmt.Printf("label: \"%v\"\n", label)

	tokens, err := tokenise(line)
	if err != nil {
		return err
	}

	// fmt.Printf("tokens  : %q\n", tokens)
	// fmt.Printf("len(tokens): %v\n", len(tokens))
	// fmt.Printf("a.Labels: %q\n", a.Labels)

	// Check the first token in the list for a matching instruction
	instruction, exists := instructionSet[tokens[0]]
	if !exists {
		return fmt.Errorf("unknown instruction: %v", tokens[0])
	}

	switch instruction.Length {
	case 1: // Single byte, so return the opcode and carry on
		a.ByteCode = append(a.ByteCode, instruction.Opcode)
	case 2: // Two bytes (a word), so return the opcode and the next instruction converted to uint8
		if len(tokens) < 2 {
			return fmt.Errorf("missing operand for instruction: %v", tokens[0])
		}
		// parseOperand returns two bytes, but we can ignore the first if we're only expecting one
		_, lowByte, err := a.parseOperand(tokens[1])
		if err != nil {
			return err
		}

		a.ByteCode = append(a.ByteCode, instruction.Opcode, byte(lowByte))
	case 3: // Three bytes, so return the opcode and the next instruction converted to two uint8s
		if len(tokens) < 2 {
			return fmt.Errorf("missing operands for instruction: %v", tokens[0])
		}

		highByte, lowByte, err := a.parseOperand(tokens[1])
		if err != nil {
			return err
		}

		// Split the 16-bit int into high and low order bytes.  The 8080 CPU is little endian,
		// so we split the bytes and return the low order byte first.
		a.ByteCode = append(a.ByteCode, instruction.Opcode, byte(lowByte), byte(highByte))
	default:
		return fmt.Errorf("invalid instruction length for: %v", tokens[0])
	}

	return nil
}

// Assemble takes a newline separated string of code, parses the input into tokens, and then converts each instruction to a valid opcode from the instructionSet map.
func (a *Assembler) Assemble(code string) error {
	lines := strings.Split(code, "\n")
	for _, line := range lines {
		err := a.parseLine(line)
		if err != nil {
			return err
		}
	}

	return nil
}

func removeComment(line string) string {
	index := strings.Index(line, ";")
	if index != -1 {
		line = line[:index]
	}

	return line
}

// parseHex takes a string with an H suffix or 0x prefix and parses it into a 16 bit integer, returning the result as two int8s.  This has a nice side effect of being able to take a one byte string and also returning the result as two 8-bit integers, which is required for our two-byte instructions.  For example: "4AH" -> 0x00, 0x4A.
func (a Assembler) parseOperand(token string) (uint8, uint8, error) {

	// TODO: This function needs to first check whether the token is a label

	// fmt.Printf("checking operand token: %v\n", token)

	address, exists := a.Labels[token]
	if exists {
		// fmt.Printf("FOUND TOKEN IN MAP... return addresss %v", address)
		return 0, uint8(address), nil
	}

	// Check the first token in the list for a matching instruction
	// instruction, exists := instructionSet[tokens[0]]
	// if !exists {
	// 	return fmt.Errorf("unknown instruction: %v", tokens[0])
	// }

	token = strings.TrimSuffix(token, "H")
	token = strings.TrimPrefix(token, "0X")
	token = strings.TrimPrefix(token, "0x")
	hex, err := strconv.ParseUint(token, 16, 16)
	if err != nil {

		return 0, 0, err
	}

	highByte := uint8(hex >> 8)
	lowByte := uint8(hex & 0x00FF)

	return highByte, lowByte, nil
}
