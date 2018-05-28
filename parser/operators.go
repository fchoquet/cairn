package parser

import "github.com/fchoquet/cairn/tokens"

type associativity string

const (
	AssocLeft  associativity = "left"
	AssocRight associativity = "right"
)

var BinaryOpPrecedence = map[tokens.TokenType]int{
	// operating on numbers
	tokens.PLUS:  1,
	tokens.MINUS: 1,
	tokens.MULT:  2,
	tokens.DIV:   2,

	// operating on strings
	tokens.CONCAT: 1,

	// operating on booleans

	// general purpose
	tokens.EQ:  3,
	tokens.NEQ: 3,
}

var BinaryOpAssociativity = map[tokens.TokenType]associativity{
	// operating on numbers
	tokens.PLUS:  AssocLeft,
	tokens.MINUS: AssocLeft,
	tokens.MULT:  AssocLeft,
	tokens.DIV:   AssocLeft,

	// operating on strings
	tokens.CONCAT: AssocLeft,

	// operating on booleans

	// general purpose
	tokens.EQ:  AssocLeft,
	tokens.NEQ: AssocLeft,
}

func isUnaryOp(tk *tokens.Token) bool {
	switch tk.Type {
	case tokens.PLUS,
		tokens.MINUS,
		tokens.NOT:
		return true
	default:
		return false
	}
}

func isBinaryOp(tk *tokens.Token) bool {
	switch tk.Type {
	case tokens.PLUS,
		tokens.MINUS,
		tokens.MULT,
		tokens.DIV,
		tokens.EQ,
		tokens.NEQ,
		tokens.CONCAT:
		return true
	default:
		return false
	}
}
