package main

import (
	"bufio"
	"io"
	"os"
	"testing"
)

func TestHackAssembler(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"test", "./test/Pong.asm"}

	main()

	outputFile, err := os.Open("./test/Pong.hack")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer func() {
		if err := outputFile.Close(); err != nil {
			t.Fatal(err.Error())
		}
	}()
	outputFileReader := bufio.NewReader(outputFile)

	expectedFile, err := os.Open("./test/ExpectedPong.hack")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer func() {
		if err := expectedFile.Close(); err != nil {
			t.Fatal(err.Error())
		}
	}()
	expectedFileReader := bufio.NewReader(expectedFile)

	for {
		expectedLine, err1 := expectedFileReader.ReadString('\n')
		outputLine, err2 := outputFileReader.ReadString('\n')
		if err1 != nil && err1 != io.EOF {
			t.Fatalf("error when reading file: %s", err1.Error())
		}
		if err2 != nil && err2 != io.EOF {
			t.Fatalf("error when reading file: %s", err2.Error())
		}
		if (err1 == io.EOF) != (err2 == io.EOF) {
			t.Fatal("unexpected number of lines")
		}
		if expectedLine != outputLine {
			t.Fatal("error when comparing files")
		}
		if (err1 == io.EOF) && (err2 == io.EOF) {
			break
		}
	}
}
