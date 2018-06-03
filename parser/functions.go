package parser

import (
	"github.com/fchoquet/cairn/ast"
	"github.com/fchoquet/cairn/tokens"
)

func looksLikeFunctionDecl(tk *tokens.Token) bool {
	return tk.Type == tokens.FUNC
}

func (p *Parser) functionDecl() (*ast.FuncDecl, error) {
	tk, err := p.consume(tokens.FUNC)
	if err != nil {
		return nil, err
	}

	name, err := p.consume(tokens.IDENTIFIER)
	if err != nil {
		return nil, err
	}

	sign, err := p.signature()
	if err != nil {
		return nil, err
	}

	body, err := p.block()
	if err != nil {
		return nil, err
	}

	return &ast.FuncDecl{
		Token:     tk,
		Name:      name,
		Signature: sign,
		Body:      body,
	}, nil
}

func (p *Parser) signature() (*ast.Signature, error) {
	pl, err := p.parameterList()
	if err != nil {
		return nil, err
	}

	returnType, err := p.typeId()
	if err != nil {
		return nil, err
	}

	return &ast.Signature{
		Token:      pl.Token,
		Parameters: pl,
		ReturnType: returnType,
	}, nil
}

func (p *Parser) parameterList() (*ast.ParameterList, error) {
	lparen, err := p.consume(tokens.LPAREN)
	if err != nil {
		return nil, err
	}

	parameters := []*ast.Parameter{}

	index := 0
	for tk := p.current(); tk != nil && tk.Type != tokens.RPAREN; tk = p.current() {
		// we expect a comma between each parameter
		if index > 0 {
			if _, err := p.consume(tokens.COMMA); err != nil {
				return nil, err
			}
		}

		param, err := p.parameterDecl()
		if err != nil {
			return nil, err
		}
		parameters = append(parameters, param)
		index++
	}

	if _, err := p.consume(tokens.RPAREN); err != nil {
		return nil, err
	}

	return &ast.ParameterList{
		Token:      lparen,
		Parameters: parameters,
	}, nil
}

func (p *Parser) parameterDecl() (*ast.Parameter, error) {
	name, err := p.consume(tokens.IDENTIFIER)
	if err != nil {
		return nil, err
	}

	typeId, err := p.typeId()
	if err != nil {
		return nil, err
	}

	return &ast.Parameter{
		Name:  name.Value,
		Token: name,
		Type:  typeId,
	}, nil
}
