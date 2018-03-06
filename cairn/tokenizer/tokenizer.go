package tokenizer

import (
	"errors"
	"fmt"

	"github.com/fchoquet/cairn/cairn/tokens"
)

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

func (t *Tokenizer) NextToken() (*tokens.Token, error) {
	tk := <-t.Channel

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

func readString(input string) (string, error) {
	if input == "" {
		return "", errors.New("unexpected end of string")
	}

	if input[0] != '"' {
		return "", errors.New("expected \" at the beginning of string")
	}

	if len(input) < 2 {
		return "", errors.New("a string litteral must be at least 2 chars long (including surrounding quotes)")
	}

	return readStringContents(input[1:])
}

func readStringContents(text string) (string, error) {
	if len(text) == 0 {
		return "", errors.New("could not find end of string litteral")
	}

	// let's skip the initial "
	head := text[0]
	tail := text[1:]

	var left string

	switch head {
	case '"':
		// end of string reached. stop recursion
		return "", nil
	case '\\':
		var err error
		left, err = readEscapeSequence(tail)
		if err != nil {
			return "", err
		}
		tail = tail[len(left):]
	case '\n':
		return "", errors.New("could not find end of string litteral")
	default:
		left = string(head)
	}

	// reads the rest of the string recursively
	right, err := readStringContents(tail)
	if err != nil {
		return "", err
	}

	return (left + right), nil
}

func readEscapeSequence(text string) (string, error) {
	// starting an escape sequence
	if len(text) == 0 {
		return "", errors.New("invalid escape sequence. did you mean \\\\?")
	}

	switch text[0] {
	case '\\':
		return "\\", nil
	case '"':
		return "\"", nil
	case 'n':
		return "\n", nil
	default:
		return "", errors.New("invalid escape sequence")
	}
}

func isWhiteSpace(char byte) bool {
	return char <= ' '
}

func isDigit(char byte) bool {
	return char >= '0' && char <= '9'
}
