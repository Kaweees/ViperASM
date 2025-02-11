package main

import "fmt"

// Symbol lookup table
var symbolTable = map[string]int32{}

func parseTokens(dfa *DFA) error {
	lineCount := 0
	for _, tokenList := range dfa.totalTokens {
		// fmt.Println("Line: ", lineCount)
		fmt.Println("Tokens: ", tokenList)
		for _, token := range tokenList {
			if token.tokenType == "DotIdentifier" {
				if _, ok := symbolTable[token.tokenValue]; ok {
					return fmt.Errorf("duplicate label definition: %s", token.tokenValue)
				} else {
					symbolTable[token.tokenValue] = int32(lineCount)
				}
			}
		}
		lineCount += 1
	}
	return nil
}
