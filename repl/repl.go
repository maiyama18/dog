package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/maiyama18/dog/evaluate"

	"github.com/maiyama18/dog/lex"
	"github.com/maiyama18/dog/parse"
)

const PROMPT = "~> "

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print(PROMPT)
	for scanner.Scan() {
		line := scanner.Text()
		lexer := lex.NewLexer(line)
		parser := parse.NewParser(lexer)

		program := parser.ParseProgram()
		if len(parser.Errors()) > 0 {
			for _, err := range parser.Errors() {
				fmt.Println(err.Error())
			}
		} else {
			fmt.Println(evaluate.Eval(program).Inspect())
		}

		fmt.Print(PROMPT)
	}
}
