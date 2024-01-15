package compiler

import (
	"Voca-2/ast"
	"Voca-2/lexer"
	"os"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

var stringType = types.NewArray(20, types.I8)

var boolType = types.I1
var Functions = make(map[string]Function)

type Variable struct {
	Type      string
	Value     value.Value
	hasstring bool
	length    int
}
type Function struct {
	Type  string
	Value value.Value
}

type Program struct {
	Input       string
	Errs        []error
	Tokens      []lexer.Token
	Program     ast.Program
	GenerateAST bool
	Args        []string
	File        string
	Output      string
	Arch        string
	OS          string
	LoadAST     bool
	Ir          bool
}

func New(program Program) []error {
	tmp, err := os.Create(program.File + ".ll")
	program.Errs = append(program.Errs, err)
	_, err = tmp.Write([]byte(GenerateIR(program.Program)))

	program.Errs = append(program.Errs, err)
	tmp.Close()
	if !program.Ir {
		err = os.Remove(program.File + ".ll")
		program.Errs = append(program.Errs, err)
	}
	return program.Errs
}
func GenerateIR(program ast.Program) string {

	var module = ir.NewModule()
	i := 0
	for i < len(program.Statements) {
		statement := program.Statements[i].Node
		switch statement.(type) {
		case ast.FuncDeclaration:

			var fn *ir.Func
			par := statement.(ast.FuncDeclaration).Arguments
			var params []*ir.Param
			for i := 0; i < len(par); i++ {
				if par[i].(ast.VariableDeclaration).Type.Value == "int" {
					params = append(params, ir.NewParam(par[i].(ast.VariableDeclaration).Name.Value.(string), types.I32))
				} else if par[i].(ast.VariableDeclaration).Type.Value == "string" {
					params = append(params, ir.NewParam(par[i].(ast.VariableDeclaration).Name.Value.(string), types.I8Ptr))
				}
			}
			if statement.(ast.FuncDeclaration).Type.Value.(string) == "int" {
				fn = module.NewFunc(statement.(ast.FuncDeclaration).Name.Value.(string), types.I32)
			} else if statement.(ast.FuncDeclaration).Type.Value.(string) == "string" {
				fn = module.NewFunc(statement.(ast.FuncDeclaration).Name.Value.(string), types.I8Ptr)
			} else if statement.(ast.FuncDeclaration).Type.Value.(string) == "void" {
				fn = module.NewFunc(statement.(ast.FuncDeclaration).Name.Value.(string), types.Void)
			}
			fn.Params = params
			entry := fn.NewBlock("entry")
			function := Function{Type: statement.(ast.FuncDeclaration).Type.Value.(string), Value: fn}
			Functions[statement.(ast.FuncDeclaration).Name.Value.(string)] = function
			entry = Compile(entry, statement.(ast.FuncDeclaration).Body, make(map[string]Variable))
			entry.NewRet(constant.NewInt(types.I32, 0))

		}
		i++
	}

	return module.String()
}
func Compile(block *ir.Block, statements []ast.Statement, variables map[string]Variable) *ir.Block {
	i := 0
	for i < len(statements) {
		statement := statements[i].Node
		variab := Variable{}
		hasstring := false
		isstring := false
		length := 0
		switch statement.(type) {
		case ast.VariableDeclaration:
			var variable value.Value
			if statement.(ast.VariableDeclaration).Type.Value.(string) == "string" {
				isstring = true
				variab.Type = "string"

			} else if statement.(ast.VariableDeclaration).Type.Value.(string) == "bool" {
				variable = block.NewAlloca(boolType)
				variable.(*ir.InstAlloca).SetName(statement.(ast.VariableDeclaration).Name.Value.(string))
				variab.Type = "bool"
			} else if statement.(ast.VariableDeclaration).Type.Value.(string) == "int" {
				variable = block.NewAlloca(types.I32)
				variable.(*ir.InstAlloca).SetName(statement.(ast.VariableDeclaration).Name.Value.(string))
				variab.Type = "int"
			}
			if statement.(ast.VariableDeclaration).Value != nil {
				switch statement.(ast.VariableDeclaration).Value.(type) {
				case ast.ExpressionStatement:
					if isstring {
						hasstring = true

						str := CompileExpression(block, statement.(ast.VariableDeclaration).Value.(ast.ExpressionStatement), variables)
						switch str.(type) {
						case *constant.CharArray:
							variable = block.NewAlloca(types.NewArray(str.(*constant.CharArray).Typ.Len, types.I8))
							variable.(*ir.InstAlloca).SetName(statement.(ast.VariableDeclaration).Name.Value.(string))
							length = int(str.(*constant.CharArray).Typ.Len)
							block.NewStore(str, variable)
						case *ir.InstCall:
							variable = str
							variable.(*ir.InstCall).SetName(statement.(ast.VariableDeclaration).Name.Value.(string))

						}

					} else {
						block.NewStore(CompileExpression(block, statement.(ast.VariableDeclaration).Value.(ast.ExpressionStatement), variables), variable)
					}
				case lexer.Token:
					if statement.(ast.VariableDeclaration).Value.(lexer.Token).Type == lexer.Int {
						block.NewStore(constant.NewInt(types.I32, statement.(ast.VariableDeclaration).Value.(lexer.Token).Value.(int64)), variable)
					} else if statement.(ast.VariableDeclaration).Value.(lexer.Token).Type == lexer.String {
						hasstring = true

						str := constant.NewCharArrayFromString(statement.(ast.VariableDeclaration).Value.(lexer.Token).Value.(string))
						variable = block.NewAlloca(types.NewArray(str.Typ.Len, types.I8))
						variable.(*ir.InstAlloca).SetName(statement.(ast.VariableDeclaration).Name.Value.(string))
						length = int(str.Typ.Len)
						block.NewStore(str, variable)
					} else if statement.(ast.VariableDeclaration).Value.(lexer.Token).Type == lexer.Identifier {
						if isstring {
							hasstring = true
							value := variables[statement.(ast.VariableDeclaration).Value.(lexer.Token).Value.(string)].Value.(*constant.CharArray)
							variable = block.NewAlloca(types.NewArray(value.Typ.Len, types.I8))
							variable.(*ir.InstAlloca).SetName(statement.(ast.VariableDeclaration).Name.Value.(string))
							length = int(value.Typ.Len)
							block.NewStore(value, variable)
						} else {
							block.NewStore(variables[statement.(ast.VariableDeclaration).Value.(lexer.Token).Value.(string)].Value, variable)
						}
					}
				case ast.FuncCall:
					name := statement.(ast.VariableDeclaration).Value.(ast.FuncCall).Name.Value.(string)
					args := statement.(ast.VariableDeclaration).Value.(ast.FuncCall).Arguments
					var arguments []value.Value
					for _, arg := range args {
						switch arg.(type) {
						case ast.ExpressionStatement:
							arguments = append(arguments, CompileExpression(block, arg.(ast.ExpressionStatement), variables))
						case lexer.Token:
							if arg.(lexer.Token).Type == lexer.Int {
								arguments = append(arguments, constant.NewInt(types.I32, arg.(lexer.Token).Value.(int64)))
							} else if arg.(lexer.Token).Type == lexer.Identifier {
								arguments = append(arguments, variables[arg.(lexer.Token).Value.(string)].Value)
							} else if arg.(lexer.Token).Type == lexer.String {
								str := constant.NewCharArrayFromString(arg.(lexer.Token).Value.(string))
								arguments = append(arguments, str)
							}
						}
					}
					if Functions[name].Type == "int" {
						block.NewStore(block.NewCall(Functions[name].Value, arguments...), variable)
					} else if Functions[name].Type == "string" {
						hasstring = true

						str := block.NewCall(Functions[name].Value, arguments...)

						variable = str
						variable.(*ir.InstCall).SetName(statement.(ast.VariableDeclaration).Name.Value.(string))

						block.NewStore(str, variable)
					}
				}
				if isstring && !hasstring {
					variab.hasstring = false
				} else if isstring && hasstring {
					variab.hasstring = true
					variab.Value = variable
					variab.length = length
				} else {

					variab.Value = variable
				}

				//block.NewStore(variables[statement.(ast.VariableDeclaration).Name.Value.(string)], variable)
			}
			variables[statement.(ast.VariableDeclaration).Name.Value.(string)] = variab
		case ast.FuncCall:
			name := statement.(ast.FuncCall).Name.Value.(string)
			args := statement.(ast.FuncCall).Arguments
			var arguments []value.Value
			for _, arg := range args {
				switch arg.(type) {
				case ast.ExpressionStatement:
					arguments = append(arguments, CompileExpression(block, arg.(ast.ExpressionStatement), variables))
				case lexer.Token:
					if arg.(lexer.Token).Type == lexer.Int {
						arguments = append(arguments, constant.NewInt(types.I32, arg.(lexer.Token).Value.(int64)))
					} else if arg.(lexer.Token).Type == lexer.Identifier {
						arguments = append(arguments, variables[arg.(lexer.Token).Value.(string)].Value)
					} else if arg.(lexer.Token).Type == lexer.String {
						str := constant.NewCharArrayFromString(arg.(lexer.Token).Value.(string))
						arguments = append(arguments, str)
					}
				}
			}
			block.NewCall(Functions[name].Value, arguments...)

		}
		i++
	}
	return block
}

func CompileExpression(block *ir.Block, expression ast.ExpressionStatement, variables map[string]Variable) value.Value {
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
					return block.NewAdd(CompileExpression(block, expression.Left.(ast.ExpressionStatement), variables), variables[expression.Right.(lexer.Token).Value.(string)].Value)
				}
			case ast.FuncCall:
				name := expression.Right.(ast.FuncCall).Name.Value.(string)
				args := expression.Right.(ast.FuncCall).Arguments
				var arguments []value.Value
				for _, arg := range args {
					switch arg.(type) {
					case ast.ExpressionStatement:
						arguments = append(arguments, CompileExpression(block, arg.(ast.ExpressionStatement), variables))
					case lexer.Token:
						if arg.(lexer.Token).Type == lexer.Int {
							arguments = append(arguments, constant.NewInt(types.I32, arg.(lexer.Token).Value.(int64)))
						} else if arg.(lexer.Token).Type == lexer.Identifier {
							arguments = append(arguments, variables[arg.(lexer.Token).Value.(string)].Value)
						} else if arg.(lexer.Token).Type == lexer.String {
							str := constant.NewCharArrayFromString(arg.(lexer.Token).Value.(string))
							arguments = append(arguments, str)
						}
					}
				}
				if Functions[name].Type == "int" {
					return block.NewAdd(CompileExpression(block, expression.Left.(ast.ExpressionStatement), variables), block.NewCall(Functions[name].Value, arguments...))
				}
			}
		case lexer.Token:
			switch expression.Right.(type) {
			case ast.ExpressionStatement:
				if expression.Left.(lexer.Token).Type == lexer.Int {

					return block.NewAdd(constant.NewInt(types.I32, int64(expression.Left.(lexer.Token).Value.(int))), CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables))
				} else if expression.Left.(lexer.Token).Type == lexer.Identifier {
					if variables[expression.Left.(lexer.Token).Value.(string)].Type == "int" {
						return block.NewAdd(variables[expression.Left.(lexer.Token).Value.(string)].Value, CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables))
					} else if variables[expression.Left.(lexer.Token).Value.(string)].Type == "string" {
						str1 := variables[expression.Left.(lexer.Token).Value.(string)].Value.(*constant.CharArray)
						str2 := CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables)
						str := constant.NewCharArray(append(str1.X, str2.(*constant.CharArray).X...))
						return str
					}
				} else if expression.Left.(lexer.Token).Type == lexer.String {
					str1 := constant.NewCharArrayFromString(expression.Left.(lexer.Token).Value.(string))
					str2 := CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables).(*constant.CharArray)
					str := constant.NewCharArray(append(str1.X, str2.X...))
					return str

				}
			case lexer.Token:
				if expression.Left.(lexer.Token).Type == lexer.Int && expression.Right.(lexer.Token).Type == lexer.Int {
					return block.NewAdd(constant.NewInt(types.I32, int64(expression.Left.(lexer.Token).Value.(int))), constant.NewInt(types.I32, int64(expression.Right.(lexer.Token).Value.(int))))
				} else if expression.Left.(lexer.Token).Type == lexer.Identifier && expression.Right.(lexer.Token).Type == lexer.Int {
					return block.NewAdd(variables[expression.Left.(lexer.Token).Value.(string)].Value, constant.NewInt(types.I32, int64(expression.Right.(lexer.Token).Value.(int))))
				} else if expression.Left.(lexer.Token).Type == lexer.Int && expression.Right.(lexer.Token).Type == lexer.Identifier {
					return block.NewAdd(constant.NewInt(types.I32, int64(expression.Left.(lexer.Token).Value.(int))), variables[expression.Right.(lexer.Token).Value.(string)].Value)
				} else if expression.Left.(lexer.Token).Type == lexer.Identifier && expression.Right.(lexer.Token).Type == lexer.Identifier {
					if (variables[expression.Left.(lexer.Token).Value.(string)].Type == "int") && (variables[expression.Right.(lexer.Token).Value.(string)].Type == "int") {
						return block.NewAdd(variables[expression.Left.(lexer.Token).Value.(string)].Value, variables[expression.Right.(lexer.Token).Value.(string)].Value)
					} else if (variables[expression.Left.(lexer.Token).Value.(string)].Type == "string") && (variables[expression.Right.(lexer.Token).Value.(string)].Type == "string") {

						return block.NewCall(Functions["append"].Value, block.NewGetElementPtr(types.I8, variables[expression.Left.(lexer.Token).Value.(string)].Value, constant.NewInt(types.I32, 0)), block.NewGetElementPtr(types.I8, variables[expression.Right.(lexer.Token).Value.(string)].Value, constant.NewInt(types.I32, 0)))
					}
				} else if expression.Left.(lexer.Token).Type == lexer.String && expression.Right.(lexer.Token).Type == lexer.String {
					str1 := constant.NewCharArrayFromString(expression.Left.(lexer.Token).Value.(string))
					str2 := constant.NewCharArrayFromString(expression.Right.(lexer.Token).Value.(string))
					str := constant.NewCharArray(append(str1.X, str2.X...))
					return str
				} else if expression.Left.(lexer.Token).Type == lexer.String && expression.Right.(lexer.Token).Type == lexer.Identifier {
					str1 := constant.NewCharArrayFromString(expression.Left.(lexer.Token).Value.(string))
					str2 := variables[expression.Right.(lexer.Token).Value.(string)].Value.(*constant.CharArray)
					str := constant.NewCharArray(append(str1.X, str2.X...))
					return str
				} else if expression.Left.(lexer.Token).Type == lexer.Identifier && expression.Right.(lexer.Token).Type == lexer.String {
					str1 := variables[expression.Left.(lexer.Token).Value.(string)].Value.(*constant.CharArray)
					str2 := constant.NewCharArrayFromString(expression.Right.(lexer.Token).Value.(string))
					str := constant.NewCharArray(append(str1.X, str2.X...))
					return str
				}
			case ast.FuncCall:
				name2 := expression.Right.(ast.FuncCall).Name.Value.(string)
				args2 := expression.Right.(ast.FuncCall).Arguments
				var arguments2 []value.Value
				for _, arg := range args2 {
					switch arg.(type) {
					case ast.ExpressionStatement:
						arguments2 = append(arguments2, CompileExpression(block, arg.(ast.ExpressionStatement), variables))
					case lexer.Token:
						if arg.(lexer.Token).Type == lexer.Int {
							arguments2 = append(arguments2, constant.NewInt(types.I32, arg.(lexer.Token).Value.(int64)))
						} else if arg.(lexer.Token).Type == lexer.Identifier {
							arguments2 = append(arguments2, variables[arg.(lexer.Token).Value.(string)].Value)
						} else if arg.(lexer.Token).Type == lexer.String {
							str := constant.NewCharArrayFromString(arg.(lexer.Token).Value.(string))
							arguments2 = append(arguments2, str)
						}
					}
				}
				if Functions[name2].Type == "int" && expression.Left.(lexer.Token).Type == lexer.Int {
					return block.NewAdd(constant.NewInt(types.I32, int64(expression.Left.(lexer.Token).Value.(int))), block.NewCall(Functions[name2].Value, arguments2...))
				} else if Functions[name2].Type == "int" && expression.Left.(lexer.Token).Type == lexer.Identifier {
					return block.NewAdd(variables[expression.Left.(lexer.Token).Value.(string)].Value, block.NewCall(Functions[name2].Value, arguments2...))
				}

			}
		case ast.FuncCall:
			name := expression.Left.(ast.FuncCall).Name.Value.(string)
			args := expression.Left.(ast.FuncCall).Arguments
			var arguments []value.Value
			for _, arg := range args {
				switch arg.(type) {
				case ast.ExpressionStatement:
					arguments = append(arguments, CompileExpression(block, arg.(ast.ExpressionStatement), variables))
				case lexer.Token:
					if arg.(lexer.Token).Type == lexer.Int {
						arguments = append(arguments, constant.NewInt(types.I32, arg.(lexer.Token).Value.(int64)))
					} else if arg.(lexer.Token).Type == lexer.Identifier {
						arguments = append(arguments, variables[arg.(lexer.Token).Value.(string)].Value)
					} else if arg.(lexer.Token).Type == lexer.String {
						str := constant.NewCharArrayFromString(arg.(lexer.Token).Value.(string))
						arguments = append(arguments, str)
					}
				}
			}
			switch expression.Right.(type) {
			case ast.ExpressionStatement:
				if Functions[name].Type == "int" {
					return block.NewAdd(block.NewCall(Functions[name].Value, arguments...), CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables))
				}
			case lexer.Token:
				if expression.Right.(lexer.Token).Type == lexer.Int {
					return block.NewAdd(block.NewCall(Functions[name].Value, arguments...), constant.NewInt(types.I32, int64(expression.Right.(lexer.Token).Value.(int))))
				} else if expression.Right.(lexer.Token).Type == lexer.Identifier {
					return block.NewAdd(block.NewCall(Functions[name].Value, arguments...), variables[expression.Right.(lexer.Token).Value.(string)].Value)
				} else if expression.Right.(lexer.Token).Type == lexer.String {
					str := constant.NewCharArrayFromString(expression.Right.(lexer.Token).Value.(string))
					return block.NewCall(Functions["append"].Value, block.NewCall(Functions[name].Value, arguments...), str)
				}
			case ast.FuncCall:
				name2 := expression.Right.(ast.FuncCall).Name.Value.(string)
				args2 := expression.Right.(ast.FuncCall).Arguments
				var arguments2 []value.Value
				for _, arg := range args2 {
					switch arg.(type) {
					case ast.ExpressionStatement:
						arguments2 = append(arguments2, CompileExpression(block, arg.(ast.ExpressionStatement), variables))
					case lexer.Token:
						if arg.(lexer.Token).Type == lexer.Int {
							arguments2 = append(arguments2, constant.NewInt(types.I32, arg.(lexer.Token).Value.(int64)))
						} else if arg.(lexer.Token).Type == lexer.Identifier {
							arguments2 = append(arguments2, variables[arg.(lexer.Token).Value.(string)].Value)
						} else if arg.(lexer.Token).Type == lexer.String {
							str := constant.NewCharArrayFromString(arg.(lexer.Token).Value.(string))
							arguments2 = append(arguments2, str)
						}
					}
				}
				if Functions[name].Type == "int" && Functions[name2].Type == "int" {
					return block.NewAdd(block.NewCall(Functions[name].Value, arguments...), block.NewCall(Functions[name2].Value, arguments2...))
				} else if Functions[name].Type == "string" && Functions[name2].Type == "string" {
					return block.NewCall(Functions["append"].Value, block.NewCall(Functions[name].Value, arguments...), block.NewCall(Functions[name2].Value, arguments2...))
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
					return block.NewSub(CompileExpression(block, expression.Left.(ast.ExpressionStatement), variables), variables[expression.Right.(lexer.Token).Value.(string)].Value)
				}
			case ast.FuncCall:
				name := expression.Right.(ast.FuncCall).Name.Value.(string)
				args := expression.Right.(ast.FuncCall).Arguments
				var arguments []value.Value
				for _, arg := range args {
					switch arg.(type) {
					case ast.ExpressionStatement:
						arguments = append(arguments, CompileExpression(block, arg.(ast.ExpressionStatement), variables))
					case lexer.Token:
						if arg.(lexer.Token).Type == lexer.Int {
							arguments = append(arguments, constant.NewInt(types.I32, arg.(lexer.Token).Value.(int64)))
						} else if arg.(lexer.Token).Type == lexer.Identifier {
							arguments = append(arguments, variables[arg.(lexer.Token).Value.(string)].Value)
						} else if arg.(lexer.Token).Type == lexer.String {
							str := constant.NewCharArrayFromString(arg.(lexer.Token).Value.(string))
							arguments = append(arguments, str)
						}
					}
				}
				if Functions[name].Type == "int" {
					return block.NewSub(CompileExpression(block, expression.Left.(ast.ExpressionStatement), variables), block.NewCall(Functions[name].Value, arguments...))
				}
			}
		case lexer.Token:
			switch expression.Right.(type) {
			case ast.ExpressionStatement:
				if expression.Left.(lexer.Token).Type == lexer.Int {
					return block.NewSub(constant.NewInt(types.I32, int64(expression.Left.(lexer.Token).Value.(int))), CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables))
				} else if expression.Left.(lexer.Token).Type == lexer.Identifier {
					return block.NewSub(variables[expression.Left.(lexer.Token).Value.(string)].Value, CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables))
				}
			case lexer.Token:
				if expression.Left.(lexer.Token).Type == lexer.Int && expression.Right.(lexer.Token).Type == lexer.Int {
					return block.NewSub(constant.NewInt(types.I32, int64(expression.Left.(lexer.Token).Value.(int))), constant.NewInt(types.I32, int64(expression.Right.(lexer.Token).Value.(int))))
				} else if expression.Left.(lexer.Token).Type == lexer.Identifier && expression.Right.(lexer.Token).Type == lexer.Int {
					return block.NewSub(variables[expression.Left.(lexer.Token).Value.(string)].Value, constant.NewInt(types.I32, int64(expression.Right.(lexer.Token).Value.(int))))
				} else if expression.Left.(lexer.Token).Type == lexer.Int && expression.Right.(lexer.Token).Type == lexer.Identifier {
					return block.NewSub(constant.NewInt(types.I32, int64(expression.Left.(lexer.Token).Value.(int))), variables[expression.Right.(lexer.Token).Value.(string)].Value)
				} else if expression.Left.(lexer.Token).Type == lexer.Identifier && expression.Right.(lexer.Token).Type == lexer.Identifier {
					return block.NewSub(variables[expression.Left.(lexer.Token).Value.(string)].Value, variables[expression.Right.(lexer.Token).Value.(string)].Value)
				}
			case ast.FuncCall:
				name2 := expression.Right.(ast.FuncCall).Name.Value.(string)
				args2 := expression.Right.(ast.FuncCall).Arguments
				var arguments2 []value.Value
				for _, arg := range args2 {
					switch arg.(type) {
					case ast.ExpressionStatement:
						arguments2 = append(arguments2, CompileExpression(block, arg.(ast.ExpressionStatement), variables))
					case lexer.Token:
						if arg.(lexer.Token).Type == lexer.Int {
							arguments2 = append(arguments2, constant.NewInt(types.I32, arg.(lexer.Token).Value.(int64)))
						} else if arg.(lexer.Token).Type == lexer.Identifier {
							arguments2 = append(arguments2, variables[arg.(lexer.Token).Value.(string)].Value)
						} else if arg.(lexer.Token).Type == lexer.String {
							str := constant.NewCharArrayFromString(arg.(lexer.Token).Value.(string))
							arguments2 = append(arguments2, str)
						}
					}
				}
				if Functions[name2].Type == "int" && expression.Left.(lexer.Token).Type == lexer.Int {
					return block.NewSub(constant.NewInt(types.I32, int64(expression.Left.(lexer.Token).Value.(int))), block.NewCall(Functions[name2].Value, arguments2...))
				} else if Functions[name2].Type == "int" && expression.Left.(lexer.Token).Type == lexer.Identifier {
					return block.NewSub(variables[expression.Left.(lexer.Token).Value.(string)].Value, block.NewCall(Functions[name2].Value, arguments2...))
				}
			}
		case ast.FuncCall:
			name := expression.Left.(ast.FuncCall).Name.Value.(string)
			args := expression.Left.(ast.FuncCall).Arguments
			var arguments []value.Value
			for _, arg := range args {
				switch arg.(type) {
				case ast.ExpressionStatement:
					arguments = append(arguments, CompileExpression(block, arg.(ast.ExpressionStatement), variables))
				case lexer.Token:
					if arg.(lexer.Token).Type == lexer.Int {
						arguments = append(arguments, constant.NewInt(types.I32, arg.(lexer.Token).Value.(int64)))
					} else if arg.(lexer.Token).Type == lexer.Identifier {
						arguments = append(arguments, variables[arg.(lexer.Token).Value.(string)].Value)
					} else if arg.(lexer.Token).Type == lexer.String {
						str := constant.NewCharArrayFromString(arg.(lexer.Token).Value.(string))
						arguments = append(arguments, str)
					}
				}
			}
			switch expression.Right.(type) {
			case ast.ExpressionStatement:
				if Functions[name].Type == "int" {
					return block.NewSub(block.NewCall(Functions[name].Value, arguments...), CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables))
				}
			case lexer.Token:
				if expression.Right.(lexer.Token).Type == lexer.Int {
					return block.NewSub(block.NewCall(Functions[name].Value, arguments...), constant.NewInt(types.I32, int64(expression.Right.(lexer.Token).Value.(int))))
				} else if expression.Right.(lexer.Token).Type == lexer.Identifier {
					return block.NewSub(block.NewCall(Functions[name].Value, arguments...), variables[expression.Right.(lexer.Token).Value.(string)].Value)
				}
			case ast.FuncCall:
				name2 := expression.Right.(ast.FuncCall).Name.Value.(string)
				args2 := expression.Right.(ast.FuncCall).Arguments
				var arguments2 []value.Value
				for _, arg := range args2 {
					switch arg.(type) {
					case ast.ExpressionStatement:
						arguments2 = append(arguments2, CompileExpression(block, arg.(ast.ExpressionStatement), variables))
					case lexer.Token:
						if arg.(lexer.Token).Type == lexer.Int {
							arguments2 = append(arguments2, constant.NewInt(types.I32, arg.(lexer.Token).Value.(int64)))
						} else if arg.(lexer.Token).Type == lexer.Identifier {
							arguments2 = append(arguments2, variables[arg.(lexer.Token).Value.(string)].Value)
						} else if arg.(lexer.Token).Type == lexer.String {
							str := constant.NewCharArrayFromString(arg.(lexer.Token).Value.(string))
							arguments2 = append(arguments2, str)
						}
					}
				}
				if Functions[name].Type == "int" && Functions[name2].Type == "int" {
					return block.NewSub(block.NewCall(Functions[name].Value, arguments...), block.NewCall(Functions[name2].Value, arguments2...))
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
					return block.NewMul(CompileExpression(block, expression.Left.(ast.ExpressionStatement), variables), variables[expression.Right.(lexer.Token).Value.(string)].Value)
				}
			case ast.FuncCall:
				name := expression.Right.(ast.FuncCall).Name.Value.(string)
				args := expression.Right.(ast.FuncCall).Arguments
				var arguments []value.Value
				for _, arg := range args {
					switch arg.(type) {
					case ast.ExpressionStatement:
						arguments = append(arguments, CompileExpression(block, arg.(ast.ExpressionStatement), variables))
					case lexer.Token:
						if arg.(lexer.Token).Type == lexer.Int {
							arguments = append(arguments, constant.NewInt(types.I32, arg.(lexer.Token).Value.(int64)))
						} else if arg.(lexer.Token).Type == lexer.Identifier {
							arguments = append(arguments, variables[arg.(lexer.Token).Value.(string)].Value)
						} else if arg.(lexer.Token).Type == lexer.String {
							str := constant.NewCharArrayFromString(arg.(lexer.Token).Value.(string))
							arguments = append(arguments, str)
						}
					}
				}
				if Functions[name].Type == "int" {
					return block.NewMul(CompileExpression(block, expression.Left.(ast.ExpressionStatement), variables), block.NewCall(Functions[name].Value, arguments...))
				}
			}
		case lexer.Token:
			switch expression.Right.(type) {
			case ast.ExpressionStatement:
				if expression.Left.(lexer.Token).Type == lexer.Int {
					return block.NewMul(constant.NewInt(types.I32, int64(expression.Left.(lexer.Token).Value.(int))), CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables))
				} else if expression.Left.(lexer.Token).Type == lexer.Identifier {
					return block.NewMul(variables[expression.Left.(lexer.Token).Value.(string)].Value, CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables))
				}
			case lexer.Token:
				if expression.Left.(lexer.Token).Type == lexer.Int && expression.Right.(lexer.Token).Type == lexer.Int {
					return block.NewMul(constant.NewInt(types.I32, int64(expression.Left.(lexer.Token).Value.(int))), constant.NewInt(types.I32, int64(expression.Right.(lexer.Token).Value.(int))))
				} else if expression.Left.(lexer.Token).Type == lexer.Identifier && expression.Right.(lexer.Token).Type == lexer.Int {
					return block.NewMul(variables[expression.Left.(lexer.Token).Value.(string)].Value, constant.NewInt(types.I32, int64(expression.Right.(lexer.Token).Value.(int))))
				} else if expression.Left.(lexer.Token).Type == lexer.Int && expression.Right.(lexer.Token).Type == lexer.Identifier {
					return block.NewMul(constant.NewInt(types.I32, int64(expression.Left.(lexer.Token).Value.(int))), variables[expression.Right.(lexer.Token).Value.(string)].Value)
				} else if expression.Left.(lexer.Token).Type == lexer.Identifier && expression.Right.(lexer.Token).Type == lexer.Identifier {
					return block.NewMul(variables[expression.Left.(lexer.Token).Value.(string)].Value, variables[expression.Right.(lexer.Token).Value.(string)].Value)
				}
			case ast.FuncCall:
				name2 := expression.Right.(ast.FuncCall).Name.Value.(string)
				args2 := expression.Right.(ast.FuncCall).Arguments
				var arguments2 []value.Value
				for _, arg := range args2 {
					switch arg.(type) {
					case ast.ExpressionStatement:
						arguments2 = append(arguments2, CompileExpression(block, arg.(ast.ExpressionStatement), variables))
					case lexer.Token:
						if arg.(lexer.Token).Type == lexer.Int {
							arguments2 = append(arguments2, constant.NewInt(types.I32, arg.(lexer.Token).Value.(int64)))
						} else if arg.(lexer.Token).Type == lexer.Identifier {
							arguments2 = append(arguments2, variables[arg.(lexer.Token).Value.(string)].Value)
						} else if arg.(lexer.Token).Type == lexer.String {
							str := constant.NewCharArrayFromString(arg.(lexer.Token).Value.(string))
							arguments2 = append(arguments2, str)
						}
					}
				}
				if Functions[name2].Type == "int" && expression.Left.(lexer.Token).Type == lexer.Int {
					return block.NewMul(constant.NewInt(types.I32, int64(expression.Left.(lexer.Token).Value.(int))), block.NewCall(Functions[name2].Value, arguments2...))
				} else if Functions[name2].Type == "int" && expression.Left.(lexer.Token).Type == lexer.Identifier {
					return block.NewMul(variables[expression.Left.(lexer.Token).Value.(string)].Value, block.NewCall(Functions[name2].Value, arguments2...))
				}
			}
		case ast.FuncCall:
			name := expression.Left.(ast.FuncCall).Name.Value.(string)
			args := expression.Left.(ast.FuncCall).Arguments
			var arguments []value.Value
			for _, arg := range args {
				switch arg.(type) {
				case ast.ExpressionStatement:
					arguments = append(arguments, CompileExpression(block, arg.(ast.ExpressionStatement), variables))
				case lexer.Token:
					if arg.(lexer.Token).Type == lexer.Int {
						arguments = append(arguments, constant.NewInt(types.I32, arg.(lexer.Token).Value.(int64)))
					} else if arg.(lexer.Token).Type == lexer.Identifier {
						arguments = append(arguments, variables[arg.(lexer.Token).Value.(string)].Value)
					} else if arg.(lexer.Token).Type == lexer.String {
						str := constant.NewCharArrayFromString(arg.(lexer.Token).Value.(string))
						arguments = append(arguments, str)
					}
				}
			}
			switch expression.Right.(type) {
			case ast.ExpressionStatement:
				if Functions[name].Type == "int" {
					return block.NewMul(block.NewCall(Functions[name].Value, arguments...), CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables))
				}
			case lexer.Token:
				if expression.Right.(lexer.Token).Type == lexer.Int {
					return block.NewMul(block.NewCall(Functions[name].Value, arguments...), constant.NewInt(types.I32, int64(expression.Right.(lexer.Token).Value.(int))))
				} else if expression.Right.(lexer.Token).Type == lexer.Identifier {
					return block.NewMul(block.NewCall(Functions[name].Value, arguments...), variables[expression.Right.(lexer.Token).Value.(string)].Value)
				}
			case ast.FuncCall:
				name2 := expression.Right.(ast.FuncCall).Name.Value.(string)
				args2 := expression.Right.(ast.FuncCall).Arguments
				var arguments2 []value.Value
				for _, arg := range args2 {
					switch arg.(type) {
					case ast.ExpressionStatement:
						arguments2 = append(arguments2, CompileExpression(block, arg.(ast.ExpressionStatement), variables))
					case lexer.Token:
						if arg.(lexer.Token).Type == lexer.Int {
							arguments2 = append(arguments2, constant.NewInt(types.I32, arg.(lexer.Token).Value.(int64)))
						} else if arg.(lexer.Token).Type == lexer.Identifier {
							arguments2 = append(arguments2, variables[arg.(lexer.Token).Value.(string)].Value)
						} else if arg.(lexer.Token).Type == lexer.String {
							str := constant.NewCharArrayFromString(arg.(lexer.Token).Value.(string))
							arguments2 = append(arguments2, str)
						}
					}
				}
				if Functions[name].Type == "int" && Functions[name2].Type == "int" {
					return block.NewMul(block.NewCall(Functions[name].Value, arguments...), block.NewCall(Functions[name2].Value, arguments2...))
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
					return block.NewSDiv(CompileExpression(block, expression.Left.(ast.ExpressionStatement), variables), variables[expression.Right.(lexer.Token).Value.(string)].Value)
				}
			case ast.FuncCall:
				name := expression.Right.(ast.FuncCall).Name.Value.(string)
				args := expression.Right.(ast.FuncCall).Arguments
				var arguments []value.Value
				for _, arg := range args {
					switch arg.(type) {
					case ast.ExpressionStatement:
						arguments = append(arguments, CompileExpression(block, arg.(ast.ExpressionStatement), variables))
					case lexer.Token:
						if arg.(lexer.Token).Type == lexer.Int {
							arguments = append(arguments, constant.NewInt(types.I32, arg.(lexer.Token).Value.(int64)))
						} else if arg.(lexer.Token).Type == lexer.Identifier {
							arguments = append(arguments, variables[arg.(lexer.Token).Value.(string)].Value)
						} else if arg.(lexer.Token).Type == lexer.String {
							str := constant.NewCharArrayFromString(arg.(lexer.Token).Value.(string))
							arguments = append(arguments, str)
						}
					}
				}
				if Functions[name].Type == "int" {
					return block.NewSDiv(CompileExpression(block, expression.Left.(ast.ExpressionStatement), variables), block.NewCall(Functions[name].Value, arguments...))
				}
			}
		case lexer.Token:
			switch expression.Right.(type) {
			case ast.ExpressionStatement:
				if expression.Left.(lexer.Token).Type == lexer.Int {
					return block.NewSDiv(constant.NewInt(types.I32, int64(expression.Left.(lexer.Token).Value.(int))), CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables))
				} else if expression.Left.(lexer.Token).Type == lexer.Identifier {
					return block.NewSDiv(variables[expression.Left.(lexer.Token).Value.(string)].Value, CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables))
				}
			case lexer.Token:
				if expression.Left.(lexer.Token).Type == lexer.Int && expression.Right.(lexer.Token).Type == lexer.Int {
					return block.NewSDiv(constant.NewInt(types.I32, int64(expression.Left.(lexer.Token).Value.(int))), constant.NewInt(types.I32, int64(expression.Right.(lexer.Token).Value.(int))))
				} else if expression.Left.(lexer.Token).Type == lexer.Identifier && expression.Right.(lexer.Token).Type == lexer.Int {
					return block.NewSDiv(variables[expression.Left.(lexer.Token).Value.(string)].Value, constant.NewInt(types.I32, int64(expression.Right.(lexer.Token).Value.(int))))
				} else if expression.Left.(lexer.Token).Type == lexer.Int && expression.Right.(lexer.Token).Type == lexer.Identifier {
					return block.NewSDiv(constant.NewInt(types.I32, int64(expression.Left.(lexer.Token).Value.(int))), variables[expression.Right.(lexer.Token).Value.(string)].Value)
				} else if expression.Left.(lexer.Token).Type == lexer.Identifier && expression.Right.(lexer.Token).Type == lexer.Identifier {
					return block.NewSDiv(variables[expression.Left.(lexer.Token).Value.(string)].Value, variables[expression.Right.(lexer.Token).Value.(string)].Value)
				}
			case ast.FuncCall:
				name2 := expression.Right.(ast.FuncCall).Name.Value.(string)
				args2 := expression.Right.(ast.FuncCall).Arguments
				var arguments2 []value.Value
				for _, arg := range args2 {
					switch arg.(type) {
					case ast.ExpressionStatement:
						arguments2 = append(arguments2, CompileExpression(block, arg.(ast.ExpressionStatement), variables))
					case lexer.Token:
						if arg.(lexer.Token).Type == lexer.Int {
							arguments2 = append(arguments2, constant.NewInt(types.I32, arg.(lexer.Token).Value.(int64)))
						} else if arg.(lexer.Token).Type == lexer.Identifier {
							arguments2 = append(arguments2, variables[arg.(lexer.Token).Value.(string)].Value)
						} else if arg.(lexer.Token).Type == lexer.String {
							str := constant.NewCharArrayFromString(arg.(lexer.Token).Value.(string))
							arguments2 = append(arguments2, str)
						}
					}
				}
				if Functions[name2].Type == "int" && expression.Left.(lexer.Token).Type == lexer.Int {
					return block.NewSDiv(constant.NewInt(types.I32, int64(expression.Left.(lexer.Token).Value.(int))), block.NewCall(Functions[name2].Value, arguments2...))
				} else if Functions[name2].Type == "int" && expression.Left.(lexer.Token).Type == lexer.Identifier {
					return block.NewSDiv(variables[expression.Left.(lexer.Token).Value.(string)].Value, block.NewCall(Functions[name2].Value, arguments2...))
				}
			}
		case ast.FuncCall:
			name := expression.Left.(ast.FuncCall).Name.Value.(string)
			args := expression.Left.(ast.FuncCall).Arguments
			var arguments []value.Value
			for _, arg := range args {
				switch arg.(type) {
				case ast.ExpressionStatement:
					arguments = append(arguments, CompileExpression(block, arg.(ast.ExpressionStatement), variables))
				case lexer.Token:
					if arg.(lexer.Token).Type == lexer.Int {
						arguments = append(arguments, constant.NewInt(types.I32, arg.(lexer.Token).Value.(int64)))
					} else if arg.(lexer.Token).Type == lexer.Identifier {
						arguments = append(arguments, variables[arg.(lexer.Token).Value.(string)].Value)
					} else if arg.(lexer.Token).Type == lexer.String {
						str := constant.NewCharArrayFromString(arg.(lexer.Token).Value.(string))
						arguments = append(arguments, str)
					}
				}
			}
			switch expression.Right.(type) {
			case ast.ExpressionStatement:
				if Functions[name].Type == "int" {
					return block.NewSDiv(block.NewCall(Functions[name].Value, arguments...), CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables))
				}
			case lexer.Token:
				if expression.Right.(lexer.Token).Type == lexer.Int {
					return block.NewSDiv(block.NewCall(Functions[name].Value, arguments...), constant.NewInt(types.I32, int64(expression.Right.(lexer.Token).Value.(int))))
				} else if expression.Right.(lexer.Token).Type == lexer.Identifier {
					return block.NewSDiv(block.NewCall(Functions[name].Value, arguments...), variables[expression.Right.(lexer.Token).Value.(string)].Value)
				}
			case ast.FuncCall:
				name2 := expression.Right.(ast.FuncCall).Name.Value.(string)
				args2 := expression.Right.(ast.FuncCall).Arguments
				var arguments2 []value.Value
				for _, arg := range args2 {
					switch arg.(type) {
					case ast.ExpressionStatement:
						arguments2 = append(arguments2, CompileExpression(block, arg.(ast.ExpressionStatement), variables))
					case lexer.Token:
						if arg.(lexer.Token).Type == lexer.Int {
							arguments2 = append(arguments2, constant.NewInt(types.I32, arg.(lexer.Token).Value.(int64)))
						} else if arg.(lexer.Token).Type == lexer.Identifier {
							arguments2 = append(arguments2, variables[arg.(lexer.Token).Value.(string)].Value)
						} else if arg.(lexer.Token).Type == lexer.String {
							str := constant.NewCharArrayFromString(arg.(lexer.Token).Value.(string))
							arguments2 = append(arguments2, str)
						}
					}
				}
				if Functions[name].Type == "int" && Functions[name2].Type == "int" {
					return block.NewSDiv(block.NewCall(Functions[name].Value, arguments...), block.NewCall(Functions[name2].Value, arguments2...))
				}
			}
		}
	}

	return block
}
