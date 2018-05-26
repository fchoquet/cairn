package parser

import (
	"fmt"

	"github.com/fchoquet/cairn/ast"
	"github.com/fchoquet/cairn/tokenizer"
	"github.com/fchoquet/cairn/tokens"
)

// Parser reads a text and converts it to an AST using the Tokenizer
type Parser struct {
	buffer TokenBuffer
}

// Parse builds an AST from a text
func (p *Parser) Parse(fileName, text string) (ast.Node, error) {
	p.buffer = NewTokenBuffer(tokenizer.Tokenize(fileName, text), 2)

	return p.expr()
}

func (p *Parser) current() *tokens.Token {
	// LookAhead(0) never returns an error since it does not have to preload a new token
	current, _ := p.buffer.LookAhead(0)
	return current
}

func (p *Parser) lookAhead(n int) (*tokens.Token, error) {
	return p.buffer.LookAhead(n)
}

func (p *Parser) consume(tkType tokens.TokenType) (*tokens.Token, error) {
	tk, newBuffer, err := p.buffer.Consume()
	if tk.Type != tkType {
		return nil, fmt.Errorf("wrong input type. Expected %s - got %s", tkType, p.current())
	}

	// let's use mutation for now
	p.buffer = newBuffer
	return tk, err
}

// expr : assignment | variable | arithmexpr | strexpr
func (p *Parser) expr() (ast.Node, error) {
	token := p.current()
	switch token.Type {
	case tokens.IDENTIFIER:
		next, _ := p.lookAhead(1)
		switch next.Type {
		case tokens.ASSIGN:
			return p.assignment()
		case tokens.CONCAT:
			return p.strexpr()
		case tokens.EQ, tokens.NEQ:
			return p.boolexpr()
		default:
			return p.arithmexpr()
		}
	case tokens.STRING:
		return p.strexpr()
	case tokens.BOOL, tokens.NOT:
		return p.boolexpr()
	default:
		return p.arithmexpr()
	}
}

// assignment : IDENTIFIER ASSIGN expr
func (p *Parser) assignment() (ast.Node, error) {
	idToken, err := p.consume(tokens.IDENTIFIER)
	if err != nil {
		return nil, err
	}

	assignToken, err := p.consume(tokens.ASSIGN)
	if err != nil {
		return nil, err
	}

	right, err := p.expr()
	if err != nil {
		return nil, err
	}

	return &ast.Assignment{
		Variable: ast.Variable{
			Token: idToken,
			Name:  idToken.Value,
		},
		Token: assignToken,
		Right: right,
	}, nil
}

// factor : (PLUS|MINUS)factor | INTEGER | IDENTIFIER | LPAREN arithmexpr RPAREN
func (p *Parser) factor() (ast.Node, error) {
	token := p.current()

	switch token.Type {
	case tokens.PLUS, tokens.MINUS:
		p.consume(token.Type)

		node, err := p.factor()
		if err != nil {
			return nil, err
		}

		return &ast.UnaryOp{Expr: node, Op: token}, nil
	case tokens.INTEGER:
		p.consume(tokens.INTEGER)
		return &ast.Num{Token: token, Value: token.Value}, nil
	case tokens.IDENTIFIER:
		p.consume(tokens.IDENTIFIER)
		return &ast.Variable{Token: token, Name: token.Value}, nil
	case tokens.LPAREN:
		p.consume(tokens.LPAREN)

		node, err := p.arithmexpr()
		if err != nil {
			return nil, err
		}

		if _, err := p.consume(tokens.RPAREN); err != nil {
			return nil, err
		}

		return node, nil
	}
	return nil, fmt.Errorf("unexpected factor type: %s", token)
}

// term : factor ((MUL | DIV) factor)*
func (p *Parser) term() (ast.Node, error) {
	node, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.current().Type == tokens.MULT || p.current().Type == tokens.DIV {
		token, _ := p.consume(p.current().Type)

		right, err := p.factor()
		if err != nil {
			return nil, err
		}

		node = &ast.BinOp{Left: node, Op: token, Right: right}
	}

	return node, nil
}

// arithmexpr : term ((PLUS | MINUS) term)*
func (p *Parser) arithmexpr() (ast.Node, error) {
	node, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.current().Type == tokens.PLUS || p.current().Type == tokens.MINUS {
		token, _ := p.consume(p.current().Type)

		right, err := p.term()
		if err != nil {
			return nil, err
		}

		node = &ast.BinOp{Left: node, Op: token, Right: right}
	}

	return node, nil
}

// strexpr : str (CONCAT str)*
func (p *Parser) strexpr() (ast.Node, error) {
	node, err := p.str()
	if err != nil {
		return nil, err
	}

	for p.current().Type == tokens.CONCAT {
		token, _ := p.consume(tokens.CONCAT)

		right, err := p.str()
		if err != nil {
			return nil, err
		}
		node = &ast.BinOp{Left: node, Op: token, Right: right}
	}
	return node, nil
}

// str : STRING | IDENTIFIER
func (p *Parser) str() (ast.Node, error) {
	token := p.current()
	switch token.Type {
	case tokens.STRING:
		p.consume(tokens.STRING)
		return &ast.String{Token: token, Value: token.Value}, nil
	case tokens.IDENTIFIER:
		p.consume(tokens.IDENTIFIER)
		return &ast.Variable{Token: token, Name: token.Value}, nil
	}
	return nil, fmt.Errorf("a string was expected: %s", token)
}
