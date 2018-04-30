package parser

import (
	"fmt"

	"github.com/fchoquet/cairn/cairn/tokenizer"
	"github.com/fchoquet/cairn/cairn/tokens"
)

const bufferSize = 2

// TokenBuffer allows implementation of a LL(n) Recusive Descent Parser
type TokenBuffer interface {
	// LookAhead returns the nth token (or an error if trying to read after end of file)
	LookAhead(n int) (*tokens.Token, error)

	// Consume a token and returns the new buffer (or a syntax error)
	Consume() (*tokens.Token, TokenBuffer, error)
}

// NewTokenBuffer creates a token buffer
func NewTokenBuffer(tokenizer *tokenizer.Tokenizer) TokenBuffer {
	return &buffer{
		buffer:    [bufferSize]*tokens.Token{},
		position:  0,
		tokenizer: tokenizer,
	}
}

type buffer struct {
	buffer    [bufferSize]*tokens.Token
	position  int
	tokenizer *tokenizer.Tokenizer
}

// we're using a rotating buffer. this function returns an index
// based on the current position
func (b *buffer) indexOf(n int) int {
	return (b.position + n) % bufferSize
}

func (b *buffer) load(n int) (*tokens.Token, error) {
	realIndex := b.indexOf(n)
	if b.buffer[realIndex] == nil {
		// load next token
		tk, err := b.tokenizer.NextToken()
		if err != nil {
			return nil, err
		}
		b.buffer[realIndex] = tk
	}
	return b.buffer[realIndex], nil
}

func (b *buffer) LookAhead(n int) (*tokens.Token, error) {
	if n >= bufferSize {
		// We do not return an error. The code has to be fixed so a panic is fine
		panic(fmt.Sprintf("parser buffer overflow"))
	}
	return b.load(n)
}

func (b *buffer) Consume() (*tokens.Token, TokenBuffer, error) {
	tk, err := b.LookAhead(0)
	if err != nil {
		return nil, nil, err
	}

	// create a new buffer with current position empty
	newBuffer := [bufferSize]*tokens.Token{}
	for index, token := range b.buffer {
		if index != b.position {
			newBuffer[index] = token
		}
	}

	return tk, &buffer{
		buffer:    newBuffer,
		position:  (b.position + 1) % bufferSize,
		tokenizer: b.tokenizer,
	}, nil
}
