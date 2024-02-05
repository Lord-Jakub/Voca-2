package compiler

import (
	"os"
	"path/filepath"
	"strings"
)

func parseCfg(libName string, libPath string, arch string, Os string) string {
	var tokens []string
	var addflags string
	i := 0
	cfgPath := filepath.Join(libPath, libName+".cfg")
	data, err := os.ReadFile(cfgPath)

	if err != nil {
		return ""
	}
	input := string(data)
	input = strings.Replace(input, "\r", " ", -1)
	input = strings.Replace(input, "   ", "\t", -1)

	for i < len(input) {
		switch input[i] {
		case ' ':
			i++
		case '\n':
			i++
		case '\t':
			tokens = append(tokens, " ")
			i++
		case '"':
			s := ""
			i++
			for i < len(input) && input[i] != '"' {
				s = s + string(input[i])
				i++
			}
			tokens = append(tokens, s)
			i++
		default:
			s := ""
			for i < len(input) && input[i] != ' ' && input[i] != '\n' && input[i] != '\t' {
				s = s + string(input[i])
				i++
			}
			tokens = append(tokens, s)

		}
	}
	i = 0
	for i < len(tokens) {
		if tokens[i] != " " {
			if tokens[i] == arch+":" {
				i++
				if tokens[i] == " " && tokens[i+1] == Os+":" {
					i += 2
					for i+1 < len(tokens) && tokens[i] == " " && tokens[i+1] == " " {
						i += 2
						switch tokens[i] {
						case "ADDFLAGS:":
							i++
							addflags = tokens[i]
							i++
						}

					}
				}
			}
		}
		i++
	}
	return strings.Replace(addflags, "{libDir}", libPath, -1)
}
