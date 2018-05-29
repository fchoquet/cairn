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

	t.Run("reads booleans", func(t *testing.T) {
		fixtures := []struct {
			input    string
			expected string
		}{
			{`true`, `true`},
			{`false`, `false`},
		}

		for _, f := range fixtures {
			token, err := Tokenize("test.ca", f.input).NextToken()
			if !assert.Nil(err) {
				continue
			}
			assert.Equal(tokens.BOOL, token.Type)
			assert.Equal(f.expected, token.Value)
		}

	})
}

func TestBasicExpressions(t *testing.T) {
	assert := assert.New(t)

	t.Run("valid basic expressions", func(t *testing.T) {
		fixtures := []struct {
			input    string
			expected string
		}{
			{`"foo" ++ "bar"`, `foo:STRING,++:CONCAT,bar:STRING`},
			{`12 + 34`, `12:INTEGER,+:PLUS,34:INTEGER`},
			{`12 - 34`, `12:INTEGER,-:MINUS,34:INTEGER`},
			{`12 * 34`, `12:INTEGER,*:MULT,34:INTEGER`},
			{`12 / 34`, `12:INTEGER,/:DIV,34:INTEGER`},
			{`12^34`, `12:INTEGER,^:POW,34:INTEGER`},
			{`12 * (34 + 56)`, `12:INTEGER,*:MULT,LPAREN:LPAREN,34:INTEGER,+:PLUS,56:INTEGER,RPAREN:RPAREN`},
			{`12 == 34`, `12:INTEGER,==:EQ,34:INTEGER`},
			{`12 != 34`, `12:INTEGER,!=:NEQ,34:INTEGER`},
			{`!true`, `!:NOT,true:BOOL`},
			{`true && false`, `true:BOOL,&&:AND,false:BOOL`},
			{`true || false`, `true:BOOL,||:OR,false:BOOL`},
		}

		for _, f := range fixtures {
			tks, err := Tokenize("test.ca", f.input).Flush()
			if !assert.Nil(err) {
				continue
			}

			stringTks := []string{}
			for _, tk := range tks {
				stringTks = append(stringTks, tk.String())
			}

			assert.Equal(f.expected, strings.Join(stringTks, ","))
		}

	})

	t.Run("syntaxically invalid basic expressions", func(t *testing.T) {
		fixtures := []string{
			`12 = 34`,
			`true & false`,
			`true | false`,
		}

		for _, f := range fixtures {
			_, err := Tokenize("test.ca", f).Flush()
			assert.Error(err)
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
			{`foo := 123`, `foo:IDENTIFIER,:=:ASSIGN,123:INTEGER`},
			{`bar := "bar"`, `bar:IDENTIFIER,:=:ASSIGN,bar:STRING`},
			{`foo := 1+2`, `foo:IDENTIFIER,:=:ASSIGN,1:INTEGER,+:PLUS,2:INTEGER`},
		}

		for _, f := range fixtures {
			tks, err := Tokenize("test.ca", f.input).Flush()
			if !assert.Nil(err) {
				continue
			}

			stringTks := []string{}
			for _, tk := range tks {
				stringTks = append(stringTks, tk.String())
			}

			assert.Equal(f.expected, strings.Join(stringTks, ","))
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
