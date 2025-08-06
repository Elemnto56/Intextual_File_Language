package main

import "fmt"

func ReRunParser(tokens []Tokens) []map[string]interface{} {
	i := 0
	tsa := []map[string]interface{}{}

	for i < len(tokens) {
		token := current(&i, tokens)

		if token.Type == "KEYWORD" {
			switch token.Val {
			case "let", "declare":
				meta := make(map[string]interface{})
				grandType := token.Val
				advance(&i)
				token := current(&i, tokens)
				if token.Type == "IDENTIFIER" {
					name := token.Val
					advance(&i)
					token := current(&i, tokens)
					if token.Type == "SYMBOL" && token.Val == ":" {
						advance(&i)
						token := current(&i, tokens)
						if token.Type == "TYPESYS" && Contains([]interface{}{"bool", "string", "int", "char", "float", "ord", "order"}, fmt.Sprint(token.Val)) {
							Type := token.Val
							advance(&i)
							token := current(&i, tokens)
							if token.Type == "OPERATOR" && token.Val == "=" {
								advance(&i)
								token := current(&i, tokens)
								if Contains([]interface{}{"INT", "BOOL", "STRING", "CHAR", "ORD", "IDENTIFIER"}, token.Type) {
									value := token.Val
									_type := token.Type
									meta["raw_type"] = _type
									meta["math"] = false

									temp := i + 1 // Did this in order for it to be a sight into the future
									if current(&temp, tokens).Type == "SYMBOL" && current(&temp, tokens).Val == ";" {
										advance(&i)
										tsa = append(tsa, map[string]interface{}{
											"type":      grandType,
											"var_type":  Type,
											"var_name":  name,
											"var_value": value,
											"line":      token.Line,
											"meta":      meta,
										})
									} else if (current(&temp, tokens).Type == "OPERATOR" && current(&temp, tokens).Val == "+") || (current(&temp, tokens).Type == "SYMBOL" && current(&temp, tokens).Val == ",") {
										first := current(&i, tokens).Val
										advance(&i)
										token := current(&i, tokens)
										concatCatch := []interface{}{}

										concatCatch = append(concatCatch, fmt.Sprint(first))

										for {
											if token.Type == "SYMBOL" && token.Val == ";" {
												break
											}

											if (token.Type == "SYMBOL" && token.Val == ",") || (token.Type == "OPERATOR" && token.Val == "+") {
												advance(&i)
												token = current(&i, tokens)
												continue
											}

											concatCatch = append(concatCatch, fmt.Sprint(token.Val))
											advance(&i)
											token = current(&i, tokens)
										}

										if token.Type == "SYMBOL" && token.Val == ";" {
											meta["raw_type"] = "concat"
											tsa = append(tsa, map[string]interface{}{
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
									advance(&i)
									token := current(&i, tokens)
									userList := []interface{}{} // Make list for order

									userList = append(userList, token.Val) // Put in current token
									VarRef := make(map[string]interface{})

									var tempi int = 0
									advance(&i)
									for {
										token = current(&i, tokens)

										if Contains([]interface{}{"STRING", "INT", "BOOL", "FLOAT", "CHAR"}, token.Type) {
											userList = append(userList, token.Val)
											tempi += 1
											advance(&i)
											continue
										}

										if token.Type == "IDENTIFIER" {
											userList = append(userList, token.Val)
											tempi += 1
											VarRef[fmt.Sprint(token.Val)] = tempi
											advance(&i)
											continue
										}

										if token.Type == "COMMA" {
											advance(&i)
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
										tsa = append(tsa, map[string]interface{}{
											"type":      grandType,
											"var_type":  Type,
											"var_name":  name,
											"var_value": userList,
											"line":      token.Line,
											"meta":      meta,
										})
										advance(&i)
									}
								} else if token.Type == "FUNC" {
									switch token.Val {
									case "read":
										if Type == "string" {
											advance(&i)
											token := current(&i, tokens)
											if token.Type == "PARA" {
												advance(&i)
												token = current(&i, tokens)
												file := token.Val // Grab the file trying to read
												advance(&i)
												token = current(&i, tokens)
												if token.Type == "PARA" {
													meta["raw_type"] = "FUNC"
													meta["math"] = false
													tsa = append(tsa, map[string]interface{}{
														"type":      grandType,
														"var_type":  Type,
														"var_name":  name,
														"var_value": map[string]interface{}{"read": file},
														"line":      token.Line,
														"meta":      meta,
													})
													advance(&i)
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

									temp := i + 1
									if current(&temp, tokens).Type == "SYMBOL" && current(&temp, tokens).Val == ";" {
										advance(&i)
										tsa = append(tsa, map[string]interface{}{
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
				advance(&i)
				token := current(&i, tokens)
				if true { // Added this here because the value is ambiguous
					meta := make(map[string]string)
					_type := token.Type
					val := token.Val

					temp := i + 1
					if current(&temp, tokens).Type == "SYMBOL" && current(&temp, tokens).Val == ";" {
						switch token.Type {
						case "STRING", "INT", "BOOL", "FLOAT", "ORD", "CHAR", "IDENTIFIER":
							meta["raw_type"] = _type
							meta["print_type"] = "simple"
							tsa = append(tsa, map[string]interface{}{
								"type":  "output",
								"value": val,
								"meta":  meta,
								"line":  token.Line,
							})
							advance(&i)
						case "MATH":
							meta["print_type"] = "mathematics"
							meta["raw_type"] = "none"
							tsa = append(tsa, map[string]interface{}{
								"type":  "output",
								"value": val,
								"meta":  meta,
								"line":  token.Line,
							})
							advance(&i)
						}
					} else if (current(&temp, tokens).Type == "SYMBOL" || current(&temp, tokens).Type == "COMMA") && current(&temp, tokens).Val == "," {
						spagList := []interface{}{}
						spagList = append(spagList, val) // Add the first val into the list; none left behind!
						advance(&i)

						for {
							newVal := current(&i, tokens)
							var i int = i + 1

							if (newVal.Type == "SYMBOL" || newVal.Type == "COMMA") && newVal.Val == "," {
								advance(&i)
								continue
							}

							spagList = append(spagList, interface{}(newVal.Val))
							advance(&i)

							if current(&i, tokens).Type == "SYMBOL" && current(&i, tokens).Val == ";" {
								break
							}
						}
						token := current(&i, tokens)
						if (token.Type == "SYMBOL" || token.Type == "COMMA") && token.Val == ";" {
							meta["print_type"] = "mixed"
							meta["raw_type"] = "none"
							tsa = append(tsa, map[string]interface{}{
								"type":  "output",
								"value": spagList,
								"meta":  meta,
								"line":  token.Line,
							})
							advance(&i)
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
				advance(&i)
				token := current(&i, tokens)
				if token.Type == "PARA" && token.Val == "(" {
					advance(&i)
					token = current(&i, tokens)
					if token.Type == "STRING" {
						targetFile := token.Val
						advance(&i)
						token = current(&i, tokens)
						if token.Type == "SYMBOL" && token.Val == "," {
							advance(&i)
							token = current(&i, tokens)
							if token.Type == "IDENTIFIER" {
								wrVar := token.Val
								advance(&i)
								token = current(&i, tokens)
								if token.Type == "PARA" && token.Val == ")" {
									advance(&i)
									meta["target"] = targetFile
									meta["input"] = wrVar
									meta["perms"] = 666
									tsa = append(tsa, map[string]interface{}{
										"type": "function",
										"line": token.Line,
										"call": "write",
										"meta": meta,
									})
								} else if token.Type == "SYMBOL" && token.Val == "," {
									advance(&i)
									token = current(&i, tokens)
									if token.Type == "INT" {
										perms := token.Val
										advance(&i)
										token = current(&i, tokens)
										if token.Type == "PARA" && token.Val == ")" {
											meta["target"] = targetFile
											meta["input"] = wrVar
											meta["perms"] = perms
											tsa = append(tsa, map[string]interface{}{
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
				advance(&i)
				token := current(&i, tokens)
				if token.Type == "PARA" && token.Val == "(" {
					advance(&i)
					token = current(&i, tokens)
					if token.Type == "STRING" {
						target := token.Val
						advance(&i)
						token = current(&i, tokens)
						if token.Type == "SYMBOL" && token.Val == "," {
							advance(&i)
							token = current(&i, tokens)
							if token.Type == "IDENTIFIER" {
								input := token.Val
								advance(&i)
								token = current(&i, tokens)
								if token.Type == "PARA" && token.Val == ")" {
									advance(&i)
									meta["perms"] = 0
									meta["input"] = input
									meta["target"] = target
									tsa = append(tsa, map[string]interface{}{
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
				advance(&i)
				token := current(&i, tokens)
				if token.Type == "PARA" && token.Val == "(" {
					advance(&i)
					token = current(&i, tokens)
					if token.Type == "STRING" || token.Type == "IDENTIFIER" {
						meta["raw"] = token.Type
						targetFile := token.Val
						advance(&i)
						token = current(&i, tokens)
						if token.Type == "PARA" && token.Val == ")" {
							meta["target"] = targetFile
							advance(&i)
							tsa = append(tsa, map[string]interface{}{
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
		}
	}

	return tsa
}
