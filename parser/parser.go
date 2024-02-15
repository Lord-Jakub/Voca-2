package parser

import (
	"Voca-2/ast"
	"Voca-2/lexer"
	"errors"
	"fmt"
	"strconv"
)

var parseInFunc bool = false

func New(tokens []lexer.Token) (ast.Program, error) {
	var err error = nil
	Functions, err = MapFunctions(tokens)
	steatments, err := Parse(tokens, make(map[string]string))

	return ast.Program{Statements: steatments, Externals: ParseExternals(tokens)}, err
}

func ParseExternals(tokens []lexer.Token) []string {
	var externals []string
	externals = append(externals, "vsl")
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Type == lexer.Keyword && tokens[i].Value == "external" {
			externals = append(externals, tokens[i+1].Value.(string))
		}
	}

	return externals
}

var Functions = make(map[string]string)

func MapFunctions(tokens []lexer.Token) (map[string]string, error) {
	var functions = make(map[string]string)
	functions["print"] = "void"
	functions["append"] = "string"
	functions["strlen"] = "int"
	functions["IntToString"] = "string"
	functions["Read"] = "string"
	functions["StringToFloat"] = "float"
	functions["StringToInt"] = "int"
	functions["FloatToString"] = "string"
	functions["FloatToInt"] = "int"
	functions["delay"] = "void"
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Type == lexer.Keyword && tokens[i].Value == "func" {
			name := tokens[i+1].Value
			for tokens[i].Type != lexer.OpenBrace {
				i++
			}
			if tokens[i-1].Value == ")" {
				functions[name.(string)] = "void"

			} else {
				functions[name.(string)] = tokens[i-1].Value.(string)
			}
		} else if tokens[i].Type == lexer.Keyword && tokens[i].Value == "extern_func" {
			name := tokens[i+1].Value
			for tokens[i].Type != lexer.CloseParen {
				i++
			}
			i++
			if tokens[i].Value == "string" || tokens[i].Value == "int" || tokens[i].Value == "float" || tokens[i].Value == "bool" {
				functions[name.(string)] = tokens[i].Value.(string)
			} else {
				functions[name.(string)] = "void"
			}
		}
	}
	return functions, nil
}
func Parse(tokens []lexer.Token, Variables map[string]string) ([]ast.Statement, error) {
	var program ast.Program
	if !parseInFunc {
		program.Statements = append(program.Statements, ast.Statement{Node: ast.ExternFuncDeclaration{Name: lexer.Token{Type: lexer.Keyword, Value: "print"}, Arguments: []any{ast.VariableDeclaration{Type: lexer.Token{Type: lexer.Keyword, Value: "string"}, Name: lexer.Token{Type: lexer.Identifier, Value: "s"}}}, Type: lexer.Token{Type: lexer.Keyword, Value: "void"}}})
		program.Statements = append(program.Statements, ast.Statement{Node: ast.ExternFuncDeclaration{Name: lexer.Token{Type: lexer.Keyword, Value: "append"}, Arguments: []any{ast.VariableDeclaration{Type: lexer.Token{Type: lexer.Keyword, Value: "string"}, Name: lexer.Token{Type: lexer.Identifier, Value: "s1"}}, ast.VariableDeclaration{Type: lexer.Token{Type: lexer.Keyword, Value: "string"}, Name: lexer.Token{Type: lexer.Identifier, Value: "s2"}}}, Type: lexer.Token{Type: lexer.Keyword, Value: "string"}}})
		program.Statements = append(program.Statements, ast.Statement{Node: ast.ExternFuncDeclaration{Name: lexer.Token{Type: lexer.Keyword, Value: "strlen"}, Arguments: []any{ast.VariableDeclaration{Type: lexer.Token{Type: lexer.Keyword, Value: "string"}, Name: lexer.Token{Type: lexer.Identifier, Value: "s"}}}, Type: lexer.Token{Type: lexer.Keyword, Value: "int"}}})
		program.Statements = append(program.Statements, ast.Statement{Node: ast.ExternFuncDeclaration{Name: lexer.Token{Type: lexer.Keyword, Value: "IntToString"}, Arguments: []any{ast.VariableDeclaration{Type: lexer.Token{Type: lexer.Keyword, Value: "int"}, Name: lexer.Token{Type: lexer.Identifier, Value: "num"}}}, Type: lexer.Token{Type: lexer.Keyword, Value: "string"}}})
		program.Statements = append(program.Statements, ast.Statement{Node: ast.ExternFuncDeclaration{Name: lexer.Token{Type: lexer.Keyword, Value: "Read"}, Arguments: []any{}, Type: lexer.Token{Type: lexer.Keyword, Value: "string"}}})
		program.Statements = append(program.Statements, ast.Statement{Node: ast.ExternFuncDeclaration{Name: lexer.Token{Type: lexer.Keyword, Value: "StringToFloat"}, Arguments: []any{ast.VariableDeclaration{Type: lexer.Token{Type: lexer.Keyword, Value: "string"}, Name: lexer.Token{Type: lexer.Identifier, Value: "s"}}}, Type: lexer.Token{Type: lexer.Keyword, Value: "float"}}})
		program.Statements = append(program.Statements, ast.Statement{Node: ast.ExternFuncDeclaration{Name: lexer.Token{Type: lexer.Keyword, Value: "StringToInt"}, Arguments: []any{ast.VariableDeclaration{Type: lexer.Token{Type: lexer.Keyword, Value: "string"}, Name: lexer.Token{Type: lexer.Identifier, Value: "s"}}}, Type: lexer.Token{Type: lexer.Keyword, Value: "int"}}})
		program.Statements = append(program.Statements, ast.Statement{Node: ast.ExternFuncDeclaration{Name: lexer.Token{Type: lexer.Keyword, Value: "FloatToString"}, Arguments: []any{ast.VariableDeclaration{Type: lexer.Token{Type: lexer.Keyword, Value: "float"}, Name: lexer.Token{Type: lexer.Identifier, Value: "f"}}}, Type: lexer.Token{Type: lexer.Keyword, Value: "string"}}})
		program.Statements = append(program.Statements, ast.Statement{Node: ast.ExternFuncDeclaration{Name: lexer.Token{Type: lexer.Keyword, Value: "FloatToInt"}, Arguments: []any{ast.VariableDeclaration{Type: lexer.Token{Type: lexer.Keyword, Value: "float"}, Name: lexer.Token{Type: lexer.Identifier, Value: "f"}}}, Type: lexer.Token{Type: lexer.Keyword, Value: "int"}}})
		program.Statements = append(program.Statements, ast.Statement{Node: ast.ExternFuncDeclaration{Name: lexer.Token{Type: lexer.Keyword, Value: "delay"}, Arguments: []any{ast.VariableDeclaration{Type: lexer.Token{Type: lexer.Keyword, Value: "int"}, Name: lexer.Token{Type: lexer.Identifier, Value: "ms"}}}, Type: lexer.Token{Type: lexer.Keyword, Value: "void"}}})
	}
	parseInFunc = true
	var err error = nil
	for i := 0; i < len(tokens); i++ {
		var statement ast.Statement
		switch tokens[i].Type {
		case lexer.Keyword:
			switch tokens[i].Value {
			case "int", "float":

				var variableDeclaration ast.VariableDeclaration
				variableDeclaration.Type = tokens[i]
				isarray := false
				if tokens[i+1].Type == lexer.OpenBracket {
					isarray = true
				}
				for i < len(tokens) && tokens[i].Type != lexer.Identifier {
					i++
				}
				variableDeclaration.Name = tokens[i]
				if !isarray {
					if tokens[i+1].Type == lexer.Equal {
						var expression []lexer.Token
						i += 2
						for tokens[i].Type == lexer.OpenParen {
							expression = append(expression, tokens[i])
							i++
						}
						for lexer.IsOperator(tokens[i+1]) || tokens[i+1].Type == lexer.OpenParen {
							expression = append(expression, tokens[i])
							expression = append(expression, tokens[i+1])
							i += 2
							for tokens[i].Type == lexer.OpenParen {
								expression = append(expression, tokens[i])
								i++
							}
							for tokens[i+1].Type == lexer.CloseParen {
								expression = append(expression, tokens[i])
								i++
							}
							if i+1 >= len(tokens) {
								break
							}
						}
						expression = append(expression, tokens[i])
						i++
						for tokens[i].Type == lexer.CloseParen {
							expression = append(expression, tokens[i])
							i++
						}
						value, IsNum, err := ParseExpression(expression, Variables, false)
						if err != nil {
							return program.Statements, err
						}
						variableDeclaration.Value = value

						if !IsNum {
							err = errors.New(fmt.Sprintf("Expected number on line: %d at position %d, not '%s' in file '%s'", expression[0].Line, expression[0].LinePos, expression[0].Value, expression[0].File))
							return program.Statements, err
						}
					}
					statement.Node = variableDeclaration
					if variableDeclaration.Type.Value == "int" {
						Variables[variableDeclaration.Name.Value.(string)] = "int"
					} else {
						Variables[variableDeclaration.Name.Value.(string)] = "float"
					}
				} else {
					var arrayDeclaration ast.ArrayDeclaration
					array := make([]lexer.Token, 0)
					arrayDeclaration.Name = tokens[i]
					pos := i
					for i >= 0 && tokens[i].Type != lexer.OpenBracket {
						i--
					}
					arrayDeclaration.Length, _, _ = ParseExpression(tokens[i+1:pos-1], Variables, false)
					i--
					var arraytype []lexer.Token
					for i >= 0 && tokens[i].Type != lexer.NewLine {
						arraytype = append(arraytype, tokens[i])
						i--
					}
					arrayDeclaration.Type, _ = ParseArray(arraytype, Variables)
					i = pos

					if tokens[i+1].Type == lexer.Equal {
						i += 2
						if tokens[i].Type == lexer.OpenBracket {

							array = append(array, tokens[i])
							n := 1
							i++
							for n > 0 {
								if tokens[i].Type == lexer.OpenBracket {
									n++
								} else if tokens[i].Type == lexer.CloseBracket {
									n--
								}
								array = append(array, tokens[i])
								i++
							}
							arr, _, _ := ParseExpression(array, Variables, false)
							arrayDeclaration.Value = arr
						}
					}

					statement.Node = arrayDeclaration
					Variables[arrayDeclaration.Name.Value.(string)] = "array"

				}
			case "string":
				var variableDeclaration ast.VariableDeclaration
				variableDeclaration.Type = tokens[i]
				variableDeclaration.Name = tokens[i+1]
				if tokens[i+2].Type == lexer.Equal {
					var expression []lexer.Token
					i += 3
					for tokens[i].Type == lexer.OpenParen {
						expression = append(expression, tokens[i])
						i++
					}
					for lexer.IsOperator(tokens[i+1]) {
						expression = append(expression, tokens[i])
						expression = append(expression, tokens[i+1])
						i += 2
						for tokens[i].Type == lexer.OpenParen {
							expression = append(expression, tokens[i])
							i++
						}
						for tokens[i+1].Type == lexer.CloseParen {
							expression = append(expression, tokens[i])
							i++
						}
						if i+1 >= len(tokens) {
							break
						}
					}
					expression = append(expression, tokens[i])
					for tokens[i].Type == lexer.CloseParen {
						expression = append(expression, tokens[i])
						i++
					}
					value, IsNum, err := ParseExpression(expression, Variables, false)
					if err != nil {
						return program.Statements, err
					}
					variableDeclaration.Value = value
					if IsNum {
						err = errors.New(fmt.Sprintf("Expected number on line: %d at position %d, not '%s' in file '%s'", expression[0].Line, expression[0].LinePos, expression[0].Value, expression[0].File))
						return program.Statements, err
					}
				}
				statement.Node = variableDeclaration
				Variables[variableDeclaration.Name.Value.(string)] = "string"
			case "bool":
				var variableDeclaration ast.VariableDeclaration
				variableDeclaration.Type = tokens[i]
				variableDeclaration.Name = tokens[i+1]
				if tokens[i+2].Type == lexer.Equal {
					i += 3
					variableDeclaration.Value = tokens[i]
				}
				statement.Node = variableDeclaration
				Variables[variableDeclaration.Name.Value.(string)] = "bool"
			case "func":
				var funcDeclaration ast.FuncDeclaration
				funcDeclaration.Name = tokens[i+1]
				funcDeclaration.Arguments = make([]any, 0)
				i += 3
				for tokens[i].Type != lexer.CloseParen {

					if tokens[i].Type == lexer.Keyword {
						var arg ast.VariableDeclaration
						if tokens[i].Value == "int" || tokens[i].Value == "string" || tokens[i].Value == "float" || tokens[i].Value == "bool" {
							i++
							arg = ast.VariableDeclaration{Type: tokens[i-1], Name: tokens[i]}
						} else {
							err = errors.New(fmt.Sprintf("Expected type on line: %d at position %d, not '%s'in file '%s'", tokens[i].Line, tokens[i].LinePos, tokens[i].Value, tokens[i].File))
							return program.Statements, err
						}
						funcDeclaration.Arguments = append(funcDeclaration.Arguments, arg)
					} else if tokens[i].Type == lexer.Comma {

					} else {

						err = errors.New(fmt.Sprintf("Expected argument on line: %d at position %d, not '%s'in file '%s'", tokens[i].Line, tokens[i].LinePos, tokens[i].Value, tokens[i].File))
						return program.Statements, err
					}

					if err != nil {
						return program.Statements, err
					}
					i++
				}
				i++
				if tokens[i].Value == "string" || tokens[i].Value == "int" || tokens[i].Value == "float" || tokens[i].Value == "bool" {
					funcDeclaration.Type = tokens[i]
					i++
				} else if tokens[i].Type == lexer.OpenBrace {
					funcDeclaration.Type = lexer.Token{Type: lexer.Keyword, Value: "void"}
				} else {
					err = errors.New(fmt.Sprintf("Expected type or { on line: %d at position %d, not '%s'in file '%s'", tokens[i].Line, tokens[i].LinePos, tokens[i].Value, tokens[i].File))
					return program.Statements, err
				}

				funcDeclaration.Body = make([]ast.Statement, 0)
				bodyTokens := make([]lexer.Token, 0)
				i++
				n := 1
				for n > 0 {
					if tokens[i].Type == lexer.OpenBrace {
						n++
					} else if tokens[i].Type == lexer.CloseBrace {
						n--
					}
					bodyTokens = append(bodyTokens, tokens[i])
					i++
				}
				variables := make(map[string]string)
				i2 := 0
				for i2 < len(funcDeclaration.Arguments) {
					variables[funcDeclaration.Arguments[i2].(ast.VariableDeclaration).Name.Value.(string)] = funcDeclaration.Arguments[i2].(ast.VariableDeclaration).Type.Value.(string)
					i2++
				}
				funcDeclaration.Body, err = Parse(bodyTokens, variables)
				if err != nil {
					return program.Statements, err
				}
				statement.Node = funcDeclaration
			case "extern_func":
				var funcDeclaration ast.ExternFuncDeclaration
				funcDeclaration.Name = tokens[i+1]
				funcDeclaration.Arguments = make([]any, 0)
				i += 3
				for tokens[i].Type != lexer.CloseParen {

					if tokens[i].Type == lexer.Keyword {
						var arg ast.VariableDeclaration
						if tokens[i].Value == "int" || tokens[i].Value == "string" || tokens[i].Value == "float" || tokens[i].Value == "bool" {
							i++
							arg = ast.VariableDeclaration{Type: tokens[i-1], Name: tokens[i]}
						} else {
							err = errors.New(fmt.Sprintf("Expected type on line: %d at position %d, not '%s' in file '%s'", tokens[i].Line, tokens[i].LinePos, tokens[i].Value, tokens[i].File))
							return program.Statements, err
						}
						funcDeclaration.Arguments = append(funcDeclaration.Arguments, arg)
					} else if tokens[i].Type == lexer.Comma {

					} else {

						err = errors.New(fmt.Sprintf("Expected argument on line: %d at position %d, not '%s'in file '%s'", tokens[i].Line, tokens[i].LinePos, tokens[i].Value, tokens[i].File))
						return program.Statements, err
					}

					if err != nil {
						return program.Statements, err
					}
					i++
				}
				i++
				if tokens[i].Value == "string" || tokens[i].Value == "int" || tokens[i].Value == "float" || tokens[i].Value == "bool" {
					funcDeclaration.Type = tokens[i]
					i++
				} else if tokens[i].Type == lexer.NewLine {
					funcDeclaration.Type = lexer.Token{Type: lexer.Keyword, Value: "void"}
				} else {
					err = errors.New(fmt.Sprintf("Expected type or { on line: %d at position %d, not '%s'in file '%s'", tokens[i].Line, tokens[i].LinePos, tokens[i].Value, tokens[i].File))
					return program.Statements, err
				}

				statement.Node = funcDeclaration
			case "return":
				var returnStatement ast.ReturnStatement
				var expression []lexer.Token
				i += 1
				for tokens[i].Type != lexer.NewLine {
					expression = append(expression, tokens[i])
					i++
				}
				value, _, err := ParseExpression(expression, Variables, false)
				if err != nil {
					return program.Statements, err
				}
				returnStatement.Value = value
				statement.Node = returnStatement
			case "if":
				var ifStatement ast.IfStatement
				var expression []lexer.Token
				var bodyTokens []lexer.Token
				var condition any
				invert := false
				i += 1
				if tokens[i].Type == lexer.Not {
					invert = true
					i++
				}
				i2 := i
				isexpression := false

				for tokens[i].Type != lexer.OpenBrace {
					expression = append(expression, tokens[i])
					if tokens[i].Type == lexer.DoubleEqual || tokens[i].Type == lexer.NotEqual || tokens[i].Type == lexer.LessThan || tokens[i].Type == lexer.MoreThan {
						isexpression = true
					}
					i++
				}
				i2++
				if isexpression {
					condition, err = ParseBool(expression, Variables)
				} else if tokens[i2].Type == lexer.Identifier {
					if Variables[tokens[i2].Value.(string)] == "bool" {
						condition = tokens[i2]
					} else if Functions[tokens[i2].Value.(string)] == "bool" {
						var funcCall ast.FuncCall
						funcCall.Name = tokens[i2]
						funcCall.Arguments = make([]any, 0)
						i2 += 2
						for tokens[i2].Type != lexer.CloseParen && tokens[i2].Type != lexer.NewLine {
							var expression []lexer.Token

							expression, i2 = getArgs(tokens, i2, Variables, false)

							if tokens[i2].Type == lexer.Comma {
								i2++
							}
							arg, _, err := ParseExpression(expression, Variables, false)
							if err != nil {
								return program.Statements, err
							}

							funcCall.Arguments = append(funcCall.Arguments, arg)
							if err != nil {
								return program.Statements, err
							}
						}
						condition = funcCall

					} else {
						err = errors.New(fmt.Sprintf("Expected bool on line: %d at position %d, not '%s' in file: '%s'", tokens[i2].Line, tokens[i2].LinePos, tokens[i2].Value, tokens[i2].File))
					}
				} else if tokens[i2].Value == "true" || tokens[i2].Value == "false" {
					condition = tokens[i2]
				}
				if err != nil {
					return program.Statements, err
				}
				booleon := ast.BoolStatement{Bool: condition, Invert: invert}
				ifStatement.Condition = booleon
				i++
				n := 1
				bodyTokens = make([]lexer.Token, 0)
				ifStatement.Consequence = make([]ast.Statement, 0)
				ifStatement.Alternative = make([]ast.Statement, 0)
				for n > 0 {
					if tokens[i].Type == lexer.OpenBrace {
						n++
					} else if tokens[i].Type == lexer.CloseBrace {
						n--
					}
					bodyTokens = append(bodyTokens, tokens[i])
					i++
				}
				ifStatement.Consequence, err = Parse(bodyTokens, Variables)
				if err != nil {
					return program.Statements, err
				}
				bodyTokens = make([]lexer.Token, 0)
				if tokens[i].Value == "else" {
					i += 2
					n := 1
					for n > 0 {
						if tokens[i].Type == lexer.OpenBrace {
							n++
						} else if tokens[i].Type == lexer.CloseBrace {
							n--
						}
						bodyTokens = append(bodyTokens, tokens[i])
						i++
					}
					ifStatement.Alternative, err = Parse(bodyTokens, Variables)
				}

				statement.Node = ifStatement
			case "while":
				var whileStatement ast.WhileStatement
				var expression []lexer.Token
				var bodyTokens []lexer.Token
				var condition any
				invert := false
				i += 1
				if tokens[i].Type == lexer.Not {
					invert = true
					i++
				}
				i2 := i
				isexpression := false

				for tokens[i].Type != lexer.OpenBrace {
					expression = append(expression, tokens[i])
					if tokens[i].Type == lexer.DoubleEqual || tokens[i].Type == lexer.NotEqual || tokens[i].Type == lexer.LessThan || tokens[i].Type == lexer.MoreThan {
						isexpression = true
					}
					i++
				}
				i2++
				if isexpression {
					condition, err = ParseBool(expression, Variables)
				} else if tokens[i2].Type == lexer.Identifier {
					if Variables[tokens[i2].Value.(string)] == "bool" {
						condition = tokens[i2]
					} else if Functions[tokens[i2].Value.(string)] == "bool" {
						var funcCall ast.FuncCall
						funcCall.Name = tokens[i2]
						funcCall.Arguments = make([]any, 0)
						i2 += 2
						for tokens[i2].Type != lexer.CloseParen && tokens[i2].Type != lexer.NewLine {
							var expression []lexer.Token

							expression, i2 = getArgs(tokens, i2, Variables, false)

							if tokens[i2].Type == lexer.Comma {
								i2++
							}
							arg, _, err := ParseExpression(expression, Variables, false)
							if err != nil {
								return program.Statements, err
							}

							funcCall.Arguments = append(funcCall.Arguments, arg)
							if err != nil {
								return program.Statements, err
							}
						}
						condition = funcCall

					} else {
						err = errors.New(fmt.Sprintf("Expected bool on line: %d at position %d, not '%s' in file: '%s'", tokens[i2].Line, tokens[i2].LinePos, tokens[i2].Value, tokens[i2].File))
					}
				} else if tokens[i2].Value == "true" || tokens[i2].Value == "false" {
					condition = tokens[i2]
				}
				if err != nil {
					return program.Statements, err
				}
				booleon := ast.BoolStatement{Bool: condition, Invert: invert}
				whileStatement.Condition = booleon
				i++
				n := 1
				bodyTokens = make([]lexer.Token, 0)
				whileStatement.Consequence = make([]ast.Statement, 0)

				for n > 0 {
					if tokens[i].Type == lexer.OpenBrace {
						n++
					} else if tokens[i].Type == lexer.CloseBrace {
						n--
					}
					bodyTokens = append(bodyTokens, tokens[i])
					i++
				}
				whileStatement.Consequence, err = Parse(bodyTokens, Variables)
				if err != nil {
					return program.Statements, err
				}
				statement.Node = whileStatement

			}
		case lexer.Identifier:
			if tokens[i+1].Type == lexer.OpenParen {
				//if identifier is in functions map
				if _, exist := Functions[tokens[i].Value.(string)]; exist {
					var funcCall ast.FuncCall
					funcCall.Name = tokens[i]
					funcCall.Arguments = make([]any, 0)
					i += 2
					for tokens[i].Type != lexer.CloseParen && tokens[i].Type != lexer.NewLine {
						var expression []lexer.Token

						expression, i = getArgs(tokens, i, Variables, false)

						if tokens[i].Type == lexer.Comma {
							i++
						}
						arg, _, err := ParseExpression(expression, Variables, false)
						if err != nil {
							return program.Statements, err
						}

						funcCall.Arguments = append(funcCall.Arguments, arg)
						if err != nil {
							return program.Statements, err
						}
						if tokens[i].Type == lexer.CloseParen && tokens[i+1].Type == lexer.Comma {
							i += 2
						}
					}
					statement.Node = funcCall
				} else {
					err = errors.New(fmt.Sprintf("Function '%s' not declered on line: %d at position %d in file: '%s'", tokens[i].Value, tokens[i].Line, tokens[i].LinePos, tokens[i].File))
					return program.Statements, err
				}
			} else {
				if Variables[tokens[i].Value.(string)] != "" {
					if tokens[i+1].Type == lexer.OpenBracket {
						var arrayAssignment ast.ArrayAssignment
						arrayAssignment.Name = tokens[i]
						i++
						for tokens[i].Type != lexer.Equal {
							if tokens[i].Type == lexer.OpenBracket {
								indexExpression := make([]lexer.Token, 0)
								i++
								for tokens[i].Type != lexer.CloseBracket {
									indexExpression = append(indexExpression, tokens[i])
									i++
								}
								index, _, err := ParseExpression(indexExpression, Variables, false)
								if err != nil {
									return program.Statements, err
								}
								arrayAssignment.Indexes = append(arrayAssignment.Indexes, index)
							}

							i++
						}
						i++
						var expression []lexer.Token
						for (lexer.IsOperator(tokens[i+1]) || tokens[i+1].Type == lexer.OpenParen) && tokens[i].Type != lexer.NewLine {
							if tokens[i+1].Type == lexer.OpenParen {
								expression = append(expression, tokens[i])
								expression = append(expression, tokens[i+1])
								i += 2
								for tokens[i].Type != lexer.CloseParen {
									expression = append(expression, tokens[i])
									i++
								}
								if lexer.IsOperator(tokens[i+1]) {
									expression = append(expression, tokens[i])
									expression = append(expression, tokens[i+1])
									i += 2
								}

							} else {
								expression = append(expression, tokens[i])
								expression = append(expression, tokens[i+1])
								i += 2
							}

							if i+1 >= len(tokens) {
								break
							}
						}
						if tokens[i+1].Type == lexer.NewLine {
							expression = append(expression, tokens[i])
						}
						value, _, err := ParseExpression(expression, Variables, false)
						if err != nil {
							return program.Statements, err
						}
						arrayAssignment.Value = value
						statement.Node = arrayAssignment

					} else {
						var variableAssignment ast.VariableAssignment
						variableAssignment.Name = tokens[i]
						if tokens[i+1].Type == lexer.Equal {
							var expression []lexer.Token
							i += 2
							for (lexer.IsOperator(tokens[i+1]) || tokens[i+1].Type == lexer.OpenParen) && tokens[i].Type != lexer.NewLine {
								if tokens[i+1].Type == lexer.OpenParen {
									expression = append(expression, tokens[i])
									expression = append(expression, tokens[i+1])
									i += 2
									for tokens[i].Type != lexer.CloseParen {
										expression = append(expression, tokens[i])
										i++
									}
									if lexer.IsOperator(tokens[i+1]) {
										expression = append(expression, tokens[i])
										expression = append(expression, tokens[i+1])
										i += 2
									}

								} else {
									expression = append(expression, tokens[i])
									expression = append(expression, tokens[i+1])
									i += 2
								}

								if i+1 >= len(tokens) {
									break
								}
							}
							if tokens[i+1].Type == lexer.NewLine {
								expression = append(expression, tokens[i])
							}
							value, IsNum, err := ParseExpression(expression, Variables, false)
							if err != nil {
								return program.Statements, err
							}
							variableAssignment.Value = value
							if IsNum {
								if Variables[variableAssignment.Name.Value.(string)] != "int" {
									err = errors.New(fmt.Sprintf("Expected string on line: %d at position %d, not '%s' in file '%s'", expression[0].Line, expression[0].LinePos, expression[0].Value, expression[0].File))
									return program.Statements, err
								}
							} else {
								if Variables[variableAssignment.Name.Value.(string)] != "string" {
									err = errors.New(fmt.Sprintf("Expected number on line: %d at position %d, not '%s' in file '%s'", expression[0].Line, expression[0].LinePos, expression[0].Value, expression[0].File))
									return program.Statements, err
								}
							}
						} else if tokens[i+1].Type == lexer.Plus && tokens[i+2].Type == lexer.Plus {
							variableAssignment.Value = ast.ExpressionStatement{Left: tokens[i], Operator: tokens[i+1], Right: lexer.Token{Type: lexer.Int, Value: 1}}
						} else if tokens[i+1].Type == lexer.Minus && tokens[i+2].Type == lexer.Minus {
							variableAssignment.Value = ast.ExpressionStatement{Left: tokens[i], Operator: tokens[i+1], Right: lexer.Token{Type: lexer.Int, Value: 1}}
						}
						statement.Node = variableAssignment
					}
				} else {
					err = errors.New(fmt.Sprintf("Variable '%s' not declered on line: %d at position %d in file: '%s'", tokens[i].Value, tokens[i].Line, tokens[i].LinePos, tokens[i].File))
					return program.Statements, err
				}
			}
		}
		program.Statements = append(program.Statements, statement)
	}
	return program.Statements, err
}

func ParseArray(tokens []lexer.Token, Variables map[string]string) (any, any) {
	var arraytype any
	var length any
	for i2, j := 0, len(tokens)-1; i2 < j; i2, j = i2+1, j-1 {
		tokens[i2], tokens[j] = tokens[j], tokens[i2]
	}

	i := len(tokens) - 1
	if len(tokens) == 1 {
		switch tokens[0].Value {
		case "int":
			arraytype = "int"
		case "float":
			arraytype = "float"
		case "string":
			arraytype = "string"
		case "bool":
			arraytype = "bool"
		}
	} else {

		for i >= 0 && tokens[i].Type != lexer.OpenBracket {
			i--
		}
		length, _, _ = ParseExpression(tokens[i+1:len(tokens)-1], Variables, false)

		n := i - 1
		innerArrayLenExp := make([]lexer.Token, 0)
		for n >= 0 && tokens[n].Type != lexer.OpenBracket {
			innerArrayLenExp = append(innerArrayLenExp, tokens[n])
			n--
		}
		var innerLength any
		innerLength, _, _ = ParseExpression(innerArrayLenExp, Variables, false)

		innerArray, _ := ParseArray(tokens[:i], Variables)
		arraytype = ast.ArrayType{Type: innerArray, Length: innerLength}
	}
	return arraytype, length

}

func getArgs(tokens []lexer.Token, i int, Variables map[string]string, funct bool) ([]lexer.Token, int) {

	var expression []lexer.Token
	for tokens[i].Type != lexer.CloseParen && tokens[i].Type != lexer.NewLine {
		if !funct && tokens[i].Type == lexer.Comma {
			break
		}
		switch tokens[i].Value.(type) {
		case string:
			if _, exist := Functions[tokens[i].Value.(string)]; exist {
				expression = append(expression, tokens[i])
				i++
				expression = append(expression, tokens[i])
				i++
				if tokens[i].Type != lexer.CloseParen {
					var args []lexer.Token
					args, i = getArgs(tokens, i, Variables, true)
					expression = append(expression, args...)
					i++
					expression = append(expression, tokens[i])

				} else {
					expression = append(expression, tokens[i])
				}

			} else {
				expression = append(expression, tokens[i])
			}
		default:
			expression = append(expression, tokens[i])
		}
		if tokens[i].Type == lexer.CloseParen {

			break
		}
		if i+1 < len(tokens) {
			i++
		}
	}
	return expression, i
}
func ParseBool(tokens []lexer.Token, Variables map[string]string) (ast.BoolExpression, error) {
	var boolStatement ast.BoolExpression
	var err error = nil
	i := 0
	for tokens[i].Type != lexer.DoubleEqual && tokens[i].Type != lexer.NotEqual && tokens[i].Type != lexer.LessThan && tokens[i].Type != lexer.MoreThan {
		i++
	}
	tokens2 := make([]lexer.Token, 0)
	j := 0
	for j < len(tokens) {
		tokens2 = append(tokens2, tokens[j])
		j++
	}

	boolStatement.Condition1, _, err = ParseExpression(tokens[:i], Variables, false)
	tokens = tokens2
	boolStatement.Operator = tokens[i]
	boolStatement.Condition2, _, err = ParseExpression(tokens[i+1:], Variables, false)
	return boolStatement, err
}

func IsNumber(token lexer.Token, Variables map[string]string) bool {
	switch {
	case token.Type == lexer.Int:
		return true
	case token.Type == lexer.Identifier:
		if Variables[token.Value.(string)] == "int" || Variables[token.Value.(string)] == "float" || Variables[token.Value.(string)] == "array" {
			return true
		} else if Functions[token.Value.(string)] == "int" || Functions[token.Value.(string)] == "float" {
			return true
		} else if token.Value.(string) == "int" || token.Value.(string) == "float" {
			return true
		}
		return false
	case token.Type == lexer.Float:
		return true
	default:
		return false
	}
}
func AddParentheses(tokens []lexer.Token) []lexer.Token {
	openParen := make([]lexer.Token, 1)
	closeParen := make([]lexer.Token, 1)

	openParen[0] = lexer.Token{Type: lexer.OpenParen, Value: "("}
	closeParen[0] = lexer.Token{Type: lexer.CloseParen, Value: ")"}
	//Divide
	i := 0
	for i+1 < len(tokens) {
		for i < len(tokens) && tokens[i].Type != lexer.Divide {
			i++
		}
		if i >= len(tokens) {
			break
		}
		if tokens[i-1].Type != lexer.CloseParen {
			tokens = append(tokens[:i-1], append(openParen, tokens[i-1:]...)...)
			i++
		} else {
			j := i - 2
			n := 1
			for n > 0 {
				if tokens[j].Type == lexer.CloseParen {
					n++
				} else if tokens[j].Type == lexer.OpenParen {
					n--
				}
				j--
			}
			j++

			if _, exist := Functions[tokens[i].Value.(string)]; exist {
				j--
				tokens = append(tokens[:j], append(openParen, tokens[j:]...)...)
				i++
			} else {
				tokens = append(tokens[:j], append(openParen, tokens[j:]...)...)
				i++
			}
		}
		if tokens[i+1].Type != lexer.OpenParen {
			tokens = append(tokens[:i+2], append(closeParen, tokens[i+2:]...)...)

		} else {
			j := i + 2
			n := 1
			for n > 0 {
				if tokens[j].Type == lexer.CloseParen {
					n--
				} else if tokens[j].Type == lexer.OpenParen {
					n++
				}
				j++
			}
			j--
			tokens = append(tokens[:j], append(closeParen, tokens[j:]...)...)

		}
		i++
	}
	//Multiply
	i = 0
	for i+1 < len(tokens) {
		for i < len(tokens) && tokens[i].Type != lexer.Multiply {
			i++
		}
		if i >= len(tokens) {
			break
		}
		if tokens[i-1].Type != lexer.CloseParen {
			tokens = append(tokens[:i-1], append(openParen, tokens[i-1:]...)...)
			i++
		} else {
			j := i - 2
			n := 1
			for n > 0 {
				if tokens[j].Type == lexer.CloseParen {
					n++
				} else if tokens[j].Type == lexer.OpenParen {
					n--
				}
				j--
			}
			j++
			if _, exist := Functions[tokens[i].Value.(string)]; exist {
				j--
				tokens = append(tokens[:j], append(openParen, tokens[j:]...)...)
				i++
			} else {
				tokens = append(tokens[:j], append(openParen, tokens[j:]...)...)
				i++
			}
		}
		if tokens[i+1].Type != lexer.OpenParen {
			tokens = append(tokens[:i+2], append(closeParen, tokens[i+2:]...)...)

		} else {
			j := i + 2
			n := 1
			for n > 0 {
				if tokens[j].Type == lexer.CloseParen {
					n--
				} else if tokens[j].Type == lexer.OpenParen {
					n++
				}
				j++
			}
			j--
			tokens = append(tokens[:j], append(closeParen, tokens[j:]...)...)

		}
		i++
	}
	i = 0
	//Sub
	for i+1 < len(tokens) {
		for i < len(tokens) && tokens[i].Type != lexer.Minus {
			i++
		}
		if i >= len(tokens) {
			break
		}
		if tokens[i-1].Type != lexer.CloseParen {
			tokens = append(tokens[:i-1], append(openParen, tokens[i-1:]...)...)
			i++
		} else {
			j := i - 2
			n := 1
			for n > 0 {
				if tokens[j].Type == lexer.CloseParen {
					n++
				} else if tokens[j].Type == lexer.OpenParen {
					n--
				}
				j--
			}
			j++
			if _, exist := Functions[tokens[i].Value.(string)]; exist {
				j--
				tokens = append(tokens[:j], append(openParen, tokens[j:]...)...)
				i++
			} else {
				tokens = append(tokens[:j], append(openParen, tokens[j:]...)...)
				i++
			}
		}
		if tokens[i+1].Type != lexer.OpenParen {
			tokens = append(tokens[:i+2], append(closeParen, tokens[i+2:]...)...)

		} else {
			j := i + 2
			n := 1
			for n > 0 {
				if tokens[j].Type == lexer.CloseParen {
					n--
				} else if tokens[j].Type == lexer.OpenParen {
					n++
				}
				j++
			}
			j--
			tokens = append(tokens[:j], append(closeParen, tokens[j:]...)...)

		}
		i++
	}
	i = 0
	//Add
	for i+1 < len(tokens) {
		for i < len(tokens) && tokens[i].Type != lexer.Plus {
			i++
		}
		if i >= len(tokens) {
			break
		}
		if tokens[i-1].Type != lexer.CloseParen {
			tokens = append(tokens[:i-1], append(openParen, tokens[i-1:]...)...)
			i++
		} else {
			j := i - 2
			n := 1
			for n > 0 {
				if tokens[j].Type == lexer.CloseParen {
					n++
				} else if tokens[j].Type == lexer.OpenParen {
					n--
				}
				j--
			}
			j++
			if _, exist := Functions[tokens[i].Value.(string)]; exist {
				j--
				tokens = append(tokens[:j], append(openParen, tokens[j:]...)...)
				i++
			} else {
				tokens = append(tokens[:j], append(openParen, tokens[j:]...)...)
				i++
			}
		}
		if tokens[i+1].Type != lexer.OpenParen {
			tokens = append(tokens[:i+2], append(closeParen, tokens[i+2:]...)...)

		} else {
			j := i + 2
			n := 1
			for n > 0 {

				if tokens[j].Type == lexer.CloseParen {
					n--
				} else if tokens[j].Type == lexer.OpenParen {
					n++
				}
				j++
			}
			j--
			tokens = append(tokens[:j], append(closeParen, tokens[j:]...)...)

		}
		i++
	}

	return tokens
}
func ParseExpression(express []lexer.Token, Variables map[string]string, parentheses bool) (any, bool, error) {
	var err error = nil
	var expressionStatement ast.ExpressionStatement
	i := 0
	//hasPlusOrMinus := false
	/*InFunc := 0
	for i < len(express) {
		if express[i].Type == lexer.OpenParen {
			InFunc++
		} else if express[i].Type == lexer.CloseParen {
			InFunc--
		}
		if (lexer.TokenType(express[i].Type) == lexer.Plus || lexer.TokenType(express[i].Type) == lexer.Minus) && InFunc == 0 {
			hasPlusOrMinus = true
			break
		}
		i++
	}*/
	i = 0
	for i < len(express) && (express[i].Type == lexer.OpenParen || lexer.IsOperator(express[i])) {
		i++
	}
	if IsNumber(express[i], Variables) {
		if !parentheses {
			express = AddParentheses(express)
		}
		i = 0
		if express[i].Type == lexer.OpenParen {
			i++
			if express[i].Type == lexer.OpenParen {
				i++
				n := 1
				for n > 0 {
					if express[i].Type == lexer.OpenParen {
						n++
					} else if express[i].Type == lexer.CloseParen {
						n--
					}
					i++
				}
				if express[i].Type != lexer.CloseParen {
					expressionStatement.Left, _, err = ParseExpression(express[1:i], Variables, true)
					expressionStatement.Operator = express[i]
					if express[i+1].Type == lexer.OpenParen {
						expressionStatement.Right, _, err = ParseExpression(express[i+1:len(express)-1], Variables, true)
					} else {
						var funcCall ast.FuncCall
						var fname string
						switch express[i].Value.(type) {
						case string:
							fname = express[i].Value.(string)
						case int:
							fname = strconv.Itoa(express[i].Value.(int))
						}
						if _, exist := Functions[fname]; exist && express[i].Type == lexer.Identifier {

							funcCall.Name = express[i]
							funcCall.Arguments = make([]any, i)
							i += 2
							for express[i].Type != lexer.CloseParen && express[i].Type != lexer.NewLine {
								var expression []lexer.Token

								expression, i = getArgs(express, i, Variables, false)
								if express[i].Type == lexer.Comma {
									i++
								}
								arg, _, err := ParseExpression(expression, Variables, true)
								if err != nil {
									return expressionStatement, true, err
								}
								funcCall.Arguments = append(funcCall.Arguments, arg)
								if err != nil {
									return expressionStatement, true, err
								}
							}

							expressionStatement.Right = funcCall
						} else if Variables[fname] == "array" && express[i+1].Type == lexer.OpenBracket {
							var ArrayCall ast.ArrayCall

							ArrayCall.Name = express[i]
							ArrayCall.Indexes = make([]any, 0)
							i += 2
							n := 1
							expression := make([]lexer.Token, 0)
							first := true
							for express[i].Type == lexer.OpenBracket || first {
								first = false
								for n > 0 {
									if express[i].Type == lexer.OpenBracket {
										n++
									} else if express[i].Type == lexer.CloseBracket {
										n--
									}
									expression = append(expression, express[i])

									i++
								}
								index, _, err := ParseExpression(expression, Variables, true)
								if err != nil {
									return expressionStatement, true, err
								}
								ArrayCall.Indexes = append(ArrayCall.Indexes, index)
							}

							expressionStatement.Right = ArrayCall

						} else {
							i++
							expressionStatement.Right = express[i]

						}
					}
				} else {
					return ParseExpression(express[1:len(express)-1], Variables, true)
				}

			} else {
				var funcCall ast.FuncCall
				var fname string
				switch express[i].Value.(type) {
				case string:
					fname = express[i].Value.(string)
				case int:
					fname = strconv.Itoa(express[i].Value.(int))
				}
				if _, exist := Functions[fname]; exist && express[i].Type == lexer.Identifier {

					funcCall.Name = express[i]
					funcCall.Arguments = make([]any, i)
					i += 2
					for express[i].Type != lexer.CloseParen && express[i].Type != lexer.NewLine {
						var expression []lexer.Token

						expression, i = getArgs(express, i, Variables, false)
						if express[i].Type == lexer.Comma {
							i++
						}
						arg, _, err := ParseExpression(expression, Variables, true)
						if err != nil {
							return expressionStatement, true, err
						}
						funcCall.Arguments = append(funcCall.Arguments, arg)
						if err != nil {
							return expressionStatement, true, err
						}
					}

					expressionStatement.Left = funcCall
				} else if Variables[fname] == "array" && express[i+1].Type == lexer.OpenBracket {
					var ArrayCall ast.ArrayCall

					ArrayCall.Name = express[i]
					ArrayCall.Indexes = make([]any, 0)
					i += 2
					n := 1
					expression := make([]lexer.Token, 0)
					first := true
					for express[i].Type == lexer.OpenBracket || first {
						first = false
						for n > 0 {
							if express[i].Type == lexer.OpenBracket {
								n++
							} else if express[i].Type == lexer.CloseBracket {
								n--
							}
							expression = append(expression, express[i])

							i++
						}
						index, _, err := ParseExpression(expression, Variables, true)
						if err != nil {
							return expressionStatement, true, err
						}
						ArrayCall.Indexes = append(ArrayCall.Indexes, index)
					}

					expressionStatement.Left = ArrayCall

				} else {
					expressionStatement.Left = express[i]

				}
				i++
				expressionStatement.Operator = express[i]
				if express[i+1].Type == lexer.OpenParen {
					expressionStatement.Right, _, err = ParseExpression(express[i+1:len(express)-1], Variables, true)
				} else {
					i++
					var funcCall ast.FuncCall
					var fname string
					switch express[i].Value.(type) {
					case string:
						fname = express[i].Value.(string)
					case int:
						fname = strconv.Itoa(express[i].Value.(int))
					}
					if _, exist := Functions[fname]; exist && express[i].Type == lexer.Identifier {

						funcCall.Name = express[i]
						funcCall.Arguments = make([]any, i)
						i += 2
						for express[i].Type != lexer.CloseParen && express[i].Type != lexer.NewLine {
							var expression []lexer.Token

							expression, i = getArgs(express, i, Variables, false)
							if express[i].Type == lexer.Comma {
								i++
							}
							arg, _, err := ParseExpression(expression, Variables, true)
							if err != nil {
								return expressionStatement, true, err
							}
							funcCall.Arguments = append(funcCall.Arguments, arg)
							if err != nil {
								return expressionStatement, true, err
							}
						}

						expressionStatement.Right = funcCall
					} else if Variables[fname] == "array" && express[i+1].Type == lexer.OpenBracket {
						var ArrayCall ast.ArrayCall

						ArrayCall.Name = express[i]
						ArrayCall.Indexes = make([]any, 0)
						i += 2
						n := 1
						expression := make([]lexer.Token, 0)
						first := true
						for express[i].Type == lexer.OpenBracket || first {
							first = false
							for n > 0 {
								if express[i].Type == lexer.OpenBracket {
									n++
								} else if express[i].Type == lexer.CloseBracket {
									n--
								}
								expression = append(expression, express[i])

								i++
							}
							index, _, err := ParseExpression(expression, Variables, true)
							if err != nil {
								return expressionStatement, true, err
							}
							ArrayCall.Indexes = append(ArrayCall.Indexes, index)
						}

						expressionStatement.Right = ArrayCall

					} else {
						expressionStatement.Right = express[i]
					}

				}

			}

			/*if hasPlusOrMinus {
				for lexer.TokenType(express[i].Type) != lexer.Plus && lexer.TokenType(express[i].Type) != lexer.Minus {
					i++
				}
				expressionStatement.Left, _, err = ParseExpression(express[:i], Variables)
				expressionStatement.Operator = express[i]
				expressionStatement.Right, _, err = ParseExpression(express[i+1:], Variables)
			} else {
				if len(express) > 1 {
					op := 1
					var funcCall ast.FuncCall
					var fname string
					switch express[0].Value.(type) {
					case string:
						fname = express[0].Value.(string)
					case int:
						fname = strconv.Itoa(express[0].Value.(int))
					}
					if _, exist := Functions[fname]; exist && express[0].Type == lexer.Identifier {

						funcCall.Name = express[0]
						funcCall.Arguments = make([]any, 0)
						i += 2
						for express[i].Type != lexer.CloseParen && express[i].Type != lexer.NewLine {
							var expression []lexer.Token

							expression, i = getArgs(express, i, Variables, false)
							if express[i].Type == lexer.Comma {
								i++
							}
							arg, _, err := ParseExpression(expression, Variables)
							if err != nil {
								return expressionStatement, true, err
							}
							funcCall.Arguments = append(funcCall.Arguments, arg)
							if err != nil {
								return expressionStatement, true, err
							}
						}
						expressionStatement.Left = funcCall
						op = i + 1
					} else {
						expressionStatement.Left = express[0]
					}
					if len(express) <= op {
						return funcCall, true, err

					}
					expressionStatement.Operator = express[op]
					expressionStatement.Right, _, err = ParseExpression(express[op+1:], Variables)
				} else {
					var ret any
					var fname string
					switch express[0].Value.(type) {
					case string:
						fname = express[0].Value.(string)
					case int:
						fname = strconv.Itoa(express[0].Value.(int))
					}
					if _, exist := Functions[fname]; exist && express[0].Type == lexer.Identifier {
						var funcCall ast.FuncCall
						funcCall.Name = express[0]
						funcCall.Arguments = make([]any, 0)
						i += 2
						for express[i].Type != lexer.CloseParen && express[i].Type != lexer.NewLine {
							var expression []lexer.Token

							expression, i = getArgs(express, i, Variables, false)
							if express[i].Type == lexer.Comma {
								i++
							}
							arg, _, err := ParseExpression(expression, Variables)
							if err != nil {
								return expressionStatement, true, err
							}
							funcCall.Arguments = append(funcCall.Arguments, arg)
							if err != nil {
								return expressionStatement, true, err
							}
						}
						ret = funcCall
					} else {
						ret = express[0]
					}

					return ret, true, err
				}
			}*/
		} else {
			var funcCall ast.FuncCall
			var fname string
			switch express[i].Value.(type) {
			case string:
				fname = express[i].Value.(string)
			case int:
				fname = strconv.Itoa(express[i].Value.(int))
			}
			if _, exist := Functions[fname]; exist && express[i].Type == lexer.Identifier {

				funcCall.Name = express[i]
				funcCall.Arguments = make([]any, i)
				i += 2
				for express[i].Type != lexer.CloseParen && express[i].Type != lexer.NewLine {
					var expression []lexer.Token

					expression, i = getArgs(express, i, Variables, false)
					if express[i].Type == lexer.Comma {
						i++
					}
					arg, _, err := ParseExpression(expression, Variables, true)
					if err != nil {
						return expressionStatement, true, err
					}
					funcCall.Arguments = append(funcCall.Arguments, arg)
					if err != nil {
						return expressionStatement, true, err
					}
				}

				return funcCall, true, err
			} else if Variables[fname] == "array" && i+1 < len(express) && express[i+1].Type == lexer.OpenBracket {
				var ArrayCall ast.ArrayCall

				ArrayCall.Name = express[i]
				ArrayCall.Indexes = make([]any, 0)
				i += 1

				for express[i].Type == lexer.OpenBracket {
					expression := make([]lexer.Token, 0)
					i++
					n := 1
					for n > 0 {
						if express[i].Type == lexer.OpenBracket {
							n++
						} else if express[i].Type == lexer.CloseBracket {
							n--
						}
						if n <= 0 {
							break
						}
						expression = append(expression, express[i])

						i++
					}

					index, _, err := ParseExpression(expression, Variables, true)
					if err != nil {
						return expressionStatement, true, err
					}
					ArrayCall.Indexes = append(ArrayCall.Indexes, index)
					i++
					if i >= len(express) {
						break
					}

				}

				return ArrayCall, true, err

			} else {
				return express[i], true, err
			}
		}
		return expressionStatement, true, err
	} else if express[i].Type == lexer.OpenBracket {
		i++
		array := ast.ArrayStatement{Length: 0}
		for express[i].Type != lexer.CloseBracket {
			if express[i].Type == lexer.OpenBracket {
				expression := []lexer.Token{}
				expression = append(expression, express[i])
				i++
				n := 1
				for n > 0 {
					if express[i].Type == lexer.OpenBracket {
						n++
					} else if express[i].Type == lexer.CloseBracket {
						n--
					}
					expression = append(expression, express[i])

					i++
				}

				elem, _, _ := ParseExpression(expression, Variables, true)
				array.Content = append(array.Content, elem)

				i++

			} else {
				expression := []lexer.Token{}

				for express[i].Type != lexer.Comma && express[i].Type != lexer.CloseBracket {
					expression = append(expression, express[i])
					i++

				}
				element, _, _ := ParseExpression(expression, Variables, true)
				array.Content = append(array.Content, element)
			}
			if i >= len(express) {
				array.Length++
				break
			}

			if express[i].Type == lexer.Comma {
				i++
				array.Length++
			}

		}
		array.Length++
		return array, true, err

	} else {
		i = 0
		if len(express) > 1 {
			op := 1
			var funcCall ast.FuncCall
			var fname string
			switch express[0].Value.(type) {
			case string:
				fname = express[0].Value.(string)
			case int:
				fname = strconv.Itoa(express[0].Value.(int))
			}
			if _, exist := Functions[fname]; exist && express[0].Type == lexer.Identifier {

				funcCall.Name = express[0]
				funcCall.Arguments = make([]any, 0)
				i += 2
				for express[i].Type != lexer.CloseParen && express[i].Type != lexer.NewLine {
					var expression []lexer.Token

					expression, i = getArgs(express, i, Variables, false)

					if express[i].Type == lexer.Comma {
						i++
					}
					arg, _, err := ParseExpression(expression, Variables, true)
					if err != nil {
						return expressionStatement, true, err
					}
					funcCall.Arguments = append(funcCall.Arguments, arg)
					if err != nil {
						return expressionStatement, true, err
					}
				}
				op = i + 1
				expressionStatement.Left = funcCall
			} else {
				expressionStatement.Left = express[0]
			}
			if len(express) <= op {
				return funcCall, false, err

			}
			expressionStatement.Operator = express[op]
			expressionStatement.Right, _, err = ParseExpression(express[op+1:], Variables, true)
			if expressionStatement.Operator.Type != lexer.Plus {
				err = errors.New(fmt.Sprintf("Expected '+' on line: %d at position %d, not '%s' in file '%s'", expressionStatement.Operator.Line, expressionStatement.Operator.LinePos, expressionStatement.Operator.Value, expressionStatement.Operator.File))
			}
		} else {
			var ret any
			var fname string
			switch express[0].Value.(type) {
			case string:
				fname = express[0].Value.(string)
			case int:
				fname = strconv.Itoa(express[0].Value.(int))
			}
			if _, exist := Functions[fname]; exist && express[0].Type == lexer.Identifier {
				var funcCall ast.FuncCall
				funcCall.Name = express[0]
				funcCall.Arguments = make([]any, 0)
				i += 2
				for express[i].Type != lexer.CloseParen {
					var expression []lexer.Token

					expression, i = getArgs(express, i, Variables, false)
					if express[i].Type == lexer.Comma {
						i++
					}
					arg, _, err := ParseExpression(expression, Variables, true)
					if err != nil {
						return expressionStatement, false, err
					}
					funcCall.Arguments = append(funcCall.Arguments, arg)
					if err != nil {
						return expressionStatement, false, err
					}
				}
				ret = funcCall
			} else {
				ret = express[0]
			}

			return ret, false, err
		}
		return expressionStatement, false, err
	}
}
