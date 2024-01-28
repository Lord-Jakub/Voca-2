package parser

import (
	"Voca-2/ast"
	"Voca-2/lexer"
	"errors"
	"fmt"
	"strconv"
)

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
func Parse(tokens []lexer.Token, Variables map[string]string) ([]ast.Statement, error) {
	var program ast.Program
	program.Statements = append(program.Statements, ast.Statement{Node: ast.ExternFuncDeclaration{Name: lexer.Token{Type: lexer.Keyword, Value: "print"}, Arguments: []any{ast.VariableDeclaration{Type: lexer.Token{Type: lexer.Keyword, Value: "string"}, Name: lexer.Token{Type: lexer.Identifier, Value: "s"}}}, Type: lexer.Token{Type: lexer.Keyword, Value: "void"}}})
	program.Statements = append(program.Statements, ast.Statement{Node: ast.ExternFuncDeclaration{Name: lexer.Token{Type: lexer.Keyword, Value: "append"}, Arguments: []any{ast.VariableDeclaration{Type: lexer.Token{Type: lexer.Keyword, Value: "string"}, Name: lexer.Token{Type: lexer.Identifier, Value: "s1"}}, ast.VariableDeclaration{Type: lexer.Token{Type: lexer.Keyword, Value: "string"}, Name: lexer.Token{Type: lexer.Identifier, Value: "s2"}}}, Type: lexer.Token{Type: lexer.Keyword, Value: "string"}}})
	program.Statements = append(program.Statements, ast.Statement{Node: ast.ExternFuncDeclaration{Name: lexer.Token{Type: lexer.Keyword, Value: "strlen"}, Arguments: []any{ast.VariableDeclaration{Type: lexer.Token{Type: lexer.Keyword, Value: "string"}, Name: lexer.Token{Type: lexer.Identifier, Value: "s"}}}, Type: lexer.Token{Type: lexer.Keyword, Value: "int"}}})
	program.Statements = append(program.Statements, ast.Statement{Node: ast.ExternFuncDeclaration{Name: lexer.Token{Type: lexer.Keyword, Value: "IntToString"}, Arguments: []any{ast.VariableDeclaration{Type: lexer.Token{Type: lexer.Keyword, Value: "int"}, Name: lexer.Token{Type: lexer.Identifier, Value: "num"}}}, Type: lexer.Token{Type: lexer.Keyword, Value: "string"}}})

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
					value, IsNum, err := ParseExpression(expression, Variables)
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
					value, IsNum, err := ParseExpression(expression, Variables)
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
				funcDeclaration.Body, err = Parse(bodyTokens, make(map[string]string))
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
				} else if tokens[i].Type == lexer.NewLine {
					funcDeclaration.Type = lexer.Token{Type: lexer.Keyword, Value: "void"}
				} else {
					err = errors.New(fmt.Sprintf("Expected type or { on line: %d at position %d, not '%s'", tokens[i].Line, tokens[i].LinePos, tokens[i].Value))
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
				value, _, err := ParseExpression(expression, Variables)
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
				if isexpression {
					condition, err = ParseBool(expression, Variables)
				} else if tokens[i2].Type == lexer.Identifier {
					if Variables[tokens[i2].Value.(string)] == "bool" {
						condition = tokens[i2]
					} else {
						err = errors.New(fmt.Sprintf("Expected bool on line: %d at position %d, not '%s'", tokens[i2].Line, tokens[i2].LinePos, tokens[i2].Value))
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
				if isexpression {
					condition, err = ParseBool(expression, Variables)
				} else if tokens[i2].Type == lexer.Identifier {
					if Variables[tokens[i2].Value.(string)] == "bool" {
						condition = tokens[i2]
					} else {
						err = errors.New(fmt.Sprintf("Expected bool on line: %d at position %d, not '%s'", tokens[i2].Line, tokens[i2].LinePos, tokens[i2].Value))
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
					for tokens[i].Type != lexer.CloseParen {
						var expression []lexer.Token
						for tokens[i].Type != lexer.Comma && tokens[i].Type != lexer.CloseParen {
							expression = append(expression, tokens[i])
							i++
						}
						if tokens[i].Type == lexer.Comma {
							i++
						}
						arg, _, err := ParseExpression(expression, Variables)
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
						value, IsNum, err := ParseExpression(expression, Variables)
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
func ParseBool(tokens []lexer.Token, Variables map[string]string) (ast.BoolExpression, error) {
	var boolStatement ast.BoolExpression
	var err error = nil
	i := 0
	for tokens[i].Type != lexer.DoubleEqual && tokens[i].Type != lexer.NotEqual && tokens[i].Type != lexer.LessThan && tokens[i].Type != lexer.MoreThan {
		i++
	}

	boolStatement.Condition1, _, err = ParseExpression(tokens[:i], Variables)
	boolStatement.Operator = tokens[i]
	boolStatement.Condition2, _, err = ParseExpression(tokens[i+1:], Variables)
	return boolStatement, err
}

func IsNumber(token lexer.Token, Variables map[string]string) bool {
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
func ParseExpression(tokens []lexer.Token, Variables map[string]string) (any, bool, error) {
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
	if IsNumber(tokens[i], Variables) {
		if hasPlusOrMinus {
			for lexer.TokenType(tokens[i].Type) != lexer.Plus && lexer.TokenType(tokens[i].Type) != lexer.Minus {
				i++
			}
			expressionStatement.Left, _, err = ParseExpression(tokens[:i], Variables)
			expressionStatement.Operator = tokens[i]
			expressionStatement.Right, _, err = ParseExpression(tokens[i+1:], Variables)
		} else {
			if len(tokens) > 1 {
				op := 1
				var funcCall ast.FuncCall
				var fname string
				switch tokens[0].Value.(type) {
				case string:
					fname = tokens[0].Value.(string)
				case int:
					fname = strconv.Itoa(tokens[0].Value.(int))
				}
				if _, exist := Functions[fname]; exist && tokens[0].Type == lexer.Identifier {

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
					expressionStatement.Left = tokens[0]
				}
				if len(tokens) <= op {
					return funcCall, true, err

				}
				expressionStatement.Operator = tokens[op]
				expressionStatement.Right, _, err = ParseExpression(tokens[op+1:], Variables)
			} else {
				var ret any
				var fname string
				switch tokens[0].Value.(type) {
				case string:
					fname = tokens[0].Value.(string)
				case int:
					fname = strconv.Itoa(tokens[0].Value.(int))
				}
				if _, exist := Functions[fname]; exist && tokens[0].Type == lexer.Identifier {
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
			var fname string
			switch tokens[0].Value.(type) {
			case string:
				fname = tokens[0].Value.(string)
			case int:
				fname = strconv.Itoa(tokens[0].Value.(int))
			}
			if _, exist := Functions[fname]; exist && tokens[0].Type == lexer.Identifier {

				funcCall.Name = tokens[0]
				funcCall.Arguments = make([]any, 0)
				i += 2
				for (i < len(tokens)) && (tokens[i].Type != lexer.CloseParen) {
					var expression []lexer.Token
					for (i < len(tokens)) && (tokens[i].Type != lexer.Comma && tokens[i].Type != lexer.CloseParen) {
						expression = append(expression, tokens[i])
						i++
					}
					if (i < len(tokens)) && (tokens[i].Type == lexer.Comma) {
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
				op = i + 1
				expressionStatement.Left = funcCall
			} else {
				expressionStatement.Left = tokens[0]
			}
			if len(tokens) <= op {
				return funcCall, false, err

			}
			expressionStatement.Operator = tokens[op]
			expressionStatement.Right, _, err = ParseExpression(tokens[op+1:], Variables)
			if expressionStatement.Operator.Type != lexer.Plus {
				err = errors.New(fmt.Sprintf("Expected '+' on line: %d at position %d, not '%s'", expressionStatement.Operator.Line, expressionStatement.Operator.LinePos, expressionStatement.Operator.Value))
			}
		} else {
			var ret any
			var fname string
			switch tokens[0].Value.(type) {
			case string:
				fname = tokens[0].Value.(string)
			case int:
				fname = strconv.Itoa(tokens[0].Value.(int))
			}
			if _, exist := Functions[fname]; exist && tokens[0].Type == lexer.Identifier {
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
					arg, _, err := ParseExpression(expression, Variables)
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
