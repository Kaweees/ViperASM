package main

import (
	"bufio"
	"os"
	"regexp"
	"strings"
	// "unicode"
)

// Regex patterns for the token types
var tokenRegex = map[TokenType]*regexp.Regexp{
	INSTRUCTION: regexp.MustCompile(`^(add|sub|and|or|xor|sll|srl|sra|slt|sltu|addi|lw|sw|beq|bne|jal|jalr)`),
	REGISTER:    regexp.MustCompile(`^(x[0-9]|x[1-2][0-9]|x3[0-1]|zero|ra|sp|gp|tp|t[0-6]|s[0-9]|s1[0-1]|a[0-7])`),
	IMMEDIATE:   regexp.MustCompile(`^-?\d+`),
	LABEL:       regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*:`),
	COMMENT:     regexp.MustCompile(`^#.*`),
	COMMA:       regexp.MustCompile(`^,`),
	LPAREN:      regexp.MustCompile(`^\(`),
	RPAREN:      regexp.MustCompile(`^\)`),
}

// Represents a Deterministic Finite Automaton(DFA) based lexical scanner.
type Scanner struct {
	File     *os.File
	Patterns map[TokenType]*regexp.Regexp
	Tokens   [][]Token
}

// Constructor to initialize memory for the Scanner.
func NewScanner(file *os.File) (*Scanner, error) {
	scanner := &Scanner{}
	scanner.File = file
	scanner.Patterns = tokenRegex
	scanner.Tokens = [][]Token{}
	return scanner, nil
}

// Add a line of tokens to the list of tokens.
func (scanner *Scanner) StoreLine(tokens []Token) {
	if len(tokens) > 0 {
		scanner.Tokens = append(scanner.Tokens, tokens)
	}
}

// Scan the file for tokens.
func (scanner *Scanner) ScanFile() ([][]Token, error) {
	buf := bufio.NewScanner(scanner.File)
	lineNum := 0
	for buf.Scan() {
		line := buf.Text()
		tokens, err := scanner.ScanLine(line, lineNum)
		if err != nil {
			return nil, err
		}
		scanner.StoreLine(tokens)
		lineNum++
	}
	return scanner.Tokens, nil
}

// ScanLine takes a line of assembly and returns all tokens found
func (s *Scanner) ScanLine(line string, lineNum int) ([]Token, error) {
	tokens := []Token{}
	line = strings.TrimSpace(line)

	for len(line) > 0 {
		// Skip whitespace
		line = strings.TrimLeft(line, " \t")
		if len(line) == 0 {
			break
		}

		// Try to match each pattern
		matched := false
		for tokType, pattern := range s.Patterns {
			if tokType == EOL {
				break
			}
			if match := pattern.FindString(line); match != "" {
				tokens = append(tokens, Token{
					Type:    tokType,
					Literal: match,
					LineNum: lineNum,
				})
				line = line[len(match):]
				matched = true
				break
			}
		}

		if !matched {
			// Handle invalid token
			tokens = append(tokens, Token{
				Type:    INVALID,
				Literal: string(line[0]),
				LineNum: lineNum,
			})
			line = line[1:]
		}
	}

	// Add EOL token
	tokens = append(tokens, Token{
		Type:    EOL,
		Literal: "",
		LineNum: lineNum,
	})

	return tokens, nil
}

// // Store the current state of the DFA.
// func (dfa *DFA) Store() {
// 	if dfa.State != Initial && dfa.State != Comment {
// 		dfa.AddToken(dfa.State.String(), dfa.Token)
// 	}
// 	dfa.Reset()
// }

// // Reset the DFA to its initial state.
// func (dfa *DFA) Reset() {
// 	dfa.State = Initial
// 	dfa.Token = ""
// 	dfa.currentString = 0
// }
