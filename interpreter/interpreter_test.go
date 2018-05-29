package interpreter

import (
	"strings"
	"testing"

	"github.com/fchoquet/cairn/parser"
	"github.com/stretchr/testify/assert"
)

func TestInterpreter(t *testing.T) {
	assert := assert.New(t)
	t.Run("expressions", func(t *testing.T) {
		fixtures := []struct {
			source string
			result string
		}{
			{`12`, `12`},
			{`12 + 34`, `46`},
			{`1 + 2 + 3`, `6`},
			{`1 + 2 * 3`, `7`},
			{`1 + 2 / 3`, `1`},
			{`(1 + 2) * 3`, `9`},
			{`(1 - 2) - 3`, `-4`},
			{`1 - 2 - 3`, `-4`},
			{`1 - (2 - 3)`, `2`},
			{`-1`, `-1`},
			{`----1`, `1`},
			{`-1 + -+4`, `-5`},
			{`2^4 + 2 * (3^2 - 1)`, `32`},
			// strings
			{`"foo"`, `foo`},
			{`"foo" ++ "bar"`, `foobar`},
			{`"foo" ++ "bar" ++ "baz"`, `foobarbaz`},
			{`"foo" ++ ("bar" ++ "baz")`, `foobarbaz`},
			// booleans
			{`true`, `true`},
			{`!true`, `false`},
			{`!!true`, `true`},
			{`false`, `false`},
			{`true && true`, `true`},
			{`true && false`, `false`},
			{`false && false`, `false`},
			{`true || true`, `true`},
			{`false || true`, `true`},
			{`false || false`, `false`},
			{`true && false || false`, `false`},
			{`true && (false || true)`, `true`},
			// equality
			{`true == true`, `true`},
			{`true == false`, `false`},
			{`1 == 2`, `false`},
			{`2 == 2`, `true`},
			{`"foo" == "foo"`, `true`},
			{`"foo" == "fOO"`, `false`},
			// complex operator precedence
			{`2*2==2^2 && true==(2==2)`, `true`},
		}

		i := New(&parser.Parser{})

		for _, f := range fixtures {
			result, err := i.Interpret("test.ca", f.source)
			if !assert.Nil(err) {
				break
			}
			assert.Equal(f.result, result, f.source)
		}
	})

	t.Run("variables", func(t *testing.T) {
		fixtures := []struct {
			source []string
			result string
		}{
			{[]string{`foo := 12`, `bar := 34`, `foo + bar`}, `46`},
			{[]string{`foo := "hello"`, `bar := " world"`, `foo ++ bar`}, `hello world`},
		}

		for _, f := range fixtures {
			i := New(&parser.Parser{})
			var (
				result string
				err    error
			)
			for _, s := range f.source {
				result, err = i.Interpret("test.ca", s)
				if !assert.Nil(err) {
					return
				}
			}
			assert.Equal(f.result, result, strings.Join(f.source, "\n"))
		}
	})
}
