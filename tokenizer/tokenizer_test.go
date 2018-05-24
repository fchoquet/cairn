package tokenizer

import (
	"testing"

	"github.com/fchoquet/cairn/tokens"
	"github.com/stretchr/testify/assert"
)

func TestBuffer(t *testing.T) {
	assert := assert.New(t)
	t.Run("reads tokens until the end of input", func(t *testing.T) {
		source := `12 + 34`
		tokenizer := Tokenize("test.ca", source)

		t1, err := tokenizer.NextToken()
		if !assert.Nil(err) {
			return
		}
		assert.Equal(tokens.INTEGER, t1.Type)
		assert.Equal("12", t1.Value)

		t2, err := tokenizer.NextToken()
		if !assert.Nil(err) {
			return
		}
		assert.Equal(tokens.PLUS, t2.Type)

		t3, err := tokenizer.NextToken()
		if !assert.Nil(err) {
			return
		}
		assert.Equal(tokens.INTEGER, t3.Type)
		assert.Equal("34", t3.Value)

		t4, err := tokenizer.NextToken()
		if !assert.Nil(err) {
			return
		}
		assert.Equal(tokens.EOF, t4.Type)

		_, err = tokenizer.NextToken()
		assert.Error(err)
	})
}
