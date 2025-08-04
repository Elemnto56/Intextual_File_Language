package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/expr-lang/expr"
)

func Interpreter() {
	// Grab AST
	bytes, _ := os.ReadFile("AST.json")
	nodes := []map[string]interface{}{}
	err := json.Unmarshal(bytes, &nodes)
	Check(err)

	// Banks
	variables := map[string]interface{}{}

	// Globals
	var isBinary bool
	pat2 = `^([A-Za-z]+\_*?)+\[[0-9]+\];?$`
	re2 = regexp.MustCompile(pat2)

	// Iterate through each node
	for _, node := range nodes {
		switch node["type"] {
		case "let", "declare":
			name := node["var_name"].(string)
			Type := node["var_type"]
			val := node["var_value"]
			meta := node["meta"].(map[string]interface{})

			if re2.MatchString(fmt.Sprint(val)) {
				v, _ := expr.Eval(fmt.Sprintln(val), variables)
				variables[name] = v
			} else {
				switch Type {
				case "int":
					switch meta["math"] {
					case true:
						cmp, err := expr.Compile(fmt.Sprint(val))
						Check(err)
						ans, errd := expr.Run(cmp, variables)
						Check(errd)
						variables[name] = ans
					case false:
						variables[name] = val
					}
				case "string", "char":
					switch meta["raw_type"] {
					case "STRING", "CHAR":
						val := val.(string)
						variables[name] = val

					case "FUNC":
						value := val.(map[string]interface{})
						file := value["read"]

						data, err := os.ReadFile(fmt.Sprint(file))
						Check(err)
						isBinary = binaryCheck(data)
						variables[name] = string(data)
					case "concat":
						var catch string
						for _, element := range val.([]interface{}) {
							if re2.MatchString(fmt.Sprint(element)) {
								v, _ := expr.Eval(fmt.Sprint(element), variables)
								catch += fmt.Sprint(v)
							} else {
								catch += fmt.Sprint(element)
							}
						}
						variables[name] = catch
					}
				case "bool":
					variables[name] = val
				case "float":
					variables[name] = val
				case "order", "ord":
					valList := val.([]interface{})

					for i, element := range valList {
						if _, ok := variables[fmt.Sprint(element)]; ok {
							value := variables[fmt.Sprint(element)]

							valList[i] = value
						}
					}

					variables[name] = valList
				}
			}
		case "output":
			value := node["value"]
			meta := node["meta"].(map[string]interface{})

			switch meta["print_type"] {
			case "simple":
				val, ok := variables[value.(string)]
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
					fmt.Println(NullCheck(val))
				} else {
					if re2.MatchString(fmt.Sprint(value)) {
						val, _ := expr.Eval(fmt.Sprint(value), variables)
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
						vari, ok := variables[SpagVal.(string)]
						if ok {
							SpagList = append(SpagList, vari)
						} else {
							if re2.MatchString(SpagVal.(string)) {
								val, _ := expr.Eval(SpagVal.(string), variables)
								SpagList = append(SpagList, NullCheck(val))
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
				ans, err := expr.Run(cmp, variables)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(ans)
			case "ord_index":
				OrdRef := value.(map[string]interface{})

				for key, val := range OrdRef {
					ListRaw := variables[key]
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
				val := variables[fmt.Sprint(input)]
				target := meta["target"]
				perms := meta["perms"]
				octal, _ := strconv.ParseInt(fmt.Sprint(perms), 8, 64)
				err := os.WriteFile(fmt.Sprint(target), []byte(val.(string)), 0000)
				Check(err)
				erra := os.Chmod(fmt.Sprint(target), os.FileMode(octal))
				Check(erra)
			case "append":
				fileTaget := meta["target"]
				val := fmt.Sprint(variables[fmt.Sprint(meta["input"])])
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
				if _, ok := variables[fileTarget]; ok {
					val = fmt.Sprint(variables[fileTarget])
					os.Remove(val)
				} else {
					os.Remove(fileTarget)
				}
			}
		}
	}
}

func ReadFile(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func binaryCheck(data []byte) bool {
	var nonPrintable int
	for _, b := range data {
		if (b < 32 || b > 126) && b != 9 && b != 10 && b != 13 {
			nonPrintable++
		}
	}
	return float64(nonPrintable)/float64(len(data)) > 0.3
}

func NullCheck(val interface{}) string {
	if val == nil {
		return "\033[0;35mnull\033[0m"
	} else {
		return fmt.Sprint(val)
	}
}
