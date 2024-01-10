package main

import (
	"Voca-2/ast"
	"Voca-2/compiler"
	"Voca-2/lexer"
	"Voca-2/lib"
	"Voca-2/parser"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

func astToJson(ast ast.Program) error {
	jsonData, err := json.MarshalIndent(ast, "", "  ")
	if err != nil {
		return err
	}

	soubor, err := os.Create("ast.json")
	if err != nil {
		return err
	}
	defer soubor.Close()

	_, err = soubor.Write(jsonData)
	if err != nil {
		return err
	}

	fmt.Println("AST was written in to the file ast.json:")
	return nil
}

type Program struct {
	input       string
	err         error
	tokens      []lexer.Token
	program     ast.Program
	generateAST bool
	args        []string
	file        string
	output      string
	Arch        string
	OS          string
	loadAST     bool
}

func input(program Program, i int) (Program, int) {
	if i+1 < len(program.args) {
		program.file = program.args[i+1]
		i++
	} else {
		lib.Print("No input file specified")
	}

	return program, i
}
func help() {
	lib.Print("\\ \\        / /    ______        ______         /\\")
	lib.Print(" \\ \\      / /    /  __  \\      /  __  \\       /  \\")
	lib.Print("  \\ \\    / /    /  /   \\ \\    /  /  \\ \\      / /\\ \\")
	lib.Print("   \\ \\  / /    |  |     | |  |  |           / /__\\ \\")
	lib.Print("    \\ \\/ /     |  |     | |  |  |          / /____\\ \\")
	lib.Print("     \\  /       \\ \\____/ /   \\  \\__/  /   / /      \\ \\")
	lib.Print("      \\/         \\______/     \\______/   / /        \\ \\")

	lib.Print("Voca-2 is a simple static typed, compiled programming language")
	lib.Print("Usage: voca [options] [value]")
	lib.Print("Options:")
	lib.Print("  -i, -input [file] - specify input file")
	lib.Print("  -o, -output [file] - specify output file")
	lib.Print("  -a, -ast - want generate AST?")
	lib.Print("  -h, -help - show this help")
	lib.Print("  -arch [arch] - specify architecture")
	lib.Print("  -os [os] - specify operating system")
	lib.Print("  -loadAST - load AST from ast.json")
	os.Exit(0)
}
func output(program Program, i int) (Program, int) {
	if i+1 < len(program.args) {
		program.output = program.args[i+1]
		i++
	} else {
		lib.Print("No output file specified")
	}
	return program, i
}
func arch(program Program, i int) (Program, int) {
	if i+1 < len(program.args) {
		program.Arch = program.args[i+1]
		i++
	} else {
		lib.Print("No architecture specified")
	}
	return program, i
}
func os_var(program Program, i int) (Program, int) {
	if i+1 < len(program.args) {
		program.OS = program.args[i+1]
		i++
	} else {
		lib.Print("No operating system specified")
	}
	return program, i
}
func main() {
	program := Program{args: os.Args[1:], generateAST: false, file: "main.voc", loadAST: false}
	if len(program.args) >= 1 {
		i := 0
		for i < len(program.args) {
			switch program.args[i] {
			case "-i":
				program, i = input(program, i)
			case "-input":
				program, i = input(program, i)
			case "-help":
				help()
			case "-h":
				help()
			case "-ast":
				program.generateAST = true
			case "-a":
				program.generateAST = true
			case "-output":
				program, i = output(program, i)
			case "-o":
				program, i = output(program, i)
			case "-arch":
				program, i = arch(program, i)
			case "-os":
				program, i = os_var(program, i)
			case "-loadAST":
				program.loadAST = true
			}
			i++
		}

	}

	if program.output == "" {
		program.output = strings.Replace(program.file, ".voc", "", -1)
	}
	if program.err != nil {
		panic(program.err)
	}
	if program.loadAST {
		data, err := os.ReadFile("ast.json")
		if err != nil {
			lib.Print("Can't load file: " + err.Error())
			log.Fatal(err)
		}
		err = json.Unmarshal(data, &program.program)
		if err != nil {
			panic(err)
		}
	} else {
		data, err := os.ReadFile(program.file)
		if err != nil {
			lib.Print("Can't load file: " + err.Error())
			log.Fatal(err)
		}
		program.input = string(data)
		program.input = strings.Replace(program.input, "\r", " ", -1)
		program.tokens, err = lexer.Lex(program.input)
		if err != nil {
			panic(err)
		}
		program.program, err = parser.New(program.tokens)
		if err != nil {
			panic(err)
		}
		if program.generateAST {
			err := astToJson(program.program)
			if err != nil {
				panic(err)
			}
		}
		fmt.Println(compiler.GenerateIR(program.program))
	}

}
