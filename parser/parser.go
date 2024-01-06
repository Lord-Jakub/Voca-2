package parser

import (
	"Voca-2/ast"
	"Voca-2/lexer"
	"errors"
	"fmt"
)

func New(tokens []lexer.Token) (ast.Program, error) {
	var err error = nil
	Functions, err = MapFunctions(tokens)
	steatments, err := Parse(tokens)

	return ast.Program{Statements: steatments}, err
}

var Functions = make(map[string]string)
var Variables = make(map[string]string)

func MapFunctions(tokens []lexer.Token) (map[string]string, error) {
	var functions = make(map[string]string)
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Type == lexer.Keyword && tokens[i].Value == "func" {
			name := tokens[i+1].Value
			for tokens[i].Type != lexer.OpenBrace {
				i++
			}
			functions[name.(string)] = tokens[i-1].Value.(string)
		}
	}
	return functions, nil
}
func Parse(tokens []lexer.Token) ([]ast.Statement, error) {
	var program ast.Program
	var err error = nil
	for i := 0; i < len(tokens); i++ {
		var statement ast.Statement
		switch tokens[i].Type {
		case lexer.Keyword:
			switch tokens[i].Value {
			case "int":
				var variableDeclaration ast.VariableDeclaration
				variableDeclaration.Type = tokens[i]
				variableDeclaration.Name = tokens[i+1]
				if tokens[i+2].Type == lexer.Equal {
					var expression []lexer.Token
					i += 3
					for lexer.IsOperator(tokens[i+1]) {
						expression = append(expression, tokens[i])
						expression = append(expression, tokens[i+1])
						i += 2
						if i+1 >= len(tokens) {
							break
						}
					}
					expression = append(expression, tokens[i])
					value, IsNum, err := ParseExpression(expression)
					if err != nil {
						return program.Statements, err
					}
					variableDeclaration.Value = value

					if !IsNum {
						err = errors.New(fmt.Sprintf("Expected number on line: %d at position %d, not '%s'", expression[0].Line, expression[0].LinePos, expression[0].Value))
						return program.Statements, err
					}
				}
				statement.Node = variableDeclaration
				Variables[variableDeclaration.Name.Value.(string)] = "int"
			case "string":
				var variableDeclaration ast.VariableDeclaration
				variableDeclaration.Type = tokens[i]
				variableDeclaration.Name = tokens[i+1]
				if tokens[i+2].Type == lexer.Equal {
					var expression []lexer.Token
					i += 3
					for lexer.IsOperator(tokens[i+1]) {
						expression = append(expression, tokens[i])
						expression = append(expression, tokens[i+1])
						i += 2
						if i+1 >= len(tokens) {
							break
						}
					}
					expression = append(expression, tokens[i])
					value, IsNum, err := ParseExpression(expression)
					if err != nil {
						return program.Statements, err
					}
					variableDeclaration.Value = value
					if IsNum {
						err = errors.New(fmt.Sprintf("Expected number on line: %d at position %d, not '%s'", expression[0].Line, expression[0].LinePos, expression[0].Value))
						return program.Statements, err
					}
				}
				statement.Node = variableDeclaration
				Variables[variableDeclaration.Name.Value.(string)] = "string"
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
							err = errors.New(fmt.Sprintf("Expected type on line: %d at position %d, not '%s'", tokens[i].Line, tokens[i].LinePos, tokens[i].Value))
							return program.Statements, err
						}
						funcDeclaration.Arguments = append(funcDeclaration.Arguments, arg)
					} else if tokens[i].Type == lexer.Comma {

					} else {

						err = errors.New(fmt.Sprintf("Expected argument on line: %d at position %d, not '%s'", tokens[i].Line, tokens[i].LinePos, tokens[i].Value))
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
					err = errors.New(fmt.Sprintf("Expected type or { on line: %d at position %d, not '%s'", tokens[i].Line, tokens[i].LinePos, tokens[i].Value))
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
				funcDeclaration.Body, err = Parse(bodyTokens)
				if err != nil {
					return program.Statements, err
				}
				statement.Node = funcDeclaration
			case "print":
				var printStatement ast.PrintStatement
				var expression []lexer.Token
				i += 2
				for tokens[i].Type != lexer.CloseParen {
					expression = append(expression, tokens[i])
					i++
				}
				value, _, err := ParseExpression(expression)
				if err != nil {
					return program.Statements, err
				}

				printStatement.Value = value
				statement.Node = printStatement
			case "return":
				var returnStatement ast.ReturnStatement
				var expression []lexer.Token
				i += 1
				for tokens[i].Type != lexer.NewLine {
					expression = append(expression, tokens[i])
					i++
				}
				value, _, err := ParseExpression(expression)
				if err != nil {
					return program.Statements, err
				}
				returnStatement.Value = value
				statement.Node = returnStatement

			}
		case lexer.Identifier:
			if tokens[i+1].Type == lexer.OpenParen {
				//if identifier is in functions map
				if _, exist := Functions[tokens[i].Value.(string)]; exist {
					var funcCall ast.FuncCall
					funcCall.Name = tokens[i]
					funcCall.Arguments = make([]any, 0)
					i += 2
					for tokens[i].Type != lexer.CloseParen {
						var expression []lexer.Token
						for tokens[i].Type != lexer.Comma && tokens[i].Type != lexer.CloseParen {
							expression = append(expression, tokens[i])
							i++
						}
						if tokens[i].Type == lexer.Comma {
							i++
						}
						arg, _, err := ParseExpression(expression)
						if err != nil {
							return program.Statements, err
						}

						funcCall.Arguments = append(funcCall.Arguments, arg)
						if err != nil {
							return program.Statements, err
						}
					}
					statement.Node = funcCall
				} else {
					err = errors.New(fmt.Sprintf("Function '%s' not declered on line: %d at position %d", tokens[i].Value, tokens[i].Line, tokens[i].LinePos))
					return program.Statements, err
				}
			} else {
				if Variables[tokens[i].Value.(string)] != "" {
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
						value, IsNum, err := ParseExpression(expression)
						if err != nil {
							return program.Statements, err
						}
						variableAssignment.Value = value
						if IsNum {
							if Variables[variableAssignment.Name.Value.(string)] != "int" {
								err = errors.New(fmt.Sprintf("Expected string on line: %d at position %d, not '%s'", expression[0].Line, expression[0].LinePos, expression[0].Value))
								return program.Statements, err
							}
						} else {
							if Variables[variableAssignment.Name.Value.(string)] != "string" {
								err = errors.New(fmt.Sprintf("Expected number on line: %d at position %d, not '%s'", expression[0].Line, expression[0].LinePos, expression[0].Value))
								return program.Statements, err
							}
						}
					} else if tokens[i+1].Type == lexer.Plus && tokens[i+2].Type == lexer.Plus {
						variableAssignment.Value = ast.ExpressionStatement{Left: tokens[i], Operator: tokens[i+1], Right: lexer.Token{Type: lexer.Int, Value: 1}}
					} else if tokens[i+1].Type == lexer.Minus && tokens[i+2].Type == lexer.Minus {
						variableAssignment.Value = ast.ExpressionStatement{Left: tokens[i], Operator: tokens[i+1], Right: lexer.Token{Type: lexer.Int, Value: 1}}
					}
					statement.Node = variableAssignment
				} else {
					err = errors.New(fmt.Sprintf("Variable '%s' not declered on line: %d at position %d", tokens[i].Value, tokens[i].Line, tokens[i].LinePos))
					return program.Statements, err
				}
			}
		}
		program.Statements = append(program.Statements, statement)
	}
	return program.Statements, err
}
func IsNumber(token lexer.Token) bool {
	switch {
	case token.Type == lexer.Int:
		return true
	case token.Type == lexer.Identifier:
		if Variables[token.Value.(string)] == "int" {
			return true
		} else if Functions[token.Value.(string)] == "int" {
			return true
		} else if token.Value.(string) == "int" {
			return true
		}
		return false
	default:
		return false
	}
}
func ParseExpression(tokens []lexer.Token) (any, bool, error) {
	var err error = nil
	var expressionStatement ast.ExpressionStatement
	i := 0
	hasPlusOrMinus := false
	InFunc := 0
	for i < len(tokens) {
		if tokens[i].Type == lexer.OpenParen {
			InFunc++
		} else if tokens[i].Type == lexer.CloseParen {
			InFunc--
		}
		if (lexer.TokenType(tokens[i].Type) == lexer.Plus || lexer.TokenType(tokens[i].Type) == lexer.Minus) && InFunc == 0 {
			hasPlusOrMinus = true
			break
		}
		i++
	}
	i = 0
	if IsNumber(tokens[i]) {
		if hasPlusOrMinus {
			for lexer.TokenType(tokens[i].Type) != lexer.Plus && lexer.TokenType(tokens[i].Type) != lexer.Minus {
				i++
			}
			expressionStatement.Left, _, err = ParseExpression(tokens[:i])
			expressionStatement.Operator = tokens[i]
			expressionStatement.Right, _, err = ParseExpression(tokens[i+1:])
		} else {
			if len(tokens) > 1 {
				op := 1
				var funcCall ast.FuncCall
				if _, exist := Functions[tokens[0].Value.(string)]; exist && tokens[0].Type == lexer.Identifier {

					funcCall.Name = tokens[0]
					funcCall.Arguments = make([]any, 0)
					i += 2
					for tokens[i].Type != lexer.CloseParen {
						var expression []lexer.Token
						for tokens[i].Type != lexer.Comma && tokens[i].Type != lexer.CloseParen {
							expression = append(expression, tokens[i])
							i++
						}
						if tokens[i].Type == lexer.Comma {
							i++
						}
						arg, _, err := ParseExpression(expression)
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
					expressionStatement.Left = tokens[0]
				}
				if len(tokens) <= op {
					return funcCall, true, err

				}
				expressionStatement.Operator = tokens[op]
				expressionStatement.Right, _, err = ParseExpression(tokens[op+1:])
			} else {
				var ret any
				if _, exist := Functions[tokens[0].Value.(string)]; exist && tokens[0].Type == lexer.Identifier {
					var funcCall ast.FuncCall
					funcCall.Name = tokens[0]
					funcCall.Arguments = make([]any, 0)
					i += 2
					for tokens[i].Type != lexer.CloseParen {
						var expression []lexer.Token
						for tokens[i].Type != lexer.Comma && tokens[i].Type != lexer.CloseParen {
							expression = append(expression, tokens[i])
							i++
						}
						if tokens[i].Type == lexer.Comma {
							i++
						}
						arg, _, err := ParseExpression(expression)
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
					ret = tokens[0]
				}

				return ret, true, err
			}
		}
		return expressionStatement, true, err
	} else {
		if len(tokens) > 1 {
			op := 1
			var funcCall ast.FuncCall
			if _, exist := Functions[tokens[0].Value.(string)]; exist && tokens[0].Type == lexer.Identifier {

				funcCall.Name = tokens[0]
				funcCall.Arguments = make([]any, 0)
				i += 2
				for tokens[i].Type != lexer.CloseParen {
					var expression []lexer.Token
					for tokens[i].Type != lexer.Comma && tokens[i].Type != lexer.CloseParen {
						expression = append(expression, tokens[i])
						i++
					}
					if tokens[i].Type == lexer.Comma {
						i++
					}
					arg, _, err := ParseExpression(expression)
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
				expressionStatement.Left = tokens[0]
			}
			if len(tokens) <= op {
				return funcCall, false, err

			}
			expressionStatement.Operator = tokens[op]
			expressionStatement.Right, _, err = ParseExpression(tokens[op+1:])
			if expressionStatement.Operator.Type != lexer.Plus {
				err = errors.New(fmt.Sprintf("Expected '+' on line: %d at position %d, not '%s'", expressionStatement.Operator.Line, expressionStatement.Operator.LinePos, expressionStatement.Operator.Value))
			}
		} else {
			var ret any
			if _, exist := Functions[tokens[0].Value.(string)]; exist && tokens[0].Type == lexer.Identifier {
				var funcCall ast.FuncCall
				funcCall.Name = tokens[0]
				funcCall.Arguments = make([]any, 0)
				i += 2
				for tokens[i].Type != lexer.CloseParen {
					var expression []lexer.Token
					for tokens[i].Type != lexer.Comma && tokens[i].Type != lexer.CloseParen {
						expression = append(expression, tokens[i])
						i++
					}
					if tokens[i].Type == lexer.Comma {
						i++
					}
					arg, _, err := ParseExpression(expression)
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
				ret = tokens[0]
			}

			return ret, false, err
		}
		return expressionStatement, false, err
	}
}
