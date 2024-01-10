package compiler

import (
	"Voca-2/ast"
	"Voca-2/lexer"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

var stringType = types.NewPointer(types.I8)
var boolType = types.I1
var Functions = make(map[string]value.Value)

func GenerateIR(program ast.Program) string {

	var module = ir.NewModule()
	i := 0
	for i < len(program.Statements) {
		statement := program.Statements[i].Node
		switch statement.(type) {
		case ast.FuncDeclaration:

			fn := module.NewFunc(statement.(ast.FuncDeclaration).Name.Value.(string), types.I32)
			entry := fn.NewBlock("entry")
			Functions[statement.(ast.FuncDeclaration).Name.Value.(string)] = fn
			entry = Compile(entry, statement.(ast.FuncDeclaration).Body, make(map[string]value.Value))
			entry.NewRet(constant.NewInt(types.I32, 0))

		}
		i++
	}
	return module.String()
}
func Compile(block *ir.Block, statements []ast.Statement, variables map[string]value.Value) *ir.Block {
	i := 0
	for i < len(statements) {
		statement := statements[i].Node
		switch statement.(type) {
		case ast.VariableDeclaration:
			var variable *ir.InstAlloca
			if statement.(ast.VariableDeclaration).Type.Value.(string) == "string" {
				variable = block.NewAlloca(stringType)
				variable.SetName(statement.(ast.VariableDeclaration).Name.Value.(string))

			} else if statement.(ast.VariableDeclaration).Type.Value.(string) == "bool" {
				variable = block.NewAlloca(boolType)
				variable.SetName(statement.(ast.VariableDeclaration).Name.Value.(string))
			} else if statement.(ast.VariableDeclaration).Type.Value.(string) == "int" {
				variable = block.NewAlloca(types.I32)
				variable.SetName(statement.(ast.VariableDeclaration).Name.Value.(string))
			}
			if statement.(ast.VariableDeclaration).Value != nil {
				variables[statement.(ast.VariableDeclaration).Name.Value.(string)] = variable
				switch statement.(ast.VariableDeclaration).Value.(type) {
				case ast.ExpressionStatement:
					block.NewStore(CompileExpression(block, statement.(ast.VariableDeclaration).Value.(ast.ExpressionStatement), variables), variable)
				case lexer.Token:
					if statement.(ast.VariableDeclaration).Value.(lexer.Token).Type == lexer.Int {
						block.NewStore(constant.NewInt(types.I32, statement.(ast.VariableDeclaration).Value.(lexer.Token).Value.(int64)), variable)
					}
				}
				//block.NewStore(variables[statement.(ast.VariableDeclaration).Name.Value.(string)], variable)
			} else {
				variables[statement.(ast.VariableDeclaration).Name.Value.(string)] = variable
			}
		}
		i++
	}
	return block
}

func CompileExpression(block *ir.Block, expression ast.ExpressionStatement, variables map[string]value.Value) value.Value {
	if expression.Operator.Type == lexer.Plus {
		switch expression.Left.(type) {
		case ast.ExpressionStatement:
			switch expression.Right.(type) {
			case ast.ExpressionStatement:
				return block.NewAdd(CompileExpression(block, expression.Left.(ast.ExpressionStatement), variables), CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables))
			case lexer.Token:
				if expression.Right.(lexer.Token).Type == lexer.Int {
					return block.NewAdd(CompileExpression(block, expression.Left.(ast.ExpressionStatement), variables), constant.NewInt(types.I32, int64(expression.Right.(lexer.Token).Value.(int))))
				} else if expression.Right.(lexer.Token).Type == lexer.Identifier {
					return block.NewAdd(CompileExpression(block, expression.Left.(ast.ExpressionStatement), variables), variables[expression.Right.(lexer.Token).Value.(string)])
				}
			}
		case lexer.Token:
			switch expression.Right.(type) {
			case ast.ExpressionStatement:
				if expression.Left.(lexer.Token).Type == lexer.Int {

					return block.NewAdd(constant.NewInt(types.I32, int64(expression.Left.(lexer.Token).Value.(int))), CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables))
				} else if expression.Left.(lexer.Token).Type == lexer.Identifier {
					return block.NewAdd(variables[expression.Left.(lexer.Token).Value.(string)], CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables))
				}
			case lexer.Token:
				if expression.Left.(lexer.Token).Type == lexer.Int && expression.Right.(lexer.Token).Type == lexer.Int {
					return block.NewAdd(constant.NewInt(types.I32, int64(expression.Left.(lexer.Token).Value.(int))), constant.NewInt(types.I32, int64(expression.Right.(lexer.Token).Value.(int))))
				} else if expression.Left.(lexer.Token).Type == lexer.Identifier && expression.Right.(lexer.Token).Type == lexer.Int {
					return block.NewAdd(variables[expression.Left.(lexer.Token).Value.(string)], constant.NewInt(types.I32, int64(expression.Right.(lexer.Token).Value.(int))))
				} else if expression.Left.(lexer.Token).Type == lexer.Int && expression.Right.(lexer.Token).Type == lexer.Identifier {
					return block.NewAdd(constant.NewInt(types.I32, int64(expression.Left.(lexer.Token).Value.(int))), variables[expression.Right.(lexer.Token).Value.(string)])
				} else if expression.Left.(lexer.Token).Type == lexer.Identifier && expression.Right.(lexer.Token).Type == lexer.Identifier {
					return block.NewAdd(variables[expression.Left.(lexer.Token).Value.(string)], variables[expression.Right.(lexer.Token).Value.(string)])
				}
			}
		}
	} else if expression.Operator.Type == lexer.Minus {
		switch expression.Left.(type) {
		case ast.ExpressionStatement:
			switch expression.Right.(type) {
			case ast.ExpressionStatement:
				return block.NewSub(CompileExpression(block, expression.Left.(ast.ExpressionStatement), variables), CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables))
			case lexer.Token:
				if expression.Right.(lexer.Token).Type == lexer.Int {
					return block.NewSub(CompileExpression(block, expression.Left.(ast.ExpressionStatement), variables), constant.NewInt(types.I32, int64(expression.Right.(lexer.Token).Value.(int))))
				} else if expression.Right.(lexer.Token).Type == lexer.Identifier {
					return block.NewSub(CompileExpression(block, expression.Left.(ast.ExpressionStatement), variables), variables[expression.Right.(lexer.Token).Value.(string)])
				}
			}
		case lexer.Token:
			switch expression.Right.(type) {
			case ast.ExpressionStatement:
				if expression.Left.(lexer.Token).Type == lexer.Int {
					return block.NewSub(constant.NewInt(types.I32, int64(expression.Left.(lexer.Token).Value.(int))), CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables))
				} else if expression.Left.(lexer.Token).Type == lexer.Identifier {
					return block.NewSub(variables[expression.Left.(lexer.Token).Value.(string)], CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables))
				}
			case lexer.Token:
				if expression.Left.(lexer.Token).Type == lexer.Int && expression.Right.(lexer.Token).Type == lexer.Int {
					return block.NewSub(constant.NewInt(types.I32, int64(expression.Left.(lexer.Token).Value.(int))), constant.NewInt(types.I32, int64(expression.Right.(lexer.Token).Value.(int))))
				} else if expression.Left.(lexer.Token).Type == lexer.Identifier && expression.Right.(lexer.Token).Type == lexer.Int {
					return block.NewSub(variables[expression.Left.(lexer.Token).Value.(string)], constant.NewInt(types.I32, int64(expression.Right.(lexer.Token).Value.(int))))
				} else if expression.Left.(lexer.Token).Type == lexer.Int && expression.Right.(lexer.Token).Type == lexer.Identifier {
					return block.NewSub(constant.NewInt(types.I32, int64(expression.Left.(lexer.Token).Value.(int))), variables[expression.Right.(lexer.Token).Value.(string)])
				} else if expression.Left.(lexer.Token).Type == lexer.Identifier && expression.Right.(lexer.Token).Type == lexer.Identifier {
					return block.NewSub(variables[expression.Left.(lexer.Token).Value.(string)], variables[expression.Right.(lexer.Token).Value.(string)])
				}
			}
		}
	} else if expression.Operator.Type == lexer.Multiply {
		switch expression.Left.(type) {
		case ast.ExpressionStatement:
			switch expression.Right.(type) {
			case ast.ExpressionStatement:
				return block.NewMul(CompileExpression(block, expression.Left.(ast.ExpressionStatement), variables), CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables))
			case lexer.Token:
				if expression.Right.(lexer.Token).Type == lexer.Int {
					return block.NewMul(CompileExpression(block, expression.Left.(ast.ExpressionStatement), variables), constant.NewInt(types.I32, int64(expression.Right.(lexer.Token).Value.(int))))
				} else if expression.Right.(lexer.Token).Type == lexer.Identifier {
					return block.NewMul(CompileExpression(block, expression.Left.(ast.ExpressionStatement), variables), variables[expression.Right.(lexer.Token).Value.(string)])
				}
			}
		case lexer.Token:
			switch expression.Right.(type) {
			case ast.ExpressionStatement:
				if expression.Left.(lexer.Token).Type == lexer.Int {
					return block.NewMul(constant.NewInt(types.I32, int64(expression.Left.(lexer.Token).Value.(int))), CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables))
				} else if expression.Left.(lexer.Token).Type == lexer.Identifier {
					return block.NewMul(variables[expression.Left.(lexer.Token).Value.(string)], CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables))
				}
			case lexer.Token:
				if expression.Left.(lexer.Token).Type == lexer.Int && expression.Right.(lexer.Token).Type == lexer.Int {
					return block.NewMul(constant.NewInt(types.I32, int64(expression.Left.(lexer.Token).Value.(int))), constant.NewInt(types.I32, int64(expression.Right.(lexer.Token).Value.(int))))
				} else if expression.Left.(lexer.Token).Type == lexer.Identifier && expression.Right.(lexer.Token).Type == lexer.Int {
					return block.NewMul(variables[expression.Left.(lexer.Token).Value.(string)], constant.NewInt(types.I32, int64(expression.Right.(lexer.Token).Value.(int))))
				} else if expression.Left.(lexer.Token).Type == lexer.Int && expression.Right.(lexer.Token).Type == lexer.Identifier {
					return block.NewMul(constant.NewInt(types.I32, int64(expression.Left.(lexer.Token).Value.(int))), variables[expression.Right.(lexer.Token).Value.(string)])
				} else if expression.Left.(lexer.Token).Type == lexer.Identifier && expression.Right.(lexer.Token).Type == lexer.Identifier {
					return block.NewMul(variables[expression.Left.(lexer.Token).Value.(string)], variables[expression.Right.(lexer.Token).Value.(string)])
				}
			}
		}
	} else if expression.Operator.Type == lexer.Divide {
		switch expression.Left.(type) {
		case ast.ExpressionStatement:
			switch expression.Right.(type) {
			case ast.ExpressionStatement:
				return block.NewSDiv(CompileExpression(block, expression.Left.(ast.ExpressionStatement), variables), CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables))
			case lexer.Token:
				if expression.Right.(lexer.Token).Type == lexer.Int {
					return block.NewSDiv(CompileExpression(block, expression.Left.(ast.ExpressionStatement), variables), constant.NewInt(types.I32, int64(expression.Right.(lexer.Token).Value.(int))))
				} else if expression.Right.(lexer.Token).Type == lexer.Identifier {
					return block.NewSDiv(CompileExpression(block, expression.Left.(ast.ExpressionStatement), variables), variables[expression.Right.(lexer.Token).Value.(string)])
				}
			}
		case lexer.Token:
			switch expression.Right.(type) {
			case ast.ExpressionStatement:
				if expression.Left.(lexer.Token).Type == lexer.Int {
					return block.NewSDiv(constant.NewInt(types.I32, int64(expression.Left.(lexer.Token).Value.(int))), CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables))
				} else if expression.Left.(lexer.Token).Type == lexer.Identifier {
					return block.NewSDiv(variables[expression.Left.(lexer.Token).Value.(string)], CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables))
				}
			case lexer.Token:
				if expression.Left.(lexer.Token).Type == lexer.Int && expression.Right.(lexer.Token).Type == lexer.Int {
					return block.NewSDiv(constant.NewInt(types.I32, int64(expression.Left.(lexer.Token).Value.(int))), constant.NewInt(types.I32, int64(expression.Right.(lexer.Token).Value.(int))))
				} else if expression.Left.(lexer.Token).Type == lexer.Identifier && expression.Right.(lexer.Token).Type == lexer.Int {
					return block.NewSDiv(variables[expression.Left.(lexer.Token).Value.(string)], constant.NewInt(types.I32, int64(expression.Right.(lexer.Token).Value.(int))))
				} else if expression.Left.(lexer.Token).Type == lexer.Int && expression.Right.(lexer.Token).Type == lexer.Identifier {
					return block.NewSDiv(constant.NewInt(types.I32, int64(expression.Left.(lexer.Token).Value.(int))), variables[expression.Right.(lexer.Token).Value.(string)])
				} else if expression.Left.(lexer.Token).Type == lexer.Identifier && expression.Right.(lexer.Token).Type == lexer.Identifier {
					return block.NewSDiv(variables[expression.Left.(lexer.Token).Value.(string)], variables[expression.Right.(lexer.Token).Value.(string)])
				}
			}
		}
	}

	return block
}
