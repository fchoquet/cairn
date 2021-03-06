package interpreter

import (
	"fmt"
	"math"
	"strconv"

	"github.com/fchoquet/cairn/ast"
	"github.com/fchoquet/cairn/parser"
	"github.com/fchoquet/cairn/tokens"
)

// Interpreter traverses the AST returned by the parser and yields results
type Interpreter struct {
	Parser      *parser.Parser
	SymbolTable SymbolTable
}

// New creates a new interpreter
func New(parser *parser.Parser) *Interpreter {
	return &Interpreter{
		Parser:      parser,
		SymbolTable: SymbolTable{},
	}
}

type Symbol struct {
	Scope      string
	Identifier string
}

type SymbolTable map[Symbol]string

func (i *Interpreter) Interpret(fileName, text string) (string, error) {
	ast, err := i.Parser.Parse(fileName, text)

	// DEBUG CODE
	// display all the tokens
	fmt.Println("--- AST: ---")
	fmt.Println(ast)
	fmt.Println("------------")
	// END DEBUG CODE

	if err != nil {
		return "", fmt.Errorf("Parser error: %s", err)
	}

	return i.visit(ast)
}

func (i *Interpreter) visit(node ast.Node) (string, error) {
	switch n := node.(type) {
	case *ast.SourceFile:
		return i.visitSourceFile(n)
	case *ast.FuncDecl:
		return i.visitFuncDecl(n)
	case *ast.StatementList:
		return i.visitStatementList(n)
	case *ast.BlockStmt:
		return i.visitBlockStmt(n)
	case *ast.Num:
		return i.visitNum(n)
	case *ast.String:
		return i.visitString(n)
	case *ast.Bool:
		return i.visitBool(n)
	case *ast.UnaryOp:
		return i.visitUnaryOp(n)
	case *ast.BinOp:
		return i.visitBinOp(n)
	case *ast.Assignment:
		return i.visitAssignment(n)
	case *ast.Variable:
		return i.visitVariable(n)
	default:
		return "", fmt.Errorf("unexpected node type: %v", node)
	}
}

func (i *Interpreter) visitSourceFile(node *ast.SourceFile) (string, error) {
	for _, f := range node.Functions {
		if _, err := i.visit(f); err != nil {
			return "", err
		}
	}
	return i.visitStatementList(node.Statements)
}

func (i *Interpreter) visitFuncDecl(node *ast.FuncDecl) (string, error) {
	// TODO implement
	return "", nil
}

func (i *Interpreter) visitStatementList(node *ast.StatementList) (string, error) {
	output := ""
	for _, st := range node.Statements {
		s, err := i.visit(st)
		if err != nil {
			return output, err
		}
		// the output is the output of the last statement
		output = s
	}

	return output, nil
}

func (i *Interpreter) visitBlockStmt(node *ast.BlockStmt) (string, error) {
	return i.visitStatementList(node.Statements)
}

func (i *Interpreter) visitNum(node *ast.Num) (string, error) {
	return node.Value, nil
}

func (i *Interpreter) visitString(node *ast.String) (string, error) {
	return node.Value, nil
}

func (i *Interpreter) visitBool(node *ast.Bool) (string, error) {
	return node.Value, nil
}

func (i *Interpreter) visitUnaryOp(node *ast.UnaryOp) (string, error) {
	expr, err := i.visit(node.Expr)
	if err != nil {
		return "", err
	}

	switch node.Op.Type {
	case tokens.NOT:
		val, err := strconv.ParseBool(expr)
		if err != nil {
			return "", err
		}
		return strconv.FormatBool(!val), err
	default:
		val, err := strconv.Atoi(expr)
		if err != nil {
			return "", err
		}
		switch node.Op.Type {
		case tokens.PLUS:
			return strconv.Itoa(+val), nil
		case tokens.MINUS:
			return strconv.Itoa(-val), nil
		default:
			return "", fmt.Errorf("unexpected binary operator: %s", node.Op)
		}
	}
}

func (i *Interpreter) visitBinOp(node *ast.BinOp) (string, error) {
	left, err := i.visit(node.Left)
	if err != nil {
		return "", err
	}

	right, err := i.visit(node.Right)
	if err != nil {
		return "", err
	}

	switch node.Op.Type {
	case tokens.CONCAT:
		return (left + right), nil
	case tokens.EQ:
		return strconv.FormatBool(left == right), nil
	case tokens.NEQ:
		return strconv.FormatBool(left != right), nil
	case tokens.POW:
		// we need floats to use math.Pow but we expect ints only for now
		leftVal, err := strconv.ParseFloat(left, 64)
		if err != nil {
			return "", err
		}
		rightVal, err := strconv.ParseFloat(right, 64)
		if err != nil {
			return "", err
		}
		// watch the int conversion here
		return strconv.Itoa(int(math.Pow(leftVal, rightVal))), nil
	case tokens.AND, tokens.OR:
		leftVal, err := strconv.ParseBool(left)
		if err != nil {
			return "", err
		}
		rightVal, err := strconv.ParseBool(right)
		if err != nil {
			return "", err
		}
		switch node.Op.Type {
		case tokens.AND:
			return strconv.FormatBool(leftVal && rightVal), nil
		case tokens.OR:
			return strconv.FormatBool(leftVal || rightVal), nil
		default:
			return "", fmt.Errorf("unexpected binary operator: %s", node.Op)
		}
	default:
		// string to int conversions
		leftVal, err := strconv.Atoi(left)
		if err != nil {
			return "", err
		}

		rightVal, err := strconv.Atoi(right)
		if err != nil {
			return "", err
		}

		switch node.Op.Type {
		case tokens.PLUS:
			return strconv.Itoa(leftVal + rightVal), nil
		case tokens.MINUS:
			return strconv.Itoa(leftVal - rightVal), nil
		case tokens.MULT:
			return strconv.Itoa(leftVal * rightVal), nil
		case tokens.DIV:
			return strconv.Itoa(leftVal / rightVal), nil
		default:
			return "", fmt.Errorf("unexpected binary operator: %s", node.Op)
		}
	}
}

func (i *Interpreter) visitAssignment(node *ast.Assignment) (string, error) {
	right, err := i.visit(node.Right)
	if err != nil {
		return "", err
	}

	i.SymbolTable[Symbol{Scope: "global", Identifier: node.Variable.Name}] = right

	// DEBUG code
	fmt.Printf("%+v\n", i.SymbolTable)
	// END DEBUG code

	return right, nil
}

func (i *Interpreter) visitVariable(node *ast.Variable) (string, error) {
	value, ok := i.SymbolTable[Symbol{Scope: "global", Identifier: node.Name}]
	if !ok {
		return "", fmt.Errorf("unknown identifier: %s", node.Name)
	}
	return value, nil
}
