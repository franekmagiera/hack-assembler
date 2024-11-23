package assembler

import (
	"slices"
	"strings"
	"testing"
)

func TestAddProgram(t *testing.T) {
	input := ` 
		@2
		D=A
		@3
		D=D+A
		@0
		M=D
		`

	output, err := Assemble(strings.NewReader(input))
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	expectedOutput := []string{
		"0000000000000010",
		"1110110000010000",
		"0000000000000011",
		"1110000010010000",
		"0000000000000000",
		"1110001100001000",
	}

	if !slices.Equal(output, expectedOutput) {
		t.Log("expected:\n")
		for i, code := range expectedOutput {
			t.Logf("\t%d %s\n", i, code)
		}
		t.Log("got:\n")
		for i, code := range output {
			t.Logf("\t%d %s\n", i, code)
		}
		t.FailNow()
	}
}

func TestRectProgram(t *testing.T) {
	input := ` 
		// If (R0 <= 0) goto END else n = R0"
		@R0
		D=M
		@END
		D;JLE
		@n
		M=D
		// addr = base address of first screen row
		@SCREEN
		D=A
		@addr
		M=D
		(LOOP)
		// RAM[addr] = -1
		@addr
		A=M
		M=-1
		// addr = base address of next screen row
		@addr
		D=M
		@32
		D=D+A
		@addr
		M=D
		// decrements n and loops
		@n
		MD=M-1
		@LOOP
		D;JGT
		(END)
		@END
		0;JMP
	`

	output, err := Assemble(strings.NewReader(input))
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	expectedOutput := []string{
		"0000000000000000",
		"1111110000010000",
		"0000000000010111",
		"1110001100000110",
		"0000000000010000",
		"1110001100001000",
		"0100000000000000",
		"1110110000010000",
		"0000000000010001",
		"1110001100001000",
		"0000000000010001",
		"1111110000100000",
		"1110111010001000",
		"0000000000010001",
		"1111110000010000",
		"0000000000100000",
		"1110000010010000",
		"0000000000010001",
		"1110001100001000",
		"0000000000010000",
		"1111110010011000",
		"0000000000001010",
		"1110001100000001",
		"0000000000010111",
		"1110101010000111",
	}

	if !slices.Equal(output, expectedOutput) {
		t.Log("expected:\n")
		for i, code := range expectedOutput {
			t.Logf("\t%d %s\n", i, code)
		}
		t.Log("got:\n")
		for i, code := range output {
			t.Logf("\t%d %s\n", i, code)
		}
		t.FailNow()
	}
}
