package main

import (
	"fmt"
	"os"
)

func main() {
	// Parse the arguments
	cli, _ := GetCliArgs()

	// Initialize the logger
	initalizeLogger()

	file, err := os.Open(cli.args.FileName)
	if err != nil {
		Log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close() // Close the file after the function returns

	viperASM(file) // assemble the file
}

func viperASM(file *os.File) {
	// Initialize the Scanner
	scanner, err := NewScanner(file)
	if err != nil {
		Log.Fatalf("Error initializing the Scanner: %v", err)
		return
	} else {
		Log.Info(fmt.Sprintf("Scanner of %s initialized", file.Name()))
	}

	// Scan the file for tokens
	tokens, err := scanner.ScanFile(file)
	if err != nil {
		Log.Fatalf("Error scanning file: %v", err)
		return
	}

	// // Process each line of tokens
	// for lineNum, lineTokens := range tokens {
	// 	fmt.Printf("Line %d:\n", lineNum+1)
	// 	for _, token := range lineTokens {
	// 		fmt.Printf("  %s\n", token)
	// 	}
	// }

	// Initialize the Parser
	parser, err := NewParser()
	if err != nil {
		Log.Fatalf("Error initializing the Parser: %v", err)
		return
	}

	// Parse the tokens to generate the instructions
	Log.Info(fmt.Sprintf("Parsing tokens of %s initalized", file.Name()))
	instructions, err := parser.ParseTokens(tokens)
	if err != nil {
		Log.Fatalf("Error parsing tokens: %v", err)
		return
	}

	// Print the instructions
	for _, instruction := range instructions {
		fmt.Printf("% s\n", instruction)
	}
}
