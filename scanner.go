package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

// Represents a Deterministic Finite Automaton(DFA) based lexical scanner.
type Scanner struct {
	File   *os.File
	Token  Token
	Tokens []Token
	Reader *bufio.Reader
	State  TokenType
}

// Constructor to initialize memory for the Scanner.
func NewScanner(file *os.File) (*Scanner, error) {
	scanner := &Scanner{}
	scanner.File = file
	scanner.Token = Token{}
	scanner.Tokens = []Token{}
	scanner.Reader = nil
	scanner.State = Initial
	return scanner, nil
}

// Add a token to the list of tokens.
func (scanner *Scanner) AddToken() {
	scanner.Tokens = append(scanner.Tokens, scanner.Token)
}

func (scanner *Scanner) ReadChar() (rune, error) {
	b, err := scanner.Reader.ReadByte()
	if err != nil {
		return 0, err
	}
	return rune(b), nil
}

// Scan the file for tokens.
func (scanner *Scanner) ScanFile() ([]Token, error) {
	scanner.Reader = bufio.NewReader(scanner.File)
	// Scan the file character by character to generate tokens
	for {
		b, err := scanner.Reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				return scanner.Tokens, nil
			}
			return nil, err
		}
		scanner.Transition(r)
		ch := rune(b)
		if ch == 'h' {
			fmt.Print("h")
		}
	}
}

// Transition the scanner to a new state based on the input.
func (scanner *Scanner) Transition(input byte) {
	// fmt.Printf("State: %s, Rune: '%c'\n", dfa.State.String(), input)
	switch scanner.State {
	case Initial:
		if input == 'x' {
			scanner.State = DotIdentifier
		} else if input == '$' {
			scanner.Token = string(input)
			scanner.State = Register
		} else if input == '0' {
			scanner.State = Zero
		} else if unicode.IsDigit(input) || input == '-' {
			scanner.Token = string(input)
			scanner.State = Decimal
		} else if input == ',' {
			scanner.AddToken(scanner.State.String(), ",")
			scanner.Reset()
		} else if input == '(' {
			scanner.AddToken(scanner.State.String(), "(")
			scanner.Reset()
		} else if input == ')' {
			scanner.AddToken(scanner.State.String(), ")")
			scanner.Reset()
		} else if input == '#' || input == ';' {
			scanner.Store()
			scanner.State = Comment
		} else if input == '"' || input == '\'' {
			scanner.currentString = input
			scanner.Token = string(input)
			scanner.State = String
		} else if !unicode.IsSpace(input) {
			scanner.Token = string(input)
			scanner.State = Identifier
		}
	case Identifier:
		if input == ':' {
			scanner.State = LabelDef
			scanner.Store()
		} else if !unicode.IsSpace(input) {
			scanner.Token += string(input)
		} else {
			scanner.Store()
		}
	case DotIdentifier:
		if unicode.IsLetter(input) {
			scanner.Token += string(input)
		} else {
			scanner.Store()
		}
	case Register:
		if input == ',' {
			scanner.Store()
			scanner.AddToken(scanner.State.String(), ",")
		} else if unicode.IsDigit(input) || unicode.IsLetter(input) {
			scanner.Token += string(input)
		} else {
			scanner.Store()
		}
	case Zero:
		if input == 'x' {
			scanner.State = Hexadecimal
		} else if unicode.IsDigit(input) {
			scanner.Token = string(input)
			scanner.State = Decimal
		} else {
			scanner.Store()
		}
	case Decimal:
		if unicode.IsDigit(input) {
			scanner.Token += string(input)
		} else {
			scanner.Store()
		}
	case Hexadecimal:
		if strings.ContainsAny(string(input), "0123456789abcdefABCDEF") {
			scanner.Token += string(input)
		} else {
			scanner.Store()
		}
	case String:
		scanner.Token += string(input)
		if input == scanner.currentString {
			scanner.Store()
		}
	}
}

// // Store the current state of the DFA.
// func (dfa *DFA) Store() {
// 	if dfa.State != Initial && dfa.State != Comment {
// 		dfa.AddToken(dfa.State.String(), dfa.Token)
// 	}
// 	dfa.Reset()
// }

// func (dfa *DFA) StoreLine() {
// 	if len(dfa.tokens) > 0 {
// 		dfa.totalTokens = append(dfa.totalTokens, dfa.tokens)
// 		dfa.tokens = []Token{}
// 	}
// }

// // Reset the DFA to its initial state.
// func (dfa *DFA) Reset() {
// 	dfa.State = Initial
// 	dfa.Token = ""
// 	dfa.currentString = 0
// }
