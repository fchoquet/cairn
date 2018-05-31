package tokens

import "fmt"

// TokenType represents a token type
type TokenType string

// Token types
const (
	ERROR      TokenType = "ERROR"
	EOF        TokenType = "EOF"
	EOL        TokenType = "EOL"
	BEGIN      TokenType = "BEGIN"
	END        TokenType = "END"
	LPAREN     TokenType = "LPAREN"
	RPAREN     TokenType = "RPAREN"
	ASSIGN     TokenType = "ASSIGN"
	IDENTIFIER TokenType = "IDENTIFIER"
	FUNC       TokenType = "FUNC"
	COLUMN     TokenType = "COLUMN"
	COMMA      TokenType = "COMMA"

	// primaty type litterals
	INTEGER TokenType = "INTEGER"
	STRING  TokenType = "STRING"
	BOOL    TokenType = "BOOL"

	// operators
	PLUS   TokenType = "PLUS"
	MINUS  TokenType = "MINUS"
	MULT   TokenType = "MULT"
	DIV    TokenType = "DIV"
	POW    TokenType = "POW"
	CONCAT TokenType = "CONCAT"
	NOT    TokenType = "NOT"
	AND    TokenType = "AND"
	OR     TokenType = "OR"
	EQ     TokenType = "EQ"
	NEQ    TokenType = "NEQ"
)

// Token reprensents the result of a lexical analysis
type Token struct {
	Type     TokenType
	Value    string
	Position Position
}

func (t *Token) String() string {
	return fmt.Sprintf("%s:%s", t.Value, t.Type)
}

func (t *Token) Debug() string {
	return fmt.Sprintf("%s:%s@%s", t.Value, t.Type, t.Position)
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
