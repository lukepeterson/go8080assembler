package assembler

import (
	"reflect"
	"testing"
)

func TestAssemblerParseLine(t *testing.T) {
	tests := []struct {
		name      string
		line      string
		wantBytes []byte
		wantError bool
	}{
		{
			name:      "single byte instruction without comma",
			line:      "HLT",
			wantBytes: []byte{0x76},
			wantError: false,
		},
		{
			name:      "single byte instruction with comma",
			line:      "MOV B, B",
			wantBytes: []byte{0x40},
			wantError: false,
		},
		{
			name:      "two byte instruction",
			line:      "MVI B, 55H",
			wantBytes: []byte{0x06, 0x55},
			wantError: false,
		},
		{
			name:      "three byte instruction with explicit high byte",
			line:      "LDA 1234H",
			wantBytes: []byte{0x3A, 0x34, 0x12},
			wantError: false,
		},
		{
			name:      "three byte instruction with implicit high byte",
			line:      "LDA 34H",
			wantBytes: []byte{0x3A, 0x34, 0x00},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Assembler{}
			err := a.parseLine(tt.line)
			if (err != nil) != tt.wantError {
				t.Errorf("Assembler.parseLine() error = %v, wantErr %v", err, tt.wantError)
				return
			}
			if !reflect.DeepEqual(a.ByteCode, tt.wantBytes) {
				t.Errorf("Assembler.parseLine() = 0x%02X, want 0x%02X", a.ByteCode, tt.wantBytes)
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
			name:         "single line code",
			code:         `MOV B, M`,
			wantByteCode: []byte{0x46},
			wantErr:      false,
		},
		{
			name: "multi line code",
			code: `
				MVI A, 34h
				MOV B, C
				LDA 1234h
				HLT
			`,
			wantByteCode: []byte{0x3E, 0x34, 0x41, 0x3A, 0x34, 0x12, 0x76},
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
