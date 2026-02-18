package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/walonCode/code-lang/evaluator"
	"github.com/walonCode/code-lang/lexer"
	"github.com/walonCode/code-lang/object"
	"github.com/walonCode/code-lang/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer){
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	
	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		
		if !scanned {
			return
		}
		
		line := scanner.Text()
		
		l := lexer.New(line)
		p := parser.New(l)
		
		programe := p.ParsePrograme()
		if len(p.Errors()) != 0 {
			printParserError(out, p.Errors())
			continue
		}
		
		evaluated := evaluator.Eval(programe, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

const CODE_LANG = `
______          __                 __                         
/ ____/  ____   / /  ___           / /   ____ _   ____    ____ _
/ /      / __ \ / /  / _ \  ______ / /   / __ /  / __ \  / __ /
/ /___   / /_/ // /  /  __/ /_____// /___/ /_/ /  / / / / / /_/ / 
\____/   \____//_/   \___/        /_____/\__,_/  /_/ /_/  \__, /  
                                                      /____/   

`                                                    

func printParserError(out io.Writer, errors []string){
	io.WriteString(out, CODE_LANG)
	io.WriteString(out, "Whoops! We can in to some Code_lang business!\n")
	io.WriteString(out," parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}


func Execute(source string, out io.Writer){
	l := lexer.New(source)
	p := parser.New(l)
	program := p.ParsePrograme()
	
	if len(p.Errors()) != 0 {
		printParserError(out, p.Errors())
		return
	}
	
	evaluated := evaluator.Eval(program, object.NewEnvironment())
	if evaluated != nil {
		io.WriteString(out, evaluated.Inspect())
        io.WriteString(out, "\n")
	}
}