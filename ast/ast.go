package ast

import (
	"Voca-2/lexer"
)

type Program struct {
	Externals  []string
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
	Invert bool
	Bool   any
}
type BoolExpression struct {
	Condition1 any
	Operator   lexer.Token
	Condition2 any
}
type WhileStatement struct {
	Condition   BoolStatement
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

type ArrayStatement struct {
	Content []any
	Length  int
}

type ArrayDeclaration struct {
	Type   any
	Name   lexer.Token
	Value  any
	Length any
}
type ArrayAssignment struct {
	Name    lexer.Token
	Indexes []any
	Value   any
}
type ArrayType struct {
	Type   any
	Length any
}
