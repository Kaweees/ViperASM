package main

import "fmt"

// Represents the type of a token.
type TokenType int

// Represents the possible types of tokens.
const (
	INVALID TokenType = iota

	// Instructions
	INSTRUCTION
	// Operands and symbols
	REGISTER  // x0-x31 or aliases
	IMMEDIATE // Numeric values
	LABEL     // Symbol names
	LABEL_DEF // Symbol definitions (ending with :)
	// Punctuation
	COMMA
	LPAREN
	RPAREN
	// Directives and others
	DIRECTIVE // .text, .data, etc.
	MACRO     // .macro
	COMMENT   // # comments
	EOL       // End of line
)

// Represents a lexical token with position information
type Token struct {
	Type    TokenType
	Literal string
}

// String method to convert the current token type to a string.
func (t TokenType) String() string {
	return [...]string{"Invalid", "Instruction", "Register", "Immediate", "Label", "LabelDef", "Comma", "LParen", "RParen", "Directive", "Macro", "Comment", "EOL"}[t]
}

// String provides a readable representation of the token
func (t Token) String() string {
	return fmt.Sprintf("Token{Type: %s, Literal: %q}",
		t.Type, t.Literal)
}
