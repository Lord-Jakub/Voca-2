package main

import (
	"Voca-2/ast"
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

func main() {
	data, err := os.ReadFile("code.voc")
	if err != nil {
		lib.Print("Can't load file: " + err.Error())
		log.Fatal(err)
	}

	input := string(data)
	input = strings.Replace(input, "\r", " ", -1)
	tokens, err := lexer.Lex(input)
	if err != nil {
		panic(err)
	}
	program, err := parser.New(tokens)
	if err != nil {
		panic(err)
	}
	err = astToJson(program)
	if err != nil {
		panic(err)
	}

	for _, token := range tokens {
		fmt.Println(token)
	}
}
