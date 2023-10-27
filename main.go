package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/procode2/lena/ln_parser/evaluator"
	"github.com/procode2/lena/ln_parser/lexer"
	"github.com/procode2/lena/ln_parser/parser"
	"github.com/procode2/lena/ln_parser/repl"
)

func parseMetaSection(metaString string) {
	// Split the meta section into lines
	lines := strings.Split(metaString, "\n")

	// Create a map to store the key-value pairs
	metaInfo := make(map[string]string)

	for _, line := range lines {
		// Split each line by ":" to separate key and value
		parts := strings.SplitN(line, ":", 2)

		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			// Store the key-value pair in the map
			metaInfo[key] = value
		}
	}
	fmt.Println(metaInfo)
}

func parseDataSection(dataString string) {

	l := lexer.New(dataString)
	p := parser.New(l)

	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		repl.PrintParserErrors(os.Stdout, p.Errors())
	}

	evaluated := evaluator.Eval(program)
	fmt.Println(evaluated.Code)
}

func main() {
	input := `title: HelloBrother whe are u doing
		description: demo
		how: you`

	input2 := `let todo = ["this is task 1", "this is task 2"];
	puts(todo);
	let lobby = fn(something) {puts(something);};
	lobby(todo[1]);`
	parseMetaSection(input)
	parseDataSection(input2)
}
