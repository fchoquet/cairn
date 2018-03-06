package ast

import "github.com/fchoquet/cairn/cairn/tokens"

type Node interface {
}

type UnaryOp struct {
	Op   *tokens.Token
	Expr Node
}

type BinOp struct {
	Left  Node
	Op    *tokens.Token
	Right Node
}

type Num struct {
	Token *tokens.Token
	Value string
}
