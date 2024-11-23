package assembler

import (
	"bytes"
	"errors"
	"io"
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
	output := Assemble(strings.NewReader(input))
	expectedOutput :=
		`0000000000000010
1110110000010000
0000000000000011
1110000010010000
0000000000000000
1110001100001000`
	if err := compare(output, strings.NewReader(expectedOutput)); err != nil {
		t.Fatal(err)
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
	output := Assemble(strings.NewReader(input))
	expectedOutput :=
		`0000000000000000
1111110000010000
0000000000010111
1110001100000110
0000000000010000
1110001100001000
0100000000000000
1110110000010000
0000000000010001
1110001100001000
0000000000010001
1111110000100000
1110111010001000
0000000000010001
1111110000010000
0000000000100000
1110000010010000
0000000000010001
1110001100001000
0000000000010000
1111110010011000
0000000000001010
1110001100000001
0000000000010111
1110101010000111`

	if err := compare(output, strings.NewReader(expectedOutput)); err != nil {
		t.Fatal(err)
	}
}

func compare(reader io.Reader, expectedReader io.Reader) error {
	outputBytes := make([]byte, 1024)
	expectedBytes := make([]byte, 1024)
	for {
		bytesRead, err1 := reader.Read(outputBytes)
		expectedBytesRead, err2 := expectedReader.Read(expectedBytes)
		if err1 != nil && err1 != io.EOF {
			return err1
		}
		if err2 != nil && err2 != io.EOF {
			return err2
		}
		if bytesRead != expectedBytesRead {
			return errors.New("output of a different size than expected")
		}
		if !bytes.Equal(outputBytes, expectedBytes) {
			return errors.New("unexpected output")
		}
		if (err1 == io.EOF) && (err2 == io.EOF) {
			return nil
		}
	}
}
