package tokenizer

import "errors"

func readString(input string) (string, error) {
	if input == "" {
		return "", errors.New("unexpected end of string")
	}

	if input[0] != '"' {
		return "", errors.New("expected \" at the beginning of string")
	}

	if len(input) < 2 {
		return "", errors.New("a string litteral must be at least 2 chars long (including surrounding quotes)")
	}

	return readStringContents(input[1:])
}

func readStringContents(text string) (string, error) {
	if len(text) == 0 {
		return "", errors.New("could not find end of string litteral")
	}

	// let's skip the initial "
	head := text[0]
	tail := text[1:]

	var left string

	switch head {
	case '"':
		// end of string reached. stop recursion
		return "", nil
	case '\\':
		var err error
		left, err = readEscapeSequence(tail)
		if err != nil {
			return "", err
		}
		tail = tail[len(left):]
	case '\n':
		return "", errors.New("could not find end of string litteral")
	default:
		left = string(head)
	}

	// reads the rest of the string recursively
	right, err := readStringContents(tail)
	if err != nil {
		return "", err
	}

	return (left + right), nil
}

func readEscapeSequence(text string) (string, error) {
	// starting an escape sequence
	if len(text) == 0 {
		return "", errors.New("invalid escape sequence. did you mean \\\\?")
	}

	switch text[0] {
	case '\\':
		return "\\", nil
	case '"':
		return "\"", nil
	case 'n':
		return "\n", nil
	default:
		return "", errors.New("invalid escape sequence")
	}
}
