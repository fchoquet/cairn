package ast

import (
	"fmt"

	"github.com/fchoquet/cairn/tokens"
)

type Node interface {
	// GetToken returns the token associated with this node
	GetToken() *tokens.Token
}

type UnaryOp struct {
	Op   *tokens.Token
	Expr Node
}

// GetToken implements the Node interface
func (op *UnaryOp) GetToken() *tokens.Token {
	return op.Op
}

func (op *UnaryOp) String() string {
	return fmt.Sprintf("UnaryOp(%s, %s)\n", op.Op, op.Expr)
}

type BinOp struct {
	Left  Node
	Op    *tokens.Token
	Right Node
}

// GetToken implements the Node interface
func (op *BinOp) GetToken() *tokens.Token {
	return op.Op
}

func (op *BinOp) String() string {
	return fmt.Sprintf("BinOp(\n%s,\n%s,\n%s\n)\n", op.Left, op.Op, op.Right)
}

type Num struct {
	Token *tokens.Token
	Value string
}

// GetToken implements the Node interface
func (num *Num) GetToken() *tokens.Token {
	return num.Token
}

func (num *Num) String() string {
	return fmt.Sprintf("Num(%s)", num.Token)
}

type String struct {
	Token *tokens.Token
	Value string
}

// GetToken implements the Node interface
func (str *String) GetToken() *tokens.Token {
	return str.Token
}

func (str *String) String() string {
	return fmt.Sprintf("String(%s)", str.Token)
}

// Assignment represent and assignment in an AST
type Assignment struct {
	Token    *tokens.Token
	Variable Variable
	Right    Node
}

// GetToken implements the Node interface
func (asgn *Assignment) GetToken() *tokens.Token {
	return asgn.Token
}

func (asgn *Assignment) String() string {
	return fmt.Sprintf("Assign(%s := %s)", asgn.Variable, asgn.Right)
}

// Variable represents a variable in an AST
type Variable struct {
	Token *tokens.Token
	Name  string
}

// GetToken implements the Node interface
func (v *Variable) GetToken() *tokens.Token {
	return v.Token
}

func (v *Variable) String() string {
	return fmt.Sprintf("Identifier(%s)", v.Name)
}