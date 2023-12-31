package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/procode2/lena/ln_parser/evaluator"
	"github.com/procode2/lena/ln_parser/lexer"
	"github.com/procode2/lena/ln_parser/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			PrintParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program)
		/*if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}*/

		io.WriteString(out, evaluated.Code)
		io.WriteString(out, "\n")
	}
}

func PrintParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, fmt.Sprintf("\t%s\n", msg))
	}
}
