package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	assert := assert.New(t)
	t.Run("expressions", func(t *testing.T) {
		fixtures := []struct {
			source string
			ast    string
		}{
			{
				`12`,
				`Num(12:INTEGER)`,
			},
			{
				`12 + 34`,
				`BinOp(+:PLUS Num(12:INTEGER) Num(34:INTEGER))`,
			},
			{
				`1 + 2 + 3`,
				`BinOp(+:PLUS BinOp(+:PLUS Num(1:INTEGER) Num(2:INTEGER)) Num(3:INTEGER))`,
			},
			{
				`1 + 2 * 3`,
				`BinOp(+:PLUS Num(1:INTEGER) BinOp(*:MULT Num(2:INTEGER) Num(3:INTEGER)))`,
			},
			{
				`(1 + 2) * 3`,
				`BinOp(*:MULT BinOp(+:PLUS Num(1:INTEGER) Num(2:INTEGER)) Num(3:INTEGER))`,
			},
			{
				`(1 - 2) - 3`,
				`BinOp(-:MINUS BinOp(-:MINUS Num(1:INTEGER) Num(2:INTEGER)) Num(3:INTEGER))`,
			},
			{
				`1 - 2 - 3`,
				`BinOp(-:MINUS BinOp(-:MINUS Num(1:INTEGER) Num(2:INTEGER)) Num(3:INTEGER))`,
			},
			{
				`1 - (2 - 3)`,
				`BinOp(-:MINUS Num(1:INTEGER) BinOp(-:MINUS Num(2:INTEGER) Num(3:INTEGER)))`,
			},
			{
				`-1`,
				`UnaryOp(-:MINUS Num(1:INTEGER))`,
			},
			{
				`----1`,
				`UnaryOp(-:MINUS UnaryOp(-:MINUS UnaryOp(-:MINUS UnaryOp(-:MINUS Num(1:INTEGER)))))`,
			},
			{
				`-1 + -+4`,
				`BinOp(+:PLUS UnaryOp(-:MINUS Num(1:INTEGER)) UnaryOp(-:MINUS UnaryOp(+:PLUS Num(4:INTEGER))))`,
			},
			{
				`2^4 + 2 * (3^2 - 1)`,
				`BinOp(+:PLUS BinOp(^:POW Num(2:INTEGER) Num(4:INTEGER)) BinOp(*:MULT Num(2:INTEGER) BinOp(-:MINUS BinOp(^:POW Num(3:INTEGER) Num(2:INTEGER)) Num(1:INTEGER))))`,
			},
			{
				`"foo"`,
				`String(foo:STRING)`,
			},
			{
				`"foo" ++ "bar"`,
				`BinOp(++:CONCAT String(foo:STRING) String(bar:STRING))`,
			},
			{
				`"foo" ++ "bar" ++ "baz"`,
				`BinOp(++:CONCAT BinOp(++:CONCAT String(foo:STRING) String(bar:STRING)) String(baz:STRING))`,
			},
			{
				`"foo" ++ ("bar" ++ "baz")`,
				`BinOp(++:CONCAT String(foo:STRING) BinOp(++:CONCAT String(bar:STRING) String(baz:STRING)))`,
			},
		}

		for _, f := range fixtures {
			parser := Parser{}
			node, err := parser.Parse("test.ca", f.source)
			if !assert.Nil(err) {
				break
			}
			assert.Equal(f.ast, node.String())
		}
	})
}
