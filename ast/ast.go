package ast

import (
	"fmt"
	"strings"

	"github.com/fchoquet/cairn/tokens"
)

type Node interface {
	fmt.Stringer
}

type SourceFile struct {
	Functions  []*FuncDecl
	Statements *StatementList
}

func (s SourceFile) String() string {
	functions := []string{}
	for _, f := range s.Functions {
		functions = append(functions, f.String())
	}
	return fmt.Sprintf("SourceFile(%s %s)", strings.Join(functions, "; "), s.Statements)
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

type TypeId struct {
	Token *tokens.Token
	Name  string
}

func (t *TypeId) String() string {
	return fmt.Sprintf("Type(%s)", t.Name)
}

type Parameter struct {
	Token *tokens.Token
	Name  string
	Type  *TypeId
}

func (p *Parameter) String() string {
	return fmt.Sprintf("Parameter(%s %s)", p.Name, p.Type)
}

type ParameterList struct {
	Token      *tokens.Token
	Parameters []*Parameter
}

func (pl ParameterList) String() string {
	params := []string{}
	for _, p := range pl.Parameters {
		params = append(params, p.String())
	}
	return fmt.Sprintf("ParameterList(%s)", strings.Join(params, " "))
}

type FuncDecl struct {
	Token     *tokens.Token
	Name      *tokens.Token
	Signature *Signature
	Body      *BlockStmt
}

func (f *FuncDecl) String() string {
	return fmt.Sprintf("FuncDecl(%s %s %s)", f.Name, f.Signature, f.Body)
}

type Signature struct {
	Token      *tokens.Token
	Parameters *ParameterList
	ReturnType *TypeId
}

func (s *Signature) String() string {
	return fmt.Sprintf("Signature(%s %s)", s.Parameters, s.ReturnType)
}
