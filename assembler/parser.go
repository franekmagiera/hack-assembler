package assembler

import (
	"bufio"
	"fmt"
	"strings"
)

type (
	Command interface {
		command()
	}

	ACommand struct {
		symbol string
	}

	CCommand struct {
		dest string
		comp string
		jump string
	}

	LCommand struct {
		symbol string
	}
)

// command() ensures that only commands can be assigned to Command.
func (*ACommand) command() {}
func (*CCommand) command() {}
func (*LCommand) command() {}

// Can return nil if result should be ignored (for example a comment).
func parseLine(line string, lineNumber int) (Command, error) {
	scanner := bufio.NewScanner(strings.NewReader(line))
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("error reading input at line %d", lineNumber)
		}
		// Ignore an empty line.
		return nil, nil
	}
	firstWord := scanner.Text()
	if strings.HasPrefix(firstWord, "//") {
		// Ignore comment.
		return nil, nil
	}
	var command Command
	if strings.HasPrefix(firstWord, "@") {
		// It has to be an "A_COMMAND".
		symbol := firstWord[1:]
		if len(symbol) == 0 {
			return nil, fmt.Errorf("empty symbol not allowed at line %d", lineNumber)
		}
		command = &ACommand{symbol}
	} else if strings.HasPrefix(firstWord, "(") {
		// It has to be an "L_COMMAND".
		symbol := firstWord[1 : len(firstWord)-1]
		if len(symbol) == 0 {
			return nil, fmt.Errorf("empty symbol not allowed at line %d", lineNumber)
		}
		command = &LCommand{symbol}
	} else {
		// Anything else has to be a "C_COMMAND".
		equalSignIndex, semicolonIndex := strings.Index(firstWord, "="), strings.Index(firstWord, ";")
		if equalSignIndex == len(firstWord)-1 ||
			semicolonIndex == len(firstWord)-1 ||
			equalSignIndex == 0 ||
			semicolonIndex == 0 {
			return nil, fmt.Errorf("could not parse C_COMMAND at line %d", lineNumber)
		}
		if equalSignIndex != -1 && semicolonIndex != -1 {
			if equalSignIndex > semicolonIndex {
				return nil, fmt.Errorf("; should not appear before = in C_COMMAND at line %d", lineNumber)
			}
			if semicolonIndex-equalSignIndex == 1 {
				return nil, fmt.Errorf("empty 'comp' in C_COMMAND at line %d", lineNumber)
			}
			dest := firstWord[:equalSignIndex]
			comp := firstWord[equalSignIndex+1 : semicolonIndex]
			jump := firstWord[semicolonIndex+1:]
			command = &CCommand{dest, comp, jump}
		}
		if equalSignIndex == -1 && semicolonIndex != -1 {
			// If dest is empty, "=" is omitted.
			comp := firstWord[:semicolonIndex]
			jump := firstWord[semicolonIndex+1:]
			command = &CCommand{"", comp, jump}
		}
		if equalSignIndex != -1 && semicolonIndex == -1 {
			// If jump is empty, ";" is ommited.
			dest := firstWord[:equalSignIndex]
			comp := firstWord[equalSignIndex+1:]
			command = &CCommand{dest, comp, ""}
		}
		if equalSignIndex == -1 && semicolonIndex == -1 {
			return nil, fmt.Errorf("could not parse C_COMMAND at line %d", lineNumber)
		}
	}
	if scanner.Scan() && !strings.HasPrefix(scanner.Text(), "//") {
		return nil, fmt.Errorf("only comments are allowed after the command at line %d", lineNumber)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading input at line %d", lineNumber)
	}
	return command, nil
}
