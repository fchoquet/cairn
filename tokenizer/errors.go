package tokenizer

import "github.com/fchoquet/cairn/tokens"

// Error implements error and contains contextual information
type Error struct {
	Message  string
	Code     []byte
	Position tokens.Position
	Hint     string
}

func (e Error) Error() string {
	return e.Message
}
