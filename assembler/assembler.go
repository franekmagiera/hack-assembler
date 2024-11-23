package assembler

func Assemble(input []string) ([]string, error) {
	symbolTable := initSymbolTable()

	romAddress := 0
	for lineNumber, line := range input {
		command, err := parseLine(line, lineNumber)
		if err != nil {
			return nil, err
		}
		if command == nil {
			continue
		}
		switch command := command.(type) {
		case *CCommand, *ACommand:
			romAddress += 1
		case *LCommand:
			symbolTable[command.symbol] = romAddress
		}
	}

	output := make([]string, 0, len(input))
	for lineNumber, line := range input {
		command, err := parseLine(line, lineNumber)
		if command == nil && err == nil {
			continue
		}
		if err != nil {
			return nil, err
		}
		machineCode, err := processCommand(command, lineNumber, symbolTable)
		if machineCode == "" && err == nil {
			continue
		}
		if err != nil {
			return nil, err
		}
		output = append(output, machineCode)
	}
	return output, nil
}

func initSymbolTable() map[string]int {
	symbolTable := map[string]int{
		"NEXT_AVAILABLE_ADDRESS": 16, // Hacky, but let's see how it works.
		"SP":                     0,
		"LCL":                    1,
		"ARG":                    2,
		"THIS":                   3,
		"THAT":                   4,
		"R0":                     0,
		"R1":                     1,
		"R2":                     2,
		"R3":                     3,
		"R4":                     4,
		"R5":                     5,
		"R6":                     6,
		"R7":                     7,
		"R8":                     8,
		"R9":                     9,
		"R10":                    10,
		"R11":                    11,
		"R12":                    12,
		"R13":                    13,
		"R14":                    14,
		"R15":                    15,
		"SCREEN":                 16384,
		"KBD":                    24576,
	}
	return symbolTable
}
