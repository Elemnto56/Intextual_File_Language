package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Tokens struct {
	Type    string                 `json:"TYPE"`
	SubType string                 `json:"SUB-TYPE"`
	Meta    map[string]interface{} `json:"META"`
	Val     interface{}            `json:"VAL"`
	Line    int                    `json:"LINE"`
}

func advance(index *int) {
	*index++
}

func current(index *int, tokens []Tokens) Tokens {
	if *index >= len(tokens) {
		panic("Tokens ranged out in parser. Dev, fix it NOW!")
	}
	return tokens[*index]
}

func Parser(givenTokens []Tokens) []map[string]interface{} {
	// Get tokens from JSON
	bytes, _ := os.ReadFile("./.intext/cache/Tokens.json")
	var tokens []Tokens
	err := json.Unmarshal(bytes, &tokens)
	Check(err)

	// AST
	ast := []map[string]interface{}{}

	// Iterators
	index := 0
	UUID := 0 // For if statements

	if givenTokens != nil {
		tokens = givenTokens
	}

	for index < len(tokens) {
		token := current(&index, tokens)

		if token.Type == "KEYWORD" {
			switch token.Val {
			case "let":
				meta := make(map[string]interface{})
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
								if Contains([]interface{}{"INT", "BOOL", "STRING", "CHAR", "IDENTIFIER"}, token.Type) {
									value := token.Val
									_type := token.Type
									meta["raw_type"] = _type

									temp := index + 1 // Did this in order for it to be a sight into the future
									if token.Meta["assignment"] == "math" {
										meta["math"] = true
										meta["assignment-type"] = "basic"

										ast = append(ast, map[string]interface{}{
											"type":      "let",
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
										meta["assignment-type"] = "basic"
										ast = append(ast, map[string]interface{}{
											"type":      "let",
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
											meta["assignment-type"] = "basic"
											ast = append(ast, map[string]interface{}{
												"type":      "let",
												"var_type":  Type,
												"var_name":  name,
												"var_value": concatCatch,
												"line":      token.Line,
												"meta":      meta,
											})
										}
									} else {
										err := NewError("MissingBreaker", token.Line, fmt.Sprintf("%v %v: %v = %v <-", "let", name, Type, value), "This line is missing a semicolon", true, "")
										err.Throw()
									}
								} else if token.Type == "TXT BLK" {
									var captureStr string

									for i, line := range token.Val.([]interface{}) {
										if i == len(token.Val.([]interface{}))-1 {
											captureStr += fmt.Sprint(line)
										} else {
											captureStr += fmt.Sprint(line) + "\n"
										}
									}

									advance(&index)
									token = current(&index, tokens)
									if token.Type == "SYMBOL" && token.Val == ";" {
										meta["math"] = false
										meta["assignment-type"] = "basic"
										meta["raw_type"] = "TXT BLK"
										ast = append(ast, map[string]interface{}{
											"type":      "let",
											"var_type":  Type,
											"var_name":  name,
											"var_value": captureStr,
											"line":      token.Line,
											"meta":      meta,
										})
									}
									// I won't be calling a missing semicolon error, since it'd be the lexer's fault; not user's
								} else if token.Type == "ORDER" {
									tokenList := token.Val.([]interface{})
									var orderList []map[string]interface{}

									for _, element := range tokenList {
										newEle := element.(map[string]interface{})
										orderList = append(orderList, map[string]interface{}{"val": newEle["VAL"], "type": newEle["TYPE"]})
									}

									meta["raw_type"] = "ORDER"
									meta["math"] = false
									meta["assignment-type"] = "basic"
									ast = append(ast, map[string]interface{}{
										"type":      "let",
										"var_type":  Type,
										"var_name":  name,
										"var_value": orderList,
										"line":      token.Line,
										"meta":      meta,
									})
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
													meta["assignment-type"] = "basic"
													ast = append(ast, map[string]interface{}{
														"type":      "let",
														"var_type":  Type,
														"var_name":  name,
														"var_value": map[string]interface{}{"read": file},
														"line":      token.Line,
														"meta":      meta,
													})
													advance(&index)
												} else {
													err := NewError("MalformedSyntax", token.Line, fmt.Sprintf("%v %v: %v = %sread(...%s;", "let", name, Type, Red, Reset), "The following read function was not properly closed", true, "Add a paranthesis after the string that calls the file (e.g. read(...))")
													err.Throw()
												}
											} else {
												err := NewError("MalformedSyntax", token.Line, fmt.Sprintf("%v %v: %v = %sread...)%s;", "let", name, Type, Red, Reset), "The following read function was not properly closed", true, "Add a paranthesis before the string that calls the file (e.g. read(...))")
												err.Throw()
											}
										} else {
											err := NewError("TypeMismatch", token.Line, fmt.Sprintf("%v %v: %s%v%s = read(...);", "let", name, Red, Type, Reset), "The wrong type was assigned to read()", true, fmt.Sprintf("Change the type to %sstring%s", Yellow, Reset))
											err.Throw()
										}
									}
								} else if token.Type == "MATH" {
									value := token.Val
									meta["math"] = true
									meta["raw_type"] = "none"
									meta["assignment-type"] = "basic"
									temp := index + 1
									if current(&temp, tokens).Type == "SYMBOL" && current(&temp, tokens).Val == ";" {
										advance(&index)
										ast = append(ast, map[string]interface{}{
											"type":      "let",
											"var_type":  Type,
											"var_name":  name,
											"var_value": value,
											"line":      token.Line,
											"meta":      meta,
										})
									}
								} else {
									err := NewError("UnknownValue", token.Line, fmt.Sprintf("%v %v: %v = %s%v%s", "let", name, Type, Red, token.Val, Reset), "The following value could not be correctly parsed", true, "Did you forget any quotes or accidently put a variable in the statement?")
									err.Throw()
								}
							} else {
								err := NewError("MalformedSyntax", token.Line, fmt.Sprintf("%v %v:%v %s??%s", "let", name, Type, Red, Reset), "The following statement failed to abide by Intext's syntax rules", true, "The \"??\" is where an \"=\" is expected")
								err.Throw()
							}
						} else {
							err := NewError("TypeMismatch", token.Line, fmt.Sprintf("%v %v: %v <-", "let", name, token.Val), "The following line does not include a valid type", true, "")
							err.Throw()
						}
					} else if token.Val == "," {
						varSet := []interface{}{}

						varSet = append(varSet, name)
						for {
							if token.Val == "=" {
								break
							}

							if token.Val == "," {
								advance(&index)
								token = current(&index, tokens)
							}

							if token.Type == "IDENTIFIER" {
								varSet = append(varSet, token.Val)
							}

							advance(&index)
							token = current(&index, tokens)
						}

						if token.Type == "OPERATOR" && token.Val == "=" {
							advance(&index)
							token = current(&index, tokens)

							// Possibly add if statment for expansion?
							// Macro calls
							meta["assignment-type"] = "multi"
							meta["macro-args"] = token.Meta["args"]
							meta["macro-name"] = token.Meta["name"]
							meta["macro-expect"] = len(varSet)
							meta["macro-type"] = "assignment"
							ast = append(ast, map[string]interface{}{
								"type": "let",
								"line": token.Line,
								"meta": meta,
								"vars": varSet,
							})
							advance(&index)
						}

					} else {
						err := NewError("LexerErr", token.Line, fmt.Sprintf("%v %v%v <-", "let", name, token.Val), "Invalid character on this line", true, "Did you mean ':'?")
						err.Throw()
					}
				}
			case "declare":
				advance(&index)
				token := current(&index, tokens)
				meta := make(map[string]interface{})

				// Possibly add if-statement to check for var (IDENTIFIER)
				varSet := []interface{}{}

				for {
					if token.Val == "=" {
						break
					}

					if token.Val == "," {
						advance(&index)
						token = current(&index, tokens)
					}

					if token.Type == "IDENTIFIER" {
						varSet = append(varSet, token.Val)
					}

					advance(&index)
					token = current(&index, tokens)
				}

				if token.Type == "OPERATOR" && token.Val == "=" {
					advance(&index)
					token = current(&index, tokens)

					// Possibly add if statment for expansion?
					// Macro calls
					meta["macro-args"] = token.Meta["args"]
					meta["macro-name"] = token.Meta["name"]
					meta["macro-expect"] = len(varSet)
					meta["macro-type"] = "assignment"
					ast = append(ast, map[string]interface{}{
						"type": "declare",
						"line": token.Line,
						"meta": meta,
						"vars": varSet,
					})
					advance(&index)
				}
			case "output":
				advance(&index)
				token := current(&index, tokens)
				if true { // Added this here because the value is ambiguous
					meta := make(map[string]interface{})
					_type := token.Type
					val := token.Val

					temp := index + 1
					if tokens[temp].Type == "SYMBOL" && tokens[temp].Val == ";" {
						switch token.Type {
						case "STRING", "INT", "BOOL", "FLOAT", "CHAR", "IDENTIFIER":
							if cmpRegEx(fmt.Sprint(val), `^([A-Za-z]+\_*?)+\[([0-9]+|\w)\]`) {
								rawIndext := fmt.Sprint(val)

								a := strings.IndexRune(rawIndext, '[')
								b := strings.IndexRune(rawIndext, ']')

								vari := rawIndext[:a]
								num := rawIndext[a+1 : b]

								meta["print_type"] = "ord_index"
								meta["raw_type"] = "none"
								ast = append(ast, map[string]interface{}{
									"type":  "output",
									"line":  token.Line,
									"meta":  meta,
									"value": map[string]interface{}{vari: num},
								})
							} else {
								meta["raw_type"] = _type
								meta["print_type"] = "simple"
								ast = append(ast, map[string]interface{}{
									"type":  "output",
									"value": val,
									"meta":  meta,
									"line":  token.Line,
								})
							}
						case "MATH":
							meta["print_type"] = "mathematics"
							meta["raw_type"] = "none"
							ast = append(ast, map[string]interface{}{
								"type":  "output",
								"value": val,
								"meta":  meta,
								"line":  token.Line,
							})
						case "ORDER":
							tokenList := token.Val.([]interface{})
							var orderList []map[string]interface{}

							for _, element := range tokenList {
								forUse := element.(map[string]interface{})
								orderList = append(orderList, map[string]interface{}{"val": forUse["VAL"], "type": forUse["TYPE"]})
							}

							meta["print_type"] = "order"
							meta["raw_type"] = "none"
							ast = append(ast, map[string]interface{}{
								"type":  "output",
								"value": orderList,
								"meta":  meta,
								"line":  token.Line,
							})
						}
					} else if (tokens[temp].Type == "SYMBOL" || tokens[temp].Type == "COMMA") && tokens[temp].Val == "," {
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
					if token.Type == "STRING" || token.Type == "IDENTIFIER" {
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
						err := NewError("TypeMismatch", token.Line, fmt.Sprintf("write(%s%v%s, ...);", Red, token.Val, Reset), "The value for the file was not a string or variable", true, "")
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
					if Contains([]interface{}{"STRING", "IDENTIFIER", "ORDER"}, token.Type) {
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

					astBody := Parser(toks)

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

					astBody := Parser(toks)

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

						astBody := Parser(toks)

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

						astBody := Parser(toks)

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
					advance(&index)
					var exprString string

					for {
						token = current(&index, tokens)

						if token.Type == "SYMBOL" && token.Val == ";" {
							break
						}

						switch token.Type {
						case "STRING":
							exprString += fmt.Sprintf("\"%v\"", token.Val)
							advance(&index)
						case "TXT BLK":
							for i, line := range token.Val.([]interface{}) {
								if i == len(token.Val.([]interface{}))-1 {
									exprString += fmt.Sprint(line)
								} else {
									exprString += fmt.Sprint(line) + "\n"
								}
							}
							exprString = fmt.Sprintf("\"%v\"", exprString)
							advance(&index)
						default:
							exprString += fmt.Sprint(token.Val)
							advance(&index)
						}

					}

					meta["value"] = strings.TrimSpace(exprString)
					meta["target"] = exprVal
					ast = append(ast, map[string]interface{}{
						"type":     "expr",
						"sub_type": "reassign",
						"meta":     meta,
						"line":     token.Line,
					})
				}
			}
		} else if token.Type == "MACRO" {
			meta := make(map[string]interface{})
			metaJSON := token.Meta
			macLine := token.Line

			switch token.SubType {
			case "declaration":

				callStr := fmt.Sprint(metaJSON["call"])
				callList := []map[string]interface{}{}

				name := strings.TrimSpace(fmt.Sprint(metaJSON["name"]))
				returns := metaJSON["returns"].([]interface{})

				paramStr := fmt.Sprint(metaJSON["param"])
				paramList := []map[string]interface{}{}

				matchForParam := `\w+\:\s?(string|int|bool|float|char|(ord|order))`

				for _, c := range sliceOfRegex(callStr, matchForParam) {
					parts := strings.Split(c, ":")
					cName := strings.TrimSpace(parts[0])
					cType := strings.TrimSpace(parts[1])

					callList = append(callList, map[string]interface{}{cName: cType})
				}

				for _, p := range sliceOfRegex(paramStr, matchForParam) {
					parts := strings.Split(p, ":")
					pName := strings.TrimSpace(parts[0])
					pType := strings.TrimSpace(parts[1])

					paramList = append(paramList, map[string]interface{}{pName: pType})
				}

				advance(&index)
				token = current(&index, tokens)

				if token.Type == "BODY" {
					var toks []Tokens
					macBody := token.Val.([]interface{})

					b, _ := json.Marshal(macBody)
					json.Unmarshal(b, &toks)

					astBody := Parser(toks)

					meta["call"] = callList
					meta["param"] = paramList
					meta["returns"] = returns
					meta["macro-type"] = "declaration"
					ast = append(ast, map[string]interface{}{
						"meta": meta,
						"type": "macro",
						"line": macLine,
						"name": name,
						"body": astBody,
					})
				} else {
					err := NewError("MalformedSyntax", macLine, fmt.Sprintf("macro ... %v ... %s???%s", name, Red, Reset), "This macro declaration is missing either a left-standing curly brace, or right-standing curly brace", true, "")
					err.Throw()
				}
			case "call":
				args := token.Meta["args"].([]interface{})
				name := token.Meta["name"]

				if len(args) != 0 {
					var capture []interface{}

					for _, element := range args {
						forUse := element.(map[string]interface{})
						capture = append(capture, map[string]interface{}{"type": forUse["TYPE"], "val": forUse["VAL"]})
					}

					meta["macro-type"] = "standalone"
					meta["args"] = capture
					ast = append(ast, map[string]interface{}{
						"type": "macro",
						"meta": meta,
						"line": macLine,
						"name": name,
					})
					advance(&index)
				} else {
					meta["macro-type"] = "standalone"
					meta["args"] = nil
					ast = append(ast, map[string]interface{}{
						"type": "macro",
						"meta": meta,
						"line": macLine,
						"name": name,
					})
					advance(&index)
				}
			}
		}

		advance(&index)
	}

	// Append to AST.json
	b, err := json.MarshalIndent(ast, "", "  ")
	Check(err)
	os.WriteFile("./.intext/cache/AST.json", b, 0666)

	return ast
}
