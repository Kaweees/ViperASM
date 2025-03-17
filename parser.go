package main

import "fmt"

// import "fmt"

// Represents an address
type Address int32

// Symbol lookup table
var symbolTable = map[string]Address{}

// Represents a lexical parser
type Parser struct {
	SymbolTable map[string]Address
}

// Constructor to initialize memory for the Parser.
func NewParser() (*Parser, error) {
	parser := &Parser{}
	parser.SymbolTable = symbolTable
	return parser, nil
}

// Parse the tokens to generate the instructions
func (p *Parser) ParseTokens(tokens [][]Token) ([]Instruction, error) {
	lineCount := 0
	instructions := []Instruction{}
	for i, tokenList := range tokens {
		// fmt.Println("Tokens: ", tokenList)
		for j, token := range tokenList {
			if token.Type == DIRECTIVE {
				// remove the leading dot from directive
				tokens[i][j].Literal = token.Literal[1:]
			} else if token.Type == LABEL_DEF {
				// remove the colon from label definition
				tokens[i][j].Literal = token.Literal[:len(token.Literal)-1]
				if _, ok := symbolTable[tokens[i][j].Literal]; ok {
					return nil, fmt.Errorf("duplicate label definition: %s", tokens[i][j].Literal)
				} else {
					symbolTable[tokens[i][j].Literal] = Address(lineCount)
				}
			}
		}
		fmt.Println("Tokens: ", tokenList)
		lineCount += 1
	}
	return instructions, nil
}
