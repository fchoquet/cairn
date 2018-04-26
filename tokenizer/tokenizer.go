package tokenizer

import (
	"bytes"

	"github.com/fchoquet/cairn/tokens"
)

//Tokenize converts source code into tokens
func Tokenize(fileName string, sourceCode []byte) ([]*tokens.Token, error) {
	tokens := []*tokens.Token{}

	// splits into lines
	lines := bytes.Split(sourceCode, []byte("\n"))

	for lineNumber, lineCode := range lines {
		if lineNumber < len(lines)-1 {
			// restores the trailing \n
			lineCode = append(lineCode, "\n"...)
		}

		tks, err := readLine(fileName, lineNumber, lineCode)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, tks...)
	}

	return tokens, nil
}

// readLine extracts all the tokens from a line
func readLine(fileName string, line int, code []byte) ([]*tokens.Token, error) {
	return read(tokens.Position{
		File: fileName,
		Line: line,
	}, code)
}

// read recursively extracts tokens from a string
// this function is never called from the middle of a multibyte token
func read(pos tokens.Position, code []byte) ([]*tokens.Token, error) {
	tks := []*tokens.Token{}

	if len(code) == 0 {
		return tks, nil
	}

	head := code[0]
	tail := code[1:]

	switch {
	// new line
	case head == '\n':
		tks = append(tks, createToken(pos, tokens.NewLine, ""))
	// indentation
	case head == ' ':
		tk, newTail, err := readIndentationToken(pos, code)
		if err != nil {
			return nil, err
		}
		if tk != nil {
			tks = append(tks, tk)
		}
		tail = newTail
	// number
	case head >= '0' && head <= '9':
		tk, newTail, err := readNumberToken(pos, code)
		if err != nil {
			return nil, err
		}
		if tk != nil {
			tks = append(tks, tk)
		}
		tail = newTail
	// string litteral
	case head == '"':
		tk, newTail, err := readStringToken(pos, code)
		if err != nil {
			return nil, err
		}
		if tk != nil {
			tks = append(tks, tk)
		}
		tail = newTail
	// plus operator
	case head == '+':
		tks = append(tks, createToken(pos, tokens.Operator, "+"))
	case head == '*':
		tks = append(tks, createToken(pos, tokens.Operator, "*"))
	default:

	}

	// recursively consumes the tail
	nextTks, err := read(tokens.Position{
		File: pos.File,
		Line: pos.Line,
		Col:  pos.Col + len(code) - len(tail),
	}, tail)
	tks = append(tks, nextTks...)
	return tks, err
}

func createToken(pos tokens.Position, tokenType tokens.TokenType, value string) *tokens.Token {
	return &tokens.Token{
		Type:     tokenType,
		Value:    value,
		Position: pos,
	}
}

// reads an indentation and returns the token along with the remaining of the line
func readIndentationToken(pos tokens.Position, code []byte) (*tokens.Token, []byte, error) {
	if string(code[0:3]) == "    " {
		return createToken(pos, tokens.Indentation, ""), code[:4], nil
	}
	// the space is consumed without geneating a token. This is just a space, no big deal
	return nil, code[1:], nil
}

// reads a number and returns the token along with the remaining of the line
func readNumberToken(pos tokens.Position, code []byte) (*tokens.Token, []byte, error) {
	num := readNumber(code)
	return createToken(pos, tokens.Number, num), code[len(num):], nil
}

func readNumber(code []byte) string {
	if len(code) == 0 {
		return ""
	}

	head := code[0]
	tail := code[1:]

	if head >= '0' && head <= '9' {
		return string(head) + readNumber(tail)
	}
	return ""
}

// reads a string and returns the token along with the remaining of the line
func readStringToken(pos tokens.Position, code []byte) (*tokens.Token, []byte, error) {
	s := ""

	// the first char is a ", let's skip it
	for i := 1; i < len(code); i++ {
		switch code[i] {
		case '\\':
			// next char is escaped. Let's consume it now and skip it in the iteration
			if i == len(code)-1 {
				// \ should not be the last char
				return nil, nil, Error{
					Message:  "\\ in string litteral does not escape anything",
					Code:     code,
					Position: pos,
					Hint:     "did you mean \\\\ ?",
				}
			}
			switch code[i+1] {
			case '\\':
				s += "\\"
			case '"':
				s += `"`
			case 'n':
				s += "\n"
			default:
				return nil, nil, Error{
					Message:  "invalid escape sequence in string",
					Code:     code,
					Position: pos,
					Hint:     `valid escape sequences are: \\ \" \n`,
				}
			}
			// we've consumed 2 chars
			i++
		case '"':
			return createToken(pos, tokens.String, s), code[i+1:], nil
		default:
			// let's add this char to the string
			s += string(code[i])
		}
	}
	// could not find the end of the string
	// multiline strings are not allowed for now
	return nil, nil, Error{
		Message:  "could not find end of string",
		Code:     code,
		Position: pos,
		Hint:     `you must have forgotten the closing " or forgotten to escape one`,
	}
}
