package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/fchoquet/cairn/cairn/interpreter"
	"github.com/fchoquet/cairn/cairn/parser"
)

func main() {
	i := interpreter.New(&parser.Parser{})

	for {
		fmt.Print("cairn> ")

		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')

		if input == "" {
			continue
		}
		output, err := i.Interpret("stdin", input)
		if err != nil {
			fmt.Println("!!! " + err.Error())
			continue
		}
		fmt.Println("--> " + output)
	}
}
