package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/fchoquet/cairn/interpreter"
)

func main() {
	for {
		fmt.Print("cairn> ")

		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')

		if input == "" {
			continue
		}

		i := &interpreter.Interpreter{
			Parser: &interpreter.Parser{
				Lexer: &interpreter.Lexer{
					Text: input,
				},
			},
		}

		output, err := i.Interpret()
		if err != nil {
			fmt.Println("!!! " + err.Error())
			continue
		}
		fmt.Println("--> " + output)
	}
}
