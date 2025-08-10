package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Tokens struct {
	Type    string      `json:"TYPE"`
	SubType string      `json:"SUB-TYPE"`
	Val     interface{} `json:"VAL"`
	Line    int         `json:"LINE"`
}

func advance(index *int) {
	*index++
}

func current(index *int, tokens []Tokens) Tokens {
	if *index >= len(tokens) {
		panic("Tokens ranged out. Semicolon or paranthesis may not have been closed.")
	}
	return tokens[*index]
}

func Parser() {
	// Get tokens from JSON
	bytes, _ := os.ReadFile("./.intext/cache/Tokens.json")
	var tokens []Tokens
	err := json.Unmarshal(bytes, &tokens)
	Check(err)

	// AST
	ast := []map[string]interface{}{}

	index := 0
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
								if Contains([]interface{}{"INT", "BOOL", "STRING", "CHAR", "ORD", "IDENTIFIER"}, token.Type) {
									value := token.Val
									_type := token.Type
									meta["raw_type"] = _type
									meta["math"] = false

									temp := index + 1 // Did this in order for it to be a sight into the future
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
					if current(&temp, tokens).Type == "SYMBOL" && current(&temp, tokens).Val == ";" {
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
							advance(&index)
						case "MATH":
							meta["print_type"] = "mathematics"
							meta["raw_type"] = "none"
							ast = append(ast, map[string]interface{}{
								"type":  "output",
								"value": val,
								"meta":  meta,
								"line":  token.Line,
							})
							advance(&index)
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
							advance(&index)
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
			case "if":
				condition := token.Val
				ifLine := token.Line

				advance(&index)
				token := current(&index, tokens)
				if token.Type == "LCURL" {
					advance(&index)
					token = current(&index, tokens)

					// Capturing loop starts
					capture := []Tokens{}
					capture = append(capture, token) // put in current token

					for {
						if token.Type == "RCURL" {
							break
						}

						advance(&index)
						token = current(&index, tokens)
						capture = append(capture, token)
					}

					logicAST := ReRunParser(capture)

					if token.Type == "RCURL" {
						meta["sub_type"] = "if"
						ast = append(ast, map[string]interface{}{
							"type":      "logic",
							"meta":      meta,
							"line":      ifLine,
							"condition": condition,
							"body":      logicAST,
						})
					} else {
						err := NewError("MalformedSyntax", token.Line, fmt.Sprintf("if %v { \n ... %s???%s", condition, Red, Reset), "This if-statement is missing a right-standing curly brace", true, "")
						err.Throw()
					}
				} else {
					err := NewError("MalformedSyntax", token.Line, fmt.Sprintf("if %v %s???%s \n ... }", condition, Red, Reset), "This if-statement is missing a left-standing curly brace", true, "")
					err.Throw()
				}
			}
		}

		advance(&index)
	}

	// Append to AST.json
	b, err := json.MarshalIndent(ast, "", "  ")
	Check(err)
	os.WriteFile("./.intext/cache/AST.json", b, 0666)
}
