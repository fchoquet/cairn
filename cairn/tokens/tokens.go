package tokens

import "fmt"

// TokenType represents a token type
type TokenType string

// Token types
const (
	ERROR   TokenType = "ERROR"
	INTEGER TokenType = "INTEGER"
	PLUS    TokenType = "PLUS"
	MINUS   TokenType = "MINUS"
	MULT    TokenType = "MULT"
	DIV     TokenType = "DIV"
	EOF     TokenType = "EOF"
	LPAREN  TokenType = "LPAREN"
	RPAREN  TokenType = "RPAREN"
)

// Token reprensents the result of a lexical analysis
type Token struct {
	Type  TokenType
	Value string
}

func (t *Token) String() string {
	return fmt.Sprintf("Token(%s, %s)", t.Type, t.Value)
}

// New builds a new Token
func New(tokenType TokenType, value string) *Token {
	return &Token{
		Type:  tokenType,
		Value: value,
	}
}

// Error returns an error Token
func Error(err error) *Token {
	return &Token{
		Type:  ERROR,
		Value: err.Error(),
	}
}
