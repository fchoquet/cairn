package tokenizer

import (
	"strings"
	"testing"

	"github.com/fchoquet/cairn/tokens"
	"github.com/stretchr/testify/assert"
)

func TestPrimaryTypes(t *testing.T) {
	assert := assert.New(t)

	t.Run("reads strings", func(t *testing.T) {
		// I choose avoid any interpolation on the "expected" side to avoid ambiguity
		fixtures := []struct {
			input    string
			expected string
		}{
			{`""`, ``},
			{`"foo"`, `foo`},
			{`"foo bar baz"`, `foo bar baz`},
			{`"foo\nbar\nbaz"`, `foo` + "\n" + `bar` + "\n" + `baz`},
			{`"foo\\bar"`, `foo\bar`},
			{`"foo\"bar\"baz"`, `foo"bar"baz`},
		}

		for _, f := range fixtures {
			token, err := Tokenize("test.ca", f.input).NextToken()
			if !assert.Nil(err) {
				continue
			}
			assert.Equal(tokens.STRING, token.Type)
			assert.Equal(f.expected, token.Value)
		}

	})

	t.Run("detects invalid strings and escape sequences", func(t *testing.T) {
		// I choose avoid any interpolation on the "expected" side to avoid ambiguity
		fixtures := []string{
			`"\"`,
			`"foo\bar"`,
			`"foo\*bar"`,
			`"foo`,
			`"foo` + "\n" + `bar"`,
		}

		for _, f := range fixtures {
			_, err := Tokenize("test.ca", f).NextToken()
			assert.Error(err)
		}
	})

	t.Run("reads integers", func(t *testing.T) {
		fixtures := []struct {
			input    string
			expected string
		}{
			{`0`, `0`},
			{`123`, `123`},
		}

		for _, f := range fixtures {
			token, err := Tokenize("test.ca", f.input).NextToken()
			if !assert.Nil(err) {
				continue
			}
			assert.Equal(tokens.INTEGER, token.Type)
			assert.Equal(f.expected, token.Value)
		}

	})
}

func TestArithmeticExpressions(t *testing.T) {
	assert := assert.New(t)

	t.Run("basic arithmentic expressions", func(t *testing.T) {
		fixtures := []struct {
			input    string
			expected string
		}{
			{`12 + 34`, `Token(INTEGER, 12),Token(PLUS, +),Token(INTEGER, 34)`},
			{`12 - 34`, `Token(INTEGER, 12),Token(MINUS, -),Token(INTEGER, 34)`},
			{`12 * 34`, `Token(INTEGER, 12),Token(MULT, *),Token(INTEGER, 34)`},
			{`12 / 34`, `Token(INTEGER, 12),Token(DIV, /),Token(INTEGER, 34)`},
			{`12 * (34 + 56)`, `Token(INTEGER, 12),Token(MULT, *),Token(LPAREN, LPAREN),Token(INTEGER, 34),Token(PLUS, +),Token(INTEGER, 56),Token(RPAREN, RPAREN)`},
		}

		for _, f := range fixtures {
			tokenizer := Tokenize("test.ca", f.input)

			tks := []string{}
			tk, err := tokenizer.NextToken()
			for ; err == nil && tk != nil && tk.Type != tokens.EOF; tk, err = tokenizer.NextToken() {
				tks = append(tks, tk.String())
			}
			assert.Equal(f.expected, strings.Join(tks, ","))
		}

	})
}

func TestAssignments(t *testing.T) {
	assert := assert.New(t)

	t.Run("assignments", func(t *testing.T) {
		fixtures := []struct {
			input    string
			expected string
		}{
			{`foo := 123`, `Token(IDENTIFIER, foo),Token(ASSIGN, :=),Token(INTEGER, 123)`},
			{`bar := "bar"`, `Token(IDENTIFIER, bar),Token(ASSIGN, :=),Token(STRING, bar)`},
			{`foo := 1+2`, `Token(IDENTIFIER, foo),Token(ASSIGN, :=),Token(INTEGER, 1),Token(PLUS, +),Token(INTEGER, 2)`},
		}

		for _, f := range fixtures {
			tokenizer := Tokenize("test.ca", f.input)

			tks := []string{}
			tk, err := tokenizer.NextToken()
			for ; err == nil && tk != nil && tk.Type != tokens.EOF; tk, err = tokenizer.NextToken() {
				tks = append(tks, tk.String())
			}
			assert.Equal(f.expected, strings.Join(tks, ","))
		}

	})
}

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
