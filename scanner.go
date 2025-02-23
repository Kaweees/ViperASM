package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Add this ordered processing list above the tokenRegex declaration
var tokenPriority = []TokenType{
	INSTRUCTION,
	PSEUDO_INSTRUCTION,
	LABEL_DEF,
	DIRECTIVE,
	REGISTER,
	IMMEDIATE,
	COMMA,
	LPAREN,
	RPAREN,
	COMMENT,
	STRING,
	CHAR,
	LABEL, // Regular labels come after label definitions
}

// Generate a regex for all instructions
func instructionRegex() string {
	instructions := []string{}
	for name, _ := range instructionSet {
		instructions = append(instructions, name)
	}
	return strings.Join(instructions, "|")
}

// Generate a regex for all pseudo instructions
func pseudoInstructionRegex() string {
	pseudoInstructions := []string{}
	for name, _ := range pseudoInstructionSet {
		pseudoInstructions = append(pseudoInstructions, name)
	}
	return strings.Join(pseudoInstructions, "|")
}

// Updated token patterns with proper priority
var tokenRegex = map[TokenType]*regexp.Regexp{
	LABEL_DEF: regexp.MustCompile(`^[a-zA-Z_.][a-zA-Z0-9_\.]*:`), // Modified to include dots
	DIRECTIVE: regexp.MustCompile(`^\.(text|data|globl|global|extern|byte|half|word|dword|string|align|section|macro|endm|ifdef|ifndef|endif|include|equ|set)`),
	// Instructions
	INSTRUCTION: regexp.MustCompile(
		fmt.Sprintf("^(%s)", instructionRegex()),
	),
	PSEUDO_INSTRUCTION: regexp.MustCompile(
		fmt.Sprintf("^(%s)", pseudoInstructionRegex()),
	),
	// Registers (including aliases)
	REGISTER: regexp.MustCompile(`^(x[0-9]|x[1-2][0-9]|x3[0-1]|zero|ra|sp|gp|tp|t[0-6]|s[0-9]|s1[0-1]|a[0-7])`),
	// Immediates (decimal, hex, binary)
	IMMEDIATE: regexp.MustCompile(`^(-?[0-9]+|0x[0-9a-fA-F]+|0b[01]+)`),
	// Labels
	LABEL: regexp.MustCompile(`^[a-zA-Z_.][a-zA-Z0-9_]*`),
	// Other tokens
	COMMA:   regexp.MustCompile(`^,`),
	LPAREN:  regexp.MustCompile(`^\(`),
	RPAREN:  regexp.MustCompile(`^\)`),
	COMMENT: regexp.MustCompile(`^#.*`),
	STRING:  regexp.MustCompile(`^"([^"\\]|\\.)*"`), // Matches quoted strings with escape support
	CHAR:    regexp.MustCompile(`^'([^'\\]|\\.)'`),
}

// Represents a lexical scanner.
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

		// Special handling for string literals
		if line[0] == '"' {
			if match := sc.Patterns[STRING].FindString(line); match != "" {
				// Remove the quotes and handle escapes
				literal := match[1 : len(match)-1] // Remove surrounding quotes
				literal = strings.ReplaceAll(literal, `\"`, `"`)
				literal = strings.ReplaceAll(literal, `\\`, `\`)
				literal = strings.ReplaceAll(literal, `\n`, "\n")
				literal = strings.ReplaceAll(literal, `\t`, "\t")

				tokens = append(tokens, Token{
					Type:    STRING,
					Literal: literal,
				})
				line = line[len(match):]
				continue
			}
		}

		// Handle character literals
		if line[0] == '\'' {
			if match := sc.Patterns[CHAR].FindString(line); match != "" {
				// Remove the quotes and handle escapes
				literal := match[1 : len(match)-1] // Remove surrounding quotes
				literal = strings.ReplaceAll(literal, `\'`, `'`)
				literal = strings.ReplaceAll(literal, `\\`, `\`)
				literal = strings.ReplaceAll(literal, `\n`, "\n")
				literal = strings.ReplaceAll(literal, `\t`, "\t")

				tokens = append(tokens, Token{
					Type:    CHAR,
					Literal: literal,
				})
				line = line[len(match):]
				continue
			}
		}

		// Try to match each pattern in priority order
		matched := false
		for _, tokType := range tokenPriority {
			pattern := sc.Patterns[tokType]
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

// ScanFile reads the input file line by line and returns all tokens
func (sc *Scanner) ScanFile(file *os.File) ([][]Token, error) {
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}

		// Scan the line and store tokens
		tokens := sc.ScanLine(line)

		// Only store non-empty token lists
		if len(tokens) > 0 {
			// Add line number to each token
			sc.StoreLine(tokens)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return sc.Tokens, nil
}
