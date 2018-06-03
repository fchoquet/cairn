package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fchoquet/cairn/interpreter"
	"github.com/fchoquet/cairn/parser"
)

func main() {
	i := interpreter.New(&parser.Parser{})

	args := os.Args[1:]
	if len(args) > 0 {
		file := os.Args[1]
		input, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Println("!!! " + err.Error())
			return
		}
		output, err := i.Interpret("stdin", string(input))
		if err != nil {
			fmt.Println("!!! " + err.Error())
			return
		}
		fmt.Println(output)
		return
	}

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
