package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func Contains(slice []interface{}, look interface{}) bool {
	for _, v := range slice {
		if look == v {
			return true
		}
	}
	return false
}

var pat2 string = `([A-Za-z]+\_*?)+\[[0-9]+\];?`
var re2 *regexp.Regexp = regexp.MustCompile(pat2)

func Lexer(filename string) {

	file, err := os.Open(filename)
	Check(err)
	defer file.Close()

	var lines []string

	scan := bufio.NewScanner(file)
	for scan.Scan() {
		lines = append(lines, scan.Text())
	}
	fileErr := scan.Err()
	Check(fileErr)

	// Banks
	allTokens := []map[string]interface{}{}

	for index := 0; index < len(lines); index++ {
		line := strings.TrimSpace(lines[index])

		if strings.Contains(line, "//") { // Checks if line is a comment
			line = strings.Split(line, "//")[0]
		}

		if line == "" || line == " " { // Checks if line is space
			continue
		}

	outer: // Label for loop
		for i := 0; i < len(line); i++ { // NOTE: Some are place early so the others behind don't get triggered beforehand
			char := rune(line[i])

			// Checking if the current char is a space
			if unicode.IsSpace(char) || char == ' ' || string(char) == " " {
				continue // Adding these so the loop continues; it'll get stuck here not knowing what to do
			}

			// Checks for math operators
			if i+1 < len(line) && Contains([]interface{}{">=", "<=", "==", "+=", "*=", "-=", "/="}, string(line[i:i+2])) {
				allTokens = append(allTokens, map[string]interface{}{
					"TYPE": "OPERATOR",
					"VAL":  string(line[i : i+2]),
					"LINE": index + 1,
				})
				continue
			}

			// Multi-line comment support
			if i+1 < len(line) && string(line[i:i+2]) == "/*" {
				var multiComment string

				index += 1
				for index < len(lines) {
					if strings.TrimSpace(lines[index]) == "*/" {
						break
					}
					multiComment += lines[index]
					index += 1
				}

				continue
			}

			// Checks for either a semicolon or equal sign
			if Contains([]interface{}{";", ",", ":"}, string(char)) {
				if string(char) == "," {
					for j := 0; j < len(allTokens)-1; j++ {
						prev := allTokens[j]
						if prev["TYPE"] == "LBRACKET" && prev["LINE"] == index+1 {
							allTokens = append(allTokens, map[string]interface{}{
								"TYPE": "COMMA",
								"VAL":  string(char),
								"LINE": index + 1,
							})
							continue outer
						}
					}
				}
				allTokens = append(allTokens, map[string]interface{}{
					"TYPE": "SYMBOL",
					"VAL":  string(char),
					"LINE": index + 1,
				})
				continue
			}

			// Checks for single char operators
			if Contains([]interface{}{"+", "-", "*", "/", ">", "<", "="}, string(char)) {
				allTokens = append(allTokens, map[string]interface{}{
					"TYPE": "OPERATOR",
					"VAL":  string(char),
					"LINE": index + 1,
				})
				continue
			}

			// Checks for words; very complex; stuff like output, declare, and types go here
			if unicode.IsLetter(char) {
				temp := ""
				var indexCheck bool = false

				if re2.MatchString(line) {
					indexCheck = true
				}

				if !indexCheck {
					for i < len(line) && (unicode.IsLetter(rune(line[i])) || string(line[i]) == "_") {
						temp += string(line[i])
						i++
					}
				} else {
					for i < len(line) && (unicode.IsLetter(rune(line[i])) || Contains([]interface{}{"_", "[", "]"}, string(line[i])) || unicode.IsDigit(rune(line[i]))) {
						temp += string(line[i])
						i++
					}
				}
				if Contains([]interface{}{"output", "declare", "let"}, temp) {
					allTokens = append(allTokens, map[string]interface{}{
						"TYPE": "KEYWORD",
						"VAL":  temp,
						"LINE": index + 1,
					})
				} else if Contains([]interface{}{"if"}, temp) {
					var rawCatch string
					for tempI := 0; tempI < len(line) && string(line[i]) != "{"; {
						tempI = i

						rawCatch += string(line[i])

						i++
					}

					cleanCatch := strings.TrimSpace(rawCatch)
					allTokens = append(allTokens, map[string]interface{}{
						"TYPE":     "LOGIC",
						"VAL":      cleanCatch,
						"SUB-TYPE": "if",
						"LINE":     index + 1,
					})
				} else if Contains([]interface{}{"read", "write", "append", "del", "remove"}, temp) {
					allTokens = append(allTokens, map[string]interface{}{
						"TYPE": "FUNC",
						"VAL":  temp,
						"LINE": index + 1,
					})
				} else if Contains([]interface{}{"bool", "string", "int", "char", "float", "ord", "order"}, temp) {
					allTokens = append(allTokens, map[string]interface{}{
						"TYPE": "TYPESYS",
						"VAL":  temp,
						"LINE": index + 1,
					})
				} else if temp == "true" || temp == "false" {
					allTokens = append(allTokens, map[string]interface{}{
						"TYPE": "BOOL",
						"VAL":  temp,
						"LINE": index + 1,
					})
				} else {
					allTokens = append(allTokens, map[string]interface{}{
						"TYPE": "IDENTIFIER",
						"VAL":  temp,
						"LINE": index + 1,
					})
				}
				i -= 1
				continue
			}

			if i+2 < len(line) && string(line[i:i+3]) == "[[[" {
				//textCapture := []interface{}{}

			}

			// Switch statement for the single employeed bums
			switch string(char) {
			case "[":
				allTokens = append(allTokens, map[string]interface{}{
					"TYPE": "LBRACKET",
					"VAL":  string(char),
					"LINE": index + 1,
				})
				continue
			case "]":
				allTokens = append(allTokens, map[string]interface{}{
					"TYPE": "RBRACKET",
					"VAL":  string(char),
					"LINE": index + 1,
				})

				continue
			case "'":
				if i+1 < len(line) && string(line[i+2]) == "'" {
					i += 1
					chr := string(line[i])
					allTokens = append(allTokens, map[string]interface{}{
						"TYPE": "CHAR",
						"VAL":  chr,
						"LINE": index + 1,
					})
					i += 1
					continue
				} else {
					err := NewError("LexerErr", index+1, line, fmt.Sprintf("Invalid %schar%s", Yellow, Reset), false, fmt.Sprintf("A %schar%s is structured as, 'A'", Yellow, Reset))
					err.Throw()
				}
			case "(", ")":
				allTokens = append(allTokens, map[string]interface{}{
					"TYPE": "PARA",
					"VAL":  string(char),
					"LINE": index + 1,
				})

				continue
			case "{":
				allTokens = append(allTokens, map[string]interface{}{
					"TYPE": "LCURL",
					"VAL":  string(char),
					"LINE": index + 1,
				})

				continue
			case "}":
				allTokens = append(allTokens, map[string]interface{}{
					"TYPE": "RCURL",
					"VAL":  string(char),
					"LINE": index + 1,
				})

				continue
			}

			// Checks for strings, since they start with quotes
			if char == '"' {
				i++
				var stringVal string

				for i < len(line) {
					if rune(line[i]) == '"' {
						break
					}
					stringVal += string(line[i])
					i++
				}

				if i >= len(line) {
					panic("String ranged out")
				}

				allTokens = append(allTokens, map[string]interface{}{
					"TYPE": "STRING",
					"VAL":  stringVal,
					"LINE": index + 1,
				})
				continue
			}

			// Checks if it's strictly a number from 0 to 9, not some arabic or roman numeral
			if unicode.IsDigit(char) {
				var num string
				var floatCatch bool = false

			number:
				for i < len(line) && (unicode.IsDigit(char) || char == '.') {
					if i+1 < len(line) && Contains([]interface{}{")", "]", ","}, string(line[i+1])) {
						num += string(line[i])
						break number

					}

					switch string(line[i]) {
					case ";":
						break number
					case ".":
						floatCatch = true
					}

					num += string(line[i])
					i += 1
				}
				num = strings.TrimSpace(num)
				if floatCatch {
					nflot, err := strconv.ParseFloat(num, 64)
					Check(err)
					allTokens = append(allTokens, map[string]interface{}{
						"TYPE": "FLOAT",
						"VAL":  nflot,
						"LINE": index + 1,
					})
					continue
				} else {
					number, err := strconv.Atoi(num)
					if err != nil {
						allTokens = append(allTokens, map[string]interface{}{
							"TYPE": "MATH",
							"VAL":  num,
							"LINE": index + 1,
						})
						continue
					}
					allTokens = append(allTokens, map[string]interface{}{
						"TYPE": "INT",
						"VAL":  number,
						"LINE": index + 1,
					})
					continue
				}
			}

			err := NewError("LexerErr", index+1, line, "Invalid character somewhere in expression", false, "")
			err.Throw()
		}

		if len(allTokens) > 0 {
			last := allTokens[len(allTokens)-1]
			if !Contains([]interface{}{"SYMBOL", "LCURL", "RCURL"}, last["TYPE"]) || !Contains([]interface{}{";", "{", "}"}, last["VAL"]) {
				allTokens = append(allTokens, map[string]interface{}{
					"TYPE": "SYMBOL",
					"VAL":  ";",
					"LINE": index + 1,
				})
			}
		}

	}
	// Finally add to tokens.json
	b, err := json.MarshalIndent(allTokens, "", "  ")
	Check(err)
	cacheDir := filepath.Join(".intext", "cache")
	os.MkdirAll(cacheDir, 0766)
	os.WriteFile("./.intext/cache/Tokens.json", b, 0666)
}
