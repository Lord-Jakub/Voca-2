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
	Condition   BoolStatement
	Consequence []Statement
	Alternative []Statement
}
type BoolStatement struct {
	Condition1 any
	Operator   lexer.Token
	Condition2 any
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
type ExternFuncDeclaration struct {
	Name      lexer.Token
	Arguments []any
	Type      lexer.Token
}
type ReturnStatement struct {
	Value any
}
type VariableDeclaration struct {
	Name  lexer.Token
	Type  lexer.Token
	Value any
}
type VariableAssignment struct {
	Name  lexer.Token
	Value any
}

type FuncCall struct {
	Name      lexer.Token
	Arguments []any
}
