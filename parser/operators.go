package parser

import (
	"github.com/fchoquet/cairn/tokens"
)

type associativity string

const (
	AssocLeft  associativity = "left"
	AssocRight associativity = "right"
)

var BinaryOpPrecedence = map[tokens.TokenType]int{
	// operating on numbers
	tokens.PLUS:  4,
	tokens.MINUS: 4,
	tokens.MULT:  5,
	tokens.DIV:   5,
	tokens.POW:   6,

	// operating on strings
	tokens.CONCAT: 4,

	// operating on booleans
	tokens.OR:  1,
	tokens.AND: 2,

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
	tokens.POW:   AssocRight,

	// operating on strings
	tokens.CONCAT: AssocLeft,

	// operating on booleans
	tokens.OR:  AssocLeft,
	tokens.AND: AssocLeft,

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
		tokens.POW,
		tokens.EQ,
		tokens.NEQ,
		tokens.CONCAT,
		tokens.AND,
		tokens.OR:
		return true
	default:
		return false
	}
}
