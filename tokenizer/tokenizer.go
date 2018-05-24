package tokenizer

import (
	"errors"
	"fmt"

	"github.com/fchoquet/cairn/tokens"
)

// Tokenizer transforms a string into a stream of tokens
type Tokenizer struct {
	Channel chan *tokens.Token
}

// Tokenize returns a Tokenizer ready to return tokens
func Tokenize(fileName, text string) *Tokenizer {
	t := &Tokenizer{
		Channel: make(chan *tokens.Token),
	}

	go func() {
		t.tokenize(text, tokens.Position{
			File: fileName,
			Line: 0,
			Col:  0,
		})

		// close the channel to notify completion
		close(t.Channel)

		return
	}()

	return t
}

// NextToken yields a new token
// since the channel is blocking, we only yield a new token when this function is called
// the tokenize function does not have to worry about memory usage
func (t *Tokenizer) NextToken() (*tokens.Token, error) {
	tk, ok := <-t.Channel
	if !ok {
		return nil, errors.New("can not read after end of file")
	}

	if tk.Type == tokens.ERROR {
		return nil, errors.New(tk.Value)
	}

	return tk, nil
}

func (t *Tokenizer) yieldToken(tkType tokens.TokenType, value string, pos tokens.Position) {
	t.Channel <- &tokens.Token{
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
	case isIdentifier(head):
		value := readIdentifier(text)
		tail = text[len(value):]
		t.yieldToken(tokens.IDENTIFIER, value, pos)
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
	case head == '"':
		value, err := readString(text)
		if err != nil {
			t.yieldToken(tokens.ERROR, err.Error(), pos)
			return
		}

		// value does not contain the surrounding quotes, so let's add 2
		length := len(value) + 2
		tail = text[length:]
		t.yieldToken(tokens.STRING, value, pos)
		pos.Col += length
	case head == ':':
		// might be an assignment
		if len(tail) == 0 || tail[0] != '=' {
			t.yieldToken(tokens.ERROR, "Unexpected :", pos)
		}

		tail = tail[1:]
		t.yieldToken(tokens.ASSIGN, ":=", pos)
		pos.Col += 2
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

func readIdentifier(input string) string {
	if input == "" {
		return ""
	}

	head := input[0]
	tail := input[1:]

	if isIdentifier(head) {
		return string(head) + readIdentifier(tail)
	}

	return ""
}

func isWhiteSpace(char byte) bool {
	return char <= ' '
}

func isDigit(char byte) bool {
	return char >= '0' && char <= '9'
}

func isIdentifier(char byte) bool {
	return (char >= 'A' && char <= 'z') || char == '_'
}
