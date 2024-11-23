package assembler

import (
	"fmt"
	"testing"
)

func TestUnknownValues(t *testing.T) {
	lineNumber := 1
	tests := []struct {
		command       Command
		expectedError string
	}{
		{&CCommand{"D", "D-B", ""}, fmt.Sprintf("unknown value for comp: D-B at line: %d", lineNumber)},
		{&CCommand{"B", "M", ""}, fmt.Sprintf("unknown value for dest: B at line: %d", lineNumber)},
		{&CCommand{"", "0", "JUMP"}, fmt.Sprintf("unknown value for jump: JUMP at line: %d", lineNumber)},
	}
	for _, test := range tests {
		result, err := processCommand(test.command, lineNumber, make(map[string]int))
		if result != "" {
			t.Fatalf("expected empty result, got: %s", result)
		}
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != test.expectedError {
			t.Fatalf("expected error: %s, got %s", test.expectedError, err.Error())
		}
	}
}

func TestTooBigNumberInACommand(t *testing.T) {
	lineNumber := 1
	result, err := processCommand(&ACommand{"32768"}, lineNumber, make(map[string]int))
	if result != "" {
		t.Fatalf("expected empty result, got: %s", result)
	}
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	expectedError := fmt.Sprintf("number too big at line: %d", lineNumber)
	if err.Error() != expectedError {
		t.Fatalf("expected error: %s, got %s", expectedError, err.Error())
	}
}

func TestIgnoreLCommand(t *testing.T) {
	lineNumber := 1
	result, err := processCommand(&LCommand{"LOOP"}, lineNumber, make(map[string]int))
	if err != nil {
		t.Fatalf("expected nil error, got: %s", err.Error())
	}
	if result != "" {
		t.Fatalf("expected empty result, got: %s", result)
	}
}
