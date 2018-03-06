package interpreter

import (
	"fmt"
	"strconv"

	"github.com/fchoquet/cairn/cairn/ast"
	"github.com/fchoquet/cairn/cairn/parser"
	"github.com/fchoquet/cairn/cairn/tokens"
)

// Interpreter traverses the AST returned by the parser and yields results
type Interpreter struct {
	Parser *parser.Parser
}

func (i *Interpreter) Interpret(fileName, text string) (string, error) {
	ast, err := i.Parser.Parse(fileName, text)

	if err != nil {
		return "", err
	}

	return visit(ast)
}

func visit(node ast.Node) (string, error) {

	if num, ok := node.(*ast.Num); ok {
		return visitNum(num)
	}

	if unaryOp, ok := node.(*ast.UnaryOp); ok {
		return visitUnaryOp(unaryOp)
	}

	if binOp, ok := node.(*ast.BinOp); ok {
		return visitBinOp(binOp)
	}

	return "", fmt.Errorf("unexpected node type: %v", node)
}

func visitNum(node *ast.Num) (string, error) {
	return node.Value, nil
}

func visitUnaryOp(node *ast.UnaryOp) (string, error) {
	expr, err := visit(node.Expr)
	if err != nil {
		return "", err
	}

	val, err := strconv.Atoi(expr)
	if err != nil {
		return "", err
	}

	switch node.Op.Type {
	case tokens.PLUS:
		return strconv.Itoa(+val), nil
	case tokens.MINUS:
		return strconv.Itoa(-val), nil
	}

	return "", fmt.Errorf("unexpected binary operator: %s", node.Op)
}

func visitBinOp(node *ast.BinOp) (string, error) {
	left, err := visit(node.Left)
	if err != nil {
		return "", err
	}

	leftVal, err := strconv.Atoi(left)
	if err != nil {
		return "", err
	}

	right, err := visit(node.Right)
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
	}

	return "", fmt.Errorf("unexpected binary operator: %s", node.Op)
}
