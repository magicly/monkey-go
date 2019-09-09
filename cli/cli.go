package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/magicly/monkey-go/evaluator"
	"github.com/magicly/monkey-go/lexer"
	"github.com/magicly/monkey-go/object"
	"github.com/magicly/monkey-go/parser"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}
func main() {
	filename := os.Args[1]
	dat, err := ioutil.ReadFile(filename)
	script := string(dat)
	check(err)

	l := lexer.New(script)
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printParseError(os.Stdout, p.Errors())
	}

	env := object.NewEnv()
	result := evaluator.Eval(program, env)
	if result != nil {
		fmt.Println(result.Inspect())
	}
}

func printParseError(out io.Writer, errors []string) {
	for _, err := range errors {
		io.WriteString(out, "\t"+err+"\n")
	}
}
