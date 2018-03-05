package interpreter

import "fmt"

type Lexer struct {
	Text string
	Pos  int
}

func (l *Lexer) currentChar() (char byte, ok bool) {
	if l.Pos > len(l.Text)-1 {
		return
	}
	return l.Text[l.Pos], true
}

// Advance the 'Pos' pointer and set the 'CurrentChar' variable.
func (l *Lexer) advance() {
	l.Pos++
}

func (l *Lexer) skipWhitespace() {
	for {
		ch, ok := l.currentChar()
		if !ok || ch != ' ' {
			break
		}
		fmt.Println("whitespace skipped")
		l.advance()
	}
}

// Return a (multidigit) integer consumed from the input
func (l *Lexer) integer() string {
	result := ""

	for {
		ch, ok := l.currentChar()
		if !ok || !isDigit(ch) {
			break
		}
		result += string(ch)
		l.advance()
	}
	return result
}

func isWhiteSpace(char byte) bool {
	return char <= ' '
}

func isDigit(char byte) bool {
	return char >= '0' && char <= '9'
}

// Lexical analyzer (also known as scanner or tokenizer)
// This method is responsible for breaking a sentence
// apart into tokens. One token at a time.
func (l *Lexer) GetNextToken() (*Token, error) {
	for {
		ch, ok := l.currentChar()
		if !ok {
			break
		}

		if isWhiteSpace(ch) {
			l.advance()
			continue
		}

		if isDigit(ch) {
			return &Token{INTEGER, l.integer()}, nil
		}

		if ch == '+' {
			l.advance()
			return &Token{PLUS, "+"}, nil
		}

		if ch == '-' {
			l.advance()
			return &Token{MINUS, "-"}, nil
		}

		if ch == '*' {
			l.advance()
			return &Token{MULT, "*"}, nil
		}

		if ch == '/' {
			l.advance()
			return &Token{DIV, "/"}, nil
		}

		if ch == '(' {
			l.advance()
			return &Token{LPAREN, "("}, nil
		}

		if ch == ')' {
			l.advance()
			return &Token{RPAREN, ")"}, nil
		}

		return nil, fmt.Errorf("syntax error in %s", l.Text)
	}

	return &Token{EOF, ""}, nil
}
