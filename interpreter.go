package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/expr-lang/expr"
)

// Globals
type MarcoDef struct {
	calls      map[string]interface{}
	parameters map[string]interface{}
	returns    []interface{}
	body       []map[string]interface{}
}

var isBinary bool
var pat2 string = `^([A-Za-z]+\_*?)+\[([0-9]+|\w)\];?$` // For order indexes (i.e. x[2])
var re2 *regexp.Regexp = regexp.MustCompile(pat2)
var breakFlag bool = false
var contFlag bool = false
var BackupVariables = make(map[string]interface{}) // Might use for later things?
var macroTable = make(map[string]MarcoDef)

func Interpreter(givenNodes []map[string]interface{}, InterpreterVariables map[string]interface{}) {
	// Func Globals
	var uuid int = 0
	MacroVariables := make(map[string]interface{})

	// Grab AST
	bytes, _ := os.ReadFile("./.intext/cache/AST.json")
	nodes := []map[string]interface{}{}
	err := json.Unmarshal(bytes, &nodes)
	Check(err)

	if givenNodes != nil {
		nodes = givenNodes
	}

	if InterpreterVariables == nil {
		InterpreterVariables = BackupVariables
	}

	// Iterate through each node
	for index := 0; index < len(nodes); index++ {
		node := nodes[index]
		line, _ := strconv.Atoi(fmt.Sprint(node["line"]))

		rawUUID := node["meta"]
		if rawUUID != nil {
			uuid, _ = strconv.Atoi(fmt.Sprint(rawUUID.(map[string]interface{})["UUID"]))
		}

		switch node["statement"] {
		case "break":
			breakFlag = true
		case "continue":
			contFlag = true
		}

		switch node["type"] {
		case "let":
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
					case "TXT BLK":
						var rawCatch string
						for _, line := range val.([]interface{}) {
							rawCatch += fmt.Sprintf("%s \n", fmt.Sprint(line))
						}
						catch := strings.TrimSpace(rawCatch)
						InterpreterVariables[name] = catch
					}
				case "bool":
					InterpreterVariables[name] = val
				case "float":
					InterpreterVariables[name] = val
				case "order", "ord":
					valList := val.([]interface{})
					var capture []interface{}

					for _, element := range valList {
						forUse := element.(map[string]interface{})
						if forUse["type"] == "IDENTIFIER" {
							a := InterpreterVariables[fmt.Sprint(forUse["val"])]
							capture = append(capture, a)
						} else {
							capture = append(capture, forUse["val"])
						}
					}

					InterpreterVariables[name] = capture
				}
			}

		case "declare":
			// TODO: Add macro giving return values `declare x = myMacro()`.
			// Also, add standalone macro calls, and variable macro calls (x.myMacro())

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

				for _, SpagVal := range value.([]interface{}) {
					switch SpagVal := SpagVal.(type) {
					case string:
						vari, ok := InterpreterVariables[SpagVal]
						if ok {
							SpagList = append(SpagList, vari)
						} else {
							if re2.MatchString(SpagVal) { // Order index check
								val, _ := expr.Eval(SpagVal, InterpreterVariables)
								SpagList = append(SpagList, NullCheck(val, false))
							} else {
								SpagList = append(SpagList, SpagVal)
							}
						}
					case []interface{}:
						var capture []interface{}
						for _, element := range SpagVal {
							forUse := element.(map[string]interface{})
							if forUse["TYPE"] == "IDENTIFIER" {
								a := InterpreterVariables[fmt.Sprint(forUse["VAL"])]
								capture = append(capture, a)
							} else {
								capture = append(capture, forUse["VAL"])
							}
						}
						SpagList = append(SpagList, capture)
					default:
						SpagList = append(SpagList, SpagVal)
					}
				}

				for _, i := range SpagList {
					fmt.Print(i)
				}
				fmt.Println()
			case "mathematics":
				val := fmt.Sprint(value)
				ans, err := EvaluateExpression(val, InterpreterVariables)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(ans)
			case "ord_index":
				OrdRef := value.(map[string]interface{})

				for key, val := range OrdRef {
					ListRaw := InterpreterVariables[key]
					List := ListRaw.([]interface{})
					var i int
					if _, ok := InterpreterVariables[fmt.Sprint(val)]; ok {
						i, _ = strconv.Atoi(fmt.Sprint(InterpreterVariables[fmt.Sprint(val)]))
					} else {
						i, _ = strconv.Atoi(fmt.Sprint(val))
					}

					if i >= len(List) {
						fmt.Println(NullCheck(nil, false))
					} else {
						fmt.Println(List[i])
					}
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
				f.Close()
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
				meta["UUID"] = uuid

				cond := fmt.Sprint(node["condition"])
				val, err := EvaluateExpression(cond, InterpreterVariables)
				Check(err)
				captureAST := []map[string]interface{}{}

				for _, element := range body {
					captureAST = append(captureAST, element.(map[string]interface{}))
				}

				v, _ := strconv.ParseBool(fmt.Sprint(val))

				if v == true {
					if i, _ := strconv.Atoi(fmt.Sprint(meta["UUID"])); i == uuid {
						index++
					}
					Interpreter(captureAST, nil)
					if i, _ := strconv.Atoi(fmt.Sprint(meta["UUID"])); i == uuid {
						index++
					}
				}
				index -= 2
			case "while":
				cond := fmt.Sprint(node["condition"])
				body := node["body"].([]interface{})

				whileCapture := []map[string]interface{}{}

				for _, element := range body {
					whileCapture = append(whileCapture, element.(map[string]interface{}))
				}

				rawLogic, err := EvaluateExpression(cond, InterpreterVariables)
				Check(err)
				logic, _ := strconv.ParseBool(strings.TrimSpace(fmt.Sprint(rawLogic)))

				for logic { // Updates logic
					if contFlag {
						contFlag = false
						Interpreter(whileCapture, nil)
						continue
					}

					if breakFlag {
						breakFlag = false
						Interpreter(whileCapture, nil)
						break
					}

					rawLogic, err = EvaluateExpression(cond, InterpreterVariables)
					Check(err)
					logic, _ = strconv.ParseBool(strings.TrimSpace(fmt.Sprint(rawLogic)))
					Interpreter(whileCapture, nil)
				}
			case "repeat":
				rawItr := meta["iterator_var"]
				itr := fmt.Sprint(rawItr)

				rawTimes := fmt.Sprint(meta["times"])
				rawTimesTwo, _ := EvaluateExpression(rawTimes, InterpreterVariables) // Incase math is present
				times, err := strconv.Atoi(fmt.Sprint(rawTimesTwo))
				Check(err)

				body := node["body"].([]interface{})
				repeatCapture := []map[string]interface{}{}

				for _, element := range body {
					repeatCapture = append(repeatCapture, element.(map[string]interface{}))
				}
				if rawItr != nil {
				inner:
					for i := 0; i < times; i++ {
						InterpreterVariables[itr] = i

						if contFlag {
							contFlag = false
							Interpreter(repeatCapture, nil)
							continue inner
						}

						if breakFlag {
							breakFlag = false
							break
						}
						Interpreter(repeatCapture, nil)
					}
				} else if rawItr == nil {
					for i := 0; i < times; i++ {
						if contFlag {
							contFlag = false
							Interpreter(repeatCapture, nil)
							continue
						}

						if breakFlag {
							breakFlag = false
							break
						}

						Interpreter(repeatCapture, nil)
					}
				}
			}
		case "expr":
			meta := node["meta"].(map[string]interface{})

			switch node["sub_type"] {
			case "incr":
				line, _ := strconv.Atoi(fmt.Sprint(node["line"]))

				incrType := meta["incr_type"]
				newValue := meta["new_val"]

				rawTarget := fmt.Sprint(meta["target"])
				target := InterpreterVariables[rawTarget]

				var strCheck bool = false
				var intTarget int
				if cmpRegEx(fmt.Sprint(target), `[\W|A-Za-z|\_]+`) {
					strCheck = true
				} else {
					var err error
					intTarget, err = strconv.Atoi(fmt.Sprint(target))
					if err != nil {
						err0 := NewError("TypeMismatch", line, fmt.Sprintf("%s%v%s %v %v", Red, rawTarget, Reset, incrType, newValue), "The following variable was not an int", true, typemismatch)
						err0.Throw()

					}
				}

				var intNewVal int
				if cmpRegEx(fmt.Sprint(newValue), `[\W|A-Za-z|\_]+`) {
					strCheck = true
				} else {
					var erra error
					if val, ok := InterpreterVariables[fmt.Sprint(newValue)]; ok {
						intNewVal, _ = strconv.Atoi(fmt.Sprint(val))
					} else {
						intNewVal, erra = strconv.Atoi(fmt.Sprint(newValue))
						if erra != nil {
							err0 := NewError("TypeMismatch", line, fmt.Sprintf("%v %v %s%v%s", rawTarget, incrType, Red, newValue, Reset), "The following value was not an int", true, typemismatch)
							err0.Throw()

						}
					}
				}

				switch incrType {
				case "+=":
					if strCheck {
						val, err := EvaluateExpression(fmt.Sprintf("\"%s\" + \"%s\"", fmt.Sprint(target), fmt.Sprint(newValue)), InterpreterVariables)
						Check(err)
						InterpreterVariables[rawTarget] = val
					} else {
						val, err := EvaluateExpression(fmt.Sprintf("%v + %v", target, newValue), InterpreterVariables)
						Check(err)
						InterpreterVariables[rawTarget] = val
					}
				case "-=":
					val := intTarget - intNewVal

					InterpreterVariables[rawTarget] = val
				case "*=":
					val := intTarget * intNewVal

					InterpreterVariables[rawTarget] = val
				case "/=":
					val := intTarget / intNewVal

					InterpreterVariables[rawTarget] = val
				}
			case "reassign":
				target := fmt.Sprint(meta["target"])
				value := fmt.Sprint(meta["value"])

				evald, err := EvaluateExpression(value, InterpreterVariables)
				Check(err)
				InterpreterVariables[target] = evald
			}

		case "macro":
			meta := node["meta"].(map[string]interface{})

			switch meta["macro-type"] {
			case "declaration":
				rawBody := node["body"].([]interface{})
				body := []map[string]interface{}{}

				for _, element := range rawBody {
					body = append(body, element.(map[string]interface{}))
				}

				macroTable[fmt.Sprint(node["name"])] = MarcoDef{
					calls:      meta["call"].(map[string]interface{}),
					parameters: meta["param"].(map[string]interface{}),
					returns:    meta["returns"].([]interface{}),
					body:       body,
				}

			case "standalone":
				args := meta["args"]
				name := fmt.Sprint(node["name"])

				if _, ok := macroTable[name]; ok {
					if args == nil {
						mac := macroTable[name]

						Interpreter(mac.body, MacroVariables)
					} else {
						args := args.([]interface{})
						calledMacro := macroTable[name]

						for topIndex, topElement := range args {
							topUse := topElement.(map[string]interface{})
							if topUse["type"] == "ORDER" {
								for lowerIndex, lowerElement := range topUse["val"].([]interface{}) {
									lowerUSe := lowerElement.(map[string]interface{})

									if lowerUSe["TYPE"] == "IDENTIFIER" {
										if a, ok := InterpreterVariables[fmt.Sprint(lowerUSe["VAL"])]; ok {
											topUse["val"].([]interface{})[lowerIndex] = a
										} else {
											err := NewError("VariableNotFound", line, fmt.Sprintf("%v(... [... %s%v%s ...] ...)", name, Red, lowerUSe["VAL"], Reset), "The variable inside this order was not found", true, "")
											err.Throw()
										}
									}
								}
							}

							if topUse["type"] == "IDENTIFIER" {
								if a, ok := InterpreterVariables[fmt.Sprint(topUse["val"])]; ok {
									args[topIndex] = a
								} else {
									err := NewError("VariableNotFound", line, fmt.Sprintf("%v(... %s%v%s ...)", name, Red, topUse["VAL"], Reset), "The variable inside this macro call was not found", true, "")
									err.Throw()
								}
							}
						}

						var count int = 0
						for paramName, expectedType := range calledMacro.parameters {
							if fmt.Sprint(expectedType) == strings.ToLower(fmt.Sprint(args[count].(map[string]interface{})["type"])) {
								MacroVariables[paramName] = args[count].(map[string]interface{})["val"]
							} else {
								err := NewError("TypeMismatch", line, fmt.Sprintf("%v(... %s%v%s(%v) ...) -> %v(%s%v%s)",
									name,
									Red,
									strings.ToLower(fmt.Sprint(args[count].(map[string]interface{})["type"])),
									Reset,
									fmt.Sprint(args[count].(map[string]interface{})["val"]),
									name,
									Green,
									expectedType,
									Reset,
								), "A data type of a value given during the macro call was unexpected", true, "")
								err.Throw()
							}

							count++
						}

						Validator(calledMacro.body, MacroVariables)
						Interpreter(calledMacro.body, MacroVariables)
					}
				}
			}
		}
	}

}
