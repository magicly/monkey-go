package repl

import (
	"bufio"
	"io"

	"github.com/magicly/monkey-go/lexer"
	"github.com/magicly/monkey-go/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParseError(out, p.Errors())
		}

		io.WriteString(out, program.String())
		io.WriteString(out, "\n")
	}

}
func printParseError(out io.Writer, errors []string) {
	for _, err := range errors {
		io.WriteString(out, "\t"+err+"\n")
	}
}
