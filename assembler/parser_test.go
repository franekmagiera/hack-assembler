package assembler

import (
	"fmt"
	"testing"
)

func TestParsingEmptyLine(t *testing.T) {
	command, err := parseLine("", 1)
	if command != nil {
		t.Errorf("expected nil command for empty line, got: %s", command)
	}
	if err != nil {
		t.Errorf("expected nil error for empty line, got: %s", err)
	}
}

func TestParsingComment(t *testing.T) {
	command, err := parseLine("// some comment", 1)
	if command != nil {
		t.Errorf("expected nil command for a comment, got: %s", command)
	}
	if err != nil {
		t.Errorf("expected nil error for a comment, got: %s", err)
	}
}

func TestParsingACommand(t *testing.T) {
	command, err := parseLine("@100", 1)
	if err != nil {
		t.Fatalf("expected nil error for a valid A_COMMAND, got: %s", err)
	}
	switch c := command.(type) {
	case *ACommand:
		if c.symbol != "100" {
			t.Errorf("expected symbol value 100, got %s", c.symbol)
		}
	default:
		t.Errorf("expected to get an A_COMMAND, got %T", c)
	}
}

func TestParsingLCommand(t *testing.T) {
	command, err := parseLine("(LOOP)", 1)
	if err != nil {
		t.Fatalf("expected nil error for a valid L_COMMAND, got: %s", err)
	}
	switch c := command.(type) {
	case *LCommand:
		if c.symbol != "LOOP" {
			t.Errorf("expected symbol value LOOP, got %s", c.symbol)
		}
	default:
		t.Errorf("expected to get an L_COMMAND, got %T", c)
	}
}

func TestParsingCCommand(t *testing.T) {
	lineToExpectedCCommand := []struct {
		line             string
		expectedCCommand CCommand
	}{
		{"dest=comp;jump", CCommand{"dest", "comp", "jump"}},
		{"dest=comp", CCommand{"dest", "comp", ""}},
		{"comp;jump", CCommand{"", "comp", "jump"}},
	}

	for _, test := range lineToExpectedCCommand {
		line := test.line
		expectedCCommand := test.expectedCCommand
		command, err := parseLine(line, 1)
		if err != nil {
			t.Fatalf("expected nil error for a valid C_COMMAND, got: %s", err)
		}
		switch c := command.(type) {
		case *CCommand:
			if c.dest != expectedCCommand.dest {
				t.Errorf("expected dest value %s, got %s", expectedCCommand.dest, c.dest)
			}
			if c.comp != expectedCCommand.comp {
				t.Errorf("expected comp value %s, got %s", expectedCCommand.comp, c.comp)
			}
			if c.jump != expectedCCommand.jump {
				t.Errorf("expected jump value %s, got %s", expectedCCommand.jump, c.jump)
			}
		default:
			t.Errorf("expected to get a C_COMMAND, got %T", c)
		}
	}
}

func TestCommentsAreAllowedAfterCommands(t *testing.T) {
	testLines := []string{"@100 // comment", "(LOOP) // comment", "0;JMP // comment"}
	for _, testLine := range testLines {
		if _, err := parseLine(testLine, 1); err != nil {
			t.Errorf("should allow comments after command for line: %s", testLine)
		}
	}
}

func TestEmptySymbolInACommand(t *testing.T) {
	command, err := parseLine("@", 1)
	if command != nil {
		t.Errorf("expected nil command, got %s", command)
	}
	if err == nil {
		t.Fatal("expected an error, got a nil")
	}
	if expectedMessage := "empty symbol not allowed at line 1"; err.Error() != expectedMessage {
		t.Errorf("Expected '%s' message, got: %s", expectedMessage, err.Error())
	}
}

func TestEmptySymbolInLCommand(t *testing.T) {
	command, err := parseLine("()", 1)
	if command != nil {
		t.Errorf("expected nil command, got %s", command)
	}
	if err == nil {
		t.Fatal("expected an error, got a nil")
	}
	if expectedMessage := "empty symbol not allowed at line 1"; err.Error() != expectedMessage {
		t.Errorf("expected '%s' message, got: %s", expectedMessage, err.Error())
	}
}

func TestEmptyValuesInCCommand(t *testing.T) {
	lineNumber := 1
	couldNotParse := fmt.Sprintf("could not parse C_COMMAND at line %d", lineNumber)
	tests := []struct {
		line          string
		expectedError string
	}{
		{"dest=", couldNotParse},
		{"=comp;", couldNotParse},
		{"=;", couldNotParse},
		{"=", couldNotParse},
		{";", couldNotParse},
		{"dest=;jump", fmt.Sprintf("empty 'comp' in C_COMMAND at line %d", lineNumber)},
		{"dest;comp=jump", fmt.Sprintf("; should not appear before = in C_COMMAND at line %d", lineNumber)},
		{"nonsense", couldNotParse},
		{"dest=comp;jump (LOOP)", fmt.Sprintf("only comments are allowed after the command at line %d", lineNumber)},
	}

	for _, test := range tests {
		command, err := parseLine(test.line, lineNumber)
		if command != nil {
			t.Errorf("expected a nil command, got %s", command)
		}
		if err == nil {
			t.Fatal("expected an error, got a nil")
		}
		if err.Error() != test.expectedError {
			t.Errorf("for input: %s, expected error: %s, got: %s", test.line, test.expectedError, err.Error())
		}
	}
}
