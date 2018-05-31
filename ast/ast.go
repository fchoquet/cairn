package ast

import (
	"fmt"
	"strings"

	"github.com/fchoquet/cairn/tokens"
)

type Node interface {
	fmt.Stringer
}

type Statement interface {
	Node
}

type StatementList struct {
	Statements []Statement
}

func (sl StatementList) String() string {
	statements := []string{}
	for _, st := range sl.Statements {
		statements = append(statements, st.String())
	}
	return fmt.Sprintf("StatementList(%s)", strings.Join(statements, "; "))
}

type BlockStmt struct {
	Begin      *tokens.Token
	Statements *StatementList
	End        *tokens.Token
}

func (bs *BlockStmt) String() string {
	return fmt.Sprintf("BlockStmt(%s %s %s)", bs.Begin, bs.Statements, bs.End)
}

type UnaryOp struct {
	Op   *tokens.Token
	Expr Node
}

func (op *UnaryOp) String() string {
	return fmt.Sprintf("UnaryOp(%s %s)", op.Op, op.Expr)
}

type BinOp struct {
	Left  Node
	Op    *tokens.Token
	Right Node
}

func (op *BinOp) String() string {
	return fmt.Sprintf("BinOp(%s %s %s)", op.Op, op.Left, op.Right)
}

type Num struct {
	Token *tokens.Token
	Value string
}

func (num *Num) String() string {
	return fmt.Sprintf("Num(%s)", num.Token)
}

type String struct {
	Token *tokens.Token
	Value string
}

func (str *String) String() string {
	return fmt.Sprintf("String(%s)", str.Token)
}

type Bool struct {
	Token *tokens.Token
	Value string
}

func (b *Bool) String() string {
	return fmt.Sprintf("Bool(%s)", b.Token)
}

// Assignment represent and assignment in an AST
type Assignment struct {
	Token    *tokens.Token
	Variable Variable
	Right    Node
}

func (asgn *Assignment) String() string {
	return fmt.Sprintf("Assign(%s %s)", asgn.Variable, asgn.Right)
}

// Variable represents a variable in an AST
type Variable struct {
	Token *tokens.Token
	Name  string
}

func (v *Variable) String() string {
	return fmt.Sprintf("Variable(%s)", v.Name)
}
