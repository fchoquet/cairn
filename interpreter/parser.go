package interpreter

// see https://ruslanspivak.com/lsbasi-part4/

import (
	"errors"
	"fmt"
)

type Parser struct {
	Lexer        *Lexer
	CurrentToken *Token
}

func (p *Parser) Parse() (Node, error) {
	// set current token to the first token taken from the input
	currentToken, err := p.Lexer.GetNextToken()
	if err != nil {
		return "", err
	}
	p.CurrentToken = currentToken

	return p.expr()
}

// compare the current token type with the passed token
// type and if they match then "eat" the current token
// and assign the next token to the self.current_token,
// otherwise raise an exception.
func (p *Parser) eat(tokenType TokenType) (err error) {
	if p.CurrentToken == nil {
		return errors.New("no more input available")
	}

	if p.CurrentToken.Type == tokenType {
		currentToken, err := p.Lexer.GetNextToken()
		if err == nil {
			p.CurrentToken = currentToken
		}
		return err
	}

	return fmt.Errorf("expected %s but could not find it", tokenType)
}

// factor : (PLUS|MINUS)factor | INTEGER | LPAREN expr RPAREN
func (p *Parser) factor() (Node, error) {
	token := p.CurrentToken

	switch token.Type {
	case PLUS, MINUS:
		if err := p.eat(token.Type); err != nil {
			return nil, err
		}
		node, err := p.factor()
		if err != nil {
			return nil, err
		}

		return &UnaryOp{Expr: node, Op: token}, nil
	case INTEGER:
		if err := p.eat(INTEGER); err != nil {
			return nil, err
		}
		return &Num{token, token.Value}, nil
	case LPAREN:
		if err := p.eat(LPAREN); err != nil {
			return nil, err
		}

		node, err := p.expr()
		if err != nil {
			return nil, err
		}

		if err := p.eat(RPAREN); err != nil {
			return nil, err
		}

		return node, nil
	}
	return nil, fmt.Errorf("unexpected factor type: %s", token)
}

// term : factor ((MUL | DIV) factor)*
func (p *Parser) term() (Node, error) {
	node, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.CurrentToken.Type == MULT || p.CurrentToken.Type == DIV {
		token := p.CurrentToken
		p.eat(p.CurrentToken.Type)

		right, err := p.factor()
		if err != nil {
			return nil, err
		}

		node = &BinOp{Left: node, Op: token, Right: right}
	}

	return node, nil
}

// Arithmetic expression parser / interpreter.
//
// expr   : term ((PLUS | MINUS) term)*
// term   : factor ((MUL | DIV) factor)*
// factor : INTEGER | LPAREN expr RPAREN
func (p *Parser) expr() (Node, error) {
	node, err := p.term()
	if err != nil {
		return "", err
	}

	for p.CurrentToken.Type == PLUS || p.CurrentToken.Type == MINUS {
		token := p.CurrentToken
		p.eat(p.CurrentToken.Type)

		right, err := p.term()
		if err != nil {
			return "", err
		}

		node = &BinOp{Left: node, Op: token, Right: right}
	}

	return node, nil
}
