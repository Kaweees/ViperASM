package main

import "fmt"

// Represents the type of a token.
type TokenType int

// Represents the possible types of tokens.
const (
	Initial TokenType = iota
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

// Represents a token in the scanner.
type Token struct {
	Type    TokenType
	Literal string
}

// String method to convert the current token type to a string.
func (t TokenType) String() string {
	return [...]string{"Initial", "Identifier", "DotIdentifier", "Register", "Zero", "Decimal", "Hexadecimal", "Comma", "LParen", "RParen", "LabelDef", "Comment", "String"}[t]
}

// String method to convert the current token to a string.
func (t Token) String() string {
	return fmt.Sprintf("Token{Type: %s, Literal: %s}", t.Type, t.Literal)
}
