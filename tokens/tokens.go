package tokens

import "fmt"

// TokenType represents a token type
type TokenType string

// Token types
const (
	ERROR      TokenType = "ERROR"
	INTEGER    TokenType = "INTEGER"
	PLUS       TokenType = "PLUS"
	MINUS      TokenType = "MINUS"
	MULT       TokenType = "MULT"
	DIV        TokenType = "DIV"
	EOF        TokenType = "EOF"
	LPAREN     TokenType = "LPAREN"
	RPAREN     TokenType = "RPAREN"
	STRING     TokenType = "STRING"
	IDENTIFIER TokenType = "IDENTIFIER"
	ASSIGN     TokenType = "ASSIGN"
)

// Token reprensents the result of a lexical analysis
type Token struct {
	Type     TokenType
	Value    string
	Position Position
}

func (t *Token) String() string {
	return fmt.Sprintf("Token(%s, %s, %s)", t.Type, t.Value, t.Position)
}

// Position represents the position of a token in the source code
type Position struct {
	File string
	Line int
	Col  int
}

func (p Position) String() string {
	return fmt.Sprintf("Pos(%s, %d, %d)", p.File, p.Line, p.Col)
}
