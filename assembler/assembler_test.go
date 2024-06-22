package assembler

import (
	"reflect"
	"testing"
)

func TestAssemblerParseLine(t *testing.T) {
	tests := []struct {
		name         string
		line         string
		wantByteCode []byte
		wantErr      bool
	}{
		{
			name:         "single byte instruction without comma",
			line:         "HLT",
			wantByteCode: []byte{0x76},
			wantErr:      false,
		},
		{
			name:         "single byte instruction with comma",
			line:         "MOV B, B",
			wantByteCode: []byte{0x40},
			wantErr:      false,
		},
		{
			name:         "two byte instruction",
			line:         "MVI B, 55H",
			wantByteCode: []byte{0x06, 0x55},
			wantErr:      false,
		},
		{
			name:         "two byte instruction with missing operand",
			line:         "MVI B",
			wantByteCode: nil,
			wantErr:      true,
		},
		{
			name:         "three byte instruction with missing operand",
			line:         "JNZ",
			wantByteCode: nil,
			wantErr:      true,
		},
		{
			name:         "three byte instruction with explicit high byte",
			line:         "LDA 1234H",
			wantByteCode: []byte{0x3A, 0x34, 0x12},
			wantErr:      false,
		},
		{
			name:         "three byte instruction with implicit high byte",
			line:         "LDA 34H",
			wantByteCode: []byte{0x3A, 0x34, 0x00},
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Assembler{}
			err := a.parseLine(tt.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("Assembler.parseLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(a.ByteCode, tt.wantByteCode) {
				t.Errorf("Assembler.parseLine() ByteCode = 0x%02X, wantByteCode 0x%02X", a.ByteCode, tt.wantByteCode)
			}
		})
	}
}

func TestAssemblerParseHex(t *testing.T) {
	tests := []struct {
		name     string
		token    string
		wantHigh uint8
		wantLow  uint8
		wantErr  bool
	}{
		{
			name:     "one byte input with suffix",
			token:    "1AH",
			wantHigh: 0x00,
			wantLow:  0x1A,
			wantErr:  false,
		},
		{
			name:     "one byte input with prefix",
			token:    "0x2A",
			wantHigh: 0x00,
			wantLow:  0x2A,
			wantErr:  false,
		},
		{
			name:     "one byte input without suffix",
			token:    "1B",
			wantHigh: 0x00,
			wantLow:  0x1B,
			wantErr:  false,
		},
		{
			name:     "two byte input with suffix",
			token:    "1D2DH",
			wantHigh: 0x1D,
			wantLow:  0x2D,
			wantErr:  false,
		},
		{
			name:     "two byte input with prefix",
			token:    "0x3D4D",
			wantHigh: 0x3D,
			wantLow:  0x4D,
			wantErr:  false,
		},
		{
			name:     "two byte input without prefix or suffix",
			token:    "3A4A",
			wantHigh: 0x3A,
			wantLow:  0x4A,
			wantErr:  false,
		},
		{
			name:     "one byte input without leading zero, prefix or suffix",
			token:    "7",
			wantHigh: 0x00,
			wantLow:  0x07,
			wantErr:  false,
		},
		{
			name:     "one byte input with suffix and without leading zero",
			token:    "8H",
			wantHigh: 0x00,
			wantLow:  0x08,
			wantErr:  false,
		},
		{
			name:     "one byte input with prefix and without leading zero",
			token:    "0xA",
			wantHigh: 0x00,
			wantLow:  0x0A,
			wantErr:  false,
		},
		{
			name:     "two byte input with suffix and without leading zero",
			token:    "5A6H",
			wantHigh: 0x05,
			wantLow:  0xA6,
			wantErr:  false,
		},
		{
			name:     "two byte input with prefix and without leading zero",
			token:    "0xDB9",
			wantHigh: 0x0D,
			wantLow:  0xB9,
			wantErr:  false,
		},
		{
			name:     "two bytes at end of 16 bit address space (65,535 bytes)",
			token:    "0xFFFF",
			wantHigh: 0xFF,
			wantLow:  0xFF,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHigh, gotLow, err := parseHex(tt.token)
			if gotHigh != tt.wantHigh {
				t.Errorf("parseHex() gotHigh = %02X, wantHigh %02X", gotHigh, tt.wantHigh)
			}
			if gotLow != tt.wantLow {
				t.Errorf("parseHex() gotLow = %02X, wantLow %02X", gotLow, tt.wantLow)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("parseHex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestAssemblerAssemble(t *testing.T) {
	tests := []struct {
		name         string
		code         string
		wantByteCode []byte
		wantErr      bool
	}{
		{
			name:         "single line of code",
			code:         `MOV B, M`,
			wantByteCode: []byte{0x46},
			wantErr:      false,
		},
		{
			name:         "empty code",
			code:         "",
			wantByteCode: nil,
			wantErr:      false,
		},
		{
			name:         "whitespace only",
			code:         "    ",
			wantByteCode: nil,
			wantErr:      false,
		},
		{
			name:         "invalid opcode",
			code:         "FOO",
			wantByteCode: nil,
			wantErr:      true,
		},
		{
			name:         "single line with comment",
			code:         "MVI A, 34h ;Load immediate value",
			wantByteCode: []byte{0x3E, 0x34},
			wantErr:      false,
		},
		{
			name:         "single line with comment and no instruction",
			code:         " ; Do nothing",
			wantByteCode: nil,
			wantErr:      false,
		},
		{
			name:         "correct decoding of INR A",
			code:         "INR A",
			wantByteCode: []byte{0x3C},
			wantErr:      false,
		},
		{
			name: "correct decoding of INR A, avoiding collision with IN",
			code: `
				INR A
				IN 44H`,
			wantByteCode: []byte{0x3C, 0xDB, 0x44},
			wantErr:      false,
		},
		{
			name:         "correct decoding of IN",
			code:         "IN 33H",
			wantByteCode: []byte{0xDB, 0x33},
			wantErr:      false,
		},
		{
			name: "correct decoding of IN, avoiding collision with INR A",
			code: `
				IN 33H
				INR A`,
			wantByteCode: []byte{0xDB, 0x33, 0x3C},
			wantErr:      false,
		},
		{
			name: "multiple lines with comments",
			code: `
				MVI A, 34H ;Load immediate value
				           ;This is a comment on an empty line
				MOV B, C   ;Move C to B
				HLT        ;Halt`,
			wantByteCode: []byte{0x3E, 0x34, 0x41, 0x76},
			wantErr:      false,
		},
		{
			name: "multi line code",
			code: `
				MVI A, 34H
				MOV B, C
				LDA 1234H
				HLT
			`,
			wantByteCode: []byte{0x3E, 0x34, 0x41, 0x3A, 0x34, 0x12, 0x76},
			wantErr:      false,
		},
		{
			name: "end of address space",
			code: `
				LXI H, FFFFH
				HLT
			`,
			wantByteCode: []byte{0x21, 0xFF, 0xFF, 0x76},
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Assembler{}
			err := a.Assemble(tt.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("Assembler.Assemble() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(a.ByteCode, tt.wantByteCode) {
				t.Errorf("Assembler.Assemble() ByteCode = %02X, wantByteCode %02X", a.ByteCode, tt.wantByteCode)
			}
		})
	}
}
