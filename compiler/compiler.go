package compiler

import (
	"Voca-2/ast"
	"Voca-2/lexer"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	Type      any
	Value     value.Value
	hasstring bool
	length    int
}
type Function struct {
	Type  string
	Value value.Value
}

type Program struct {
	Input          string
	Errs           []error
	Tokens         []lexer.Token
	Program        ast.Program
	GenerateAST    bool
	Args           []string
	File           string
	Output         string
	Arch           string
	OS             string
	LoadAST        bool
	Ir             bool
	Obj            bool
	Exec           bool
	Optimalisation int
	JustParse      bool
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
	err = CompileToObj(program)
	if err != nil {
		err2 := errors.New("lld")
		program.Errs = append(program.Errs, fmt.Errorf("%s: %w", err2, err))
	}
	if program.Exec {
		err = CompileToExecutable(program)
		if err != nil {
			err2 := errors.New("clang")
			program.Errs = append(program.Errs, fmt.Errorf("%s: %w", err2, err))
		}
	}
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

		cmd = exec.Command(llc, "-mtriple="+program.Arch+"-w64-"+program.OS+"-gnu", "-O"+strconv.Itoa(program.Optimalisation), "-filetype=obj", program.File[:len(program.File)-4]+".ll", "-o", program.File[:len(program.File)-4]+".obj")
	} else if program.OS == "linux" {
		cmd = exec.Command(llc, "-mtriple="+program.Arch+"-pc-"+program.OS+"-gnu", "-O"+strconv.Itoa(program.Optimalisation), "-filetype=obj", program.File[:len(program.File)-4]+".ll", "-o", program.File[:len(program.File)-4]+".o")
	}
	fmt.Println(cmd.Args)
	return cmd.Run()

}
func CompileToExecutable(program Program) error {
	var cmd *exec.Cmd
	var lld string
	exePath, err := os.Executable()
	program.Errs = append(program.Errs, err)
	libs := filepath.Join(filepath.Dir(exePath), "libs")
	//exeDir := filepath.Dir(exePath)
	lld = "clang"

	if program.Arch == "x86_64" {
		if program.OS == "windows" {
			i := 0
			var libsObjs []string
			libsObjs = append(libsObjs, program.File[:len(program.File)-4]+".obj")
			for i < len(program.Program.Externals) {
				addflags := parseCfg(program.Program.Externals[i], filepath.Join(libs, program.Program.Externals[i]), "amd64", "windows")

				libsObjs = append(libsObjs, filepath.Join(libs, program.Program.Externals[i], program.Program.Externals[i]+"_amd64-windows.obj"))
				libsObjs = append(libsObjs, strings.Split(addflags, " ")...)
				i++
			}
			libsObjs = append(libsObjs, "-o", program.Output, "-lgcc")
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

			cmd = exec.Command("gcc", libsObjs...)
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
	err = cmd.Run()

	return err

}

var module = ir.NewModule()

func GenerateIR(program ast.Program) string {

	i := 0
	for i < len(program.Statements) {
		statement := program.Statements[i].Node
		switch statement.(type) {
		case ast.FuncDeclaration:

			var fn *ir.Func
			var vars map[string]Variable
			vars = make(map[string]Variable)
			if statement.(ast.FuncDeclaration).Type.Value.(string) == "int" {
				fn = module.NewFunc(statement.(ast.FuncDeclaration).Name.Value.(string), types.I32)
			} else if statement.(ast.FuncDeclaration).Type.Value.(string) == "string" {
				fn = module.NewFunc(statement.(ast.FuncDeclaration).Name.Value.(string), types.I8Ptr)
			} else if statement.(ast.FuncDeclaration).Type.Value.(string) == "void" {
				fn = module.NewFunc(statement.(ast.FuncDeclaration).Name.Value.(string), types.Void)
			} else if statement.(ast.FuncDeclaration).Type.Value.(string) == "bool" {
				fn = module.NewFunc(statement.(ast.FuncDeclaration).Name.Value.(string), types.I1)
			} else if statement.(ast.FuncDeclaration).Type.Value.(string) == "float" {
				fn = module.NewFunc(statement.(ast.FuncDeclaration).Name.Value.(string), types.Float)
			}
			entry := fn.NewBlock("entry")
			par := statement.(ast.FuncDeclaration).Arguments
			var params []*ir.Param
			for i := 0; i < len(par); i++ {
				variable := Variable{}
				if par[i].(ast.VariableDeclaration).Type.Value == "int" {
					params = append(params, ir.NewParam(par[i].(ast.VariableDeclaration).Name.Value.(string), types.I32))
					variab := entry.NewAlloca(types.I32)
					entry.NewStore(params[len(params)-1], variab)
					variable.Value = variab
					variable.Type = "int"
				} else if par[i].(ast.VariableDeclaration).Type.Value == "string" {
					params = append(params, ir.NewParam(par[i].(ast.VariableDeclaration).Name.Value.(string), types.I8Ptr))
					variab := entry.NewAlloca(types.I8Ptr)
					entry.NewStore(params[len(params)-1], variab)
					variable.Value = variab
					variable.Type = "string"
				} else if statement.(ast.FuncDeclaration).Arguments[i].(ast.VariableDeclaration).Type.Value == "bool" {
					params = append(params, ir.NewParam(par[i].(ast.VariableDeclaration).Name.Value.(string), types.I1))
					variab := entry.NewAlloca(types.I1)
					entry.NewStore(params[len(params)-1], variab)
					variable.Value = variab
					variable.Type = "bool"
				} else if statement.(ast.FuncDeclaration).Arguments[i].(ast.VariableDeclaration).Type.Value == "float" {
					params = append(params, ir.NewParam(par[i].(ast.VariableDeclaration).Name.Value.(string), types.Float))
					variab := entry.NewAlloca(types.Float)
					entry.NewStore(params[len(params)-1], variab)
					variable.Value = variab
					variable.Type = "float"
				}
				vars[statement.(ast.FuncDeclaration).Arguments[i].(ast.VariableDeclaration).Name.Value.(string)] = variable
			}
			if statement.(ast.FuncDeclaration).Name.Value.(string) == "main" {
				params = append(params, ir.NewParam("argc", types.I32), ir.NewParam("argv", types.NewPointer(types.I8Ptr)))
			}

			fn.Params = params

			function := Function{Type: statement.(ast.FuncDeclaration).Type.Value.(string), Value: fn}
			Functions[statement.(ast.FuncDeclaration).Name.Value.(string)] = function
			entry = Compile(entry, statement.(ast.FuncDeclaration).Body, vars)
			if entry.Term == nil {
				if statement.(ast.FuncDeclaration).Type.Value.(string) == "int" {
					entry.NewRet(constant.NewInt(types.I32, 0))
				} else if statement.(ast.FuncDeclaration).Type.Value.(string) == "string" {
					entry.NewRet(constant.NewPtrToInt(constant.NewInt(types.I32, 0), types.I8Ptr))
				} else if statement.(ast.FuncDeclaration).Type.Value.(string) == "void" {
					entry.NewRet(nil)
				} else if statement.(ast.FuncDeclaration).Type.Value.(string) == "bool" {
					entry.NewRet(constant.NewBool(false))
				} else if statement.(ast.FuncDeclaration).Type.Value.(string) == "float" {
					entry.NewRet(constant.NewFloat(types.Float, 0))
				} else {
					entry.NewRet(nil)
				}
			}
		case ast.ExternFuncDeclaration:
			var fn *ir.Func
			if statement.(ast.ExternFuncDeclaration).Type.Value.(string) == "int" {
				fn = module.NewFunc(statement.(ast.ExternFuncDeclaration).Name.Value.(string), types.I32)
			} else if statement.(ast.ExternFuncDeclaration).Type.Value.(string) == "string" {
				fn = module.NewFunc(statement.(ast.ExternFuncDeclaration).Name.Value.(string), types.I8Ptr)
			} else if statement.(ast.ExternFuncDeclaration).Type.Value.(string) == "void" {
				fn = module.NewFunc(statement.(ast.ExternFuncDeclaration).Name.Value.(string), types.Void)
			} else if statement.(ast.ExternFuncDeclaration).Type.Value.(string) == "bool" {
				fn = module.NewFunc(statement.(ast.ExternFuncDeclaration).Name.Value.(string), types.I1)
			} else if statement.(ast.ExternFuncDeclaration).Type.Value.(string) == "float" {
				fn = module.NewFunc(statement.(ast.ExternFuncDeclaration).Name.Value.(string), types.Float)
			}
			par := statement.(ast.ExternFuncDeclaration).Arguments
			var params []*ir.Param
			for i := 0; i < len(par); i++ {
				variable := Variable{}
				if par[i].(ast.VariableDeclaration).Type.Value == "int" {
					params = append(params, ir.NewParam(par[i].(ast.VariableDeclaration).Name.Value.(string), types.I32))
					variable.Type = "int"
				} else if par[i].(ast.VariableDeclaration).Type.Value == "string" {
					params = append(params, ir.NewParam(par[i].(ast.VariableDeclaration).Name.Value.(string), types.I8Ptr))
					variable.Type = "string"
				} else if statement.(ast.ExternFuncDeclaration).Arguments[i].(ast.VariableDeclaration).Type.Value == "bool" {
					params = append(params, ir.NewParam(par[i].(ast.VariableDeclaration).Name.Value.(string), types.I1))
					variable.Type = "bool"
				} else if statement.(ast.ExternFuncDeclaration).Arguments[i].(ast.VariableDeclaration).Type.Value == "float" {
					params = append(params, ir.NewParam(par[i].(ast.VariableDeclaration).Name.Value.(string), types.Float))
					variable.Type = "float"
				}
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

var strCounter int = 0

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
			} else if statement.(ast.VariableDeclaration).Type.Value.(string) == "float" {
				variable = block.NewAlloca(types.Float)
				variable.(*ir.InstAlloca).SetName(statement.(ast.VariableDeclaration).Name.Value.(string))
				variab.Type = "float"
			}
			if statement.(ast.VariableDeclaration).Value != nil {
				switch statement.(ast.VariableDeclaration).Value.(type) {
				case ast.ExpressionStatement:
					if isstring {
						hasstring = true

						str, _ := CompileExpression(block, statement.(ast.VariableDeclaration).Value.(ast.ExpressionStatement), variables)
						switch str.(type) {
						case *ir.InstGetElementPtr:

							block.NewStore(str, variable)
							//variab.length = int(str.(*constant.CharArray).Typ.Len)
						case *ir.InstCall:
							block.NewStore(str, variable)

						}

					} else {
						exp, _ := CompileExpression(block, statement.(ast.VariableDeclaration).Value.(ast.ExpressionStatement), variables)
						block.NewStore(exp, variable)
					}
				case lexer.Token:
					if statement.(ast.VariableDeclaration).Value.(lexer.Token).Type == lexer.Int {
						block.NewStore(constant.NewInt(types.I32, int64(statement.(ast.VariableDeclaration).Value.(lexer.Token).Value.(int))), variable)
					} else if statement.(ast.VariableDeclaration).Value.(lexer.Token).Type == lexer.String {
						hasstring = true

						str := constant.NewCharArrayFromString(statement.(ast.VariableDeclaration).Value.(lexer.Token).Value.(string) + "\x00")

						length = int(str.Typ.Len)
						//strvar := block.NewAlloca(types.NewArray(uint64(str.Typ.Len), types.I8))
						strvar := module.NewGlobalDef("str"+strconv.Itoa(strCounter), str)
						strCounter++
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
					} else if statement.(ast.VariableDeclaration).Value.(lexer.Token).Value == "true" {
						block.NewStore(constant.NewBool(true), variable)
					} else if statement.(ast.VariableDeclaration).Value.(lexer.Token).Value == "false" {
						block.NewStore(constant.NewBool(false), variable)
					} else if statement.(ast.VariableDeclaration).Value.(lexer.Token).Type == lexer.Float {
						block.NewStore(constant.NewFloat(types.Float, statement.(ast.VariableDeclaration).Value.(lexer.Token).Value.(float64)), variable)
					}

				case ast.FuncCall:
					name := statement.(ast.VariableDeclaration).Value.(ast.FuncCall).Name.Value.(string)
					arguments := FuncCall(block, statement.(ast.VariableDeclaration).Value.(ast.FuncCall), variables)
					if Functions[name].Type == "int" {
						block.NewStore(block.NewCall(Functions[name].Value, arguments...), variable)
					} else if Functions[name].Type == "string" {
						hasstring = true

						str := block.NewCall(Functions[name].Value, arguments...)

						variable = str

						block.NewStore(str, variable)
					} else if Functions[name].Type == "float" {
						block.NewStore(block.NewCall(Functions[name].Value, arguments...), variable)
					} else if Functions[name].Type == "bool" {
						block.NewStore(block.NewCall(Functions[name].Value, arguments...), variable)
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
		case ast.ArrayDeclaration:
			var lenght value.Value
			variab := Variable{}
			if statement.(ast.ArrayDeclaration).Length != nil {
				switch statement.(ast.ArrayDeclaration).Length.(type) {
				case ast.ExpressionStatement:
					exp, _ := CompileExpression(block, statement.(ast.ArrayDeclaration).Length.(ast.ExpressionStatement), variables)
					lenght = exp
				case lexer.Token:
					if statement.(ast.ArrayDeclaration).Length.(lexer.Token).Type == lexer.Int {
						lenght = constant.NewInt(types.I32, int64(statement.(ast.ArrayDeclaration).Length.(lexer.Token).Value.(int)))
					} else if statement.(ast.ArrayDeclaration).Length.(lexer.Token).Type == lexer.Identifier {
						lenght = block.NewLoad(types.I32, variables[statement.(ast.ArrayDeclaration).Length.(lexer.Token).Value.(string)].Value)
					}
				case ast.FuncCall:
					name := statement.(ast.ArrayDeclaration).Length.(ast.FuncCall).Name.Value.(string)
					arguments := FuncCall(block, statement.(ast.ArrayDeclaration).Length.(ast.FuncCall), variables)
					lenght = block.NewCall(Functions[name].Value, arguments...)

				}
			}
			var variable value.Value
			switch statement.(ast.ArrayDeclaration).Type.(type) {
			case string:
				if statement.(ast.ArrayDeclaration).Type.(string) == "int" {
					variable = block.NewAlloca(types.NewArray(uint64(lenght.(*constant.Int).X.Int64()), types.I32))
					variab.Type = types.NewArray(uint64(lenght.(*constant.Int).X.Int64()), types.I32)
				}
			case ast.ArrayType:
				var ArrType = ArrayType(statement.(ast.ArrayDeclaration).Type.(ast.ArrayType))
				variable = block.NewAlloca(types.NewArray(uint64(lenght.(*constant.Int).X.Int64()), ArrType))
				variab.Type = types.NewArray(uint64(lenght.(*constant.Int).X.Int64()), ArrType)
			}
			variable.(*ir.InstAlloca).SetName(statement.(ast.ArrayDeclaration).Name.Value.(string))

			array := statement.(ast.ArrayDeclaration).Value
			indexes := []int{0}
			variable = CompileArrayAssigment(block, array, variables, variable, variab, indexes, 1) //.(*ir.InstAlloca)

			variab.Value = variable
			variab.length = int(lenght.(*constant.Int).X.Int64())
			variables[statement.(ast.ArrayDeclaration).Name.Value.(string)] = variab

		case ast.VariableAssignment:
			variable := variables[statement.(ast.VariableAssignment).Name.Value.(string)].Value
			isstring = variables[statement.(ast.VariableAssignment).Name.Value.(string)].Type == "string"
			hasstring = variables[statement.(ast.VariableAssignment).Name.Value.(string)].hasstring
			switch statement.(ast.VariableAssignment).Value.(type) {
			case ast.ExpressionStatement:
				if isstring {
					hasstring = true

					str, _ := CompileExpression(block, statement.(ast.VariableAssignment).Value.(ast.ExpressionStatement), variables)
					switch str.(type) {
					case *ir.InstGetElementPtr:
						//strvar := block.NewAlloca(types.NewArray(uint64(str.(*constant.CharArray).Typ.Len), types.I8))
						//strvar := module.NewGlobalDef("str"+strconv.Itoa(strCounter), str)
						//strCounter++
						//block.NewStore(str, strvar)
						//strin := block.NewGetElementPtr(types.NewArray(uint64(str.(*constant.CharArray).Typ.Len), types.I8), strvar, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
						block.NewStore(str, variable)
						//variab.length = int(str.(*constant.CharArray).Typ.Len)
					case *ir.InstCall:
						block.NewStore(str, variable)

					}

				} else {
					exp, _ := CompileExpression(block, statement.(ast.VariableAssignment).Value.(ast.ExpressionStatement), variables)

					block.NewStore(exp, variable)
				}
			case lexer.Token:
				if statement.(ast.VariableAssignment).Value.(lexer.Token).Type == lexer.Int {
					block.NewStore(constant.NewInt(types.I32, statement.(ast.VariableAssignment).Value.(lexer.Token).Value.(int64)), variable)
				} else if statement.(ast.VariableAssignment).Value.(lexer.Token).Type == lexer.String {
					hasstring = true

					str := constant.NewCharArrayFromString(statement.(ast.VariableAssignment).Value.(lexer.Token).Value.(string) + "\x00")

					length = int(str.Typ.Len)
					strvar := module.NewGlobalDef("str"+strconv.Itoa(strCounter), str)
					strCounter++
					block.NewStore(str, strvar)
					strin := block.NewGetElementPtr(types.NewArray(uint64(length), types.I8), strvar, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
					block.NewStore(strin, variable)
				} else if statement.(ast.VariableAssignment).Value.(lexer.Token).Type == lexer.Identifier {
					if isstring {
						hasstring = true
						value := variables[statement.(ast.VariableAssignment).Value.(lexer.Token).Value.(string)].Value
						length = variables[statement.(ast.VariableAssignment).Value.(lexer.Token).Value.(string)].length
						block.NewStore(value, variable)
					} else {
						block.NewStore(variables[statement.(ast.VariableAssignment).Value.(lexer.Token).Value.(string)].Value, variable)
					}
				} else if statement.(ast.VariableAssignment).Value.(lexer.Token).Value == "true" {
					block.NewStore(constant.NewBool(true), variable)
				} else if statement.(ast.VariableAssignment).Value.(lexer.Token).Value == "false" {
					block.NewStore(constant.NewBool(false), variable)
				} else if statement.(ast.VariableAssignment).Value.(lexer.Token).Type == lexer.Float {
					block.NewStore(constant.NewFloat(types.Float, statement.(ast.VariableAssignment).Value.(lexer.Token).Value.(float64)), variable)
				}

			case ast.FuncCall:
				name := statement.(ast.VariableAssignment).Value.(ast.FuncCall).Name.Value.(string)
				arguments := FuncCall(block, statement.(ast.VariableAssignment).Value.(ast.FuncCall), variables)
				if Functions[name].Type == "int" {
					block.NewStore(block.NewCall(Functions[name].Value, arguments...), variable)
				} else if Functions[name].Type == "string" {
					hasstring = true

					str := block.NewCall(Functions[name].Value, arguments...)

					variable = str

					block.NewStore(str, variable)
				} else if Functions[name].Type == "float" {
					block.NewStore(block.NewCall(Functions[name].Value, arguments...), variable)
				} else if Functions[name].Type == "bool" {
					block.NewStore(block.NewCall(Functions[name].Value, arguments...), variable)
				}
			}

		case ast.FuncCall:
			name := statement.(ast.FuncCall).Name.Value.(string)
			arguments := FuncCall(block, statement.(ast.FuncCall), variables)
			block.NewCall(Functions[name].Value, arguments...)
		case ast.IfStatement:
			ifs++
			condition := CompileBool(block, statement.(ast.IfStatement).Condition, variables)
			if statement.(ast.IfStatement).Condition.Invert {
				//booleon := block.NewAlloca(boolType)
				//block.NewStore(constant.NewBool(true), booleon)
				cond := block.NewLoad(boolType, condition)
				//booleon2 := block.NewLoad(boolType, booleon)

				condition = block.NewXor(cond, constant.NewBool(true))

			}
			trueBlock := block.Parent.NewBlock("true" + strconv.Itoa(ifs))
			trueBlock = Compile(trueBlock, statement.(ast.IfStatement).Consequence, variables)
			falseBlock := block.Parent.NewBlock("false" + strconv.Itoa(ifs))
			falseBlock = Compile(falseBlock, statement.(ast.IfStatement).Alternative, variables)
			afterBlock := block.Parent.NewBlock("after" + strconv.Itoa(ifs))
			trueBlock.NewBr(afterBlock)
			falseBlock.NewBr(afterBlock)
			block.NewCondBr(condition, trueBlock, falseBlock)

			block = afterBlock
		case ast.WhileStatement:
			ifs++

			afterBlock := block.Parent.NewBlock("after" + strconv.Itoa(ifs))

			loopBlock := block.Parent.NewBlock("loop" + strconv.Itoa(ifs))

			loopBodyBlock := block.Parent.NewBlock("loop_body" + strconv.Itoa(ifs))
			/*loopBodyBlock = Compile(loopBodyBlock, statement.(ast.WhileStatement).Consequence, variables)
			condition := CompileBool(loopBodyBlock, statement.(ast.WhileStatement).Condition, variables)
			if statement.(ast.WhileStatement).Condition.Invert {
				cond := block.NewLoad(boolType, condition)
				condition = block.NewXor(cond, constant.NewBool(true))
			}
			loopBodyBlock.NewCondBr(condition, loopBlock, afterBlock)*/
			condition := CompileBool(loopBlock, statement.(ast.WhileStatement).Condition, variables)
			if statement.(ast.WhileStatement).Condition.Invert {
				cond := block.NewLoad(boolType, condition)
				condition = block.NewXor(cond, constant.NewBool(true))
			}
			loopBlock.NewCondBr(condition, loopBodyBlock, afterBlock)

			loopBodyBlock = Compile(loopBodyBlock, statement.(ast.WhileStatement).Consequence, variables)
			loopBodyBlock = Compile(loopBodyBlock, statement.(ast.WhileStatement).Consequence, variables)
			condition = CompileBool(loopBodyBlock, statement.(ast.WhileStatement).Condition, variables)
			if statement.(ast.WhileStatement).Condition.Invert {
				cond := block.NewLoad(boolType, condition)
				condition = block.NewXor(cond, constant.NewBool(true))
			}
			loopBodyBlock.NewCondBr(condition, loopBlock, afterBlock)

			condition = CompileBool(block, statement.(ast.WhileStatement).Condition, variables)
			if statement.(ast.WhileStatement).Condition.Invert {
				cond := block.NewLoad(boolType, condition)
				condition = block.NewXor(cond, constant.NewBool(true))
			}
			block.NewCondBr(condition, loopBlock, afterBlock)

			block = afterBlock
		case ast.ReturnStatement:
			value := statement.(ast.ReturnStatement).Value
			switch value.(type) {
			case ast.ExpressionStatement:
				exp, _ := CompileExpression(block, value.(ast.ExpressionStatement), variables)
				block.NewRet(exp)
			case lexer.Token:
				if value.(lexer.Token).Type == lexer.Int {
					block.NewRet(constant.NewInt(types.I32, int64(value.(lexer.Token).Value.(int))))
				} else if value.(lexer.Token).Type == lexer.String {
					str := constant.NewCharArrayFromString(value.(lexer.Token).Value.(string) + "\x00")
					strvar := module.NewGlobalDef("str"+strconv.Itoa(strCounter), str)
					strCounter++
					block.NewStore(str, strvar)
					strPtr := block.NewGetElementPtr(types.NewArray(uint64(str.Typ.Len), types.I8), strvar, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
					block.NewRet(strPtr)
				} else if value.(lexer.Token).Type == lexer.Identifier {
					if variables[value.(lexer.Token).Value.(string)].Type == "string" {
						block.NewRet(block.NewLoad(types.NewArray(uint64(variables[value.(lexer.Token).Value.(string)].length), types.I8), variables[value.(lexer.Token).Value.(string)].Value))
					} else if variables[value.(lexer.Token).Value.(string)].Type == "int" {
						value := block.NewLoad(types.I32, variables[value.(lexer.Token).Value.(string)].Value)
						block.NewRet(value)
					} else if variables[value.(lexer.Token).Value.(string)].Type == "bool" {
						value := block.NewLoad(types.I1, variables[value.(lexer.Token).Value.(string)].Value)
						block.NewRet(value)
					} else if variables[value.(lexer.Token).Value.(string)].Type == "float" {
						value := block.NewLoad(types.Float, variables[value.(lexer.Token).Value.(string)].Value)
						block.NewRet(value)
					}
				} else if value.(lexer.Token).Value == "true" {
					block.NewRet(constant.NewBool(true))
				} else if value.(lexer.Token).Value == "false" {
					block.NewRet(constant.NewBool(false))
				} else if value.(lexer.Token).Type == lexer.Float {
					block.NewRet(constant.NewFloat(types.Float, value.(lexer.Token).Value.(float64)))
				}
			case ast.FuncCall:
				name := value.(ast.FuncCall).Name.Value.(string)
				arguments := FuncCall(block, value.(ast.FuncCall), variables)
				block.NewRet(block.NewCall(Functions[name].Value, arguments...))
			}
		case ast.ArrayAssignment:
			arrayAssignment := statement.(ast.ArrayAssignment)
			indexesForGep := []value.Value{constant.NewInt(types.I32, 0)}
			for j := 0; j < len(arrayAssignment.Indexes); j++ {
				indexesForGep = append(indexesForGep, constant.NewInt(types.I32, int64(arrayAssignment.Indexes[j].(lexer.Token).Value.(int))))
			}
			var val value.Value
			switch arrayAssignment.Value.(type) {
			case ast.ExpressionStatement:
				val, _ = CompileExpression(block, arrayAssignment.Value.(ast.ExpressionStatement), variables)
			case lexer.Token:
				if arrayAssignment.Value.(lexer.Token).Type == lexer.Identifier {
					val = block.NewLoad(types.I32, variables[arrayAssignment.Value.(lexer.Token).Value.(string)].Value)
				} else if arrayAssignment.Value.(lexer.Token).Type == lexer.Int {
					val = constant.NewInt(types.I32, int64(arrayAssignment.Value.(lexer.Token).Value.(int)))
				}
			case ast.FuncCall:
				name := arrayAssignment.Value.(ast.FuncCall).Name.Value.(string)
				arguments := FuncCall(block, arrayAssignment.Value.(ast.FuncCall), variables)
				val = block.NewCall(Functions[name].Value, arguments...)

			}
			indexPtr := block.NewGetElementPtr(variables[arrayAssignment.Name.Value.(string)].Type.(types.Type), variables[arrayAssignment.Name.Value.(string)].Value, indexesForGep...)
			block.NewStore(val, indexPtr)

		}
		i++
	}

	return block
}
func CompileArrayAssigment(block *ir.Block, array ast.ArrayStatement, variables map[string]Variable, variable value.Value, variab Variable, indexes []int, curIndex int) value.Value /*, []int*/ {
	index := 0
	indexes = append(indexes, index)
	for index < len(array.Content) {

		switch array.Content[index].(type) {
		case ast.ExpressionStatement:
			//variableLoad := block.NewLoad(variab.Type.(types.Type), variable)
			inexesForGep := make([]value.Value, 0)
			for i := 0; i < len(indexes); i++ {
				inexesForGep = append(inexesForGep, constant.NewInt(types.I32, int64(indexes[i])))
			}

			indexPtr := block.NewGetElementPtr(variab.Type.(types.Type), variable, inexesForGep...)
			exp, _ := CompileExpression(block, array.Content[index].(ast.ExpressionStatement), variables)
			block.NewStore(exp, indexPtr)
		case lexer.Token:
			//variableLoad := block.NewLoad(variab.Type.(types.Type), variable)

			//Ok, musím tohle nějak vyřešit, protože momentálně to bere jen array dvou dymenzí...

			inexesForGep := make([]value.Value, 0)
			for i := 0; i < len(indexes); i++ {
				inexesForGep = append(inexesForGep, constant.NewInt(types.I32, int64(indexes[i])))
			}

			indexPtr := block.NewGetElementPtr(variab.Type.(types.Type), variable, inexesForGep...)
			if array.Content[index].(lexer.Token).Type == lexer.Int {
				block.NewStore(constant.NewInt(types.I32, int64(array.Content[index].(lexer.Token).Value.(int))), indexPtr)
			}
		case ast.FuncCall:
			//variableLoad := block.NewLoad(variab.Type.(types.Type), variable)
			inexesForGep := make([]value.Value, 0)
			for i := 0; i < len(indexes); i++ {
				inexesForGep = append(inexesForGep, constant.NewInt(types.I32, int64(indexes[i])))
			}

			indexPtr := block.NewGetElementPtr(variab.Type.(types.Type), variable, inexesForGep...)
			name := array.Content[index].(ast.FuncCall).Name.Value.(string)
			arguments := FuncCall(block, array.Content[index].(ast.FuncCall), variables)
			block.NewStore(block.NewCall(Functions[name].Value, arguments...), indexPtr)
		case int:
			//variableLoad := block.NewLoad(variab.Type.(types.Type), variable)
			inexesForGep := make([]value.Value, 0)
			for i := 0; i < len(indexes); i++ {
				inexesForGep = append(inexesForGep, constant.NewInt(types.I32, int64(indexes[i])))
			}

			indexPtr := block.NewGetElementPtr(variab.Type.(types.Type), variable, inexesForGep...)
			block.NewStore(constant.NewInt(types.I32, int64(array.Content[index].(int))), indexPtr)
		case ast.ArrayStatement:
			//variableLoad := block.NewLoad(variab.Type.(types.Type), variable)
			inexesForGep := make([]value.Value, 0)
			for i := 0; i < len(indexes); i++ {
				inexesForGep = append(inexesForGep, constant.NewInt(types.I32, int64(indexes[i])))
			}

			indexPtr := block.NewGetElementPtr(variab.Type.(types.Type), variable, inexesForGep...)
			//arr := block.NewLoad(variab.Type.(types.Type), indexPtr)

			variable /*, indexes*/ = CompileArrayAssigment(block, array.Content[index].(ast.ArrayStatement), variables, indexPtr, variab, indexes, curIndex+1)

		}
		indexes[curIndex]++
		index++

	}
	indexes[curIndex] = 0
	return variable /*, indexes*/
}

func ArrayType(array ast.ArrayType) types.Type {
	var ArrType types.Type
	switch array.Type.(type) {
	case string:
		if array.Type.(string) == "int" {
			ArrType = types.NewArray(uint64(array.Length.(lexer.Token).Value.(int)), types.I32)
		}
	case ast.ArrayType:
		ArrType = ArrayType(array.Type.(ast.ArrayType))
	}
	return ArrType
}

func FuncCall(block *ir.Block, function ast.FuncCall, variables map[string]Variable) []value.Value {
	args := function.Arguments
	var arguments []value.Value
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg.(type) {
		case ast.ExpressionStatement:
			express, _ := CompileExpression(block, arg.(ast.ExpressionStatement), variables)
			arguments = append(arguments, express)
		case lexer.Token:
			if arg.(lexer.Token).Type == lexer.Int {
				arguments = append(arguments, constant.NewInt(types.I32, int64(arg.(lexer.Token).Value.(int))))
			} else if arg.(lexer.Token).Type == lexer.Identifier {
				if variables[arg.(lexer.Token).Value.(string)].Type == "string" {
					s := block.NewLoad(types.I8Ptr, variables[arg.(lexer.Token).Value.(string)].Value)
					arguments = append(arguments, s)
				} else if variables[arg.(lexer.Token).Value.(string)].Type == "int" {
					i := block.NewLoad(types.I32, variables[arg.(lexer.Token).Value.(string)].Value)
					arguments = append(arguments, i)
				} else if variables[arg.(lexer.Token).Value.(string)].Type == "bool" {
					b := block.NewLoad(types.I1, variables[arg.(lexer.Token).Value.(string)].Value)
					arguments = append(arguments, b)

				} else if variables[arg.(lexer.Token).Value.(string)].Type == "float" {
					f := block.NewLoad(types.Float, variables[arg.(lexer.Token).Value.(string)].Value)
					arguments = append(arguments, f)
				}
			} else if arg.(lexer.Token).Type == lexer.String {
				str := constant.NewCharArrayFromString(arg.(lexer.Token).Value.(string) + "\x00")
				strvar := module.NewGlobalDef("str"+strconv.Itoa(strCounter), str)
				strCounter++
				block.NewStore(str, strvar)
				strPtr := block.NewGetElementPtr(types.NewArray(uint64(str.Typ.Len), types.I8), strvar, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
				arguments = append(arguments, strPtr)
			} else if arg.(lexer.Token).Value == "true" {
				arguments = append(arguments, constant.NewBool(true))
			} else if arg.(lexer.Token).Value == "false" {
				arguments = append(arguments, constant.NewBool(false))
			} else if arg.(lexer.Token).Type == lexer.Float {
				arguments = append(arguments, constant.NewFloat(types.Float, arg.(lexer.Token).Value.(float64)))
			}
		case ast.FuncCall:
			name2 := arg.(ast.FuncCall).Name.Value.(string)
			args2 := FuncCall(block, arg.(ast.FuncCall), variables)
			arguments = append(arguments, block.NewCall(Functions[name2].Value, args2...))

		}
	}
	return arguments
}
func CompileBool(block *ir.Block, condition ast.BoolStatement, variables map[string]Variable) value.Value {
	switch condition.Bool.(type) {
	case ast.BoolExpression:
		var con1 value.Value
		var con2 value.Value
		switch condition.Bool.(ast.BoolExpression).Condition1.(type) {
		case ast.ExpressionStatement:
			con1, _ = CompileExpression(block, condition.Bool.(ast.BoolExpression).Condition1.(ast.ExpressionStatement), variables)
		case lexer.Token:
			if condition.Bool.(ast.BoolExpression).Condition1.(lexer.Token).Type == lexer.Int {
				con1 = constant.NewInt(types.I32, int64(condition.Bool.(ast.BoolExpression).Condition1.(lexer.Token).Value.(int)))
			} else if condition.Bool.(ast.BoolExpression).Condition1.(lexer.Token).Type == lexer.Identifier {
				con1 = variables[condition.Bool.(ast.BoolExpression).Condition1.(lexer.Token).Value.(string)].Value
				if variables[condition.Bool.(ast.BoolExpression).Condition1.(lexer.Token).Value.(string)].Type == "string" {
					con1 = block.NewLoad(types.NewArray(uint64(variables[condition.Bool.(ast.BoolExpression).Condition1.(lexer.Token).Value.(string)].length), types.I8), con1)
				} else if variables[condition.Bool.(ast.BoolExpression).Condition1.(lexer.Token).Value.(string)].Type == "int" {
					con1 = block.NewLoad(types.I32, con1)
				} else if variables[condition.Bool.(ast.BoolExpression).Condition1.(lexer.Token).Value.(string)].Type == "bool" {
					con1 = block.NewLoad(types.I1, con1)
				}
			} else if condition.Bool.(ast.BoolExpression).Condition1.(lexer.Token).Type == lexer.String {
				con1 = constant.NewCharArrayFromString(condition.Bool.(ast.BoolExpression).Condition1.(lexer.Token).Value.(string) + "\x00")
			} else if condition.Bool.(ast.BoolExpression).Condition1.(lexer.Token).Value == "true" {
				con1 = constant.NewBool(true)
			} else if condition.Bool.(ast.BoolExpression).Condition1.(lexer.Token).Value == "false" {
				con1 = constant.NewBool(false)
			} else if condition.Bool.(ast.BoolExpression).Condition1.(lexer.Token).Type == lexer.Float {
				con1 = constant.NewFloat(types.Float, condition.Bool.(ast.BoolExpression).Condition1.(lexer.Token).Value.(float64))
			}
		}
		switch condition.Bool.(ast.BoolExpression).Condition2.(type) {
		case ast.ExpressionStatement:
			con2, _ = CompileExpression(block, condition.Bool.(ast.BoolExpression).Condition2.(ast.ExpressionStatement), variables)
		case lexer.Token:
			if condition.Bool.(ast.BoolExpression).Condition2.(lexer.Token).Type == lexer.Int {
				con2 = constant.NewInt(types.I32, int64(condition.Bool.(ast.BoolExpression).Condition2.(lexer.Token).Value.(int)))
			} else if condition.Bool.(ast.BoolExpression).Condition2.(lexer.Token).Type == lexer.Identifier {
				con2 = variables[condition.Bool.(ast.BoolExpression).Condition2.(lexer.Token).Value.(string)].Value
				if variables[condition.Bool.(ast.BoolExpression).Condition2.(lexer.Token).Value.(string)].Type == "string" {
					con2 = block.NewLoad(types.NewArray(uint64(variables[condition.Bool.(ast.BoolExpression).Condition2.(lexer.Token).Value.(string)].length), types.I8), con2)
				}
			} else if condition.Bool.(ast.BoolExpression).Condition2.(lexer.Token).Type == lexer.String {
				con2 = constant.NewCharArrayFromString(condition.Bool.(ast.BoolExpression).Condition2.(lexer.Token).Value.(string) + "\x00")
			} else if condition.Bool.(ast.BoolExpression).Condition2.(lexer.Token).Value == "true" {
				con2 = constant.NewBool(true)
			} else if condition.Bool.(ast.BoolExpression).Condition2.(lexer.Token).Value == "false" {
				con2 = constant.NewBool(false)
			} else if condition.Bool.(ast.BoolExpression).Condition2.(lexer.Token).Type == lexer.Float {
				con2 = constant.NewFloat(types.Float, condition.Bool.(ast.BoolExpression).Condition2.(lexer.Token).Value.(float64))
			}
		}
		if condition.Bool.(ast.BoolExpression).Operator.Type == lexer.DoubleEqual {
			return block.NewICmp(enum.IPredEQ, con1, con2)
		} else if condition.Bool.(ast.BoolExpression).Operator.Type == lexer.NotEqual {
			return block.NewICmp(enum.IPredNE, con1, con2)
		} else if condition.Bool.(ast.BoolExpression).Operator.Type == lexer.MoreThan {
			return block.NewICmp(enum.IPredSGT, con1, con2)
		} else if condition.Bool.(ast.BoolExpression).Operator.Type == lexer.LessThan {
			return block.NewICmp(enum.IPredSLT, con1, con2)
		}
	case lexer.Token:
		if condition.Bool.(lexer.Token).Value.(string) == "true" {
			return constant.NewBool(true)
		} else if condition.Bool.(lexer.Token).Value.(string) == "false" {
			return constant.NewBool(false)
		} else if condition.Bool.(lexer.Token).Type == lexer.Identifier {
			if variables[condition.Bool.(lexer.Token).Value.(string)].Type == "bool" {
				variab := block.NewLoad(types.I1, variables[condition.Bool.(lexer.Token).Value.(string)].Value)
				return variab
			} else if Functions[condition.Bool.(lexer.Token).Value.(string)].Type == "bool" {
				return block.NewCall(Functions[condition.Bool.(lexer.Token).Value.(string)].Value)
			}
		}
	case ast.FuncCall:
		name := condition.Bool.(ast.FuncCall).Name.Value.(string)
		arguments := FuncCall(block, condition.Bool.(ast.FuncCall), variables)
		return block.NewCall(Functions[name].Value, arguments...)
	}
	return nil
}

func CompileExpression(block *ir.Block, expression ast.ExpressionStatement, variables map[string]Variable) (value.Value, string) {
	if expression.Operator.Type == lexer.Plus {
		switch expression.Left.(type) {
		case ast.ExpressionStatement:
			switch expression.Right.(type) {
			case ast.ExpressionStatement:
				exp, exp_type := CompileExpression(block, expression.Left.(ast.ExpressionStatement), variables)
				exp2, exp2_type := CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables)
				if exp_type == "int" && exp2_type == "int" {
					return block.NewAdd(exp, exp2), "int"
				} else if exp_type == "float" && exp2_type == "float" {
					return block.NewFAdd(exp, exp2), "float"
				}
			case lexer.Token:
				if expression.Right.(lexer.Token).Type == lexer.Int {
					exp, _ := CompileExpression(block, expression.Left.(ast.ExpressionStatement), variables)
					return block.NewAdd(exp, constant.NewInt(types.I32, int64(expression.Right.(lexer.Token).Value.(int)))), "int"
				} else if expression.Right.(lexer.Token).Type == lexer.Identifier {
					if variables[expression.Right.(lexer.Token).Value.(string)].Type == "int" {
						exp, _ := CompileExpression(block, expression.Left.(ast.ExpressionStatement), variables)
						return block.NewAdd(exp, block.NewLoad(types.I32, variables[expression.Right.(lexer.Token).Value.(string)].Value)), "int"
					} else if variables[expression.Right.(lexer.Token).Value.(string)].Type == "float" {
						exp, _ := CompileExpression(block, expression.Left.(ast.ExpressionStatement), variables)
						return block.NewFAdd(exp, block.NewLoad(types.Float, variables[expression.Right.(lexer.Token).Value.(string)].Value)), "float"
					}

				} else if expression.Right.(lexer.Token).Type == lexer.Float {
					exp, _ := CompileExpression(block, expression.Left.(ast.ExpressionStatement), variables)
					return block.NewFAdd(exp, constant.NewFloat(types.Float, expression.Right.(lexer.Token).Value.(float64))), "float"
				}
			case ast.FuncCall:
				name := expression.Right.(ast.FuncCall).Name.Value.(string)
				arguments := FuncCall(block, expression.Right.(ast.FuncCall), variables)
				if Functions[name].Type == "int" {
					exp, _ := CompileExpression(block, expression.Left.(ast.ExpressionStatement), variables)
					return block.NewAdd(exp, block.NewCall(Functions[name].Value, arguments...)), "int"
				} else if Functions[name].Type == "float" {
					exp, _ := CompileExpression(block, expression.Left.(ast.ExpressionStatement), variables)
					return block.NewFAdd(exp, block.NewCall(Functions[name].Value, arguments...)), "float"
				}
			}
		case lexer.Token:
			switch expression.Right.(type) {
			case ast.ExpressionStatement:
				if expression.Left.(lexer.Token).Type == lexer.Int {
					exp, _ := CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables)
					return block.NewAdd(constant.NewInt(types.I32, int64(expression.Left.(lexer.Token).Value.(int))), exp), "int"
				} else if expression.Left.(lexer.Token).Type == lexer.Identifier {
					if variables[expression.Left.(lexer.Token).Value.(string)].Type == "int" {
						variable := block.NewLoad(types.I32, variables[expression.Left.(lexer.Token).Value.(string)].Value)
						exp, _ := CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables)
						return block.NewAdd(variable, exp), "int"
					} else if variables[expression.Left.(lexer.Token).Value.(string)].Type == "string" {
						str1 := variables[expression.Left.(lexer.Token).Value.(string)].Value
						str2, _ := CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables)

						return block.NewCall(Functions["append"].Value, str1, str2), "string"
					} else if variables[expression.Left.(lexer.Token).Value.(string)].Type == "float" {
						variable := block.NewLoad(types.Float, variables[expression.Left.(lexer.Token).Value.(string)].Value)
						exp, _ := CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables)
						return block.NewFAdd(variable, exp), "float"
					}
				} else if expression.Left.(lexer.Token).Type == lexer.String {
					length := len(expression.Left.(lexer.Token).Value.(string) + "\x00")
					str1 := module.NewGlobalDef("str"+strconv.Itoa(strCounter), constant.NewCharArrayFromString(expression.Left.(lexer.Token).Value.(string)+"\x00"))
					strCounter++
					str2, _ := CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables)

					return block.NewCall(Functions["append"].Value, block.NewGetElementPtr(types.NewArray(uint64(length), types.I8), str1, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0)), str2), "string"

				} else if expression.Left.(lexer.Token).Type == lexer.Float {
					exp, _ := CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables)
					return block.NewFAdd(constant.NewFloat(types.Float, expression.Left.(lexer.Token).Value.(float64)), exp), "float"
				}
			case lexer.Token:
				if expression.Left.(lexer.Token).Type == lexer.Int && expression.Right.(lexer.Token).Type == lexer.Int {
					return block.NewAdd(constant.NewInt(types.I32, int64(expression.Left.(lexer.Token).Value.(int))), constant.NewInt(types.I32, int64(expression.Right.(lexer.Token).Value.(int)))), "int"
				} else if expression.Left.(lexer.Token).Type == lexer.Identifier && expression.Right.(lexer.Token).Type == lexer.Int {
					variable := block.NewLoad(types.I32, variables[expression.Left.(lexer.Token).Value.(string)].Value)
					return block.NewAdd(variable, constant.NewInt(types.I32, int64(expression.Right.(lexer.Token).Value.(int)))), "int"
				} else if expression.Left.(lexer.Token).Type == lexer.Int && expression.Right.(lexer.Token).Type == lexer.Identifier {
					variable := block.NewLoad(types.I32, variables[expression.Right.(lexer.Token).Value.(string)].Value)
					return block.NewAdd(constant.NewInt(types.I32, int64(expression.Left.(lexer.Token).Value.(int))), variable), "int"
				} else if expression.Left.(lexer.Token).Type == lexer.Identifier && expression.Right.(lexer.Token).Type == lexer.Identifier {
					if (variables[expression.Left.(lexer.Token).Value.(string)].Type == "int") && (variables[expression.Right.(lexer.Token).Value.(string)].Type == "int") {
						variable := block.NewLoad(types.I32, variables[expression.Left.(lexer.Token).Value.(string)].Value)
						variable2 := block.NewLoad(types.I32, variables[expression.Right.(lexer.Token).Value.(string)].Value)
						return block.NewAdd(variable, variable2), "int"
					} else if (variables[expression.Left.(lexer.Token).Value.(string)].Type == "string") && (variables[expression.Right.(lexer.Token).Value.(string)].Type == "string") {

						return block.NewCall(Functions["append"].Value, block.NewLoad(types.I8Ptr, variables[expression.Left.(lexer.Token).Value.(string)].Value), block.NewLoad(types.I8Ptr, variables[expression.Right.(lexer.Token).Value.(string)].Value)), "string"
					} else if (variables[expression.Left.(lexer.Token).Value.(string)].Type == "float") && (variables[expression.Right.(lexer.Token).Value.(string)].Type == "float") {
						variable := block.NewLoad(types.Float, variables[expression.Left.(lexer.Token).Value.(string)].Value)
						variable2 := block.NewLoad(types.Float, variables[expression.Right.(lexer.Token).Value.(string)].Value)
						return block.NewFAdd(variable, variable2), "float"
					}
				} else if expression.Left.(lexer.Token).Type == lexer.String && expression.Right.(lexer.Token).Type == lexer.String {
					str1 := expression.Left.(lexer.Token).Value.(string)
					str2 := expression.Right.(lexer.Token).Value.(string) + "\x00"
					str := constant.NewCharArrayFromString(str1 + str2)
					strvar := module.NewGlobalDef("str"+strconv.Itoa(strCounter), str)
					strCounter++

					return block.NewGetElementPtr(types.NewArray(uint64(str.Typ.Len), types.I8), strvar, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0)), "string"
				} else if expression.Left.(lexer.Token).Type == lexer.String && expression.Right.(lexer.Token).Type == lexer.Identifier {
					str1 := constant.NewCharArrayFromString(expression.Left.(lexer.Token).Value.(string) + "\x00")
					str1var := module.NewGlobalDef("str"+strconv.Itoa(strCounter), str1)
					strCounter++
					block.NewStore(str1, str1var)
					str1in := block.NewGetElementPtr(types.NewArray(uint64(str1.Typ.Len), types.I8), str1var, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
					str2 := block.NewLoad(types.I8Ptr, variables[expression.Right.(lexer.Token).Value.(string)].Value)
					return block.NewCall(Functions["append"].Value, str1in, str2), "string"
				} else if expression.Left.(lexer.Token).Type == lexer.Identifier && expression.Right.(lexer.Token).Type == lexer.String {
					str1 := block.NewLoad(types.I8Ptr, variables[expression.Left.(lexer.Token).Value.(string)].Value)
					str2 := constant.NewCharArrayFromString(expression.Right.(lexer.Token).Value.(string) + "\x00")
					str2var := module.NewGlobalDef("str"+strconv.Itoa(strCounter), str2)
					strCounter++
					block.NewStore(str2, str2var)
					str2in := block.NewGetElementPtr(types.NewArray(uint64(str2.Typ.Len), types.I8), str2var, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
					return block.NewCall(Functions["append"].Value, str1, str2in), "string"
				} else if expression.Left.(lexer.Token).Type == lexer.Float && expression.Right.(lexer.Token).Type == lexer.Float {
					return block.NewFAdd(constant.NewFloat(types.Float, expression.Left.(lexer.Token).Value.(float64)), constant.NewFloat(types.Float, expression.Right.(lexer.Token).Value.(float64))), "float"
				} else if expression.Left.(lexer.Token).Type == lexer.Float && expression.Right.(lexer.Token).Type == lexer.Identifier {
					variable := block.NewLoad(types.Float, variables[expression.Right.(lexer.Token).Value.(string)].Value)
					return block.NewFAdd(constant.NewFloat(types.Float, expression.Left.(lexer.Token).Value.(float64)), variable), "float"
				} else if expression.Left.(lexer.Token).Type == lexer.Identifier && expression.Right.(lexer.Token).Type == lexer.Float {
					variable := block.NewLoad(types.Float, variables[expression.Left.(lexer.Token).Value.(string)].Value)
					return block.NewFAdd(variable, constant.NewFloat(types.Float, expression.Right.(lexer.Token).Value.(float64))), "float"
				}
			case ast.FuncCall:
				name2 := expression.Right.(ast.FuncCall).Name.Value.(string)

				arguments2 := FuncCall(block, expression.Right.(ast.FuncCall), variables)
				if Functions[name2].Type == "int" && expression.Left.(lexer.Token).Type == lexer.Int {
					return block.NewAdd(constant.NewInt(types.I32, int64(expression.Left.(lexer.Token).Value.(int))), block.NewCall(Functions[name2].Value, arguments2...)), "int"
				} else if Functions[name2].Type == "int" && expression.Left.(lexer.Token).Type == lexer.Identifier {
					variable := block.NewLoad(types.I32, variables[expression.Left.(lexer.Token).Value.(string)].Value)
					return block.NewAdd(variable, block.NewCall(Functions[name2].Value, arguments2...)), "int"
				} else if Functions[name2].Type == "float" && expression.Left.(lexer.Token).Type == lexer.Float {
					return block.NewFAdd(constant.NewFloat(types.Float, expression.Left.(lexer.Token).Value.(float64)), block.NewCall(Functions[name2].Value, arguments2...)), "float"
				} else if Functions[name2].Type == "float" && expression.Left.(lexer.Token).Type == lexer.Identifier {
					variable := block.NewLoad(types.Float, variables[expression.Left.(lexer.Token).Value.(string)].Value)
					return block.NewFAdd(variable, block.NewCall(Functions[name2].Value, arguments2...)), "float"
				} else if Functions[name2].Type == "string" && expression.Left.(lexer.Token).Type == lexer.String {
					str := constant.NewCharArrayFromString(expression.Left.(lexer.Token).Value.(string) + "\x00")
					strvar := module.NewGlobalDef("str"+strconv.Itoa(strCounter), str)
					strCounter++
					block.NewStore(str, strvar)
					strin := block.NewGetElementPtr(types.NewArray(uint64(str.Typ.Len), types.I8), strvar, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
					return block.NewCall(Functions["append"].Value, strin, block.NewCall(Functions[name2].Value, arguments2...)), "string"
				} else if Functions[name2].Type == "string" && expression.Left.(lexer.Token).Type == lexer.Identifier {
					str := block.NewLoad(types.I8Ptr, variables[expression.Left.(lexer.Token).Value.(string)].Value)
					return block.NewCall(Functions["append"].Value, str, block.NewCall(Functions[name2].Value, arguments2...)), "string"
				}

			}
		case ast.FuncCall:
			name := expression.Left.(ast.FuncCall).Name.Value.(string)
			arguments := FuncCall(block, expression.Left.(ast.FuncCall), variables)
			switch expression.Right.(type) {
			case ast.ExpressionStatement:
				if Functions[name].Type == "int" {
					exp, _ := CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables)
					return block.NewAdd(block.NewCall(Functions[name].Value, arguments...), exp), "int"
				} else if Functions[name].Type == "float" {
					exp, _ := CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables)
					return block.NewFAdd(block.NewCall(Functions[name].Value, arguments...), exp), "float"
				}
			case lexer.Token:
				if expression.Right.(lexer.Token).Type == lexer.Int {
					return block.NewAdd(block.NewCall(Functions[name].Value, arguments...), constant.NewInt(types.I32, int64(expression.Right.(lexer.Token).Value.(int)))), "int"
				} else if expression.Right.(lexer.Token).Type == lexer.Identifier {
					if variables[expression.Right.(lexer.Token).Value.(string)].Type == "int" {
						return block.NewAdd(block.NewCall(Functions[name].Value, arguments...), block.NewLoad(types.I32, variables[expression.Right.(lexer.Token).Value.(string)].Value)), "int"
					} else if variables[expression.Right.(lexer.Token).Value.(string)].Type == "float" {
						return block.NewFAdd(block.NewCall(Functions[name].Value, arguments...), block.NewLoad(types.Float, variables[expression.Right.(lexer.Token).Value.(string)].Value)), "float"
					}
				} else if expression.Right.(lexer.Token).Type == lexer.String {
					str := constant.NewCharArrayFromString(expression.Right.(lexer.Token).Value.(string) + "\x00")
					return block.NewCall(Functions["append"].Value, block.NewCall(Functions[name].Value, arguments...), str), "string"
				} else if expression.Right.(lexer.Token).Type == lexer.Float {
					return block.NewFAdd(block.NewCall(Functions[name].Value, arguments...), constant.NewFloat(types.Float, expression.Right.(lexer.Token).Value.(float64))), "float"
				}
			case ast.FuncCall:
				name2 := expression.Right.(ast.FuncCall).Name.Value.(string)
				arguments2 := FuncCall(block, expression.Right.(ast.FuncCall), variables)
				if Functions[name].Type == "int" && Functions[name2].Type == "int" {
					return block.NewAdd(block.NewCall(Functions[name].Value, arguments...), block.NewCall(Functions[name2].Value, arguments2...)), "int"
				} else if Functions[name].Type == "string" && Functions[name2].Type == "string" {
					return block.NewCall(Functions["append"].Value, block.NewCall(Functions[name].Value, arguments...), block.NewCall(Functions[name2].Value, arguments2...)), "string"
				} else if Functions[name].Type == "float" && Functions[name2].Type == "float" {
					return block.NewFAdd(block.NewCall(Functions[name].Value, arguments...), block.NewCall(Functions[name2].Value, arguments2...)), "float"
				}
			}
		}
	} else if expression.Operator.Type == lexer.Minus {
		switch expression.Left.(type) {
		case ast.ExpressionStatement:
			switch expression.Right.(type) {
			case ast.ExpressionStatement:
				exp, exp_type := CompileExpression(block, expression.Left.(ast.ExpressionStatement), variables)
				exp2, exp2_type := CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables)
				if exp_type == "int" && exp2_type == "int" {
					return block.NewSub(exp, exp2), "int"
				} else if exp_type == "float" && exp2_type == "float" {
					return block.NewFSub(exp, exp2), "float"
				}
			case lexer.Token:
				if expression.Right.(lexer.Token).Type == lexer.Int {
					exp, _ := CompileExpression(block, expression.Left.(ast.ExpressionStatement), variables)
					return block.NewSub(exp, constant.NewInt(types.I32, int64(expression.Right.(lexer.Token).Value.(int)))), "int"
				} else if expression.Right.(lexer.Token).Type == lexer.Identifier {
					if variables[expression.Right.(lexer.Token).Value.(string)].Type == "int" {
						exp, _ := CompileExpression(block, expression.Left.(ast.ExpressionStatement), variables)
						return block.NewSub(exp, block.NewLoad(types.I32, variables[expression.Right.(lexer.Token).Value.(string)].Value)), "int"
					} else if variables[expression.Right.(lexer.Token).Value.(string)].Type == "float" {
						exp, _ := CompileExpression(block, expression.Left.(ast.ExpressionStatement), variables)
						return block.NewFSub(exp, block.NewLoad(types.Float, variables[expression.Right.(lexer.Token).Value.(string)].Value)), "float"
					}

				} else if expression.Right.(lexer.Token).Type == lexer.Float {
					exp, _ := CompileExpression(block, expression.Left.(ast.ExpressionStatement), variables)
					return block.NewFSub(exp, constant.NewFloat(types.Float, expression.Right.(lexer.Token).Value.(float64))), "float"
				}
			case ast.FuncCall:
				name := expression.Right.(ast.FuncCall).Name.Value.(string)
				arguments := FuncCall(block, expression.Right.(ast.FuncCall), variables)
				if Functions[name].Type == "int" {
					exp, _ := CompileExpression(block, expression.Left.(ast.ExpressionStatement), variables)
					return block.NewSub(exp, block.NewCall(Functions[name].Value, arguments...)), "int"
				} else if Functions[name].Type == "float" {
					exp, _ := CompileExpression(block, expression.Left.(ast.ExpressionStatement), variables)
					return block.NewFSub(exp, block.NewCall(Functions[name].Value, arguments...)), "float"
				}
			}
		case lexer.Token:
			switch expression.Right.(type) {
			case ast.ExpressionStatement:
				if expression.Left.(lexer.Token).Type == lexer.Int {
					exp, _ := CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables)
					return block.NewSub(constant.NewInt(types.I32, int64(expression.Left.(lexer.Token).Value.(int))), exp), "int"
				} else if expression.Left.(lexer.Token).Type == lexer.Identifier {
					if variables[expression.Left.(lexer.Token).Value.(string)].Type == "int" {
						variable := block.NewLoad(types.I32, variables[expression.Left.(lexer.Token).Value.(string)].Value)
						exp, _ := CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables)
						return block.NewSub(variable, exp), "int"
					} else if variables[expression.Left.(lexer.Token).Value.(string)].Type == "float" {
						variable := block.NewLoad(types.Float, variables[expression.Left.(lexer.Token).Value.(string)].Value)
						exp, _ := CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables)
						return block.NewFSub(variable, exp), "float"
					}
				} else if expression.Left.(lexer.Token).Type == lexer.Float {
					exp, _ := CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables)
					return block.NewFSub(constant.NewFloat(types.Float, expression.Left.(lexer.Token).Value.(float64)), exp), "float"
				}
			case lexer.Token:
				if expression.Left.(lexer.Token).Type == lexer.Int && expression.Right.(lexer.Token).Type == lexer.Int {
					return block.NewSub(constant.NewInt(types.I32, int64(expression.Left.(lexer.Token).Value.(int))), constant.NewInt(types.I32, int64(expression.Right.(lexer.Token).Value.(int)))), "int"
				} else if expression.Left.(lexer.Token).Type == lexer.Identifier && expression.Right.(lexer.Token).Type == lexer.Int {
					variable := block.NewLoad(types.I32, variables[expression.Left.(lexer.Token).Value.(string)].Value)
					return block.NewSub(variable, constant.NewInt(types.I32, int64(expression.Right.(lexer.Token).Value.(int)))), "int"
				} else if expression.Left.(lexer.Token).Type == lexer.Int && expression.Right.(lexer.Token).Type == lexer.Identifier {
					variable := block.NewLoad(types.I32, variables[expression.Right.(lexer.Token).Value.(string)].Value)
					return block.NewSub(constant.NewInt(types.I32, int64(expression.Left.(lexer.Token).Value.(int))), variable), "int"
				} else if expression.Left.(lexer.Token).Type == lexer.Identifier && expression.Right.(lexer.Token).Type == lexer.Identifier {
					if (variables[expression.Left.(lexer.Token).Value.(string)].Type == "int") && (variables[expression.Right.(lexer.Token).Value.(string)].Type == "int") {
						variable := block.NewLoad(types.I32, variables[expression.Left.(lexer.Token).Value.(string)].Value)
						variable2 := block.NewLoad(types.I32, variables[expression.Right.(lexer.Token).Value.(string)].Value)
						return block.NewSub(variable, variable2), "int"
					} else if (variables[expression.Left.(lexer.Token).Value.(string)].Type == "float") && (variables[expression.Right.(lexer.Token).Value.(string)].Type == "float") {
						variable := block.NewLoad(types.Float, variables[expression.Left.(lexer.Token).Value.(string)].Value)
						variable2 := block.NewLoad(types.Float, variables[expression.Right.(lexer.Token).Value.(string)].Value)
						return block.NewFSub(variable, variable2), "float"
					}
				} else if expression.Left.(lexer.Token).Type == lexer.Float && expression.Right.(lexer.Token).Type == lexer.Float {
					return block.NewFSub(constant.NewFloat(types.Float, expression.Left.(lexer.Token).Value.(float64)), constant.NewFloat(types.Float, expression.Right.(lexer.Token).Value.(float64))), "float"
				} else if expression.Left.(lexer.Token).Type == lexer.Float && expression.Right.(lexer.Token).Type == lexer.Identifier {
					variable := block.NewLoad(types.Float, variables[expression.Right.(lexer.Token).Value.(string)].Value)
					return block.NewFSub(constant.NewFloat(types.Float, expression.Left.(lexer.Token).Value.(float64)), variable), "float"
				} else if expression.Left.(lexer.Token).Type == lexer.Identifier && expression.Right.(lexer.Token).Type == lexer.Float {
					variable := block.NewLoad(types.Float, variables[expression.Left.(lexer.Token).Value.(string)].Value)
					return block.NewFSub(variable, constant.NewFloat(types.Float, expression.Right.(lexer.Token).Value.(float64))), "float"
				}
			case ast.FuncCall:
				name2 := expression.Right.(ast.FuncCall).Name.Value.(string)

				arguments2 := FuncCall(block, expression.Right.(ast.FuncCall), variables)
				if Functions[name2].Type == "int" && expression.Left.(lexer.Token).Type == lexer.Int {
					return block.NewSub(constant.NewInt(types.I32, int64(expression.Left.(lexer.Token).Value.(int))), block.NewCall(Functions[name2].Value, arguments2...)), "int"
				} else if Functions[name2].Type == "int" && expression.Left.(lexer.Token).Type == lexer.Identifier {
					variable := block.NewLoad(types.I32, variables[expression.Left.(lexer.Token).Value.(string)].Value)
					return block.NewSub(variable, block.NewCall(Functions[name2].Value, arguments2...)), "int"
				} else if Functions[name2].Type == "float" && expression.Left.(lexer.Token).Type == lexer.Float {
					return block.NewFSub(constant.NewFloat(types.Float, expression.Left.(lexer.Token).Value.(float64)), block.NewCall(Functions[name2].Value, arguments2...)), "float"
				} else if Functions[name2].Type == "float" && expression.Left.(lexer.Token).Type == lexer.Identifier {
					variable := block.NewLoad(types.Float, variables[expression.Left.(lexer.Token).Value.(string)].Value)
					return block.NewFSub(variable, block.NewCall(Functions[name2].Value, arguments2...)), "float"
				}

			}
		case ast.FuncCall:
			name := expression.Left.(ast.FuncCall).Name.Value.(string)
			arguments := FuncCall(block, expression.Left.(ast.FuncCall), variables)
			switch expression.Right.(type) {
			case ast.ExpressionStatement:
				if Functions[name].Type == "int" {
					exp, _ := CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables)
					return block.NewSub(block.NewCall(Functions[name].Value, arguments...), exp), "int"
				} else if Functions[name].Type == "float" {
					exp, _ := CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables)
					return block.NewFSub(block.NewCall(Functions[name].Value, arguments...), exp), "float"
				}
			case lexer.Token:
				if expression.Right.(lexer.Token).Type == lexer.Int {
					return block.NewSub(block.NewCall(Functions[name].Value, arguments...), constant.NewInt(types.I32, int64(expression.Right.(lexer.Token).Value.(int)))), "int"
				} else if expression.Right.(lexer.Token).Type == lexer.Identifier {
					if variables[expression.Right.(lexer.Token).Value.(string)].Type == "int" {
						return block.NewSub(block.NewCall(Functions[name].Value, arguments...), block.NewLoad(types.I32, variables[expression.Right.(lexer.Token).Value.(string)].Value)), "int"
					} else if variables[expression.Right.(lexer.Token).Value.(string)].Type == "float" {
						return block.NewFSub(block.NewCall(Functions[name].Value, arguments...), block.NewLoad(types.Float, variables[expression.Right.(lexer.Token).Value.(string)].Value)), "float"
					}
				} else if expression.Right.(lexer.Token).Type == lexer.Float {
					return block.NewFSub(block.NewCall(Functions[name].Value, arguments...), constant.NewFloat(types.Float, expression.Right.(lexer.Token).Value.(float64))), "float"
				}
			case ast.FuncCall:
				name2 := expression.Right.(ast.FuncCall).Name.Value.(string)
				arguments2 := FuncCall(block, expression.Right.(ast.FuncCall), variables)
				if Functions[name].Type == "int" && Functions[name2].Type == "int" {
					return block.NewSub(block.NewCall(Functions[name].Value, arguments...), block.NewCall(Functions[name2].Value, arguments2...)), "int"
				} else if Functions[name].Type == "float" && Functions[name2].Type == "float" {
					return block.NewFSub(block.NewCall(Functions[name].Value, arguments...), block.NewCall(Functions[name2].Value, arguments2...)), "float"
				}
			}
		}
	} else if expression.Operator.Type == lexer.Multiply {
		switch expression.Left.(type) {
		case ast.ExpressionStatement:
			switch expression.Right.(type) {
			case ast.ExpressionStatement:
				exp, exp_type := CompileExpression(block, expression.Left.(ast.ExpressionStatement), variables)
				exp2, exp2_type := CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables)
				if exp_type == "int" && exp2_type == "int" {
					return block.NewMul(exp, exp2), "int"
				} else if exp_type == "float" && exp2_type == "float" {
					return block.NewFMul(exp, exp2), "float"
				}
			case lexer.Token:
				if expression.Right.(lexer.Token).Type == lexer.Int {
					exp, _ := CompileExpression(block, expression.Left.(ast.ExpressionStatement), variables)
					return block.NewMul(exp, constant.NewInt(types.I32, int64(expression.Right.(lexer.Token).Value.(int)))), "int"
				} else if expression.Right.(lexer.Token).Type == lexer.Identifier {
					if variables[expression.Right.(lexer.Token).Value.(string)].Type == "int" {
						exp, _ := CompileExpression(block, expression.Left.(ast.ExpressionStatement), variables)
						return block.NewMul(exp, block.NewLoad(types.I32, variables[expression.Right.(lexer.Token).Value.(string)].Value)), "int"
					} else if variables[expression.Right.(lexer.Token).Value.(string)].Type == "float" {
						exp, _ := CompileExpression(block, expression.Left.(ast.ExpressionStatement), variables)
						return block.NewFMul(exp, block.NewLoad(types.Float, variables[expression.Right.(lexer.Token).Value.(string)].Value)), "float"
					}

				} else if expression.Right.(lexer.Token).Type == lexer.Float {
					exp, _ := CompileExpression(block, expression.Left.(ast.ExpressionStatement), variables)
					return block.NewFMul(exp, constant.NewFloat(types.Float, expression.Right.(lexer.Token).Value.(float64))), "float"
				}
			case ast.FuncCall:
				name := expression.Right.(ast.FuncCall).Name.Value.(string)
				arguments := FuncCall(block, expression.Right.(ast.FuncCall), variables)
				if Functions[name].Type == "int" {
					exp, _ := CompileExpression(block, expression.Left.(ast.ExpressionStatement), variables)
					return block.NewMul(exp, block.NewCall(Functions[name].Value, arguments...)), "int"
				} else if Functions[name].Type == "float" {
					exp, _ := CompileExpression(block, expression.Left.(ast.ExpressionStatement), variables)
					return block.NewFMul(exp, block.NewCall(Functions[name].Value, arguments...)), "float"
				}
			}
		case lexer.Token:
			switch expression.Right.(type) {
			case ast.ExpressionStatement:
				if expression.Left.(lexer.Token).Type == lexer.Int {
					exp, _ := CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables)
					return block.NewMul(constant.NewInt(types.I32, int64(expression.Left.(lexer.Token).Value.(int))), exp), "int"
				} else if expression.Left.(lexer.Token).Type == lexer.Identifier {
					if variables[expression.Left.(lexer.Token).Value.(string)].Type == "int" {
						variable := block.NewLoad(types.I32, variables[expression.Left.(lexer.Token).Value.(string)].Value)
						exp, _ := CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables)
						return block.NewMul(variable, exp), "int"
					} else if variables[expression.Left.(lexer.Token).Value.(string)].Type == "float" {
						variable := block.NewLoad(types.Float, variables[expression.Left.(lexer.Token).Value.(string)].Value)
						exp, _ := CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables)
						return block.NewFMul(variable, exp), "float"
					}
				} else if expression.Left.(lexer.Token).Type == lexer.Float {
					exp, _ := CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables)
					return block.NewFMul(constant.NewFloat(types.Float, expression.Left.(lexer.Token).Value.(float64)), exp), "float"
				}
			case lexer.Token:
				if expression.Left.(lexer.Token).Type == lexer.Int && expression.Right.(lexer.Token).Type == lexer.Int {
					return block.NewMul(constant.NewInt(types.I32, int64(expression.Left.(lexer.Token).Value.(int))), constant.NewInt(types.I32, int64(expression.Right.(lexer.Token).Value.(int)))), "int"
				} else if expression.Left.(lexer.Token).Type == lexer.Identifier && expression.Right.(lexer.Token).Type == lexer.Int {
					variable := block.NewLoad(types.I32, variables[expression.Left.(lexer.Token).Value.(string)].Value)
					return block.NewMul(variable, constant.NewInt(types.I32, int64(expression.Right.(lexer.Token).Value.(int)))), "int"
				} else if expression.Left.(lexer.Token).Type == lexer.Int && expression.Right.(lexer.Token).Type == lexer.Identifier {
					variable := block.NewLoad(types.I32, variables[expression.Right.(lexer.Token).Value.(string)].Value)
					return block.NewMul(constant.NewInt(types.I32, int64(expression.Left.(lexer.Token).Value.(int))), variable), "int"
				} else if expression.Left.(lexer.Token).Type == lexer.Identifier && expression.Right.(lexer.Token).Type == lexer.Identifier {
					if (variables[expression.Left.(lexer.Token).Value.(string)].Type == "int") && (variables[expression.Right.(lexer.Token).Value.(string)].Type == "int") {
						variable := block.NewLoad(types.I32, variables[expression.Left.(lexer.Token).Value.(string)].Value)
						variable2 := block.NewLoad(types.I32, variables[expression.Right.(lexer.Token).Value.(string)].Value)
						return block.NewMul(variable, variable2), "int"
					} else if (variables[expression.Left.(lexer.Token).Value.(string)].Type == "float") && (variables[expression.Right.(lexer.Token).Value.(string)].Type == "float") {
						variable := block.NewLoad(types.Float, variables[expression.Left.(lexer.Token).Value.(string)].Value)
						variable2 := block.NewLoad(types.Float, variables[expression.Right.(lexer.Token).Value.(string)].Value)
						return block.NewFMul(variable, variable2), "float"
					}
				} else if expression.Left.(lexer.Token).Type == lexer.Float && expression.Right.(lexer.Token).Type == lexer.Float {
					return block.NewFMul(constant.NewFloat(types.Float, expression.Left.(lexer.Token).Value.(float64)), constant.NewFloat(types.Float, expression.Right.(lexer.Token).Value.(float64))), "float"
				} else if expression.Left.(lexer.Token).Type == lexer.Float && expression.Right.(lexer.Token).Type == lexer.Identifier {
					variable := block.NewLoad(types.Float, variables[expression.Right.(lexer.Token).Value.(string)].Value)
					return block.NewFMul(constant.NewFloat(types.Float, expression.Left.(lexer.Token).Value.(float64)), variable), "float"
				} else if expression.Left.(lexer.Token).Type == lexer.Identifier && expression.Right.(lexer.Token).Type == lexer.Float {
					variable := block.NewLoad(types.Float, variables[expression.Left.(lexer.Token).Value.(string)].Value)
					return block.NewFMul(variable, constant.NewFloat(types.Float, expression.Right.(lexer.Token).Value.(float64))), "float"
				}
			case ast.FuncCall:
				name2 := expression.Right.(ast.FuncCall).Name.Value.(string)

				arguments2 := FuncCall(block, expression.Right.(ast.FuncCall), variables)
				if Functions[name2].Type == "int" && expression.Left.(lexer.Token).Type == lexer.Int {
					return block.NewMul(constant.NewInt(types.I32, int64(expression.Left.(lexer.Token).Value.(int))), block.NewCall(Functions[name2].Value, arguments2...)), "int"
				} else if Functions[name2].Type == "int" && expression.Left.(lexer.Token).Type == lexer.Identifier {
					variable := block.NewLoad(types.I32, variables[expression.Left.(lexer.Token).Value.(string)].Value)
					return block.NewMul(variable, block.NewCall(Functions[name2].Value, arguments2...)), "int"
				} else if Functions[name2].Type == "float" && expression.Left.(lexer.Token).Type == lexer.Float {
					return block.NewFMul(constant.NewFloat(types.Float, expression.Left.(lexer.Token).Value.(float64)), block.NewCall(Functions[name2].Value, arguments2...)), "float"
				} else if Functions[name2].Type == "float" && expression.Left.(lexer.Token).Type == lexer.Identifier {
					variable := block.NewLoad(types.Float, variables[expression.Left.(lexer.Token).Value.(string)].Value)
					return block.NewFMul(variable, block.NewCall(Functions[name2].Value, arguments2...)), "float"
				}

			}
		case ast.FuncCall:
			name := expression.Left.(ast.FuncCall).Name.Value.(string)
			arguments := FuncCall(block, expression.Left.(ast.FuncCall), variables)
			switch expression.Right.(type) {
			case ast.ExpressionStatement:
				if Functions[name].Type == "int" {
					exp, _ := CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables)
					return block.NewMul(block.NewCall(Functions[name].Value, arguments...), exp), "int"
				} else if Functions[name].Type == "float" {
					exp, _ := CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables)
					return block.NewFMul(block.NewCall(Functions[name].Value, arguments...), exp), "float"
				}
			case lexer.Token:
				if expression.Right.(lexer.Token).Type == lexer.Int {
					return block.NewMul(block.NewCall(Functions[name].Value, arguments...), constant.NewInt(types.I32, int64(expression.Right.(lexer.Token).Value.(int)))), "int"
				} else if expression.Right.(lexer.Token).Type == lexer.Identifier {
					if variables[expression.Right.(lexer.Token).Value.(string)].Type == "int" {
						return block.NewMul(block.NewCall(Functions[name].Value, arguments...), block.NewLoad(types.I32, variables[expression.Right.(lexer.Token).Value.(string)].Value)), "int"
					} else if variables[expression.Right.(lexer.Token).Value.(string)].Type == "float" {
						return block.NewFMul(block.NewCall(Functions[name].Value, arguments...), block.NewLoad(types.Float, variables[expression.Right.(lexer.Token).Value.(string)].Value)), "float"
					}
				} else if expression.Right.(lexer.Token).Type == lexer.Float {
					return block.NewFMul(block.NewCall(Functions[name].Value, arguments...), constant.NewFloat(types.Float, expression.Right.(lexer.Token).Value.(float64))), "float"
				}
			case ast.FuncCall:
				name2 := expression.Right.(ast.FuncCall).Name.Value.(string)
				arguments2 := FuncCall(block, expression.Right.(ast.FuncCall), variables)
				if Functions[name].Type == "int" && Functions[name2].Type == "int" {
					return block.NewMul(block.NewCall(Functions[name].Value, arguments...), block.NewCall(Functions[name2].Value, arguments2...)), "int"
				} else if Functions[name].Type == "float" && Functions[name2].Type == "float" {
					return block.NewFMul(block.NewCall(Functions[name].Value, arguments...), block.NewCall(Functions[name2].Value, arguments2...)), "float"
				}
			}
		}
	} else if expression.Operator.Type == lexer.Divide {
		switch expression.Left.(type) {
		case ast.ExpressionStatement:
			switch expression.Right.(type) {
			case ast.ExpressionStatement:
				exp, exp_type := CompileExpression(block, expression.Left.(ast.ExpressionStatement), variables)
				exp2, exp2_type := CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables)
				if exp_type == "int" && exp2_type == "int" {
					return block.NewSDiv(exp, exp2), "int"
				} else if exp_type == "float" && exp2_type == "float" {
					return block.NewFDiv(exp, exp2), "float"
				}
			case lexer.Token:
				if expression.Right.(lexer.Token).Type == lexer.Int {
					exp, _ := CompileExpression(block, expression.Left.(ast.ExpressionStatement), variables)
					return block.NewSDiv(exp, constant.NewInt(types.I32, int64(expression.Right.(lexer.Token).Value.(int)))), "int"
				} else if expression.Right.(lexer.Token).Type == lexer.Identifier {
					if variables[expression.Right.(lexer.Token).Value.(string)].Type == "int" {
						exp, _ := CompileExpression(block, expression.Left.(ast.ExpressionStatement), variables)
						return block.NewSDiv(exp, block.NewLoad(types.I32, variables[expression.Right.(lexer.Token).Value.(string)].Value)), "int"
					} else if variables[expression.Right.(lexer.Token).Value.(string)].Type == "float" {
						exp, _ := CompileExpression(block, expression.Left.(ast.ExpressionStatement), variables)
						return block.NewFDiv(exp, block.NewLoad(types.Float, variables[expression.Right.(lexer.Token).Value.(string)].Value)), "float"
					}

				} else if expression.Right.(lexer.Token).Type == lexer.Float {
					exp, _ := CompileExpression(block, expression.Left.(ast.ExpressionStatement), variables)
					return block.NewFDiv(exp, constant.NewFloat(types.Float, expression.Right.(lexer.Token).Value.(float64))), "float"
				}
			case ast.FuncCall:
				name := expression.Right.(ast.FuncCall).Name.Value.(string)
				arguments := FuncCall(block, expression.Right.(ast.FuncCall), variables)
				if Functions[name].Type == "int" {
					exp, _ := CompileExpression(block, expression.Left.(ast.ExpressionStatement), variables)
					return block.NewSDiv(exp, block.NewCall(Functions[name].Value, arguments...)), "int"
				} else if Functions[name].Type == "float" {
					exp, _ := CompileExpression(block, expression.Left.(ast.ExpressionStatement), variables)
					return block.NewFDiv(exp, block.NewCall(Functions[name].Value, arguments...)), "float"
				}
			}
		case lexer.Token:
			switch expression.Right.(type) {
			case ast.ExpressionStatement:
				if expression.Left.(lexer.Token).Type == lexer.Int {
					exp, _ := CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables)
					return block.NewSDiv(constant.NewInt(types.I32, int64(expression.Left.(lexer.Token).Value.(int))), exp), "int"
				} else if expression.Left.(lexer.Token).Type == lexer.Identifier {
					if variables[expression.Left.(lexer.Token).Value.(string)].Type == "int" {
						variable := block.NewLoad(types.I32, variables[expression.Left.(lexer.Token).Value.(string)].Value)
						exp, _ := CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables)
						return block.NewSDiv(variable, exp), "int"
					} else if variables[expression.Left.(lexer.Token).Value.(string)].Type == "float" {
						variable := block.NewLoad(types.Float, variables[expression.Left.(lexer.Token).Value.(string)].Value)
						exp, _ := CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables)
						return block.NewFDiv(variable, exp), "float"
					}
				} else if expression.Left.(lexer.Token).Type == lexer.Float {
					exp, _ := CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables)
					return block.NewFDiv(constant.NewFloat(types.Float, expression.Left.(lexer.Token).Value.(float64)), exp), "float"
				}
			case lexer.Token:
				if expression.Left.(lexer.Token).Type == lexer.Int && expression.Right.(lexer.Token).Type == lexer.Int {
					return block.NewSDiv(constant.NewInt(types.I32, int64(expression.Left.(lexer.Token).Value.(int))), constant.NewInt(types.I32, int64(expression.Right.(lexer.Token).Value.(int)))), "int"
				} else if expression.Left.(lexer.Token).Type == lexer.Identifier && expression.Right.(lexer.Token).Type == lexer.Int {
					variable := block.NewLoad(types.I32, variables[expression.Left.(lexer.Token).Value.(string)].Value)
					return block.NewSDiv(variable, constant.NewInt(types.I32, int64(expression.Right.(lexer.Token).Value.(int)))), "int"
				} else if expression.Left.(lexer.Token).Type == lexer.Int && expression.Right.(lexer.Token).Type == lexer.Identifier {
					variable := block.NewLoad(types.I32, variables[expression.Right.(lexer.Token).Value.(string)].Value)
					return block.NewSDiv(constant.NewInt(types.I32, int64(expression.Left.(lexer.Token).Value.(int))), variable), "int"
				} else if expression.Left.(lexer.Token).Type == lexer.Identifier && expression.Right.(lexer.Token).Type == lexer.Identifier {
					if (variables[expression.Left.(lexer.Token).Value.(string)].Type == "int") && (variables[expression.Right.(lexer.Token).Value.(string)].Type == "int") {
						variable := block.NewLoad(types.I32, variables[expression.Left.(lexer.Token).Value.(string)].Value)
						variable2 := block.NewLoad(types.I32, variables[expression.Right.(lexer.Token).Value.(string)].Value)
						return block.NewSDiv(variable, variable2), "int"
					} else if (variables[expression.Left.(lexer.Token).Value.(string)].Type == "float") && (variables[expression.Right.(lexer.Token).Value.(string)].Type == "float") {
						variable := block.NewLoad(types.Float, variables[expression.Left.(lexer.Token).Value.(string)].Value)
						variable2 := block.NewLoad(types.Float, variables[expression.Right.(lexer.Token).Value.(string)].Value)
						return block.NewFDiv(variable, variable2), "float"
					}
				} else if expression.Left.(lexer.Token).Type == lexer.Float && expression.Right.(lexer.Token).Type == lexer.Float {
					return block.NewFDiv(constant.NewFloat(types.Float, expression.Left.(lexer.Token).Value.(float64)), constant.NewFloat(types.Float, expression.Right.(lexer.Token).Value.(float64))), "float"
				} else if expression.Left.(lexer.Token).Type == lexer.Float && expression.Right.(lexer.Token).Type == lexer.Identifier {
					variable := block.NewLoad(types.Float, variables[expression.Right.(lexer.Token).Value.(string)].Value)
					return block.NewFDiv(constant.NewFloat(types.Float, expression.Left.(lexer.Token).Value.(float64)), variable), "float"
				} else if expression.Left.(lexer.Token).Type == lexer.Identifier && expression.Right.(lexer.Token).Type == lexer.Float {
					variable := block.NewLoad(types.Float, variables[expression.Left.(lexer.Token).Value.(string)].Value)
					return block.NewFDiv(variable, constant.NewFloat(types.Float, expression.Right.(lexer.Token).Value.(float64))), "float"
				}
			case ast.FuncCall:
				name2 := expression.Right.(ast.FuncCall).Name.Value.(string)

				arguments2 := FuncCall(block, expression.Right.(ast.FuncCall), variables)
				if Functions[name2].Type == "int" && expression.Left.(lexer.Token).Type == lexer.Int {
					return block.NewSDiv(constant.NewInt(types.I32, int64(expression.Left.(lexer.Token).Value.(int))), block.NewCall(Functions[name2].Value, arguments2...)), "int"
				} else if Functions[name2].Type == "int" && expression.Left.(lexer.Token).Type == lexer.Identifier {
					variable := block.NewLoad(types.I32, variables[expression.Left.(lexer.Token).Value.(string)].Value)
					return block.NewSDiv(variable, block.NewCall(Functions[name2].Value, arguments2...)), "int"
				} else if Functions[name2].Type == "float" && expression.Left.(lexer.Token).Type == lexer.Float {
					return block.NewFDiv(constant.NewFloat(types.Float, expression.Left.(lexer.Token).Value.(float64)), block.NewCall(Functions[name2].Value, arguments2...)), "float"
				} else if Functions[name2].Type == "float" && expression.Left.(lexer.Token).Type == lexer.Identifier {
					variable := block.NewLoad(types.Float, variables[expression.Left.(lexer.Token).Value.(string)].Value)
					return block.NewFDiv(variable, block.NewCall(Functions[name2].Value, arguments2...)), "float"
				}

			}
		case ast.FuncCall:
			name := expression.Left.(ast.FuncCall).Name.Value.(string)
			arguments := FuncCall(block, expression.Left.(ast.FuncCall), variables)
			switch expression.Right.(type) {
			case ast.ExpressionStatement:
				if Functions[name].Type == "int" {
					exp, _ := CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables)
					return block.NewSDiv(block.NewCall(Functions[name].Value, arguments...), exp), "int"
				} else if Functions[name].Type == "float" {
					exp, _ := CompileExpression(block, expression.Right.(ast.ExpressionStatement), variables)
					return block.NewFDiv(block.NewCall(Functions[name].Value, arguments...), exp), "float"
				}
			case lexer.Token:
				if expression.Right.(lexer.Token).Type == lexer.Int {
					return block.NewSDiv(block.NewCall(Functions[name].Value, arguments...), constant.NewInt(types.I32, int64(expression.Right.(lexer.Token).Value.(int)))), "int"
				} else if expression.Right.(lexer.Token).Type == lexer.Identifier {
					if variables[expression.Right.(lexer.Token).Value.(string)].Type == "int" {
						return block.NewSDiv(block.NewCall(Functions[name].Value, arguments...), block.NewLoad(types.I32, variables[expression.Right.(lexer.Token).Value.(string)].Value)), "int"
					} else if variables[expression.Right.(lexer.Token).Value.(string)].Type == "float" {
						return block.NewFDiv(block.NewCall(Functions[name].Value, arguments...), block.NewLoad(types.Float, variables[expression.Right.(lexer.Token).Value.(string)].Value)), "float"
					}
				} else if expression.Right.(lexer.Token).Type == lexer.Float {
					return block.NewFDiv(block.NewCall(Functions[name].Value, arguments...), constant.NewFloat(types.Float, expression.Right.(lexer.Token).Value.(float64))), "float"
				}
			case ast.FuncCall:
				name2 := expression.Right.(ast.FuncCall).Name.Value.(string)
				arguments2 := FuncCall(block, expression.Right.(ast.FuncCall), variables)
				if Functions[name].Type == "int" && Functions[name2].Type == "int" {
					return block.NewSDiv(block.NewCall(Functions[name].Value, arguments...), block.NewCall(Functions[name2].Value, arguments2...)), "int"
				} else if Functions[name].Type == "float" && Functions[name2].Type == "float" {
					return block.NewFDiv(block.NewCall(Functions[name].Value, arguments...), block.NewCall(Functions[name2].Value, arguments2...)), "float"
				}
			}
		}
	}

	return block, ""
}
