package main

import "fmt"

// Represents the type of a token.
type TokenType int

// Represents the possible types of tokens.
const (
	INVALID TokenType = iota
	INSTRUCTION
	REGISTER
	IMMEDIATE
	LABEL
	COMMENT
	COMMA
	LPAREN
	RPAREN
	EOL
)

// Represents a token in the scanner.
type Token struct {
	Type    TokenType
	Literal string
	LineNum int
}

// String method to convert the current token type to a string.
func (t TokenType) String() string {
	return [...]string{"Instruction", "Register", "Immediate", "Label", "Comment", "Comma", "LParen", "RParen", "EOL"}[t]
}

// String method to convert the current token to a string.
func (t Token) String() string {
	return fmt.Sprintf("Token{Type: %s, Literal: %s}", t.Type, t.Literal)
}
