package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

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

func sliceOfRegex(find string, regex string) []string {
	temp := regexp.MustCompile(regex)

	return temp.FindAllString(find, -1)
}

type WebFile struct {
	Type     string
	Content  string
	FileName string
}

// For the Playground
func WebAddFile(filename string, content string, Type string) WebFile {
	return WebFile{FileName: filename, Content: content, Type: Type}
}
