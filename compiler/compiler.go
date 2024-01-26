package compiler

import (
	"Voca-2/ast"
	"Voca-2/lexer"
	"fmt"
	"os"
	"path/filepath"

	"strconv"

	"os/exec"
	"runtime"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

var stringType = types.I8Ptr

var boolType = types.I1
var Functions = make(map[string]Function)
var ifs = 0

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
	Obj         bool
}

func New(program Program) []error {
	tmp, err := os.Create(program.File[:len(program.File)-4] + ".ll")
	program.Errs = append(program.Errs, err)
	_, err = tmp.Write([]byte(GenerateIR(program.Program)))

	program.Errs = append(program.Errs, err)
	tmp.Close()

	if program.OS == "" {
		program.OS = runtime.GOOS
	}
	if program.Arch == "" {
		program.Arch = runtime.GOARCH
	}
	if program.Arch == "amd64" {
		program.Arch = "x86_64"
	}
	program.Errs = append(program.Errs, CompileToObj(program))
	program.Errs = append(program.Errs, CompileToExecutable(program))
	if !program.Obj {
		if program.OS == "windows" {
			err = os.Remove(program.File[:len(program.File)-4] + ".obj")
			program.Errs = append(program.Errs, err)
		} else if program.OS == "linux" {
			err = os.Remove(program.File[:len(program.File)-4] + ".o")
			program.Errs = append(program.Errs, err)
		}
	}

	if !program.Ir {
		err = os.Remove(program.File[:len(program.File)-4] + ".ll")
		program.Errs = append(program.Errs, err)
	}
	return program.Errs
}
func CompileToObj(program Program) error {
	var cmd *exec.Cmd
	var llc string
	exePath, err := os.Executable()
	program.Errs = append(program.Errs, err)
	exeDir := filepath.Dir(exePath)
	if runtime.GOOS == "windows" {
		llc = filepath.Join(exeDir, "ll", "llc.exe")
	} else if runtime.GOOS == "linux" {
		llc = "llc"
	}

	if program.OS == "windows" {

		cmd = exec.Command(llc, "-mtriple="+program.Arch+"-pc-"+program.OS+"-msvc", "-filetype=obj", program.File[:len(program.File)-4]+".ll", "-o", program.File[:len(program.File)-4]+".obj")
	} else if program.OS == "linux" {
		cmd = exec.Command(llc, "-mtriple="+program.Arch+"-pc-"+program.OS+"-gnu", "-filetype=obj", program.File[:len(program.File)-4]+".ll", "-o", program.File[:len(program.File)-4]+".o")
	}

	return cmd.Run()

}
func CompileToExecutable(program Program) error {
	var cmd *exec.Cmd
	var lld string
	exePath, err := os.Executable()
	program.Errs = append(program.Errs, err)
	libs := filepath.Join(filepath.Dir(exePath), "libs")
	exeDir := filepath.Dir(exePath)
	if runtime.GOOS == "windows" {
		lld = filepath.Join(exeDir, "ll", "lld-link.exe")
	} else if runtime.GOOS == "linux" {
		lld = "clang"
	}

	if program.Arch == "x86_64" {
		if program.OS == "windows" {
			i := 0
			var libsObjs []string
			libsObjs = append(libsObjs, "/defaultlib:libc.lib /out:"+program.Output)
			libsObjs = append(libsObjs, program.File[:len(program.File)-4]+".obj")

			for i < len(program.Program.Externals) {
				libsObjs = append(libsObjs, filepath.Join(libs, program.Program.Externals[i], program.Program.Externals[i]+"_amd64-windows.obj"))
				i++
			}

			cmd = exec.Command(lld, libsObjs...)
		} else if program.OS == "linux" {
			i := 0
			var libsObjs []string
			libsObjs = append(libsObjs, program.File[:len(program.File)-4]+".o")
			for i < len(program.Program.Externals) {
				libsObjs = append(libsObjs, filepath.Join(libs, program.Program.Externals[i], program.Program.Externals[i]+"_amd64-linux.o"))
				i++
			}
			libsObjs = append(libsObjs, "-o", program.Output)
			cmd = exec.Command(lld, libsObjs...)
		}
	} else if program.Arch == "arm64" {
		i := 0
		var libsObjs []string
		libsObjs = append(libsObjs, program.File[:len(program.File)-4]+".o")
		for i < len(program.Program.Externals) {
			libsObjs = append(libsObjs, filepath.Join(libs, program.Program.Externals[i], program.Program.Externals[i]+"_arm64-linux.o"))
			i++
		}
		libsObjs = append(libsObjs, "-o", program.Output)
		cmd = exec.Command(lld, libsObjs...)

	}
	fmt.Println(cmd.Args)
	return cmd.Run()

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
			if statement.(ast.FuncDeclaration).Type.Value.(string) == "int" {
				entry.NewRet(constant.NewInt(types.I32, 0))
			} else if statement.(ast.FuncDeclaration).Type.Value.(string) == "string" {
				entry.NewRet(constant.NewPtrToInt(constant.NewInt(types.I32, 0), types.I8Ptr))
			} else if statement.(ast.FuncDeclaration).Type.Value.(string) == "void" {
				entry.NewRet(nil)
			}
		case ast.ExternFuncDeclaration:
			var fn *ir.Func
			par := statement.(ast.ExternFuncDeclaration).Arguments
			var params []*ir.Param
			for i := 0; i < len(par); i++ {
				if par[i].(ast.VariableDeclaration).Type.Value == "int" {
					params = append(params, ir.NewParam(par[i].(ast.VariableDeclaration).Name.Value.(string), types.I32))
				} else if par[i].(ast.VariableDeclaration).Type.Value == "string" {
					params = append(params, ir.NewParam(par[i].(ast.VariableDeclaration).Name.Value.(string), types.I8Ptr))
				}
			}
			if statement.(ast.ExternFuncDeclaration).Type.Value.(string) == "int" {
				fn = module.NewFunc(statement.(ast.ExternFuncDeclaration).Name.Value.(string), types.I32)
			} else if statement.(ast.ExternFuncDeclaration).Type.Value.(string) == "string" {
				fn = module.NewFunc(statement.(ast.ExternFuncDeclaration).Name.Value.(string), types.I8Ptr)
			} else if statement.(ast.ExternFuncDeclaration).Type.Value.(string) == "void" {
				fn = module.NewFunc(statement.(ast.ExternFuncDeclaration).Name.Value.(string), types.Void)
			}
			fn.Params = params
			function := Function{Type: statement.(ast.ExternFuncDeclaration).Type.Value.(string), Value: fn}
			Functions[statement.(ast.ExternFuncDeclaration).Name.Value.(string)] = function

		}
		i++
	}
	module_str := module.String()
	/*module_chararray := []byte(module_str)
	module_chararray = module_chararray
	n := 0
	for n < len(module_str) {

		if module_str[n] == 34 {
			n++

			for module_str[n] != 34 {

				n++
			}
			module_str = module_str[:n] + "\\00" + module_str[n:]
			n += 3
			module_chararray = []byte(module_str)

		}
		n++
	}*/
	return module_str
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
				variable = block.NewAlloca(stringType)
				variable.(*ir.InstAlloca).SetName(statement.(ast.VariableDeclaration).Name.Value.(string))
				variab.Type = "string"
				isstring = true

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
							strvar := block.NewAlloca(types.NewArray(uint64(str.(*constant.CharArray).Typ.Len), types.I8))
							block.NewStore(str, strvar)
							strin := block.NewGetElementPtr(types.NewArray(uint64(str.(*constant.CharArray).Typ.Len), types.I8), strvar, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
							block.NewStore(strin, variable)
							variab.length = int(str.(*constant.CharArray).Typ.Len)
						case *ir.InstCall:
							block.NewStore(str, variable)

						}

					} else {
						block.NewStore(CompileExpression(block, statement.(ast.VariableDeclaration).Value.(ast.ExpressionStatement), variables), variable)
					}
				case lexer.Token:
					if statement.(ast.VariableDeclaration).Value.(lexer.Token).Type == lexer.Int {
						block.NewStore(constant.NewInt(types.I32, statement.(ast.VariableDeclaration).Value.(lexer.Token).Value.(int64)), variable)
					} else if statement.(ast.VariableDeclaration).Value.(lexer.Token).Type == lexer.String {
						hasstring = true

						str := constant.NewCharArrayFromString(statement.(ast.VariableDeclaration).Value.(lexer.Token).Value.(string) + "\x00")

						length = int(str.Typ.Len)
						strvar := block.NewAlloca(types.NewArray(uint64(str.Typ.Len), types.I8))
						block.NewStore(str, strvar)
						strin := block.NewGetElementPtr(types.NewArray(uint64(length), types.I8), strvar, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
						block.NewStore(strin, variable)
					} else if statement.(ast.VariableDeclaration).Value.(lexer.Token).Type == lexer.Identifier {
						if isstring {
							hasstring = true
							value := variables[statement.(ast.VariableDeclaration).Value.(lexer.Token).Value.(string)].Value
							length = variables[statement.(ast.VariableDeclaration).Value.(lexer.Token).Value.(string)].length
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
								if variables[arg.(lexer.Token).Value.(string)].Type == "string" {
									arguments = append(arguments, block.NewLoad(types.I8Ptr, variables[arg.(lexer.Token).Value.(string)].Value))
								} else if variables[arg.(lexer.Token).Value.(string)].Type == "int" {
									arguments = append(arguments, block.NewLoad(types.I32, variables[arg.(lexer.Token).Value.(string)].Value))
								} else if variables[arg.(lexer.Token).Value.(string)].Type == "bool" {
									arguments = append(arguments, block.NewLoad(types.I1, variables[arg.(lexer.Token).Value.(string)].Value))

								}
							} else if arg.(lexer.Token).Type == lexer.String {
								str := constant.NewCharArrayFromString(arg.(lexer.Token).Value.(string) + "\x00")
								strvar := block.NewAlloca(types.NewArray(uint64(str.Typ.Len), types.I8))
								block.NewStore(str, strvar)
								strPtr := block.NewGetElementPtr(types.NewArray(uint64(str.Typ.Len), types.I8), strvar, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
								arguments = append(arguments, strPtr)
							}
						}
					}
					if Functions[name].Type == "int" {
						block.NewStore(block.NewCall(Functions[name].Value, arguments...), variable)
					} else if Functions[name].Type == "string" {
						hasstring = true

						str := block.NewCall(Functions[name].Value, arguments...)

						variable = str

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
			for i := 0; i < len(args); i++ {
				arg := args[i]
				switch arg.(type) {
				case ast.ExpressionStatement:
					arguments = append(arguments, CompileExpression(block, arg.(ast.ExpressionStatement), variables))
				case lexer.Token:
					if arg.(lexer.Token).Type == lexer.Int {
						arguments = append(arguments, constant.NewInt(types.I32, arg.(lexer.Token).Value.(int64)))
					} else if arg.(lexer.Token).Type == lexer.Identifier {
						if variables[arg.(lexer.Token).Value.(string)].Type == "string" {
							arguments = append(arguments, block.NewLoad(types.I8Ptr, variables[arg.(lexer.Token).Value.(string)].Value))
						} else if variables[arg.(lexer.Token).Value.(string)].Type == "int" {
							arguments = append(arguments, block.NewLoad(types.I32, variables[arg.(lexer.Token).Value.(string)].Value))
						} else if variables[arg.(lexer.Token).Value.(string)].Type == "bool" {
							arguments = append(arguments, block.NewLoad(types.I1, variables[arg.(lexer.Token).Value.(string)].Value))

						}
					} else if arg.(lexer.Token).Type == lexer.String {
						str := constant.NewCharArrayFromString(arg.(lexer.Token).Value.(string) + "\x00")
						strvar := block.NewAlloca(types.NewArray(uint64(str.Typ.Len), types.I8))
						block.NewStore(str, strvar)
						strPtr := block.NewGetElementPtr(types.NewArray(uint64(str.Typ.Len), types.I8), strvar, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
						arguments = append(arguments, strPtr)
					}
				}
				i++
			}
			block.NewCall(Functions[name].Value, arguments...)
		case ast.IfStatement:
			ifs++
			condition := CompileBool(block, statement.(ast.IfStatement).Condition, variables)
			trueBlock := block.Parent.NewBlock("true" + strconv.Itoa(ifs))
			trueBlock = Compile(trueBlock, statement.(ast.IfStatement).Consequence, variables)
			falseBlock := block.Parent.NewBlock("false" + strconv.Itoa(ifs))
			falseBlock = Compile(falseBlock, statement.(ast.IfStatement).Alternative, variables)
			afterBlock := block.Parent.NewBlock("after" + strconv.Itoa(ifs))
			trueBlock.NewBr(afterBlock)
			falseBlock.NewBr(afterBlock)
			if statement.(ast.IfStatement).Condition.Operator.Type == lexer.DoubleEqual {
				block.NewCondBr(condition, trueBlock, falseBlock)
			} else if statement.(ast.IfStatement).Condition.Operator.Type == lexer.NotEqual {
				block.NewCondBr(condition, trueBlock, falseBlock)
			} else if statement.(ast.IfStatement).Condition.Operator.Type == lexer.MoreThan {
				block.NewCondBr(condition, trueBlock, falseBlock)
			} else if statement.(ast.IfStatement).Condition.Operator.Type == lexer.LessThan {
				block.NewCondBr(condition, trueBlock, falseBlock)
			}

			block = afterBlock

		}
		i++
	}

	return block
}
func CompileBool(block *ir.Block, condition ast.BoolStatement, variables map[string]Variable) value.Value {
	var con1 value.Value
	var con2 value.Value
	switch condition.Condition1.(type) {
	case ast.ExpressionStatement:
		con1 = CompileExpression(block, condition.Condition1.(ast.ExpressionStatement), variables)
	case lexer.Token:
		if condition.Condition1.(lexer.Token).Type == lexer.Int {
			con1 = constant.NewInt(types.I32, int64(condition.Condition1.(lexer.Token).Value.(int)))
		} else if condition.Condition1.(lexer.Token).Type == lexer.Identifier {
			con1 = variables[condition.Condition1.(lexer.Token).Value.(string)].Value
			if variables[condition.Condition1.(lexer.Token).Value.(string)].Type == "string" {
				con1 = block.NewLoad(types.NewArray(uint64(variables[condition.Condition1.(lexer.Token).Value.(string)].length), types.I8), con1)
			}
		} else if condition.Condition1.(lexer.Token).Type == lexer.String {
			con1 = constant.NewCharArrayFromString(condition.Condition1.(lexer.Token).Value.(string) + "\x00")
		}
	}
	switch condition.Condition2.(type) {
	case ast.ExpressionStatement:
		con2 = CompileExpression(block, condition.Condition2.(ast.ExpressionStatement), variables)
	case lexer.Token:
		if condition.Condition2.(lexer.Token).Type == lexer.Int {
			con2 = constant.NewInt(types.I32, int64(condition.Condition2.(lexer.Token).Value.(int)))
		} else if condition.Condition2.(lexer.Token).Type == lexer.Identifier {
			con2 = variables[condition.Condition2.(lexer.Token).Value.(string)].Value
			if variables[condition.Condition2.(lexer.Token).Value.(string)].Type == "string" {
				con2 = block.NewLoad(types.NewArray(uint64(variables[condition.Condition2.(lexer.Token).Value.(string)].length), types.I8), con2)
			}
		} else if condition.Condition1.(lexer.Token).Type == lexer.String {
			con2 = constant.NewCharArrayFromString(condition.Condition2.(lexer.Token).Value.(string) + "\x00")
		}
	}
	if condition.Operator.Type == lexer.DoubleEqual {
		return block.NewICmp(enum.IPredEQ, con1, con2)
	} else if condition.Operator.Type == lexer.NotEqual {
		return block.NewICmp(enum.IPredNE, con1, con2)
	} else if condition.Operator.Type == lexer.MoreThan {
		return block.NewICmp(enum.IPredSGT, con1, con2)
	} else if condition.Operator.Type == lexer.LessThan {
		return block.NewICmp(enum.IPredSLT, con1, con2)
	}
	return nil
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
							if variables[arg.(lexer.Token).Value.(string)].Type == "string" {
								arguments = append(arguments, block.NewLoad(types.I8Ptr, variables[arg.(lexer.Token).Value.(string)].Value))
							} else if variables[arg.(lexer.Token).Value.(string)].Type == "int" {
								arguments = append(arguments, block.NewLoad(types.I32, variables[arg.(lexer.Token).Value.(string)].Value))
							} else if variables[arg.(lexer.Token).Value.(string)].Type == "bool" {
								arguments = append(arguments, block.NewLoad(types.I1, variables[arg.(lexer.Token).Value.(string)].Value))

							}
						} else if arg.(lexer.Token).Type == lexer.String {
							str := constant.NewCharArrayFromString(arg.(lexer.Token).Value.(string) + "\x00")
							strvar := block.NewAlloca(types.NewArray(uint64(str.Typ.Len), types.I8))
							block.NewStore(str, strvar)
							strPtr := block.NewGetElementPtr(types.NewArray(uint64(str.Typ.Len), types.I8), strvar, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
							arguments = append(arguments, strPtr)
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

						return block.NewCall(Functions["append"].Value, block.NewGetElementPtr(types.NewArray(uint64(str1.Typ.Len), types.I8), str1, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0)), block.NewGetElementPtr(types.NewArray(uint64(str2.(*constant.CharArray).Typ.Len), types.I8), str2, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0)))
					}
				} else if expression.Left.(lexer.Token).Type == lexer.String {
					str1 := constant.NewCharArrayFromString(expression.Left.(lexer.Token).Value.(string) + "\x00")
					str2 := CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables).(*constant.CharArray)
					return block.NewCall(Functions["append"].Value, block.NewGetElementPtr(types.NewArray(uint64(str1.Typ.Len), types.I8), str1, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0)), block.NewGetElementPtr(types.NewArray(uint64(str2.Typ.Len), types.I8), str2, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0)))

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

						return block.NewCall(Functions["append"].Value, block.NewLoad(types.I8Ptr, variables[expression.Left.(lexer.Token).Value.(string)].Value), block.NewLoad(types.I8Ptr, variables[expression.Right.(lexer.Token).Value.(string)].Value))
					}
				} else if expression.Left.(lexer.Token).Type == lexer.String && expression.Right.(lexer.Token).Type == lexer.String {
					str1 := constant.NewCharArrayFromString(expression.Left.(lexer.Token).Value.(string) + "\x00")
					str2 := constant.NewCharArrayFromString(expression.Right.(lexer.Token).Value.(string) + "\x00")
					str := constant.NewCharArray(append(str1.X, str2.X...))
					return block.NewGetElementPtr(types.NewArray(uint64(str.Typ.Len), types.I8), str, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
				} else if expression.Left.(lexer.Token).Type == lexer.String && expression.Right.(lexer.Token).Type == lexer.Identifier {
					str1 := constant.NewCharArrayFromString(expression.Left.(lexer.Token).Value.(string) + "\x00")
					str2 := variables[expression.Right.(lexer.Token).Value.(string)].Value.(*constant.CharArray)
					return block.NewCall(Functions["append"].Value, block.NewGetElementPtr(types.NewArray(uint64(str1.Typ.Len), types.I8), str1, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0)), str2)
				} else if expression.Left.(lexer.Token).Type == lexer.Identifier && expression.Right.(lexer.Token).Type == lexer.String {
					str1 := variables[expression.Left.(lexer.Token).Value.(string)].Value.(*constant.CharArray)
					str2 := constant.NewCharArrayFromString(expression.Right.(lexer.Token).Value.(string) + "\x00")
					return block.NewCall(Functions["append"].Value, str1, block.NewGetElementPtr(types.NewArray(uint64(str2.Typ.Len), types.I8), str2, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0)))
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
							str := constant.NewCharArrayFromString(arg.(lexer.Token).Value.(string) + "\x00")
							strvar := block.NewAlloca(types.NewArray(uint64(str.Typ.Len), types.I8))
							block.NewStore(str, strvar)
							strPtr := block.NewGetElementPtr(types.NewArray(uint64(str.Typ.Len), types.I8), strvar, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
							arguments2 = append(arguments2, strPtr)
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
						if variables[arg.(lexer.Token).Value.(string)].Type == "string" {
							arguments = append(arguments, block.NewLoad(types.I8Ptr, variables[arg.(lexer.Token).Value.(string)].Value))
						} else if variables[arg.(lexer.Token).Value.(string)].Type == "int" {
							arguments = append(arguments, block.NewLoad(types.I32, variables[arg.(lexer.Token).Value.(string)].Value))
						} else if variables[arg.(lexer.Token).Value.(string)].Type == "bool" {
							arguments = append(arguments, block.NewLoad(types.I1, variables[arg.(lexer.Token).Value.(string)].Value))

						}
					} else if arg.(lexer.Token).Type == lexer.String {
						str := constant.NewCharArrayFromString(arg.(lexer.Token).Value.(string) + "\x00")
						strvar := block.NewAlloca(types.NewArray(uint64(str.Typ.Len), types.I8))
						block.NewStore(str, strvar)
						strPtr := block.NewGetElementPtr(types.NewArray(uint64(str.Typ.Len), types.I8), strvar, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
						arguments = append(arguments, strPtr)
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
					str := constant.NewCharArrayFromString(expression.Right.(lexer.Token).Value.(string) + "\x00")
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
							str := constant.NewCharArrayFromString(arg.(lexer.Token).Value.(string) + "\x00")
							strvar := block.NewAlloca(types.NewArray(uint64(str.Typ.Len), types.I8))
							block.NewStore(str, strvar)
							strPtr := block.NewGetElementPtr(types.NewArray(uint64(str.Typ.Len), types.I8), strvar, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
							arguments2 = append(arguments2, strPtr)
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
							if variables[arg.(lexer.Token).Value.(string)].Type == "string" {
								arguments = append(arguments, block.NewLoad(types.I8Ptr, variables[arg.(lexer.Token).Value.(string)].Value))
							} else if variables[arg.(lexer.Token).Value.(string)].Type == "int" {
								arguments = append(arguments, block.NewLoad(types.I32, variables[arg.(lexer.Token).Value.(string)].Value))
							} else if variables[arg.(lexer.Token).Value.(string)].Type == "bool" {
								arguments = append(arguments, block.NewLoad(types.I1, variables[arg.(lexer.Token).Value.(string)].Value))

							}
						} else if arg.(lexer.Token).Type == lexer.String {
							str := constant.NewCharArrayFromString(arg.(lexer.Token).Value.(string) + "\x00")
							strPtr := block.NewGetElementPtr(types.NewArray(uint64(str.Typ.Len), types.I8), str, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
							arguments = append(arguments, strPtr)
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
							str := constant.NewCharArrayFromString(arg.(lexer.Token).Value.(string) + "\x00")
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
						if variables[arg.(lexer.Token).Value.(string)].Type == "string" {
							arguments = append(arguments, block.NewLoad(types.I8Ptr, variables[arg.(lexer.Token).Value.(string)].Value))
						} else if variables[arg.(lexer.Token).Value.(string)].Type == "int" {
							arguments = append(arguments, block.NewLoad(types.I32, variables[arg.(lexer.Token).Value.(string)].Value))
						} else if variables[arg.(lexer.Token).Value.(string)].Type == "bool" {
							arguments = append(arguments, block.NewLoad(types.I1, variables[arg.(lexer.Token).Value.(string)].Value))

						}
					} else if arg.(lexer.Token).Type == lexer.String {
						str := constant.NewCharArrayFromString(arg.(lexer.Token).Value.(string) + "\x00")
						strPtr := block.NewGetElementPtr(types.NewArray(uint64(str.Typ.Len), types.I8), str, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
						arguments = append(arguments, strPtr)
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
							str := constant.NewCharArrayFromString(arg.(lexer.Token).Value.(string) + "\x00")
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
							if variables[arg.(lexer.Token).Value.(string)].Type == "string" {
								arguments = append(arguments, block.NewLoad(types.I8Ptr, variables[arg.(lexer.Token).Value.(string)].Value))
							} else if variables[arg.(lexer.Token).Value.(string)].Type == "int" {
								arguments = append(arguments, block.NewLoad(types.I32, variables[arg.(lexer.Token).Value.(string)].Value))
							} else if variables[arg.(lexer.Token).Value.(string)].Type == "bool" {
								arguments = append(arguments, block.NewLoad(types.I1, variables[arg.(lexer.Token).Value.(string)].Value))

							}
						} else if arg.(lexer.Token).Type == lexer.String {
							str := constant.NewCharArrayFromString(arg.(lexer.Token).Value.(string) + "\x00")
							strPtr := block.NewGetElementPtr(types.NewArray(uint64(str.Typ.Len), types.I8), str, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
							arguments = append(arguments, strPtr)
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
							str := constant.NewCharArrayFromString(arg.(lexer.Token).Value.(string) + "\x00")
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
						if variables[arg.(lexer.Token).Value.(string)].Type == "string" {
							arguments = append(arguments, block.NewLoad(types.I8Ptr, variables[arg.(lexer.Token).Value.(string)].Value))
						} else if variables[arg.(lexer.Token).Value.(string)].Type == "int" {
							arguments = append(arguments, block.NewLoad(types.I32, variables[arg.(lexer.Token).Value.(string)].Value))
						} else if variables[arg.(lexer.Token).Value.(string)].Type == "bool" {
							arguments = append(arguments, block.NewLoad(types.I1, variables[arg.(lexer.Token).Value.(string)].Value))

						}
					} else if arg.(lexer.Token).Type == lexer.String {
						str := constant.NewCharArrayFromString(arg.(lexer.Token).Value.(string) + "\x00")
						strvar := block.NewAlloca(types.NewArray(uint64(str.Typ.Len), types.I8))
						block.NewStore(str, strvar)
						strPtr := block.NewGetElementPtr(types.NewArray(uint64(str.Typ.Len), types.I8), strvar, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
						arguments = append(arguments, strPtr)
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
							str := constant.NewCharArrayFromString(arg.(lexer.Token).Value.(string) + "\x00")
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
							if variables[arg.(lexer.Token).Value.(string)].Type == "string" {
								arguments = append(arguments, block.NewLoad(types.I8Ptr, variables[arg.(lexer.Token).Value.(string)].Value))
							} else if variables[arg.(lexer.Token).Value.(string)].Type == "int" {
								arguments = append(arguments, block.NewLoad(types.I32, variables[arg.(lexer.Token).Value.(string)].Value))
							} else if variables[arg.(lexer.Token).Value.(string)].Type == "bool" {
								arguments = append(arguments, block.NewLoad(types.I1, variables[arg.(lexer.Token).Value.(string)].Value))

							}
						} else if arg.(lexer.Token).Type == lexer.String {
							str := constant.NewCharArrayFromString(arg.(lexer.Token).Value.(string) + "\x00")
							strPtr := block.NewGetElementPtr(types.NewArray(uint64(str.Typ.Len), types.I8), str, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
							arguments = append(arguments, strPtr)
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
							str := constant.NewCharArrayFromString(arg.(lexer.Token).Value.(string) + "\x00")
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
						if variables[arg.(lexer.Token).Value.(string)].Type == "string" {
							arguments = append(arguments, block.NewLoad(types.I8Ptr, variables[arg.(lexer.Token).Value.(string)].Value))
						} else if variables[arg.(lexer.Token).Value.(string)].Type == "int" {
							arguments = append(arguments, block.NewLoad(types.I32, variables[arg.(lexer.Token).Value.(string)].Value))
						} else if variables[arg.(lexer.Token).Value.(string)].Type == "bool" {
							arguments = append(arguments, block.NewLoad(types.I1, variables[arg.(lexer.Token).Value.(string)].Value))

						}
					} else if arg.(lexer.Token).Type == lexer.String {
						str := constant.NewCharArrayFromString(arg.(lexer.Token).Value.(string) + "\x00")
						strPtr := block.NewGetElementPtr(types.NewArray(uint64(str.Typ.Len), types.I8), str, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
						arguments = append(arguments, strPtr)
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
							str := constant.NewCharArrayFromString(arg.(lexer.Token).Value.(string) + "\x00")
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
