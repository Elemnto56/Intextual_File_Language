package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"

	"github.com/expr-lang/expr"
)

func validateVal(varType interface{}, varValue interface{}, line int, meta string) bool {
	switch varType {
	case "int":
		_, err := strconv.Atoi(fmt.Sprint(varValue))
		if err != nil {
			if meta == "none" {
				return true
			}
			err0 := NewError("TypeMismatch", line, fmt.Sprintf("let x: int = %s%s%s;", Red, fmt.Sprint(varValue), Reset), "The following value was not an int", true, typemismatch)
			err0.Throw()
			return false
		}
		return true
	case "string":
		if Contains([]interface{}{"STRING", "FUNC", "IDENTIFIER", "concat"}, meta) {
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

func Validator() {
	// Take in AST.json
	b, _ := os.ReadFile("AST.json")
	nodes := []map[string]interface{}{}
	err := json.Unmarshal(b, &nodes)
	Check(err)

	variables := map[string]interface{}{}

	// Regexs
	pat0 := `((write|read|append|input)\((|.+)\)|declare|ouput|[\+|\-|\*|\\|\+\=|\-\=|\\\=|\*\=])`
	pat1 := `(true|false|[0-9]|^\[|\.)`
	re0 := regexp.MustCompile(pat0)
	re1 := regexp.MustCompile(pat1)

	for _, node := range nodes {
		line := node["line"].(float64)
		inte := int(line)
		Check(err)

		switch node["type"] {
		case "declare", "let":
			validName := node["var_name"].(string)
			validType := node["var_type"]
			validVal := node["var_value"]
			meta := node["meta"].(map[string]interface{})

			if re0.MatchString(validName) {
				err0 := NewError("UndefErr", inte, validName, "Change the variable name to not include Intext funcs", false, "")
				err0.Throw()
			} else {
				if meta["raw_type"] == "concat" && validType != "string" {
					eror := NewError("TypeMismatch", inte, fmt.Sprintf("let %v: %s%v%s = ...", validName, Red, validType, Reset), "The concated variable declaration was not a string", true, "All concats must be strings, unless the concats are all numbers, in which that would be math")
					eror.Throw()
				}

				work := validateVal(validType, validVal, inte, fmt.Sprint(meta["raw_type"]))
				if !work {
					errord := NewError("UnknownType", inte, fmt.Sprintf("let %v: %s%v%s = %v;", validName, Red, validType, Reset, validVal), "The follwing type could not be executed", true, "Did you happen to make a typo?")
					errord.Throw()
				}

				switch validType {
				case "int":
					if meta["raw_type"] == "INT" {
						valid, err := strconv.Atoi(fmt.Sprint(validVal))
						if err != nil {
							eror := NewError("TypeMismatch", inte, fmt.Sprintf("let %v: int = %s;", validName, validVal), "The following value was not an int", true, typemismatch)
							eror.Throw()
						}
						variables[validName] = int(valid)
					} else if meta["math"] == true {
						cmp, err := expr.Compile(fmt.Sprint(validVal))
						Check(err)
						ans, errd := expr.Run(cmp, variables)
						Check(errd)
						if f, ok := ans.(float64); ok && math.IsNaN(f) {
							errdd := NewError("DivisionByZero", inte, fmt.Sprintf("let %v: int = %s0/0%s;", validName, Red, Reset), "You tried to divide by zero...\x1b[3mA blackhole has opened nearby\x1b[0m", true, "")
							errdd.Throw()
						}
						variables[validName] = ans
					}
				case "string", "char":
					switch meta["raw_type"] {
					case "STRING", "CHAR":
						validVal := string(fmt.Sprint(validVal))
						variables[validName] = validVal
					case "FUNC":
						val := validVal.(map[string]interface{})
						if _, ok := val["read"]; ok {
							file := val["read"]
							_, err := os.Stat(fmt.Sprint(file))
							if err != nil {
								err := NewError("FileError", inte, fmt.Sprintf("let %v: string = read(%s%v%s);", validName, Red, file, Reset), "The following file does not exist", true, "Make sure that the directory is specific to it's location")
								err.Throw()
							}
							val, eror := ReadFile(fmt.Sprint(file))
							Check(eror)
							variables[validName] = val
						}
					case "IDENTIFIER":
						if re2.MatchString(fmt.Sprint(validVal)) {
							val, _ := expr.Eval(fmt.Sprint(validVal), variables)
							variables[validName] = val
						} else {
							variables[validName] = validVal
						}
					case "concat":
						variables[validName] = validVal
					}
				case "bool":
					validVal, err := strconv.ParseBool(fmt.Sprint(validVal))
					if err != nil {
						eror := NewError("TypeMismatch", inte, fmt.Sprintf("let %v: bool = %s;", validName, validVal), "The following value was not a boolean", true, typemismatch)
						eror.Throw()
					}
					variables[validName] = validVal
				case "float":
					if meta["raw_type"] == "FLOAT" {
						validVal, err := strconv.ParseFloat(fmt.Sprint(validVal), 64)
						if err != nil {
							eror := NewError("TypeMismatch", inte, fmt.Sprintf("let %v: float = %s;", validName, validVal), "The following value was not a float", true, typemismatch)
							eror.Throw()
						}
						variables[validName] = validVal
					} else if meta["math"] == true {
						continue
					}
				case "order", "ord":
					validList := validVal.([]interface{})
					ordRef := meta["ord-ref"].(map[string]interface{})

					for Key, element := range ordRef {
						if Contains(validList, Key) {
							i, _ := strconv.Atoi(fmt.Sprint(element))
							VarR := validList[i]

							if _, ok := variables[fmt.Sprint(VarR)]; !ok {
								err := NewError("VariableNotFound", inte, fmt.Sprintf("let %v: %v = %v %s--> %v%s", validName, validType, validList, Red, VarR, Reset), "This variable was not found", true, "The arrow points to which variable wasn't found")
								err.Throw()
							}
						}
					}

					variables[validName] = validList
				}
			}

		case "output":
			value := node["value"]
			metadata := node["meta"].(map[string]interface{})

			if metadata["print_type"] == "simple" {
				if metadata["raw_type"] == "STRING" || re1.MatchString(fmt.Sprint(value)) {
					continue
				} else {
					if _, exist := variables[fmt.Sprint(value)]; !exist { // Catches non-strings for vars
						err := NewError("VariableNotFound", inte, fmt.Sprintf("output %s%s%s;", Red, value.(string), Reset), "This variable was not found", true, "")
						err.Throw()
					}
				}
			} else if metadata["print_type"] == "mathematics" {
				cmp, err := expr.Compile(fmt.Sprint(value))
				Check(err)
				ans, errd := expr.Run(cmp, variables)
				Check(errd)
				if f, ok := ans.(float64); ok && math.IsNaN(f) {
					errdd := NewError("DivisionByZero", inte, fmt.Sprintf("output %s0/0%s;", Red, Reset), "You tried to divide by zero...\x1b[3mA blackhole has opened nearby\x1b[0m", true, "")
					errdd.Throw()
				}
			} else if metadata["print_type"] == "ord_index" {
				indexRef := value.(map[string]interface{})

				for key, val := range indexRef {
					if _, ok := variables[key]; !ok {
						err := NewError("VariableNotFound", inte, fmt.Sprintf("output %s%v%s[%v];", Red, key, Reset, val), "This variable was not found for an order", true, "")
						err.Throw() // Fix this
					}
				}
			}
		case "function":
			meta := node["meta"].(map[string]interface{})
			call := node["call"]
			switch call {
			case "write":
				if _, ok := variables[fmt.Sprint(meta["input"])]; !ok {
					err := NewError("VariableNotFound", inte, fmt.Sprintf("write(%v, %s%v%s);", meta["target"], Red, meta["input"], Reset), "This variable was not found", true, "")
					err.Throw()
				}
				if _, err := strconv.ParseInt(fmt.Sprint(meta["perms"]), 8, 64); err != nil {
					erra := NewError("FileError", inte, fmt.Sprintf("write(%v, %v, %s%v%s);", fmt.Sprint(meta["target"]), fmt.Sprint(meta["input"]), Red, fmt.Sprint(meta["perms"]), Reset), "Invalid permissions were given to create this file", true, "For Linux use the same number permission format that you'd do for chmod")
					erra.Throw()
				}
			case "append":
				if _, ok := variables[fmt.Sprint(meta["input"])]; !ok {
					err := NewError("VariableNotFound", inte, fmt.Sprintf("append(%v, %s%v%s);", meta["target"], Red, meta["input"], Reset), "This variable was not found", true, "")
					err.Throw()
				}
			case "del":
				if _, ok := variables[fmt.Sprint(meta["target"])]; !ok && meta["raw"] == "IDENTIFIER" {
					err := NewError("VariableNotFound", inte, fmt.Sprintf("del(%s%v%s);", Red, meta["target"], Reset), "This variable was not found", true, "")
					err.Throw()
				}
			}
		}
	}

}
