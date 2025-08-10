package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/expr-lang/expr"
)

// Globals
var isBinary bool
var pat2 string = `^([A-Za-z]+\_*?)+\[[0-9]+\];?$`
var re2 *regexp.Regexp = regexp.MustCompile(pat2)
var InterpreterVariables = make(map[string]interface{})

func Interpreter() {
	// Grab AST
	bytes, _ := os.ReadFile("./.intext/cache/AST.json")
	nodes := []map[string]interface{}{}
	err := json.Unmarshal(bytes, &nodes)
	Check(err)

	// Iterate through each node
	for _, node := range nodes {
		switch node["type"] {
		case "let", "declare":
			name := node["var_name"].(string)
			Type := node["var_type"]
			val := node["var_value"]
			meta := node["meta"].(map[string]interface{})

			if re2.MatchString(fmt.Sprint(val)) {
				v, _ := expr.Eval(fmt.Sprintln(val), InterpreterVariables)
				InterpreterVariables[name] = v
			} else {
				switch Type {
				case "int":
					switch meta["math"] {
					case true:
						cmp, err := expr.Compile(fmt.Sprint(val))
						Check(err)
						ans, errd := expr.Run(cmp, InterpreterVariables)
						Check(errd)
						InterpreterVariables[name] = ans
					case false:
						InterpreterVariables[name] = val
					}
				case "string", "char":
					switch meta["raw_type"] {
					case "STRING", "CHAR":
						val := val.(string)
						InterpreterVariables[name] = val

					case "FUNC":
						value := val.(map[string]interface{})
						file := value["read"]

						data, err := os.ReadFile(fmt.Sprint(file))
						Check(err)
						isBinary = binaryCheck(data)
						InterpreterVariables[name] = string(data)
					case "concat":
						var catch string
						for _, element := range val.([]interface{}) {
							if re2.MatchString(fmt.Sprint(element)) {
								v, _ := expr.Eval(fmt.Sprint(element), InterpreterVariables)
								catch += fmt.Sprint(v)
							} else {
								catch += fmt.Sprint(element)
							}
						}
						InterpreterVariables[name] = catch
					}
				case "bool":
					InterpreterVariables[name] = val
				case "float":
					InterpreterVariables[name] = val
				case "order", "ord":
					valList := val.([]interface{})

					for i, element := range valList {
						if _, ok := InterpreterVariables[fmt.Sprint(element)]; ok {
							value := InterpreterVariables[fmt.Sprint(element)]

							valList[i] = value
						}
					}

					InterpreterVariables[name] = valList
				}
			}
		case "output":
			value := node["value"]
			meta := node["meta"].(map[string]interface{})

			switch meta["print_type"] {
			case "simple":
				val, ok := InterpreterVariables[value.(string)]
				if ok {
					if isBinary {
						fmt.Println("\033[31m--- WARNING: The following variable you are about to output is linked to a file variable that is suspected to be a binary file ---\033[0m \n Ctrl + C before outputing... Or press Enter to continue regardless")

						input := make(chan byte, 1)
						go func() {
							b := make([]byte, 1)
							os.Stdin.Read(b)
							input <- b[0]
						}()

						b := <-input
						if b == '\n' || b == '\r' {
							fmt.Println(val)
						}
						fmt.Println(val)
					}
					fmt.Println(NullCheck(val, false))
				} else {
					if re2.MatchString(fmt.Sprint(value)) {
						val, _ := expr.Eval(fmt.Sprint(value), InterpreterVariables)
						fmt.Println(val)
					} else {
						fmt.Println(value)
					}
				}
			case "mixed":
				SpagList := []interface{}{}
				var catch string

				for _, SpagVal := range value.([]interface{}) {
					switch SpagVal.(type) {
					case string:
						vari, ok := InterpreterVariables[SpagVal.(string)]
						if ok {
							SpagList = append(SpagList, vari)
						} else {
							if re2.MatchString(SpagVal.(string)) {
								val, _ := expr.Eval(SpagVal.(string), InterpreterVariables)
								SpagList = append(SpagList, NullCheck(val, false))
							} else {
								SpagList = append(SpagList, vari)
							}
						}
					default:
						SpagList = append(SpagList, SpagVal)
					}
				}

				for _, i := range SpagList {
					catch += fmt.Sprint(i)
				}

				fmt.Println(catch)
			case "mathematics":
				val := fmt.Sprint(value)
				cmp, errd := expr.Compile(val)
				Check(errd)
				ans, err := expr.Run(cmp, InterpreterVariables)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(ans)
			case "ord_index":
				OrdRef := value.(map[string]interface{})

				for key, val := range OrdRef {
					ListRaw := InterpreterVariables[key]
					List := ListRaw.([]interface{})
					i, _ := strconv.Atoi(fmt.Sprint(val))

					fmt.Println(List[i])
				}
			}
		case "function":
			meta := node["meta"].(map[string]interface{})
			switch node["call"] {
			case "write":
				input := meta["input"]
				val := InterpreterVariables[fmt.Sprint(input)]
				target := meta["target"]
				perms := meta["perms"]
				octal, _ := strconv.ParseInt(fmt.Sprint(perms), 8, 64)
				err := os.WriteFile(fmt.Sprint(target), []byte(val.(string)), 0000)
				Check(err)
				erra := os.Chmod(fmt.Sprint(target), os.FileMode(octal))
				Check(erra)
			case "append":
				fileTaget := meta["target"]
				val := fmt.Sprint(InterpreterVariables[fmt.Sprint(meta["input"])])
				f, err := os.OpenFile(fmt.Sprint(fileTaget), os.O_APPEND|os.O_WRONLY, 0666)
				if err != nil {
					i, _ := strconv.Atoi(fmt.Sprint(meta["line"]))
					erra := NewError("FileError", i, fmt.Sprintf("append(%s%v%s, %v);", Red, fmt.Sprint(fileTaget), Reset, fmt.Sprint(meta["input"])), "The following file does not exist", true, "append() requires the file being appended to exist. If you wanted to create one, use write().")
					erra.Throw()
				}
				defer f.Close()
				_, errd := f.WriteString(val)
				Check(errd)
			case "del":
				fileTarget := fmt.Sprint(meta["target"])
				var val string
				if _, ok := InterpreterVariables[fileTarget]; ok {
					val = fmt.Sprint(InterpreterVariables[fileTarget])
					os.Remove(val)
				} else {
					os.Remove(fileTarget)
				}
			}

		case "logic":
			meta := node["meta"].(map[string]interface{})

			switch meta["sub_type"] {
			case "if":
				body := node["body"].([]interface{})

				cond := fmt.Sprint(node["condition"])
				val, _ := expr.Eval(cond, InterpreterVariables)

				captureAST := []map[string]interface{}{}

				for _, element := range body {
					captureAST = append(captureAST, element.(map[string]interface{}))
				}

				if v, _ := strconv.ParseBool(fmt.Sprint(val)); v {
					reRunInterpreter(captureAST)
				} else {
					continue
				}
			}
		}
	}
}

func reRunInterpreter(nodes []map[string]interface{}) {
	for _, node := range nodes {
		switch node["type"] {
		case "let", "declare":
			name := node["var_name"].(string)
			Type := node["var_type"]
			val := node["var_value"]
			meta := node["meta"].(map[string]interface{})

			if re2.MatchString(fmt.Sprint(val)) {
				v, _ := expr.Eval(fmt.Sprintln(val), InterpreterVariables)
				InterpreterVariables[name] = v
			} else {
				switch Type {
				case "int":
					switch meta["math"] {
					case true:
						cmp, err := expr.Compile(fmt.Sprint(val))
						Check(err)
						ans, errd := expr.Run(cmp, InterpreterVariables)
						Check(errd)
						InterpreterVariables[name] = ans
					case false:
						InterpreterVariables[name] = val
					}
				case "string", "char":
					switch meta["raw_type"] {
					case "STRING", "CHAR":
						val := val.(string)
						InterpreterVariables[name] = val

					case "FUNC":
						value := val.(map[string]interface{})
						file := value["read"]

						data, err := os.ReadFile(fmt.Sprint(file))
						Check(err)
						isBinary = binaryCheck(data)
						InterpreterVariables[name] = string(data)
					case "concat":
						var catch string
						for _, element := range val.([]interface{}) {
							if re2.MatchString(fmt.Sprint(element)) {
								v, _ := expr.Eval(fmt.Sprint(element), InterpreterVariables)
								catch += fmt.Sprint(v)
							} else {
								catch += fmt.Sprint(element)
							}
						}
						InterpreterVariables[name] = catch
					}
				case "bool":
					InterpreterVariables[name] = val
				case "float":
					InterpreterVariables[name] = val
				case "order", "ord":
					valList := val.([]interface{})

					for i, element := range valList {
						if _, ok := InterpreterVariables[fmt.Sprint(element)]; ok {
							value := InterpreterVariables[fmt.Sprint(element)]

							valList[i] = value
						}
					}

					InterpreterVariables[name] = valList
				}
			}
		case "output":
			value := node["value"]
			meta := node["meta"].(map[string]interface{})

			switch meta["print_type"] {
			case "simple":
				val, ok := InterpreterVariables[value.(string)]
				if ok {
					if isBinary {
						fmt.Println("\033[31m--- WARNING: The following variable you are about to output is linked to a file variable that is suspected to be a binary file ---\033[0m \n Ctrl + C before outputing... Or press Enter to continue regardless")

						input := make(chan byte, 1)
						go func() {
							b := make([]byte, 1)
							os.Stdin.Read(b)
							input <- b[0]
						}()

						b := <-input
						if b == '\n' || b == '\r' {
							fmt.Println(val)
						}
						fmt.Println(val)
					}
					fmt.Println(NullCheck(val, false))
				} else {
					if re2.MatchString(fmt.Sprint(value)) {
						val, _ := expr.Eval(fmt.Sprint(value), InterpreterVariables)
						fmt.Println(val)
					} else {
						fmt.Println(value)
					}
				}
			case "mixed":
				SpagList := []interface{}{}
				var catch string

				for _, SpagVal := range value.([]interface{}) {
					switch SpagVal.(type) {
					case string:
						vari, ok := InterpreterVariables[SpagVal.(string)]
						if ok {
							SpagList = append(SpagList, vari)
						} else {
							if re2.MatchString(SpagVal.(string)) {
								val, _ := expr.Eval(SpagVal.(string), InterpreterVariables)
								SpagList = append(SpagList, NullCheck(val, false))
							} else {
								SpagList = append(SpagList, vari)
							}
						}
					default:
						SpagList = append(SpagList, SpagVal)
					}
				}

				for _, i := range SpagList {
					catch += fmt.Sprint(i)
				}

				fmt.Println(catch)
			case "mathematics":
				val := fmt.Sprint(value)
				cmp, errd := expr.Compile(val)
				Check(errd)
				ans, err := expr.Run(cmp, InterpreterVariables)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(ans)
			case "ord_index":
				OrdRef := value.(map[string]interface{})

				for key, val := range OrdRef {
					ListRaw := InterpreterVariables[key]
					List := ListRaw.([]interface{})
					i, _ := strconv.Atoi(fmt.Sprint(val))

					fmt.Println(List[i])
				}
			}
		case "function":
			meta := node["meta"].(map[string]interface{})
			switch node["call"] {
			case "write":
				input := meta["input"]
				val := InterpreterVariables[fmt.Sprint(input)]
				target := meta["target"]
				perms := meta["perms"]
				octal, _ := strconv.ParseInt(fmt.Sprint(perms), 8, 64)
				err := os.WriteFile(fmt.Sprint(target), []byte(val.(string)), 0000)
				Check(err)
				erra := os.Chmod(fmt.Sprint(target), os.FileMode(octal))
				Check(erra)
			case "append":
				fileTaget := meta["target"]
				val := fmt.Sprint(InterpreterVariables[fmt.Sprint(meta["input"])])
				f, err := os.OpenFile(fmt.Sprint(fileTaget), os.O_APPEND|os.O_WRONLY, 0666)
				if err != nil {
					i, _ := strconv.Atoi(fmt.Sprint(meta["line"]))
					erra := NewError("FileError", i, fmt.Sprintf("append(%s%v%s, %v);", Red, fmt.Sprint(fileTaget), Reset, fmt.Sprint(meta["input"])), "The following file does not exist", true, "append() requires the file being appended to exist. If you wanted to create one, use write().")
					erra.Throw()
				}
				defer f.Close()
				_, errd := f.WriteString(val)
				Check(errd)
			case "del":
				fileTarget := fmt.Sprint(meta["target"])
				var val string
				if _, ok := InterpreterVariables[fileTarget]; ok {
					val = fmt.Sprint(InterpreterVariables[fileTarget])
					os.Remove(val)
				} else {
					os.Remove(fileTarget)
				}
			}

		case "logic":
			meta := node["meta"].(map[string]interface{})

			switch meta["sub_type"] {
			case "if":

			}
		}
	}
}
