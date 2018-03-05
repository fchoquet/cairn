package interpreter

import "fmt"

type TokenType string

const (
	INTEGER TokenType = "INTEGER"
	PLUS    TokenType = "PLUS"
	MINUS   TokenType = "MINUS"
	MULT    TokenType = "MULT"
	DIV     TokenType = "DIV"
	EOF     TokenType = "EOF"
	LPAREN  TokenType = "LPAREN"
	RPAREN  TokenType = "RPAREN"
)

type Token struct {
	Type  TokenType
	Value string
}

func (t *Token) String() string {
	return fmt.Sprintf("Token(%s, %s)", t.Type, t.Value)
}
