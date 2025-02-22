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
	tokenList, err := scanner.ScanFile()
	if err != nil {
		Log.Fatalf("Error scanning file: %v", err)
		return
	}

	// Print the tokens
	for _, token := range tokenList {
		fmt.Println(token)
	}

	// list := []int{10, 20, 30, 40, 50}
	// for i := 0; i < len(list); i++ {
	// 	fmt.Println(list[i])
	// }

	// if err := scanner.Err(); err != nil {
	// 	log.Fatal(err)
	// }

	// Synthesize the tokens
	// err = synthesizeTokens(dfa)
	// fmt.Printf("%d lines, %d bytes\n", lineCount, byteCount)
}
