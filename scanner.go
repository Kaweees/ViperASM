package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// Represents a Deterministic Finite Automaton(DFA) based lexical scanner.
type Scanner struct {
	File   *os.File
	Token  Token
	Tokens []Token
	reader *bufio.Reader
	// pos     int
	// readPos int
	// ch      byte
	// line    int
	// column  int
}

// Constructor to initialize memory for the Scanner.
func NewScanner(file *os.File) (*Scanner, error) {
	scanner := &Scanner{
		File:   file,
		Token:  Token{},
		Tokens: []Token{},
		reader: nil,
		// pos:     0,
		// readPos: 0,
		// ch:     0,
		// line:   1,
		// column: 1,
	}
	return scanner, nil
}

// Scan the file for tokens.
func (scanner *Scanner) ScanFile() ([]Token, error) {
	scanner.reader = bufio.NewReader(scanner.File)
	for {
		b, err := scanner.reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		if b == 'h' {
			fmt.Print("h")
		}
	}
	return scanner.Tokens, nil
}

// func scanFile(file *os.File, dfa *DFA) error {
// 	scanner := bufio.NewScanner(file)
//
// 	for scanner.Scan() {
// 		for _, r := range scanner.Text() {
// 			dfa.Transition(r)
// 			// fmt.Print(i, r)
// 			// fmt.Printf("Index: %d, Rune: %c\n", i, r)
// 		}
// 		dfa.Store()
// 		dfa.StoreLine()
// 	}
// 	return nil
// }


// // Represents the scanning Deterministic Finite Automaton(DFA) for the scanner.
// type DFA struct {
// 	currentState  State
// 	currentToken  string
// 	currentString rune
// 	tokens        []Token
// 	totalTokens   [][]Token
// }

// // Constructor to initialize memory for the DFA.
// func NewDFA() (*DFA, error) {
// 	dfa := &DFA{}
// 	dfa.currentState = Initial
// 	dfa.currentToken = ""
// 	dfa.currentString = 0
// 	dfa.tokens = []Token{}
// 	dfa.totalTokens = [][]Token{}
// 	return dfa, nil
// }

// // Add a token to the list of tokens.
// func (dfa *DFA) AddToken(Type string, Value string) {
// 	dfa.tokens = append(dfa.tokens, Token{Type, Value})
// }

// Transition the DFA to a new state based on the input.
func (dfa *DFA) Transition(input rune) {
	// fmt.Printf("State: %s, Rune: '%c'\n", dfa.currentState.String(), input)
	switch dfa.currentState {
	case Initial:
		if input == '.' {
			dfa.currentState = DotIdentifier
		} else if input == '$' {
			dfa.currentToken = string(input)
			dfa.currentState = Register
		} else if input == '0' {
			dfa.currentState = Zero
		} else if unicode.IsDigit(input) || input == '-' {
			dfa.currentToken = string(input)
			dfa.currentState = Decimal
		} else if input == ',' {
			dfa.AddToken(dfa.currentState.String(), ",")
			dfa.Reset()
		} else if input == '(' {
			dfa.AddToken(dfa.currentState.String(), "(")
			dfa.Reset()
		} else if input == ')' {
			dfa.AddToken(dfa.currentState.String(), ")")
			dfa.Reset()
		} else if input == '#' || input == ';' {
			dfa.Store()
			dfa.currentState = Comment
		} else if input == '"' || input == '\'' {
			dfa.currentString = input
			dfa.currentToken = string(input)
			dfa.currentState = String
		} else if !unicode.IsSpace(input) {
			dfa.currentToken = string(input)
			dfa.currentState = Identifier
		}
	case Identifier:
		if input == ':' {
			dfa.currentState = LabelDef
			dfa.Store()
		} else if !unicode.IsSpace(input) {
			dfa.currentToken += string(input)
		} else {
			dfa.Store()
		}
	case DotIdentifier:
		if unicode.IsLetter(input) {
			dfa.currentToken += string(input)
		} else {
			dfa.Store()
		}
	case Register:
		if input == ',' {
			dfa.Store()
			dfa.AddToken(dfa.currentState.String(), ",")
		} else if unicode.IsDigit(input) || unicode.IsLetter(input) {
			dfa.currentToken += string(input)
		} else {
			dfa.Store()
		}
	case Zero:
		if input == 'x' {
			dfa.currentState = Hexadecimal
		} else if unicode.IsDigit(input) {
			dfa.currentToken = string(input)
			dfa.currentState = Decimal
		} else {
			dfa.Store()
		}
	case Decimal:
		if unicode.IsDigit(input) {
			dfa.currentToken += string(input)
		} else {
			dfa.Store()
		}
	case Hexadecimal:
		if strings.ContainsAny(string(input), "0123456789abcdefABCDEF") {
			dfa.currentToken += string(input)
		} else {
			dfa.Store()
		}
	case String:
		dfa.currentToken += string(input)
		if input == dfa.currentString {
			dfa.Store()
		}
	}
}

func scanFile(file *os.File, dfa *DFA) error {
	scanner := bufio.NewScanner(file)
	// Scan the file to generate tokens
	for scanner.Scan() {
		for _, r := range scanner.Text() {
			dfa.Transition(r)
			// fmt.Print(i, r)
			// fmt.Printf("Index: %d, Rune: %c\n", i, r)
		}
		dfa.Store()
		dfa.StoreLine()
	}
	return nil
}
