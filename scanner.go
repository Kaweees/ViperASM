package main

import (
	"os"
	"regexp"
	"strings"
	// "unicode"
)

// Regex patterns for the token types
var tokenRegex = map[TokenType]*regexp.Regexp{
	// Instructions
	INSTRUCTION: regexp.MustCompile(`^(add|sub|and|or|xor|sll|srl|sra|slt|sltu|addi|lw|sw|beq|bne|jal|jalr)`),
	// Registers (including aliases)
	REGISTER: regexp.MustCompile(`^(x[0-9]|x[1-2][0-9]|x3[0-1]|zero|ra|sp|gp|tp|t[0-6]|s[0-9]|s1[0-1]|a[0-7])`),
	// Immediates (decimal, hex, binary)
	IMMEDIATE: regexp.MustCompile(`^(-?[0-9]+|0x[0-9a-fA-F]+|0b[01]+)`),
	// Labels
	LABEL:     regexp.MustCompile(`^[a-zA-Z_.][a-zA-Z0-9_]*`),
	LABEL_DEF: regexp.MustCompile(`^[a-zA-Z_.][a-zA-Z0-9_]*:`),
	// Directives
	DIRECTIVE: regexp.MustCompile(`^\.(text|data|global|extern|byte|half|word|dword|string|align|section|macro|endm|ifdef|ifndef|endif|include|equ|set)`),
	// Macros
	MACRO: regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*`),
	// Other tokens
	COMMA:   regexp.MustCompile(`^,`),
	LPAREN:  regexp.MustCompile(`^\(`),
	RPAREN:  regexp.MustCompile(`^\)`),
	COMMENT: regexp.MustCompile(`^#.*`),
}

// Represents a Deterministic Finite Automaton(DFA) based lexical scanner.
type Scanner struct {
	File     *os.File
	Patterns map[TokenType]*regexp.Regexp
	Tokens   [][]Token
}

// Constructor to initialize memory for the Scanner.
func NewScanner(file *os.File) (*Scanner, error) {
	sc := &Scanner{}
	sc.File = file
	sc.Patterns = tokenRegex
	sc.Tokens = [][]Token{}
	return sc, nil
}

// Add a line of tokens to the list of tokens.
func (sc *Scanner) StoreLine(tokens []Token) {
	if len(tokens) > 0 {
		sc.Tokens = append(sc.Tokens, tokens)
	}
}

// Scan a line of assembly and return all tokens found.
func (sc *Scanner) ScanLine(line string) []Token {
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
		for tokType, pattern := range sc.Patterns {
			if match := pattern.FindString(line); match != "" {
				if tokType != COMMENT && tokType != EOL {
					tokens = append(tokens, Token{
						Type:    tokType,
						Literal: match,
					})
				}
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
			})
			line = line[1:]
		}
	}

	return tokens
}

// Scan the file for tokens.
func (sc *Scanner) ScanFile() ([][]Token, error) {
	// Example assembly line
	line := "add x1, x2, x3  # Add registers"
	scannedTokens := sc.ScanLine(line)
	sc.StoreLine(scannedTokens)
	return sc.Tokens, nil
}
