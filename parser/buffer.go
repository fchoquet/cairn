package parser

import (
	"fmt"

	"github.com/fchoquet/cairn/tokenizer"
	"github.com/fchoquet/cairn/tokens"
)

// TokenBuffer allows implementation of a LL(n) Recusive Descent Parser
type TokenBuffer interface {
	// LookAhead returns the nth token (or an error if trying to read after end of file)
	LookAhead(n int) (*tokens.Token, error)

	// Consume a token and returns the new buffer (or a syntax error)
	Consume() (*tokens.Token, TokenBuffer, error)
}

// NewTokenBuffer creates a token buffer
func NewTokenBuffer(tokenizer *tokenizer.Tokenizer, size int) TokenBuffer {
	return &buffer{
		buffer:    make([]*tokens.Token, size, size),
		position:  0,
		tokenizer: tokenizer,
		size:      size,
	}
}

type buffer struct {
	buffer    []*tokens.Token
	position  int
	tokenizer *tokenizer.Tokenizer
	size      int
}

// we're using a rotating buffer. this function returns an index
// based on the current position
func (b *buffer) indexOf(n int) int {
	return (b.position + n) % b.size
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
	if n >= b.size {
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
	newBuffer := make([]*tokens.Token, b.size, b.size)
	for index, token := range b.buffer {
		if index != b.position {
			newBuffer[index] = token
		}
	}

	return tk, &buffer{
		buffer:    newBuffer,
		position:  (b.position + 1) % b.size,
		tokenizer: b.tokenizer,
		size:      b.size,
	}, nil
}
