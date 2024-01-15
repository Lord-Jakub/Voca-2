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

func input(program compiler.Program, i int) (compiler.Program, int) {
	if i+1 < len(program.Args) {
		program.File = program.Args[i+1]
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
	lib.Print("    \\ \\/ /     |  |     | |  |  |             / /____\\ \\")
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
func output(program compiler.Program, i int) (compiler.Program, int) {
	if i+1 < len(program.Args) {
		program.Output = program.Args[i+1]
		i++
	} else {
		lib.Print("No output file specified")
	}
	return program, i
}
func arch(program compiler.Program, i int) (compiler.Program, int) {
	if i+1 < len(program.Args) {
		program.Arch = program.Args[i+1]
		i++
	} else {
		lib.Print("No architecture specified")
	}
	return program, i
}
func os_var(program compiler.Program, i int) (compiler.Program, int) {
	if i+1 < len(program.Args) {
		program.OS = program.Args[i+1]
		i++
	} else {
		lib.Print("No operating system specified")
	}
	return program, i
}
func main() {
	program := compiler.Program{Args: os.Args, GenerateAST: false, File: "main.voc", LoadAST: false, Ir: false}
	if len(program.Args) >= 1 {
		i := 0
		for i < len(program.Args) {
			switch program.Args[i] {
			case "-i":
				program, i = input(program, i)
			case "-input":
				program, i = input(program, i)
			case "-help":
				help()
			case "-h":
				help()
			case "-ast":
				program.GenerateAST = true
			case "-a":
				program.GenerateAST = true
			case "-output":
				program, i = output(program, i)
			case "-o":
				program, i = output(program, i)
			case "-arch":
				program, i = arch(program, i)
			case "-os":
				program, i = os_var(program, i)
			case "-loadAST":
				program.LoadAST = true
			case "-ir":
				program.Ir = true

			}
			i++
		}

	}

	if program.Output == "" {
		program.Output = strings.Replace(program.File, ".voc", "", -1)
	}

	if program.LoadAST {
		data, err := os.ReadFile("ast.json")
		if err != nil {
			fmt.Print("Can't load file: " + err.Error())
			log.Fatal(err)
		}
		err = json.Unmarshal(data, &program.Program)
		program.Errs = append(program.Errs, err)
	} else {
		data, err := os.ReadFile(program.File)
		if err != nil {
			fmt.Print("Can't load file: " + err.Error())
			log.Fatal(err)
		}
		program.Input = string(data)
		program.Input = strings.Replace(program.Input, "\r", " ", -1)
		program.Tokens, err = lexer.Lex(program.Input)
		program.Errs = append(program.Errs, err)
		program.Program, err = parser.New(program.Tokens)
		program.Errs = append(program.Errs, err)
		if program.GenerateAST {
			err := astToJson(program.Program)
			program.Errs = append(program.Errs, err)
		}

		errs := compiler.New(program)
		for i := 0; i < len(errs); i++ {
			if errs[i] != nil {
				fmt.Println(errs[i])
			}
		}

	}

}
