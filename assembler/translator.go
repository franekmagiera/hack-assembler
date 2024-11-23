package assembler

import (
	"fmt"
	"strconv"
)

// Translates the command to the machine code.
// Ignores L_COMMANDs returning empty string and nil error.
// Can modify the symbolTable when a new symbol is defined as a part of an A_COMMAND.
func processCommand(command Command, lineNumber int, symbolTable map[string]int) (string, error) {
	switch command := command.(type) {
	case *LCommand:
		// Ignore the L_COMMAND.
		return "", nil
	case *CCommand:
		compCode, ok := compTranslationMap[command.comp]
		if !ok {
			return "", fmt.Errorf("unknown value for comp: %s at line: %d", command.comp, lineNumber)
		}
		destCode, ok := destTranslationMap[command.dest]
		if !ok {
			return "", fmt.Errorf("unknown value for dest: %s at line: %d", command.dest, lineNumber)
		}
		jumpCode, ok := jumpTranslationMap[command.jump]
		if !ok {
			return "", fmt.Errorf("unknown value for jump: %s at line: %d", command.jump, lineNumber)
		}
		return fmt.Sprintf("111%s%s%s", compCode, destCode, jumpCode), nil
	case *ACommand:
		num, err := strconv.Atoi(command.symbol)
		if err == nil {
			if num > 32767 {
				return "", fmt.Errorf("number too big at line: %d", lineNumber)
			}
			return fmt.Sprintf("0%015b", num), nil
		}
		value, ok := symbolTable[command.symbol]
		if !ok {
			nextAvailableAddress := symbolTable["NEXT_AVAILABLE_ADDRESS"]
			symbolTable[command.symbol] = nextAvailableAddress
			symbolTable["NEXT_AVAILABLE_ADDRESS"] += 1
			return fmt.Sprintf("0%015b", symbolTable[command.symbol]), nil
		}
		return fmt.Sprintf("0%015b", value), nil
	}
	return "", fmt.Errorf("could not translate command at line: %d", lineNumber)
}

// Never modify those maps:
var compTranslationMap = map[string]string{
	"0":   "0101010",
	"1":   "0111111",
	"-1":  "0111010",
	"D":   "0001100",
	"A":   "0110000",
	"!D":  "0001101",
	"!A":  "0110001",
	"-D":  "0001111",
	"-A":  "0110011",
	"D+1": "0011111",
	"A+1": "0110111",
	"D-1": "0001110",
	"A-1": "0110010",
	"D+A": "0000010",
	"D-A": "0010011",
	"A-D": "0000111",
	"D&A": "0000000",
	"D|A": "0010101",
	"M":   "1110000",
	"!M":  "1110001",
	"-M":  "1110011",
	"M+1": "1110111",
	"M-1": "1110010",
	"D+M": "1000010",
	"D-M": "1010011",
	"M-D": "1000111",
	"D&M": "1000000",
	"D|M": "1010101",
}

var destTranslationMap = map[string]string{
	"":    "000",
	"M":   "001",
	"D":   "010",
	"MD":  "011",
	"A":   "100",
	"AM":  "101",
	"AD":  "110",
	"AMD": "111",
}

var jumpTranslationMap = map[string]string{
	"":    "000",
	"JGT": "001",
	"JEQ": "010",
	"JGE": "011",
	"JLT": "100",
	"JNE": "101",
	"JLE": "110",
	"JMP": "111",
}
