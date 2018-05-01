package parser

import (
	"testing"

	"github.com/fchoquet/cairn/cairn/tokenizer"
	"github.com/fchoquet/cairn/cairn/tokens"
	"github.com/stretchr/testify/assert"
)

func TestBuffer(t *testing.T) {
	assert := assert.New(t)
	t.Run("look ahead up to 2 tokens", func(t *testing.T) {
		source := `12 + 34 + 56`

		tokenizer := tokenizer.Tokenize("test.ca", source)
		buffer := NewTokenBuffer(tokenizer, 2)

		tk, err := buffer.LookAhead(0)
		if !assert.Nil(err) || !assert.NotNil(tk) {
			return
		}
		assert.Equal(tokens.INTEGER, tk.Type)
		assert.Equal("12", tk.Value)

		tk, err = buffer.LookAhead(1)
		if !assert.Nil(err) || !assert.NotNil(tk) {
			return
		}
		assert.Equal(tokens.PLUS, tk.Type)

		// now consumes a token and return the current one
		tk, buffer, err = buffer.Consume()
		if !assert.Nil(err) || !assert.NotNil(tk) {
			return
		}
		// returns the current token
		assert.Equal(tokens.INTEGER, tk.Type)
		assert.Equal("12", tk.Value)

		// we can look ahead further
		tk, err = buffer.LookAhead(0)
		if !assert.Nil(err) || !assert.NotNil(tk) {
			return
		}
		assert.Equal(tokens.PLUS, tk.Type)

		tk, err = buffer.LookAhead(1)
		if !assert.Nil(err) || !assert.NotNil(tk) {
			return
		}
		assert.Equal(tokens.INTEGER, tk.Type)
		assert.Equal("34", tk.Value)

		// now reads till the end of input
		// moves to +
		tk, buffer, err = buffer.Consume()
		if !assert.Nil(err) || !assert.NotNil(tk) {
			return
		}

		// moves to 12
		tk, buffer, err = buffer.Consume()
		if !assert.Nil(err) || !assert.NotNil(tk) {
			return
		}

		// moves to +
		tk, buffer, err = buffer.Consume()
		if !assert.Nil(err) || !assert.NotNil(tk) {
			return
		}

		// moves to 56
		tk, buffer, err = buffer.Consume()
		if !assert.Nil(err) || !assert.NotNil(tk) {
			return
		}
		assert.Equal(tokens.INTEGER, tk.Type)
		assert.Equal("56", tk.Value)

		// now LookAhead after the end of the input
		tk, err = buffer.LookAhead(0)
		if !assert.Nil(err) || !assert.NotNil(tk) {
			return
		}
		assert.Equal(tokens.EOF, tk.Type)

		tk, err = buffer.LookAhead(1)
		assert.Error(err)

		// consume the EOF token
		tk, buffer, err = buffer.Consume()
		assert.Nil(err)
		assert.NotNil(tk)
		assert.Equal(tokens.EOF, tk.Type)

		// consume one token after the end of file => error
		tk, buffer, err = buffer.Consume()
		assert.Error(err)
	})

}
