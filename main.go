package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/franekmagiera/hack-assembler/assembler"
)

func main() {
	var args = os.Args
	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "hack assembler requires exactly one argument - the file path of the assembly"+
			"program: amount of arguments provided: %d\n", len(args)-1)
		return
	}
	path := args[1]
	if filepath.Ext(path) != ".asm" {
		fmt.Fprintln(os.Stderr, "only .asm files are allowed")
		return
	}
	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}()
	// Just storing the whole file in memory for simplicity.
	scanner := bufio.NewScanner(file)
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error reading file: ", err)
		return
	}
	machineCode, err := assembler.Assemble(lines)
	if err == nil {
		dir := filepath.Dir(path)
		filename := filepath.Base(path)
		filename = filename[:strings.Index(filename, ".")]
		outputFile, err := os.Create(fmt.Sprintf("%s/%s.hack", dir, filename))
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
		defer func() {
			if err := outputFile.Close(); err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
			}
		}()

		for _, code := range machineCode {
			outputFile.WriteString(fmt.Sprintln(code))
		}
	}
}
