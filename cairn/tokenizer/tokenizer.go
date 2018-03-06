package tokenizer

import (
	"errors"
	"fmt"

	"github.com/fchoquet/cairn/cairn/tokens"
)

type Tokenizer struct {
	ch chan *tokens.Token
}

// Tokenize returns a Tokenizer ready to return tokens
func Tokenize(fileName, text string) *Tokenizer {
	t := &Tokenizer{
		ch: make(chan *tokens.Token),
	}

	go func() {
		t.tokenize(text, tokens.Position{
			File: fileName,
			Line: 0,
			Col:  0,
		})

		// close the channel to notify completion
		close(t.ch)

		return
	}()

	return t
}

func (t *Tokenizer) NextToken() (*tokens.Token, error) {
	tk := <-t.ch

	if tk.Type == tokens.ERROR {
		return nil, errors.New(tk.Value)
	}

	return tk, nil
}

func (t *Tokenizer) yieldToken(tkType tokens.TokenType, value string, pos tokens.Position) {
	t.ch <- &tokens.Token{
		Type:     tkType,
		Value:    value,
		Position: pos,
	}
}

// tokenize process a string recursively
func (t *Tokenizer) tokenize(text string, pos tokens.Position) {
	if len(text) == 0 {
		t.yieldToken(tokens.EOF, "", pos)
		return
	}

	head := text[0]
	tail := text[1:]

	switch {
	case isWhiteSpace(head):
		pos.Col++
		// simply skip
	case isDigit(head):
		value := readInteger(text)
		tail = text[len(value):]
		t.yieldToken(tokens.INTEGER, value, pos)
		pos.Col += len(value)
	case head == '+':
		t.yieldToken(tokens.PLUS, "+", pos)
		pos.Col++
	case head == '-':
		t.yieldToken(tokens.MINUS, "-", pos)
		pos.Col++
	case head == '*':
		t.yieldToken(tokens.MULT, "*", pos)
		pos.Col++
	case head == '/':
		t.yieldToken(tokens.DIV, "/", pos)
		pos.Col++
	case head == '(':
		t.yieldToken(tokens.LPAREN, "(", pos)
		pos.Col++
	case head == ')':
		t.yieldToken(tokens.RPAREN, ")", pos)
		pos.Col++
	default:
		t.yieldToken(tokens.ERROR, fmt.Sprintf("syntax error in %s", text), pos)
		// stop recursion
		return
	}

	// recursively tokenize the rest of the string
	t.tokenize(tail, pos)
}

func readInteger(input string) string {
	if input == "" {
		return ""
	}

	head := input[0]
	tail := input[1:]

	if isDigit(head) {
		return string(head) + readInteger(tail)
	}

	return ""
}

func isWhiteSpace(char byte) bool {
	return char <= ' '
}

func isDigit(char byte) bool {
	return char >= '0' && char <= '9'
}
