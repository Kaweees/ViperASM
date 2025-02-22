package main

// String method to convert the current state to a string.
func (s State) String() string {
	return [...]string{"Initial", "Identifier", "DotIdentifier", "Register", "Zero", "Decimal", "Hexadecimal", "Comma", "LParen", "RParen", "LabelDef", "Comment", "String"}[s]
}

// Represents the type of a token.
type Type string

// Represents a token in the scanner.
type Token struct {
	Type  Type
	Value string
}
