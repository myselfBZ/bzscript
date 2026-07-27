package main

import (
	"fmt"
	"log"
	"os"

	"github.com/myselfBZ/bzscript/eval"
	"github.com/myselfBZ/bzscript/lexer"
	"github.com/myselfBZ/bzscript/object"
	"github.com/myselfBZ/bzscript/parser"
)


func open(path string) string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("couldn't find the file")
	}
	buff := make([]byte, 1024)
	size, err := file.Read(buff)
	if err != nil {
		log.Fatal("error reading from a file, ", err)
	}
	return string(buff[:size])
}


func main() {
	if len(os.Args) != 2 {
		fmt.Println("No file specified for the bzscript to execute")
		os.Exit(1)
	}
	input := open(os.Args[1])
	lexer := lexer.New(input)
	p := parser.New(lexer)
	program := p.ParseProgram()
	env := object.NewEnvironment()
	result := eval.Eval(program, env)
	if result.Type() == object.ERROR {
		fmt.Println(result.Inspect())
		os.Exit(1)
	}
}
