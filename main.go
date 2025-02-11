package main

import (
	"fmt"
	"os"
)

func main() {
	var err error
	var cli argsParsed

	// Parse the arguments
	cli, _ = GetCliArgs()

	// Initialize the logger
	initalizeLogger()

	var dfa *DFA

	file, err := os.Open(cli.args.FileName)
	if err != nil {
		Log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	// Initialize the dfa
	dfa, err = NewDFA()
	if err != nil {
		Log.Fatalf("Error initializing the DFA: %v", err)
		return
	}

	// Scan the file
	Log.Info(fmt.Sprintf("Scanner of %s initialized", file.Name()))
	err = scanFile(file, dfa)
	if err != nil {
		Log.Fatalf("Error scanning file: %v", err)
		return
	}

	// Parsing the tokens
	Log.Info(fmt.Sprintf("Parsing tokens of %s initalized", file.Name()))
	err = parseTokens(dfa)
	if err != nil {
		Log.Fatalf("Error parsing tokens: %v", err)
		return
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
