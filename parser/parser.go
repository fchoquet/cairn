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

	return p.sourceFile()
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

func (p *Parser) sourceFile() (*ast.SourceFile, error) {
	functions := []*ast.FuncDecl{}
	for tk := p.current(); looksLikeFunctionDecl(tk); tk = p.current() {
		f, err := p.functionDecl()
		if err != nil {
			return nil, err
		}
		functions = append(functions, f)
	}

	statements, err := p.statementList()
	if err != nil {
		return nil, err
	}

	return &ast.SourceFile{
		Functions:  functions,
		Statements: statements,
	}, nil
}

func (p *Parser) statementList() (*ast.StatementList, error) {
	statements := []ast.Statement{}

	for tk := p.current(); tk != nil && tk.Type != tokens.EOF && tk.Type != tokens.END; tk = p.current() {
		// skip end of lines
		if tk.Type == tokens.EOL {
			if _, err := p.consume(tokens.EOL); err != nil {
				return nil, err
			}
			continue
		}
		st, err := p.statement()
		if err != nil {
			return nil, err
		}
		statements = append(statements, st)
	}
	return &ast.StatementList{Statements: statements}, nil
}

func (p *Parser) statement() (ast.Node, error) {
	switch p.current().Type {
	case tokens.BEGIN:
		return p.block()
	default:
		return p.simpleStmt()
	}
}

func (p *Parser) block() (*ast.BlockStmt, error) {
	begin, err := p.consume(tokens.BEGIN)
	if err != nil {
		return nil, err
	}

	sl, err := p.statementList()
	if err != nil {
		return nil, err
	}

	end, err := p.consume(tokens.END)
	if err != nil {
		return nil, err
	}

	return &ast.BlockStmt{
		Begin:      begin,
		Statements: sl,
		End:        end,
	}, nil
}

func (p *Parser) simpleStmt() (ast.Node, error) {
	tk := p.current()
	// no need to check err here. nil is fine
	next, _ := p.lookAhead(1)
	switch {
	case looksLikeAssignment(tk, next):
		return p.assignment()
	default:
		return p.expression()
	}
}

func (p *Parser) expression() (ast.Node, error) {
	return p.computeExpression(0)
}

// uses the precedence climbing algorithm here
// https://eli.thegreenplace.net/2012/08/02/parsing-expressions-by-precedence-climbing
func (p *Parser) computeExpression(minPrec int) (ast.Node, error) {
	left, err := p.unaryExpr()
	if err != nil {
		return nil, err
	}

	for {
		op := p.current()
		if op == nil || !isBinaryOp(op) || BinaryOpPrecedence[op.Type] < minPrec {
			break
		}
		// consume this token
		p.consume(op.Type)

		// Get the operator's precedence and associativity, and compute a
		// minimal precedence for the recursive call
		nextMinPrec := BinaryOpPrecedence[op.Type]
		if BinaryOpAssociativity[op.Type] == AssocLeft {
			nextMinPrec++
		}

		// Consume the current token and prepare the next one for the recursive call
		right, err := p.computeExpression(nextMinPrec)
		if err != nil {
			return nil, err
		}

		left = &ast.BinOp{Left: left, Op: op, Right: right}
	}
	return left, nil
}

func looksLikeAssignment(tk1 *tokens.Token, tk2 *tokens.Token) bool {
	return tk1.Type == tokens.IDENTIFIER && tk2 != nil && tk2.Type == tokens.ASSIGN
}

func (p *Parser) assignment() (ast.Node, error) {
	id, err := p.consume(tokens.IDENTIFIER)
	if err != nil {
		return nil, err
	}

	op, err := p.consume(tokens.ASSIGN)
	if err != nil {
		return nil, err
	}

	right, err := p.expression()
	if err != nil {
		return nil, err
	}

	return &ast.Assignment{
		Token:    op,
		Variable: ast.Variable{Token: id, Name: id.Value},
		Right:    right,
	}, nil
}

func looksLikeUnaryExpr(tk *tokens.Token) bool {
	return isUnaryOp(tk) || looksLikePrimaryExpression(tk)
}

func (p *Parser) unaryExpr() (ast.Node, error) {
	tk := p.current()

	switch {
	case isUnaryOp(tk):
		op, err := p.consume(tk.Type)
		if err != nil {
			return nil, err
		}
		right, err := p.unaryExpr()
		if err != nil {
			return nil, err
		}

		return &ast.UnaryOp{Op: op, Expr: right}, nil
	default:
		return p.primaryExpression()
	}
}

func looksLikePrimaryExpression(tk *tokens.Token) bool {
	return looksLikeOperandName(tk) || tk.Type == tokens.LPAREN || looksLikeLitteral(tk)
}

func (p *Parser) primaryExpression() (ast.Node, error) {
	return p.operand()
}

func (p *Parser) operand() (ast.Node, error) {
	tk := p.current()

	switch {
	case looksLikeOperandName(tk):
		return p.operandName()
	case tk.Type == tokens.LPAREN:
		if _, err := p.consume(tokens.LPAREN); err != nil {
			return nil, err
		}

		nd, err := p.expression()
		if err != nil {
			return nil, err
		}

		if _, err := p.consume(tokens.RPAREN); err != nil {
			return nil, err
		}

		return nd, nil
	default:
		return p.literal()
	}
}

func looksLikeOperandName(tk *tokens.Token) bool {
	return tk.Type == tokens.IDENTIFIER
}

func (p *Parser) operandName() (ast.Node, error) {
	id, err := p.consume(tokens.IDENTIFIER)
	if err != nil {
		return nil, err
	}
	return &ast.Variable{Token: id, Name: id.Value}, nil
}

func looksLikeLitteral(tk *tokens.Token) bool {
	return tk.Type == tokens.INTEGER || tk.Type == tokens.STRING || tk.Type == tokens.BOOL
}

func (p *Parser) literal() (ast.Node, error) {
	return p.basicLit()
}

func (p *Parser) basicLit() (ast.Node, error) {
	tk := p.current()
	switch tk.Type {
	case tokens.INTEGER:
		p.consume(tk.Type)
		return &ast.Num{Token: tk, Value: tk.Value}, nil
	case tokens.STRING:
		p.consume(tk.Type)
		return &ast.String{Token: tk, Value: tk.Value}, nil
	case tokens.BOOL:
		p.consume(tk.Type)
		return &ast.Bool{Token: tk, Value: tk.Value}, nil
	default:
		return nil, fmt.Errorf("unexpected basic litteral: %s:%s", tk.Value, tk.Type)
	}
}

func (p *Parser) typeId() (*ast.TypeId, error) {
	tk, err := p.consume(tokens.COLUMN)
	if err != nil {
		return nil, err
	}

	name, err := p.consume(tokens.IDENTIFIER)
	if err != nil {
		return nil, err
	}

	return &ast.TypeId{
		Token: tk,
		Name:  name.Value,
	}, nil
}
