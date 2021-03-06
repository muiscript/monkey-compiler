package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey-compiler/compiler"
	"monkey-compiler/object"
	"monkey-compiler/vm"

	"monkey-compiler/lexer"
	"monkey-compiler/parser"
)

const prompt = ">> "

// Start starts REPL of monkey
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	constants := make([]object.Object, 0)
	symbolTable := compiler.NewSymbolTable()
	globals := make([]object.Object, vm.GlobalsSize)

	for {
		fmt.Printf(prompt)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		comp := compiler.NewWithState(symbolTable, constants)
		if err := comp.Compile(program); err != nil {
			io.WriteString(out, fmt.Sprintf("error during compilation: %v", err))
		}

		machine := vm.NewWithGlobals(comp.ByteCode(), globals)
		if err := machine.Run(); err != nil {
			io.WriteString(out, fmt.Sprintf("error during execution: %v", err))
		}

		result := machine.LastPopped()
		io.WriteString(out, result.Inspect())
		io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
