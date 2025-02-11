package main

import (
	"bufio"
	"os"
	"strings"
	"unicode"
)

// Represents the possible states of the DFA.
type State int

// Represents the possible states of the DFA.
const (
	Initial State = iota
	Identifier
	DotIdentifier
	Register
	Zero
	Decimal
	Hexadecimal
	Comma
	LParen
	RParen
	LabelDef
	Comment
	String
)

// String method to convert the current state to a string.
func (s State) String() string {
	return [...]string{"Initial", "Identifier", "DotIdentifier", "Register", "Zero", "Decimal", "Hexadecimal", "Comma", "LParen", "RParen", "LabelDef", "Comment", "String"}[s]
}

// Represents a token in the scanner.
type Token struct {
	tokenType  string
	tokenValue string
}

// Represents the scanning Deterministic Finite Automaton(DFA) for the scanner.
type DFA struct {
	currentState  State
	currentToken  string
	currentString rune
	tokens        []Token
	totalTokens   [][]Token
}

// Constructor to initialize memory for the DFA.
func NewDFA() (*DFA, error) {
	dfa := &DFA{}
	dfa.currentState = Initial
	dfa.currentToken = ""
	dfa.currentString = 0
	dfa.tokens = []Token{}
	dfa.totalTokens = [][]Token{}
	return dfa, nil
}

// Add a token to the list of tokens.
func (dfa *DFA) AddToken(tokenType string, tokenValue string) {
	dfa.tokens = append(dfa.tokens, Token{tokenType, tokenValue})
}

// Store the current state of the DFA.
func (dfa *DFA) Store() {
	if dfa.currentState != Initial && dfa.currentState != Comment {
		dfa.AddToken(dfa.currentState.String(), dfa.currentToken)
	}
	dfa.Reset()
}

func (dfa *DFA) StoreLine() {
	if len(dfa.tokens) > 0 {
		dfa.totalTokens = append(dfa.totalTokens, dfa.tokens)
		dfa.tokens = []Token{}
	}
}

// Reset the DFA to its initial state.
func (dfa *DFA) Reset() {
	dfa.currentState = Initial
	dfa.currentToken = ""
	dfa.currentString = 0
}

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
