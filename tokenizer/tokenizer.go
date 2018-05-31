package tokenizer

import (
	"errors"
	"fmt"
	"strconv"

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
		}, 0)

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
func (t *Tokenizer) tokenize(text string, pos tokens.Position, indent int) {
	if len(text) == 0 {
		t.yieldToken(tokens.EOF, "", pos)
		return
	}

	head := text[0]
	tail := text[1:]

	switch {
	case head == '\n':
		oldIndent := indent
		var consumed int
		indent, consumed = consumeTab(tail)
		diff := indent - oldIndent
		switch {
		case diff > 0:
			// indentation increased => begin block
			for i := 0; i < diff; i++ {
				t.yieldToken(tokens.BEGIN, "BEGIN"+strconv.Itoa(oldIndent+1+i), pos)
			}
		case diff < 0:
			// indentation decreased => end block
			for i := 0; i < -diff; i++ {
				t.yieldToken(tokens.END, "END"+strconv.Itoa(oldIndent-i), pos)
			}
		default:
			// no indentation change. Simply yields an EOL
			t.yieldToken(tokens.EOL, "EOL", pos)
		}

		pos.Col += 1 + consumed
	case isWhiteSpace(head):
		pos.Col++
		// simply skip
	case isDigit(head):
		value := readInteger(text)
		tail = text[len(value):]
		t.yieldToken(tokens.INTEGER, value, pos)
		pos.Col += len(value)
	case isAlpha(head):
		value := readIdentifier(text)
		tail = text[len(value):]
		// keywords should not be treated as identifiers!
		switch value {
		case "true", "false":
			t.yieldToken(tokens.BOOL, value, pos)
		case "func":
			t.yieldToken(tokens.FUNC, value, pos)
		default:
			t.yieldToken(tokens.IDENTIFIER, value, pos)
		}
		pos.Col += len(value)
	case head == ',':
		t.yieldToken(tokens.COMMA, "COMMA", pos)
		pos.Col++
	case head == '+':
		if len(tail) > 0 && tail[0] == '+' {
			tail = text[2:]
			t.yieldToken(tokens.CONCAT, "++", pos)
			pos.Col += 2
		} else {
			t.yieldToken(tokens.PLUS, "+", pos)
			pos.Col++
		}
	case head == '-':
		t.yieldToken(tokens.MINUS, "-", pos)
		pos.Col++
	case head == '*':
		t.yieldToken(tokens.MULT, "*", pos)
		pos.Col++
	case head == '/':
		t.yieldToken(tokens.DIV, "/", pos)
		pos.Col++
	case head == '^':
		t.yieldToken(tokens.POW, "^", pos)
		pos.Col++
	case head == '(':
		t.yieldToken(tokens.LPAREN, "LPAREN", pos)
		pos.Col++
	case head == ')':
		t.yieldToken(tokens.RPAREN, "RPAREN", pos)
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
		if len(tail) > 0 && tail[0] == '=' {
			tail = tail[1:]
			t.yieldToken(tokens.ASSIGN, ":=", pos)
			pos.Col += 2
		} else {
			t.yieldToken(tokens.COLUMN, "COLUMN", pos)
			pos.Col++
		}
	case head == '=':
		// might be an equality comparison
		if len(tail) == 0 || tail[0] != '=' {
			t.yieldToken(tokens.ERROR, "Unexpected :", pos)
			return
		}

		tail = tail[1:]
		t.yieldToken(tokens.EQ, "==", pos)
		pos.Col += 2
	case head == '!':
		if len(tail) > 0 && tail[0] == '=' {
			tail = tail[1:]
			t.yieldToken(tokens.NEQ, "!=", pos)
			pos.Col += 2
		} else {
			t.yieldToken(tokens.NOT, "!", pos)
			pos.Col++
		}
	case head == '|':
		if len(tail) > 0 && tail[0] == '|' {
			tail = tail[1:]
			t.yieldToken(tokens.OR, "||", pos)
			pos.Col += 2
		} else {
			t.yieldToken(tokens.ERROR, fmt.Sprintf("syntax error: unexpected | in %s", text), pos)
			return
		}
	case head == '&':
		if len(tail) > 0 && tail[0] == '&' {
			tail = tail[1:]
			t.yieldToken(tokens.AND, "&&", pos)
			pos.Col += 2
		} else {
			t.yieldToken(tokens.ERROR, fmt.Sprintf("syntax error: unexpected & in %s", text), pos)
			return
		}
	default:
		t.yieldToken(tokens.ERROR, fmt.Sprintf("syntax error in %s", text), pos)
		// stop recursion
		return
	}

	// recursively tokenize the rest of the string
	t.tokenize(tail, pos, indent)
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

	if isAlpha(head) {
		return string(head) + readIdentifier(tail)
	}

	return ""
}

func isWhiteSpace(char byte) bool {
	return char <= ' ' && char != '\n'
}

func isDigit(char byte) bool {
	return char >= '0' && char <= '9'
}

func isAlpha(char byte) bool {
	return (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z') || char == '_'
}

func consumeTab(s string) (tabs int, consumed int) {
	if len(s) == 0 {
		return
	}

	switch s[0] {
	// space
	case ' ':
		if len(s) < 4 {
			return
		}
		if string(s[0:4]) == "    " {
			tabs, consumed = consumeTab(s[4:])
			return tabs + 1, consumed + 4
		}
	// tab
	case '	':
		tabs, consumed = consumeTab(s[1:])
		return tabs + 1, consumed + 1
	}
	return
}

// Flush all remaining tokens
func (t *Tokenizer) Flush() ([]*tokens.Token, error) {
	tks := []*tokens.Token{}
	tk, err := t.NextToken()
	for ; err == nil && tk != nil && tk.Type != tokens.EOF; tk, err = t.NextToken() {
		tks = append(tks, tk)
	}
	return tks, err
}
