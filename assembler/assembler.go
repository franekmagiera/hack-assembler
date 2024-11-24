package assembler

import (
	"bufio"
	"io"
)

type SymbolTable struct {
	nextAvailableAddress int
	table                map[string]int
}

func newSymbolTable() *SymbolTable {
	nextAvailableAddress := 16
	initialSymbolTable := map[string]int{
		"SP":     0,
		"LCL":    1,
		"ARG":    2,
		"THIS":   3,
		"THAT":   4,
		"R0":     0,
		"R1":     1,
		"R2":     2,
		"R3":     3,
		"R4":     4,
		"R5":     5,
		"R6":     6,
		"R7":     7,
		"R8":     8,
		"R9":     9,
		"R10":    10,
		"R11":    11,
		"R12":    12,
		"R13":    13,
		"R14":    14,
		"R15":    15,
		"SCREEN": 16384,
		"KBD":    24576,
	}
	return &SymbolTable{nextAvailableAddress, initialSymbolTable}
}

func (symbolTable *SymbolTable) assignNextAvailableAddress(symbol string) {
	symbolTable.table[symbol] = symbolTable.nextAvailableAddress
	symbolTable.nextAvailableAddress += 1
}

type MachineCodeProvider struct {
	input           *bufio.Scanner
	inputLineNumber int
	symbolTable     *SymbolTable
}

func (machineCodeProvider *MachineCodeProvider) NextLine() (string, error) {
	line := machineCodeProvider.input.Text()
	command, err := parseLine(line, machineCodeProvider.inputLineNumber)
	if command == nil && err == nil {
		if machineCodeProvider.ScanNextLine() {
			return machineCodeProvider.NextLine()
		} else {
			return "", nil
		}
	}
	if err != nil {
		return "", err
	}
	machineCode, err := processCommand(
		command,
		machineCodeProvider.inputLineNumber,
		machineCodeProvider.symbolTable,
	)
	if machineCode == "" && err == nil {
		if machineCodeProvider.ScanNextLine() {
			return machineCodeProvider.NextLine()
		} else {
			return "", nil
		}
	}
	if err != nil {
		return "", err
	}
	return machineCode, nil
}

func (machineCodeProvider *MachineCodeProvider) ScanNextLine() bool {
	machineCodeProvider.inputLineNumber += 1
	return machineCodeProvider.input.Scan()
}

func Assemble(input io.ReadSeeker) (*MachineCodeProvider, error) {
	symbolTable := newSymbolTable()
	if err := updateSymbolTable(input, symbolTable); err != nil {
		return nil, err
	}
	input.Seek(0, io.SeekStart)
	inputScanner := bufio.NewScanner(input)
	return &MachineCodeProvider{inputScanner, 0, symbolTable}, nil
}

func updateSymbolTable(input io.Reader, symbolTable *SymbolTable) error {
	romAddress := 0
	inputScanner := bufio.NewScanner(input)
	lineNumber := 0
	for inputScanner.Scan() {
		command, err := parseLine(inputScanner.Text(), lineNumber)
		if err != nil {
			return err
		}
		if command == nil {
			continue
		}
		switch command := command.(type) {
		case *CCommand, *ACommand:
			romAddress += 1
		case *LCommand:
			symbolTable.table[command.symbol] = romAddress
		}
		lineNumber += 1
	}
	return nil
}
