package parser

import (
	"fmt"

	"github.com/fchoquet/cairn/cairn/ast"
	"github.com/fchoquet/cairn/cairn/tokenizer"
	"github.com/fchoquet/cairn/cairn/tokens"
)

// Parser reads a text and converts it to an AST using the Tokenizer
type Parser struct {
	Tokenizer    *tokenizer.Tokenizer
	CurrentToken *tokens.Token
}

// Parse builds an AST from a text
func (p *Parser) Parse(fileName, text string) (ast.Node, error) {
	// DEBUG CODE
	// display all the tokens
	fmt.Println("--- TOKENS: ---")
	tk := tokenizer.Tokenize(fileName, text)
	for token := range tk.Channel {
		fmt.Println(token)
	}
	fmt.Println("---------------")
	// END DEBUG CODE

	p.Tokenizer = tokenizer.Tokenize(fileName, text)

	// move to the 1st token
	if err := p.advance(); err != nil {
		return nil, err
	}

	return p.expr()
}

func (p *Parser) advance() error {
	token, err := p.Tokenizer.NextToken()
	if err != nil {
		return err
	}

	p.CurrentToken = token
	return nil
}

// eat test that the current token is of the expected type
// it consumes it if it's the case
// it returns an error if the types do not match
func (p *Parser) eat(tkType tokens.TokenType) error {
	if p.CurrentToken.Type != tkType {
		return fmt.Errorf("wrong input type. Expected %s - got %s", tkType, p.CurrentToken)
	}

	return p.advance()
}

// expr : assignment | arithmexpr | strexpr
func (p *Parser) expr() (ast.Node, error) {
	token := p.CurrentToken

	switch token.Type {
	case tokens.IDENTIFIER:
		return p.assignment()
	case tokens.STRING:
		return p.strexpr()
	default:
		return p.arithmexpr()
	}
}

// assignment : IDENTIFIER ASSIGN expr
func (p *Parser) assignment() (ast.Node, error) {
	idToken := p.CurrentToken
	if err := p.eat(tokens.IDENTIFIER); err != nil {
		return nil, err
	}

	assignToken := p.CurrentToken
	if err := p.eat(tokens.ASSIGN); err != nil {
		return nil, err
	}

	right, err := p.expr()
	if err != nil {
		return nil, err
	}

	return &ast.Assignment{Identifier: idToken.Value, Token: assignToken, Right: right}, nil
}

// factor : (PLUS|MINUS)factor | INTEGER | LPAREN arithmexpr RPAREN
func (p *Parser) factor() (ast.Node, error) {
	token := p.CurrentToken

	switch token.Type {
	case tokens.PLUS, tokens.MINUS:
		if err := p.eat(token.Type); err != nil {
			return nil, err
		}

		node, err := p.factor()
		if err != nil {
			return nil, err
		}

		return &ast.UnaryOp{Expr: node, Op: token}, nil
	case tokens.INTEGER:
		if err := p.eat(tokens.INTEGER); err != nil {
			return nil, err
		}
		return &ast.Num{Token: token, Value: token.Value}, nil
	case tokens.LPAREN:
		if err := p.eat(tokens.LPAREN); err != nil {
			return nil, err
		}

		node, err := p.arithmexpr()
		if err != nil {
			return nil, err
		}

		if err := p.eat(tokens.RPAREN); err != nil {
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

	for p.CurrentToken.Type == tokens.MULT || p.CurrentToken.Type == tokens.DIV {
		token := p.CurrentToken
		if err := p.eat(token.Type); err != nil {
			return nil, err
		}

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

	for p.CurrentToken.Type == tokens.PLUS || p.CurrentToken.Type == tokens.MINUS {
		token := p.CurrentToken
		if err := p.eat(token.Type); err != nil {
			return nil, err
		}

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

	for p.CurrentToken.Type == tokens.PLUS {
		token := p.CurrentToken
		if err := p.eat(tokens.PLUS); err != nil {
			return nil, err
		}

		right, err := p.str()
		if err != nil {
			return nil, err
		}
		node = &ast.BinOp{Left: node, Op: token, Right: right}
	}
	return node, nil
}

// str : STRING
func (p *Parser) str() (ast.Node, error) {
	token := p.CurrentToken
	if err := p.eat(tokens.STRING); err != nil {
		return nil, err
	}

	return &ast.String{Token: token, Value: token.Value}, nil
}
