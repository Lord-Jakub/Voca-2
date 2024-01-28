package lexer

import (
	"Voca-2/lib"
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

type Token struct {
	Type    TokenType
	Value   any
	Line    int
	LinePos int
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
var KeyWords = []string{"if", "else", "while", "func", "return", "import", "string", "int", "float", "bool", "void", "extern_func", "true", "false"}

func Lex(input string) ([]Token, error) {

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
					LinePos: linePos})
			} else {
				tokens = append(tokens, Token{
					Type:    Identifier,
					Value:   s,
					Line:    line,
					LinePos: linePos})
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
					Type:    Int,
					Value:   number,
					Line:    line,
					LinePos: linePos,
				})
			} else {
				number, _ := strconv.Atoi(num)
				tokens = append(tokens, Token{
					Type:    Int,
					Value:   number,
					Line:    line,
					LinePos: linePos,
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
			})
		case string(input[pos]) == "\n":
			//NewLine
			tokens = append(tokens, Token{
				Type:  NewLine,
				Value: "\n",
				Line:  line,
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
				})
				pos++
				linePos++

			} else {
				tokens = append(tokens, Token{
					Type:    Equal,
					Value:   "=",
					Line:    line,
					LinePos: linePos,
				})
			}
		case string(input[pos]) == "!":
			if string(input[pos+1]) == "=" {
				tokens = append(tokens, Token{
					Type:    NotEqual,
					Value:   "!=",
					Line:    line,
					LinePos: linePos,
				})
				pos++
				linePos++

			} else {
				tokens = append(tokens, Token{
					Type:    Not,
					Value:   "!",
					Line:    line,
					LinePos: linePos,
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
				})
			} else {
				err = errors.New(fmt.Sprintf("Invalid symbol: '%c' on line: %d at position: %d", input[pos], line, linePos))
			}
		}
		linePos++
		pos++
	}

	return tokens, err
}
func IsOperator(token Token) bool {
	if token.Type == Plus || token.Type == Minus || token.Type == Multiply || token.Type == Divide {
		return true
	}
	return false
}
