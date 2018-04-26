package tokenizer

import (
	"testing"

	"github.com/fchoquet/cairn/tokens"
)

func TestReadNumber(t *testing.T) {
	fixtures := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"1", "1"},
		{"1234", "1234"},
		{"12 34", "12"},
		{" 12", ""},
		{"01AZE34", "01"},
	}

	for _, f := range fixtures {
		result := readNumber([]byte(f.input))
		if result != f.expected {
			t.Errorf("expected %s - got %s", f.expected, result)
		}
	}
}

func TestTokenize(t *testing.T) {
	t.Run("basic arithmetic", func(t *testing.T) {
		fixtures := []struct {
			source   string
			expected string
		}{
			{
				source: `12 + 3`,
				expected: `
[num:12 (0:0@/test.ca)]
[op:+ (0:3@/test.ca)]
[num:3 (0:5@/test.ca)]`,
			},
			{
				source: `12 + 3 * 999
`,
				expected: `
[num:12 (0:0@/test.ca)]
[op:+ (0:3@/test.ca)]
[num:3 (0:5@/test.ca)]
[op:* (0:7@/test.ca)]
[num:999 (0:9@/test.ca)]
[nl: (0:12@/test.ca)]`,
			},
			{
				source: `12+3`,
				expected: `
[num:12 (0:0@/test.ca)]
[op:+ (0:2@/test.ca)]
[num:3 (0:3@/test.ca)]`,
			},
		}

		for _, f := range fixtures {
			tks, err := Tokenize("/test.ca", []byte(f.source))
			if err != nil {
				t.Errorf("unexpected error: %s", err.Error())
			}

			if printTokens(tks) != f.expected {
				t.Errorf("expected %s - got %s", f.expected, printTokens(tks))
			}
		}
	})

	t.Run("strings", func(t *testing.T) {
		fixtures := []struct {
			source   string
			expected string
		}{
			{
				source: `"hello" + "world"`,
				expected: `
[string:hello (0:0@/test.ca)]
[op:+ (0:8@/test.ca)]
[string:world (0:10@/test.ca)]`,
			},
			{
				source: `"hello" + "\"world\""`,
				expected: `
[string:hello (0:0@/test.ca)]
[op:+ (0:8@/test.ca)]
[string:"world" (0:10@/test.ca)]`,
			},
			{
				source: `"1234"`,
				expected: `
[string:1234 (0:0@/test.ca)]`,
			},
			{
				source: `"abcd\nefgh"`,
				expected: `
[string:abcd
efgh (0:0@/test.ca)]`,
			},
			{
				source: `"ab\\cd\\ef"`,
				expected: `
[string:ab\cd\ef (0:0@/test.ca)]`,
			},
		}

		for _, f := range fixtures {
			tks, err := Tokenize("/test.ca", []byte(f.source))
			if err != nil {
				t.Errorf("unexpected error: %s", err.Error())
			}

			if printTokens(tks) != f.expected {
				t.Errorf("expected %s - got %s", f.expected, printTokens(tks))
			}
		}
	})

	t.Run("strings with errors", func(t *testing.T) {
		fixtures := []string{
			`"hello`,
			`"hello + "world"`,
			`"hello world\`,
			`"hello world\"`,
		}

		for _, source := range fixtures {
			_, err := Tokenize("/test.ca", []byte(source))
			if err == nil {
				t.Error("an error was expected")
			}
		}
	})
}

func printTokens(tks []*tokens.Token) string {
	s := ""
	for _, tk := range tks {
		s += "\n" + tk.String()
	}
	return s
}
