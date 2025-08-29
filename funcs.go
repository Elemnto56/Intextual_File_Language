package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func ReRunParser(tokens []Tokens) []map[string]interface{} {
	ast := []map[string]interface{}{}

	// Iterators
	index := 0
	UUID := 0 // For if statements

	for index < len(tokens) {
		token := current(&index, tokens)

		if token.Type == "KEYWORD" {
			switch token.Val {
			case "let", "declare":
				meta := make(map[string]interface{})
				grandType := token.Val
				advance(&index)
				token := current(&index, tokens)
				if token.Type == "IDENTIFIER" {
					name := token.Val
					advance(&index)
					token := current(&index, tokens)
					if token.Type == "SYMBOL" && token.Val == ":" {
						advance(&index)
						token := current(&index, tokens)
						if token.Type == "TYPESYS" && Contains([]interface{}{"bool", "string", "int", "char", "float", "ord", "order"}, fmt.Sprint(token.Val)) {
							Type := token.Val
							advance(&index)
							token := current(&index, tokens)
							if token.Type == "OPERATOR" && token.Val == "=" {
								advance(&index)
								token := current(&index, tokens)
								if Contains([]interface{}{"INT", "BOOL", "STRING", "CHAR", "ORD", "IDENTIFIER", "TXT BLK"}, token.Type) {
									value := token.Val
									_type := token.Type
									meta["raw_type"] = _type

									temp := index + 1 // Did this in order for it to be a sight into the future
									if token.Meta["assignment"] == "math" {
										meta["math"] = true

										ast = append(ast, map[string]interface{}{
											"type":      grandType,
											"var_type":  Type,
											"var_name":  name,
											"var_value": value,
											"line":      token.Line,
											"meta":      meta,
										})
										advance(&index)
									} else if current(&temp, tokens).Type == "SYMBOL" && current(&temp, tokens).Val == ";" {
										advance(&index)
										meta["math"] = false
										ast = append(ast, map[string]interface{}{
											"type":      grandType,
											"var_type":  Type,
											"var_name":  name,
											"var_value": value,
											"line":      token.Line,
											"meta":      meta,
										})
									} else if (current(&temp, tokens).Type == "OPERATOR" && current(&temp, tokens).Val == "+") || (current(&temp, tokens).Type == "SYMBOL" && current(&temp, tokens).Val == ",") {
										first := current(&index, tokens).Val
										advance(&index)
										token := current(&index, tokens)
										concatCatch := []interface{}{}

										concatCatch = append(concatCatch, fmt.Sprint(first))

										for {
											if token.Type == "SYMBOL" && token.Val == ";" {
												break
											}

											if (token.Type == "SYMBOL" && token.Val == ",") || (token.Type == "OPERATOR" && token.Val == "+") {
												advance(&index)
												token = current(&index, tokens)
												continue
											}

											concatCatch = append(concatCatch, fmt.Sprint(token.Val))
											advance(&index)
											token = current(&index, tokens)
										}

										if token.Type == "SYMBOL" && token.Val == ";" {
											meta["raw_type"] = "concat"
											ast = append(ast, map[string]interface{}{
												"type":      grandType,
												"var_type":  Type,
												"var_name":  name,
												"var_value": concatCatch,
												"line":      token.Line,
												"meta":      meta,
											})
										}
									} else {
										err := NewError("MissingBreaker", token.Line, fmt.Sprintf("%v %v: %v = %v <-", grandType, name, Type, value), "This line is missing a semicolon", true, "")
										err.Throw()
									}
								} else if token.Type == "LBRACKET" {
									advance(&index)
									token := current(&index, tokens)
									userList := []interface{}{} // Make list for order

									userList = append(userList, token.Val) // Put in current token
									VarRef := make(map[string]interface{})

									var tempIndex int = 0
									advance(&index)
									for {
										token = current(&index, tokens)

										if Contains([]interface{}{"STRING", "INT", "BOOL", "FLOAT", "CHAR"}, token.Type) {
											userList = append(userList, token.Val)
											tempIndex += 1
											advance(&index)
											continue
										}

										if token.Type == "IDENTIFIER" {
											userList = append(userList, token.Val)
											tempIndex += 1
											VarRef[fmt.Sprint(token.Val)] = tempIndex
											advance(&index)
											continue
										}

										if token.Type == "COMMA" {
											advance(&index)
											continue
										}

										if token.Type == "RBRACKET" {
											break
										}
									}

									if token.Type == "RBRACKET" {
										meta["raw_type"] = "ORDER"
										meta["math"] = false
										meta["ord-ref"] = VarRef
										ast = append(ast, map[string]interface{}{
											"type":      grandType,
											"var_type":  Type,
											"var_name":  name,
											"var_value": userList,
											"line":      token.Line,
											"meta":      meta,
										})
										advance(&index)
									}
								} else if token.Type == "FUNC" {
									switch token.Val {
									case "read":
										if Type == "string" {
											advance(&index)
											token := current(&index, tokens)
											if token.Type == "PARA" {
												advance(&index)
												token = current(&index, tokens)
												file := token.Val // Grab the file trying to read
												advance(&index)
												token = current(&index, tokens)
												if token.Type == "PARA" {
													meta["raw_type"] = "FUNC"
													meta["math"] = false
													ast = append(ast, map[string]interface{}{
														"type":      grandType,
														"var_type":  Type,
														"var_name":  name,
														"var_value": map[string]interface{}{"read": file},
														"line":      token.Line,
														"meta":      meta,
													})
													advance(&index)
												} else {
													err := NewError("MalformedSyntax", token.Line, fmt.Sprintf("%v %v: %v = %sread(...%s;", grandType, name, Type, Red, Reset), "The following read function was not properly closed", true, "Add a paranthesis after the string that calls the file (e.g. read(...))")
													err.Throw()
												}
											} else {
												err := NewError("MalformedSyntax", token.Line, fmt.Sprintf("%v %v: %v = %sread...)%s;", grandType, name, Type, Red, Reset), "The following read function was not properly closed", true, "Add a paranthesis before the string that calls the file (e.g. read(...))")
												err.Throw()
											}
										} else {
											err := NewError("TypeMismatch", token.Line, fmt.Sprintf("%v %v: %s%v%s = read(...);", grandType, name, Red, Type, Reset), "The wrong type was assigned to read()", true, fmt.Sprintf("Change the type to %sstring%s", Yellow, Reset))
											err.Throw()
										}
									}
								} else if token.Type == "MATH" {
									value := token.Val
									meta["math"] = true
									meta["raw_type"] = "none"

									temp := index + 1
									if current(&temp, tokens).Type == "SYMBOL" && current(&temp, tokens).Val == ";" {
										advance(&index)
										ast = append(ast, map[string]interface{}{
											"type":      grandType,
											"var_type":  Type,
											"var_name":  name,
											"var_value": value,
											"line":      token.Line,
											"meta":      meta,
										})
									}
								} else {
									err := NewError("UnknownValue", token.Line, fmt.Sprintf("%v %v: %v = %s%v%s", grandType, name, Type, Red, token.Val, Reset), "The following value could not be correctly parsed", true, "Did you forget any quotes or accidently put a variable in the statement?")
									err.Throw()
								}
							} else {
								err := NewError("MalformedSyntax", token.Line, fmt.Sprintf("%v %v:%v %s??%s", grandType, name, Type, Red, Reset), "The following statement failed to abide by Intext's syntax rules", true, "The \"??\" is where an \"=\" is expected")
								err.Throw()
							}
						} else {
							err := NewError("TypeMismatch", token.Line, fmt.Sprintf("%v %v: %v <-", grandType, name, token.Val), "The following line does not include a valid type", true, "")
							err.Throw()
						}
					} else {
						err := NewError("LexerErr", token.Line, fmt.Sprintf("%v %v%v <-", grandType, name, token.Val), "Invalid character on this line", true, "Did you mean ':'?")
						err.Throw()
					}
				}
			case "output":
				advance(&index)
				token := current(&index, tokens)
				if true { // Added this here because the value is ambiguous
					meta := make(map[string]string)
					_type := token.Type
					val := token.Val

					temp := index + 1
					if (current(&temp, tokens).Type == "SYMBOL" && current(&temp, tokens).Val == ";") || current(&temp, tokens).Type == "LBRACKET" {
						switch token.Type {
						case "STRING", "INT", "BOOL", "FLOAT", "ORD", "CHAR", "IDENTIFIER":
							meta["raw_type"] = _type
							meta["print_type"] = "simple"
							ast = append(ast, map[string]interface{}{
								"type":  "output",
								"value": val,
								"meta":  meta,
								"line":  token.Line,
							})
						case "MATH":
							meta["print_type"] = "mathematics"
							meta["raw_type"] = "none"
							ast = append(ast, map[string]interface{}{
								"type":  "output",
								"value": val,
								"meta":  meta,
								"line":  token.Line,
							})
						}
					} else if (current(&temp, tokens).Type == "SYMBOL" || current(&temp, tokens).Type == "COMMA") && current(&temp, tokens).Val == "," {
						spagList := []interface{}{}
						spagList = append(spagList, val) // Add the first val into the list; none left behind!
						advance(&index)

						for {
							newVal := current(&index, tokens)
							var i int = index + 1

							if (newVal.Type == "SYMBOL" || newVal.Type == "COMMA") && newVal.Val == "," {
								advance(&index)
								continue
							}

							spagList = append(spagList, interface{}(newVal.Val))
							advance(&index)

							if current(&i, tokens).Type == "SYMBOL" && current(&i, tokens).Val == ";" {
								break
							}
						}
						token := current(&index, tokens)
						if (token.Type == "SYMBOL" || token.Type == "COMMA") && token.Val == ";" {
							meta["print_type"] = "mixed"
							meta["raw_type"] = "none"
							ast = append(ast, map[string]interface{}{
								"type":  "output",
								"value": spagList,
								"meta":  meta,
								"line":  token.Line,
							})
						} else {
							err := NewError("MissingBreaker", token.Line, fmt.Sprintf("output %v <-", val), "Missing semicolon", true, "")
							err.Throw()
						}

					} else {
						err := NewError("MissingBreaker", token.Line, fmt.Sprintf("output %v <-", val), "Missing semicolon", true, "")
						err.Throw()
					}
				}
			}
		} else if token.Type == "FUNC" {
			switch token.Val {
			case "write":
				meta := make(map[string]interface{})
				advance(&index)
				token := current(&index, tokens)
				if token.Type == "PARA" && token.Val == "(" {
					advance(&index)
					token = current(&index, tokens)
					if token.Type == "STRING" {
						targetFile := token.Val
						advance(&index)
						token = current(&index, tokens)
						if token.Type == "SYMBOL" && token.Val == "," {
							advance(&index)
							token = current(&index, tokens)
							if token.Type == "IDENTIFIER" {
								wrVar := token.Val
								advance(&index)
								token = current(&index, tokens)
								if token.Type == "PARA" && token.Val == ")" {
									advance(&index)
									meta["target"] = targetFile
									meta["input"] = wrVar
									meta["perms"] = 666
									ast = append(ast, map[string]interface{}{
										"type": "function",
										"line": token.Line,
										"call": "write",
										"meta": meta,
									})
								} else if token.Type == "SYMBOL" && token.Val == "," {
									advance(&index)
									token = current(&index, tokens)
									if token.Type == "INT" {
										perms := token.Val
										advance(&index)
										token = current(&index, tokens)
										if token.Type == "PARA" && token.Val == ")" {
											meta["target"] = targetFile
											meta["input"] = wrVar
											meta["perms"] = perms
											ast = append(ast, map[string]interface{}{
												"type": "function",
												"line": token.Line,
												"call": "write",
												"meta": meta,
											})
										}
									} else {
										err := NewError("TypeMismatch", token.Line, fmt.Sprintf("write(%v, %v, %s%v%s);", targetFile, wrVar, Red, token.Val, Reset), "The permission input was not a number", true, "Write the permission number as you'd do for Linux's chmod. Example: 644.")
										err.Throw()
									}
								} else {
									err := NewError("MalformedSyntax", token.Line, fmt.Sprintf("write(..., ...%s???%s;", Red, Reset), "The following function was not properly closed", true, "Add a paranthesis after the string that inputs to the file (e.g. write(..., ...))")
									err.Throw()
								}
							} else {
								err := NewError("MalformedSyntax", token.Line, fmt.Sprintf("write(%v, %s%v%s);", targetFile, Red, token.Val, Reset), "A variable was not used for the input argument", true, "Literals are not allowed in write() besides for the file call argument")
								err.Throw()
							}
						} else {
							err := NewError("MalformedSyntax", token.Line, fmt.Sprintf("write(%v%s???%s...);", targetFile, Red, Reset), "A comma was missing in the write function", true, fmt.Sprintf("The %s\"???\"%s indicates where to insert the comma", Red, Reset))
							err.Throw()
						}
					} else {
						err := NewError("TypeMismatch", token.Line, fmt.Sprintf("write(%s%v%s, ...);", Red, token.Val, Reset), "The value for the file was not a string", true, "")
						err.Throw()
					}
				} else {
					err := NewError("MalformedSyntax", token.Line, fmt.Sprintf("write%s???%s..., ...);", Red, Reset), "The following function was not properly closed", true, "Add a paranthesis before the string that calls the file (e.g. write(..., ...))")
					err.Throw()
				}

			case "append":
				meta := make(map[string]interface{})
				advance(&index)
				token := current(&index, tokens)
				if token.Type == "PARA" && token.Val == "(" {
					advance(&index)
					token = current(&index, tokens)
					if token.Type == "STRING" {
						target := token.Val
						advance(&index)
						token = current(&index, tokens)
						if token.Type == "SYMBOL" && token.Val == "," {
							advance(&index)
							token = current(&index, tokens)
							if token.Type == "IDENTIFIER" {
								input := token.Val
								advance(&index)
								token = current(&index, tokens)
								if token.Type == "PARA" && token.Val == ")" {
									advance(&index)
									meta["perms"] = 0
									meta["input"] = input
									meta["target"] = target
									ast = append(ast, map[string]interface{}{
										"type": "function",
										"line": token.Line,
										"call": "append",
										"meta": meta,
									})
								} else {
									err := NewError("MalformedSyntax", token.Line, fmt.Sprintf("append(..., ...%s???%s;", Red, Reset), "The following function was not properly closed", true, "Add a paranthesis after the string that inputs to the file (e.g. append(..., ...))")
									err.Throw()
								}
							} else {
								err := NewError("MalformedSyntax", token.Line, fmt.Sprintf("append(%v, %s%v%s);", target, Red, token.Val, Reset), "A variable was not used for the input argument", true, "Literals are not allowed in append() besides for the file call argument")
								err.Throw()
							}
						} else {
							err := NewError("MalformedSyntax", token.Line, fmt.Sprintf("append(%v%s???%s...);", target, Red, Reset), "A comma was missing in the write function", true, fmt.Sprintf("The %s\"???\"%s indicates where to insert the comma", Red, Reset))
							err.Throw()
						}
					} else {
						err := NewError("TypeMismatch", token.Line, fmt.Sprintf("append(%s%v%s, ...);", Red, token.Val, Reset), "The value for the file was not a string", true, "")
						err.Throw()
					}
				} else {
					err := NewError("MalformedSyntax", token.Line, fmt.Sprintf("append%s???%s..., ...);", Red, Reset), "The following function was not properly closed", true, "Add a paranthesis before the string that calls the file (e.g. append(..., ...))")
					err.Throw()
				}
			case "del", "remove":
				meta := make(map[string]interface{})
				grandType := token.Val
				advance(&index)
				token := current(&index, tokens)
				if token.Type == "PARA" && token.Val == "(" {
					advance(&index)
					token = current(&index, tokens)
					if token.Type == "STRING" || token.Type == "IDENTIFIER" {
						meta["raw"] = token.Type
						targetFile := token.Val
						advance(&index)
						token = current(&index, tokens)
						if token.Type == "PARA" && token.Val == ")" {
							meta["target"] = targetFile
							advance(&index)
							ast = append(ast, map[string]interface{}{
								"type": "function",
								"line": token.Line,
								"call": "del",
								"meta": meta,
							})
						} else {
							err := NewError("MalformedSyntax", token.Line, fmt.Sprintf("%v(...%s???%s;", grandType, Red, Reset), "The following function was not properly closed", true, "Add a paranthesis after the string that calls the file (e.g. del(...))")
							err.Throw()
						}
					} else {
						err := NewError("TypeMismatch", token.Line, fmt.Sprintf("%v(%s%v%s);", grandType, Red, token.Val, Reset), "The value for the file was either not a string or variable", true, "")
						err.Throw()
					}
				} else {
					err := NewError("MalformedSyntax", token.Line, fmt.Sprintf("%v%s???%s...);", grandType, Red, Reset), "The following function was not properly closed", true, "Add a paranthesis before the string that calls the file (e.g. del(...))")
					err.Throw()
				}
			}

		} else if token.Type == "LOGIC" {
			meta := make(map[string]interface{})

			switch token.SubType {
			case "if", "else", "or":
				condition := token.Val
				grandType := token.SubType

				if fmt.Sprint(condition) == "" {
					condition = true
				}

				ifLine := token.Line

				advance(&index)
				token := current(&index, tokens)
				if token.Type == "BODY" {
					var toks []Tokens
					val := token.Val.([]interface{})

					b, _ := json.Marshal(val)
					json.Unmarshal(b, &toks)

					astBody := ReRunParser(toks)

					advance(&index)
					token := current(&index, tokens)

					if token.Type == "SYMBOL" && token.Val == ";" {
						meta["sub_type"] = "if"
						if grandType == "if" {
							meta["UUID"] = UUID
							UUID++
						}
						ast = append(ast, map[string]interface{}{
							"type":      "logic",
							"meta":      meta,
							"line":      ifLine,
							"condition": condition,
							"body":      astBody,
						})
					}
				} else {
					err := NewError("MalformedSyntax", ifLine, fmt.Sprintf("if %v %s???%s", condition, Red, Reset), "This if-statement is missing either a body, left-standing curly brace, or right-standing curly brace", true, "")
					err.Throw()
				}
			case "while":
				condition := token.Val
				whileLine := token.Line

				advance(&index)
				token := current(&index, tokens)

				if token.Type == "BODY" {
					var toks []Tokens
					val := token.Val.([]interface{})

					b, _ := json.Marshal(val)
					json.Unmarshal(b, &toks)

					astBody := ReRunParser(toks)

					advance(&index)
					token := current(&index, tokens)
					if token.Type == "SYMBOL" && token.Val == ";" {
						meta["sub_type"] = "while"
						ast = append(ast, map[string]interface{}{
							"type":      "logic",
							"meta":      meta,
							"line":      whileLine,
							"condition": condition,
							"body":      astBody,
						})
					}
				} else {
					err := NewError("MalformedSyntax", whileLine, fmt.Sprintf("while %v %s???%s", condition, Red, Reset), "This while loop is missing either a body, left-standing curly brace, or right-standing curly brace", true, "")
					err.Throw()
				}
			case "repeat":
				cond := fmt.Sprint(token.Val)
				reLine := token.Line

				if cmpRegEx(cond, `\(?([A-Za-z]|\_|\d+)+\s->\s\(?.+\)?`) {
					findDash := strings.IndexRune(cond, '-')
					iterator := strings.Replace(strings.TrimSpace(cond[:findDash]), "(", "", 1)

					findGTT := strings.IndexRune(cond, '>')
					repeatValue := strings.Replace(strings.TrimSpace(cond[1+findGTT:]), ")", "", 1)

					advance(&index)
					token := current(&index, tokens)

					if token.Type == "BODY" {
						var toks []Tokens
						val := token.Val.([]interface{})

						b, _ := json.Marshal(val)
						json.Unmarshal(b, &toks)

						astBody := ReRunParser(toks)

						advance(&index)
						token := current(&index, tokens)
						if token.Type == "SYMBOL" && token.Val == ";" {
							meta["sub_type"] = "repeat"
							meta["iterator_var"] = iterator
							meta["times"] = repeatValue
							ast = append(ast, map[string]interface{}{
								"type": "logic",
								"meta": meta,
								"line": reLine,
								"body": astBody,
							})
						}
					} else {
						err := NewError("MalformedSyntax", reLine, fmt.Sprintf("repeat %v -> %v %s???%s", iterator, repeatValue, Red, Reset), "This repeat loop is missing either a body, left-standing curly brace, or right-standing curly brace", true, "")
						err.Throw()
					}
				} else if cmpRegEx(cond, `\(?\d+\)?`) {
					advance(&index)
					token := current(&index, tokens)

					if token.Type == "BODY" {
						var toks []Tokens
						val := token.Val.([]interface{})

						b, _ := json.Marshal(val)
						json.Unmarshal(b, &toks)

						astBody := ReRunParser(toks)

						advance(&index)
						token := current(&index, tokens)
						if token.Type == "SYMBOL" && token.Val == ";" {
							meta["sub_type"] = "repeat"
							meta["iterator_var"] = nil
							meta["times"] = cond
							ast = append(ast, map[string]interface{}{
								"type": "logic",
								"meta": meta,
								"line": reLine,
								"body": astBody,
							})
						}
					} else {
						err := NewError("MalformedSyntax", reLine, fmt.Sprintf("repeat %v %s???%s \n ... }", cond, Red, Reset), "This repeat loop is missing either a body, left-standing curly brace, or right-standing curly brace", true, "")
						err.Throw()
					}
				}
			}
		} else if token.Type == "IDENTIFIER" {
			meta := make(map[string]interface{})

			exprVal := token.Val

			if Contains([]interface{}{"break", "stop"}, exprVal) {
				ast = append(ast, map[string]interface{}{
					"statement": "break",
					"line":      token.Line,
				})
			} else if Contains([]interface{}{"continue", "skip"}, exprVal) {
				ast = append(ast, map[string]interface{}{
					"statement": "continue",
					"line":      token.Line,
				})
			} else {
				advance(&index)
				token := current(&index, tokens)

				if token.Type == "OPERATOR" {
					switch token.Val {
					case "+=", "*=", "-=", "/=":
						incrType := token.Val

						advance(&index)
						token = current(&index, tokens)

						incrValue := token.Val

						advance(&index)
						token = current(&index, tokens)

						if token.Type == "SYMBOL" && token.Val == ";" {
							meta["incr_type"] = incrType
							meta["target"] = exprVal
							meta["new_val"] = incrValue
							ast = append(ast, map[string]interface{}{
								"type":     "expr",
								"sub_type": "incr",
								"meta":     meta,
								"line":     token.Line,
							})
						} else {
							err := NewError("UnknownValue", token.Line, fmt.Sprintf("%v %v %v %s<--%s", exprVal, incrType, incrValue, Red, Reset), "An unknown value or known misplaced value is in this line", true, "self-incr can only have a single value to be added to itself")
							err.Throw()
						}
					case "=":
						err := NewError("PKGError", token.Line, fmt.Sprintf("%v = ...", exprVal), "The \"expression\" pkg was not included", true, "Upgrade to v0.9!")
						err.Throw()

					}
				}
			}
		}

		advance(&index)
	}

	return ast
}

func reRunLexer(lines []string) []map[string]interface{} {
	allTokens := []map[string]interface{}{}

	for index := 0; index < len(lines); index++ {
		line := strings.TrimSpace(lines[index])

		if strings.Contains(line, "//") { // Checks if line is a comment
			line = strings.Split(line, "//")[0]
		}

		if line == "" || line == " " { // Checks if line is space
			continue
		}

		if cmpRegEx(line, `crunch\(\s?.+\,\s?.+\,\s?.+\);?`) {
			err := NewError("DeprecationErr", index+1, line, "A deprecated function was found", false, "crunch() is no longer used. Do math instead. (e.g. output 5 + 5).")
			err.Throw()
		}

	outer: // Label for loop
		for i := 0; i < len(line); i++ { // NOTE: Some are place early so the others behind don't get triggered beforehand
			char := rune(line[i])

			// Checking if the current char is a space
			if unicode.IsSpace(char) || char == ' ' || string(char) == " " {
				continue // Adding these so the loop continues; it'll get stuck here not knowing what to do
			}

			// Check for arrow(s)
			if i+1 < len(line) && string(line[i:i+2]) == "->" {
				allTokens = append(allTokens, map[string]interface{}{
					"TYPE": "ARROW",
					"VAL":  string(line[i : i+2]),
					"LINE": index + 1,
				})
				i += 1
				continue
			}

			// Checks for math operators
			if i+1 < len(line) && Contains([]interface{}{">=", "<=", "==", "+=", "*=", "-=", "/=", "||", "&&"}, string(line[i:i+2])) {
				allTokens = append(allTokens, map[string]interface{}{
					"TYPE": "OPERATOR",
					"VAL":  string(line[i : i+2]),
					"LINE": index + 1,
				})
				i += 1
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
			if Contains([]interface{}{";", ",", ":", "."}, string(char)) {
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
				} else if Contains([]interface{}{"while", "repeat"}, temp) {
					var rawLogicCatch string

					for tempI := 0; tempI < len(line) && string(line[i]) != "{"; {
						tempI = i

						rawLogicCatch += string(line[i])

						i++
					}

					logicCatch := strings.TrimSpace(rawLogicCatch)

					allTokens = append(allTokens, map[string]interface{}{
						"TYPE":     "LOGIC",
						"VAL":      logicCatch,
						"SUB-TYPE": temp,
						"LINE":     index + 1,
					})
				} else if Contains([]interface{}{"if", "else", "else if", "or if"}, temp) {
					var rawCatch string

					for tempI := 0; tempI < len(line) && string(line[i]) != "{"; {
						tempI = i

						rawCatch += string(line[i])

						i++
					}

					rawwCatch := strings.Replace(rawCatch, "if ", "", 1)
					cleanCatch := strings.TrimSpace(rawwCatch)
					allTokens = append(allTokens, map[string]interface{}{
						"TYPE":     "LOGIC",
						"VAL":      cleanCatch,
						"SUB-TYPE": temp,
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
				i += 3
				captureBLK := []interface{}{}

				for index++; index < len(lines); {
					if strings.HasPrefix(lines[index], "]]]") {
						break
					}
					captureBLK = append(captureBLK, lines[index])
					index += 1
				}

				allTokens = append(allTokens, map[string]interface{}{
					"TYPE": "TXT BLK",
					"VAL":  captureBLK,
					"LINE": index + 1,
				})
				continue
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
				index++
				var rawLines []string
				var toks []map[string]interface{}

				for index < len(lines) {
					if cmpRegEx(strings.TrimSpace(lines[index]), `\}\s+(((else\s|or\s)if)|else).+\{`) || strings.TrimSpace(lines[index]) == "}" {
						break
					}
					rawLines = append(rawLines, lines[index])

					index++
				}

				toks = reRunLexer(rawLines)

				allTokens = append(allTokens, map[string]interface{}{
					"TYPE": "BODY",
					"VAL":  toks,
					"LINE": index + 1,
				})
				index--
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
					if i+1 < len(line) && Contains([]interface{}{")", "]", ",", "{"}, string(line[i+1])) {
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

			err := NewError("LexerErr", index+1, line, "Invalid character somewhere in expression", false, fmt.Sprintf("The suspected character is \"%s\"", string(char)))
			err.Throw()
		}

		if len(allTokens) > 0 {
			last := allTokens[len(allTokens)-1]
			if !(Contains([]interface{}{"SYMBOL", "LCURL", "RCURL", "OPERATOR"}, last["TYPE"]) || Contains([]interface{}{";", "{", "}", "="}, last["VAL"])) {
				allTokens = append(allTokens, map[string]interface{}{
					"TYPE": "SYMBOL",
					"VAL":  ";",
					"LINE": index + 1,
				})
			}
		}
	}

	return allTokens
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

func NullCheck(val interface{}, stringCheck bool) string {
	if !stringCheck {
		if val == nil {
			return "\033[0;35mnull\033[0m"
		} else {
			return fmt.Sprint(val)
		}
	} else if stringCheck {
		v := fmt.Sprint(val)

		NewString := strings.Replace(v, "<nil>", "\033[0;35mnull\033[0m", -1)
		return NewString
	}

	return ""
}

func ValidateVal(varType interface{}, varValue interface{}, line int, meta string) bool {
	switch varType {
	case "int":
		_, err := strconv.Atoi(fmt.Sprint(varValue))
		if err != nil {
			if meta == "none" || meta == "IDENTIFIER" {
				return true
			}
			err0 := NewError("TypeMismatch", line, fmt.Sprintf("let x: int = %s%s%s;", Red, fmt.Sprint(varValue), Reset), "The following value was not an int", true, typemismatch)
			err0.Throw()
			return false
		}
		return true
	case "string":
		if Contains([]interface{}{"STRING", "FUNC", "concat", "TXT BLK"}, meta) {
			return true
		} else {
			err2 := NewError("TypeMismatch", line, fmt.Sprintf("let x: string = %s%s%s;", Red, fmt.Sprint(varValue), Reset), "The following value was not a string", true, typemismatch)
			err2.Throw()
			return false
		}
	case "float":
		_, err := strconv.ParseFloat(fmt.Sprint(varValue), 64)
		if err != nil {
			err3 := NewError("TypeMismatch", line, fmt.Sprintf("let x: float = %s;", fmt.Sprint(varValue)), "The following value was not a float", true, typemismatch)
			err3.Throw()
		}
		return true
	case "bool":
		if fmt.Sprint(varValue) == "true" || fmt.Sprint(varValue) == "false" {
			return true
		}
		err4 := NewError("TypeMismatch", line, fmt.Sprintf("let x: bool = %s;", fmt.Sprint(varValue)), "The following value was not a boolean", true, typemismatch)
		err4.Throw()
		return false
	case "char":
		if len(fmt.Sprint(varValue)) == 3 && fmt.Sprint(varValue)[0] == '\'' && fmt.Sprint(varValue)[2] == '\'' {
			return true
		}
		err5 := NewError("TypeMismatch", line, fmt.Sprintf("let x: char = %s;", fmt.Sprint(varValue)), "The following value was not a char", true, typemismatch)
		err5.Throw()
		return false
	case "ord", "order":
		/*
			_, ok := interface{}(varValue.(string)).([]interface{})
			if !ok {
				err6 := NewError("TypeMismatch", line, fmt.Sprintf("let x: ord = %s;", varValue.(string)), "The following value was not an order", true, typemismatch)
				err6.Throw()
				return false
			}
		*/
		return true
	}
	return false
}

func cmpRegEx(find string, regex string) bool {
	temp := regexp.MustCompile(regex)

	if temp.MatchString(find) {
		return true
	}
	return false
}
