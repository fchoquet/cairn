package interpreter

// see https://ruslanspivak.com/lsbasi-part4/

import (
	"fmt"
	"strconv"
)

type Interpreter struct {
	Parser *Parser
}

func (i *Interpreter) Interpret() (string, error) {
	ast, err := i.Parser.Parse()
	if err != nil {
		return "", err
	}

	return visit(ast)
}

func visit(node Node) (string, error) {

	if num, ok := node.(*Num); ok {
		return visitNum(num)
	}

	if unaryOp, ok := node.(*UnaryOp); ok {
		return visitUnaryOp(unaryOp)
	}

	if binOp, ok := node.(*BinOp); ok {
		return visitBinOp(binOp)
	}

	return "", fmt.Errorf("unexpected node type: %v", node)
}

func visitNum(node *Num) (string, error) {
	return node.Value, nil
}

func visitUnaryOp(node *UnaryOp) (string, error) {
	expr, err := visit(node.Expr)
	if err != nil {
		return "", err
	}

	val, err := strconv.Atoi(expr)
	if err != nil {
		return "", err
	}

	switch node.Op.Type {
	case PLUS:
		return strconv.Itoa(+val), nil
	case MINUS:
		return strconv.Itoa(-val), nil
	}

	return "", fmt.Errorf("unexpected binary operator: %s", node.Op)
}

func visitBinOp(node *BinOp) (string, error) {
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
	case PLUS:
		return strconv.Itoa(leftVal + rightVal), nil
	case MINUS:
		return strconv.Itoa(leftVal - rightVal), nil
	case MULT:
		return strconv.Itoa(leftVal * rightVal), nil
	case DIV:
		return strconv.Itoa(leftVal / rightVal), nil
	}

	return "", fmt.Errorf("unexpected binary operator: %s", node.Op)
}
