package ast

import (
	"Voca-2/lexer"
)

type Program struct {
	Statements []Statement
}

type Statement struct {
	Node any
}

type ExpressionStatement struct {
	Left     any
	Operator lexer.Token
	Right    any
}
type IfStatement struct {
	Condition   []lexer.Token
	Consequence []Statement
	Alternative []Statement
}
type WhileStatement struct {
	Condition   []lexer.Token
	Consequence []Statement
}
type FuncDeclaration struct {
	Name      lexer.Token
	Arguments []any
	Type      lexer.Token
	Body      []Statement
}
type ReturnStatement struct {
	ReturnValue []lexer.Token
}
type VariableDeclaration struct {
	Name  lexer.Token
	Type  lexer.Token
	Value any
}
type PrintStatement struct {
	Value any
}
type FuncCall struct {
	Name      lexer.Token
	Arguments []ExpressionStatement
}
