package lexer

import (
	"Voca-2/lib"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

type Token struct {
	Type    TokenType
	Value   any
	Line    int
	LinePos int
	File    string
}
type TokenType int

const (
	Invalid     = iota
	OpenParen   //1
	CloseParen  //2
	OpenBrace   //3
	CloseBrace  //4
	Plus        //5
	Minus       //6
	Multiply    //7
	Divide      //8
	Backslash   //9
	NewLine     //10
	SingleQuote //11
	DoubleQuote //12
	Equal       //13
	Not         //14
	Comma       //15
	LessThan    //16
	MoreThan    //17
	Int         //18
	String      //19
	Identifier  //20
	Keyword     //21
	DoubleEqual //23
	NotEqual    //24
	Float       //25
)

var symbolMap = map[byte]TokenType{
	'(':  OpenParen,
	')':  CloseParen,
	'{':  OpenBrace,
	'}':  CloseBrace,
	'+':  Plus,
	'-':  Minus,
	'*':  Multiply,
	'/':  Divide,
	'\\': Backslash,
	'=':  Equal,
	'!':  Not,
	',':  Comma,
	'<':  LessThan,
	'>':  MoreThan,
	';':  NewLine,
}
var KeyWords = []string{"if", "else", "while", "func", "return", "string", "int", "float", "bool", "void", "extern_func", "true", "false", "import", "external"}

func AddImports(tokens []Token) []Token {
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Type == Keyword && tokens[i].Value == "import" {
			var prevTokens []Token
			if i != 0 {
				prevTokens = tokens[:i-1]
			} else {
				prevTokens = []Token{}
			}

			var s string
			i++
			file := tokens[i].Value.(string)
			data, err := os.ReadFile(file + ".voc")
			if err != nil {
				exePath, err := os.Executable()
				if err != nil {
					fmt.Print("Can't load file: " + err.Error())
				}
				lib := filepath.Join(filepath.Dir(exePath), "libs", file, file+".voc")
				data, err = os.ReadFile(lib)
				if err != nil {
					fmt.Print("Can't load file: " + err.Error())
				}

			}
			s = string(data)
			s = strings.Replace(s, "\r", " ", -1)
			newTokens, err := Lex(s, file+".voc")
			newTokens = EditImports(newTokens, file)
			if err != nil {
				fmt.Print("Lexing error: " + err.Error())
			}
			pastTokens := tokens[i+1:]
			tokens = append(prevTokens, newTokens...)
			tokens = append(tokens, pastTokens...)

		}
	}
	return tokens
}

func EditImports(tokens []Token, libname string) []Token {
	functions := []string{}
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Type == Keyword && tokens[i].Value == "func" {
			i++
			functions = append(functions, tokens[i].Value.(string))
			tokens[i].Value = libname + "." + tokens[i].Value.(string)
			i++
		} else if tokens[i].Type == Identifier {
			for j := 0; j < len(functions); j++ {
				if tokens[i].Value == functions[j] {
					tokens[i].Value = libname + "." + tokens[i].Value.(string)
				}
			}
		}

	}
	return tokens
}

func Lex(input string, file string) ([]Token, error) {

	var tokens []Token
	var err error = nil
	pos := 0
	line := 1
	linePos := 0
	for pos < len(input) {
		switch {
		case unicode.IsLetter(rune(input[pos])) && !(string(input[pos]) == " ") && !(string(input[pos]) == "\"") && !(string(input[pos]) == "'"):
			//Strings
			var s string
			for pos < len(input) && (unicode.IsLetter(rune(input[pos])) || unicode.IsDigit(rune(input[pos])) || (string(input[pos]) == ".") || (string(input[pos]) == "_")) && !(string(input[pos]) == " ") && !(string(input[pos]) == "\"") && !(string(input[pos]) == "'") {
				s += string(input[pos])
				pos++
				linePos++

			}
			pos--
			linePos--
			if lib.Contains(s, KeyWords) {
				tokens = append(tokens, Token{
					Type:    Keyword,
					Value:   s,
					Line:    line,
					LinePos: linePos,
					File:    file})
			} else {
				tokens = append(tokens, Token{
					Type:    Identifier,
					Value:   s,
					Line:    line,
					LinePos: linePos,
					File:    file})
			}
		case unicode.IsDigit(rune(input[pos])):
			//Numbers
			var num string
			isFloat := false
			for pos < len(input) && (unicode.IsDigit(rune(input[pos])) || string(input[pos]) == ".") {
				num += string(input[pos])
				if string(input[pos]) == "." {
					isFloat = true
				}
				pos++
				linePos++

			}
			pos--
			linePos--
			if isFloat {
				number, _ := strconv.ParseFloat(num, 64)
				tokens = append(tokens, Token{
					Type:    Float,
					Value:   number,
					Line:    line,
					LinePos: linePos,
					File:    file,
				})
			} else {
				number, _ := strconv.Atoi(num)
				tokens = append(tokens, Token{
					Type:    Int,
					Value:   number,
					Line:    line,
					LinePos: linePos,
					File:    file,
				})
			}

		case string(input[pos]) == "\"":
			//Strings
			var s string
			pos++
			linePos++

			for pos < len(input) && !(string(input[pos]) == "\"") {
				s += string(input[pos])
				pos++
				linePos++

			}

			tokens = append(tokens, Token{
				Type:    String,
				Value:   s,
				Line:    line,
				LinePos: linePos,
				File:    file,
			})
		case string(input[pos]) == "'":
			//Strings
			var s string
			pos++
			linePos++

			for pos < len(input) && !(string(input[pos]) == "'") {
				s += string(input[pos])
				pos++
				linePos++

			}

			tokens = append(tokens, Token{
				Type:    String,
				Value:   s,
				Line:    line,
				LinePos: linePos,
				File:    file,
			})
		case string(input[pos]) == "\n":
			//NewLine
			tokens = append(tokens, Token{
				Type:  NewLine,
				Value: "\n",
				Line:  line,
				File:  file,
			})
			line++
			linePos = 0
		case string(input[pos]) == "=":
			if string(input[pos+1]) == "=" {
				tokens = append(tokens, Token{
					Type:    DoubleEqual,
					Value:   "==",
					Line:    line,
					LinePos: linePos,
					File:    file,
				})
				pos++
				linePos++

			} else {
				tokens = append(tokens, Token{
					Type:    Equal,
					Value:   "=",
					Line:    line,
					LinePos: linePos,
					File:    file,
				})
			}
		case string(input[pos]) == "!":
			if string(input[pos+1]) == "=" {
				tokens = append(tokens, Token{
					Type:    NotEqual,
					Value:   "!=",
					Line:    line,
					LinePos: linePos,
					File:    file,
				})
				pos++
				linePos++

			} else {
				tokens = append(tokens, Token{
					Type:    Not,
					Value:   "!",
					Line:    line,
					LinePos: linePos,
					File:    file,
				})
			}
		case string(input[pos]) == " ":
			//Skip
		case string(input[pos]) == "\t":
			//Skip
		case string(input[pos]) == "/" && string(input[pos+1]) == "/":
			for pos < len(input) && !(string(input[pos]) == "\n") {
				pos++
				linePos++
			}
		case string(input[pos]) == "/" && string(input[pos+1]) == "*":
			for pos < len(input) && !(string(input[pos]) == "*" && string(input[pos+1]) == "/") {
				if string(input[pos]) == "\n" {
					line++
					linePos = 0
				}
				pos++
			}

		default:
			//Symbols
			if _, ok := symbolMap[input[pos]]; ok {
				tokens = append(tokens, Token{
					Type:    symbolMap[input[pos]],
					Value:   string(input[pos]),
					Line:    line,
					LinePos: linePos,
					File:    file,
				})
			} else if string(input[pos]) == "\n" {
				line++
				linePos = 0
				tokens = append(tokens, Token{
					Type:    NewLine,
					Value:   "\n",
					Line:    line,
					LinePos: linePos,
					File:    file,
				})
			} else {

				err = errors.New(fmt.Sprintf("Invalid symbol: '%c' on line: %d at position: %d in file: '%s'", input[pos], line, linePos, file))
			}
		}
		linePos++
		pos++
	}
	tokens = AddImports(tokens)
	//fmt.Println(tokens)
	return tokens, err
}
func IsOperator(token Token) bool {
	if token.Type == Plus || token.Type == Minus || token.Type == Multiply || token.Type == Divide {
		return true
	}
	return false
}
