package parser

import (
	"fmt"

	"github.com/fchoquet/cairn/ast"
	"github.com/fchoquet/cairn/tokens"
)

// bool: BOOL | IDENTIFIER
func (p *Parser) bool() (ast.Node, error) {
	token := p.current()
	switch token.Type {
	case tokens.BOOL:
		p.consume(tokens.BOOL)
		return &ast.Bool{Token: token, Value: token.Value}, nil
	case tokens.IDENTIFIER:
		p.consume(tokens.IDENTIFIER)
		return &ast.Variable{Token: token, Name: token.Value}, nil
	}
	return nil, fmt.Errorf("a boolean was expected: %s", token)
}

// boolterm: (NOT)boolterm | bool | LPAREN boolexpr RPAREN
func (p *Parser) boolterm() (ast.Node, error) {
	token := p.current()

	switch token.Type {
	case tokens.NOT:
		p.consume(token.Type)

		node, err := p.boolterm()
		if err != nil {
			return nil, err
		}

		return &ast.UnaryOp{Expr: node, Op: token}, nil
	case tokens.LPAREN:
		p.consume(tokens.LPAREN)

		node, err := p.boolexpr()
		if err != nil {
			return nil, err
		}

		if _, err := p.consume(tokens.RPAREN); err != nil {
			return nil, err
		}

		return node, nil
	}

	return p.bool()
}

// boolexpr: boolterm ((EQ | NEQ ) boolterm)*
func (p *Parser) boolexpr() (ast.Node, error) {
	node, err := p.boolterm()
	if err != nil {
		return nil, err
	}

	for p.current().Type == tokens.EQ || p.current().Type == tokens.NEQ {
		token, _ := p.consume(p.current().Type)

		right, err := p.boolterm()
		if err != nil {
			return nil, err
		}

		node = &ast.BinOp{Left: node, Op: token, Right: right}
	}

	return node, nil
}
