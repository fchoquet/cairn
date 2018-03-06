package tokenizer

import (
	"errors"
	"fmt"

	"github.com/fchoquet/cairn/cairn/tokens"
)

type Tokenizer struct {
	ch chan *tokens.Token
}

// Tokenize returns a Tokenizer
func Tokenize(text string) *Tokenizer {
	t := &Tokenizer{
		ch: make(chan *tokens.Token),
	}

	go func() {
		t.tokenize(text)
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

func (t *Tokenizer) yieldToken(tkType tokens.TokenType, value string) {
	t.ch <- &tokens.Token{
		Type:  tkType,
		Value: value,
	}
}

// tokenize process a string recursively
func (t *Tokenizer) tokenize(text string) {
	if len(text) == 0 {
		t.yieldToken(tokens.EOF, "")
		return
	}

	head := text[0]
	tail := text[1:]

	switch {
	case isWhiteSpace(head):
		// simply skip
	case isDigit(head):
		value := readInteger(text)
		tail = text[len(value):]
		t.yieldToken(tokens.INTEGER, value)
	case head == '+':
		t.yieldToken(tokens.PLUS, "+")
	case head == '-':
		t.yieldToken(tokens.MINUS, "-")
	case head == '*':
		t.yieldToken(tokens.MULT, "*")
	case head == '/':
		t.yieldToken(tokens.DIV, "/")
	case head == '(':
		t.yieldToken(tokens.LPAREN, "(")
	case head == ')':
		t.yieldToken(tokens.RPAREN, ")")
	default:
		t.yieldToken(tokens.ERROR, fmt.Sprintf("syntax error in %s", text))
		// stop recursion
		return
	}

	// recursively tokenize the rest of the string
	t.tokenize(tail)
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
