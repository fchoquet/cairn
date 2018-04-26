package tokens

import (
	"fmt"
)

// TokenType is an enumeration of valid token types
type TokenType string

// Token types
const (
	NewLine     TokenType = "nl"
	Indentation           = "indent"
	Number                = "num"
	String                = "string"
	Operator              = "op"
)

// Token is the product of the tokenizer
type Token struct {
	Type     TokenType
	Value    string
	Position Position
}

func (t Token) String() string {
	return fmt.Sprintf(
		"[%s:%s (%d:%d@%s)]",
		t.Type,
		t.Value,
		t.Position.Line,
		t.Position.Col,
		t.Position.File,
	)
}

// Position reprensents the position of a token in the source code
type Position struct {
	File string
	Line int
	Col  int
}
