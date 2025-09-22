package main

import (
	"fmt"
	"strconv"
)

// Simple token structure
type EvalToken struct {
	Type  string
	Value string
}

// Tokenize breaks the expression into tokens
func Tokenize(expr string) []EvalToken {
	var tokens []EvalToken
	i := 0

	for i < len(expr) {
		// Skip whitespace
		if expr[i] == ' ' {
			i++
			continue
		}

		// Numbers
		if expr[i] >= '0' && expr[i] <= '9' {
			start := i
			for i < len(expr) && (expr[i] >= '0' && expr[i] <= '9' || expr[i] == '.') {
				i++
			}
			tokens = append(tokens, EvalToken{"NUMBER", expr[start:i]})
			continue
		}

		// Variables (letters)
		if expr[i] >= 'a' && expr[i] <= 'z' || expr[i] >= 'A' && expr[i] <= 'Z' {
			start := i
			for i < len(expr) && (expr[i] >= 'a' && expr[i] <= 'z' || expr[i] >= 'A' && expr[i] <= 'Z') {
				i++
			}
			word := expr[start:i]

			if word == "true" || word == "false" {
				tokens = append(tokens, EvalToken{"BOOL", word})
			} else {
				if i+1 < len(expr) && expr[i] == '(' {
					tokens = append(tokens, EvalToken{"FUNC", word})
				} else {
					tokens = append(tokens, EvalToken{"VAR", word})
				}
			}
			continue
		}

		// Two-character operators
		if i+1 < len(expr) {
			twoChar := expr[i : i+2]
			switch twoChar {
			case "==", "!=", "<=", ">=", "&&", "||":
				tokens = append(tokens, EvalToken{"OP", twoChar})
				i += 2
				continue
			}
		}

		// Single-character operators
		switch expr[i] {
		case '+', '-', '*', '/', '<', '>', '!':
			tokens = append(tokens, EvalToken{"OP", string(expr[i])})
		case '(', ')':
			tokens = append(tokens, EvalToken{"PAREN", string(expr[i])})
		case '"':
			i++
			var capture string

			for expr[i] != '"' {
				capture += string(expr[i])
				i++
			}

			tokens = append(tokens, EvalToken{"STRING", capture})
		default:
			panic(fmt.Sprintf("Unknown character: %c", expr[i]))
		}
		i++
	}

	return tokens
}

// Simple evaluation using recursive descent
type Evaluator struct {
	tokens []EvalToken
	pos    int
	vars   map[string]interface{}
}

// Evaluate your expressions
func EvaluateExpression(expression string, variables map[string]VarManager) (interface{}, error) {
	tokens := Tokenize(expression)
	fixedVariables := make(map[string]interface{})

	if variables == nil {
		fixedVariables = make(map[string]interface{})
	}

	for key, val := range variables {
		fixedVariables[key] = val.Value
	}

	if len(tokens) == 0 {
		return nil, fmt.Errorf("empty expression")
	}

	e := &Evaluator{
		tokens: tokens,
		pos:    0,
		vars:   fixedVariables,
	}

	result := e.parseExpression()

	return result, nil
}

// Parse the full expression (handles OR)
func (e *Evaluator) parseExpression() interface{} {
	left := e.parseAnd()

	for e.pos < len(e.tokens) && e.tokens[e.pos].Value == "||" {
		e.pos++ // skip operator
		right := e.parseAnd()
		left = toBool(left) || toBool(right)
	}

	return left
}

// Parse AND expressions
func (e *Evaluator) parseAnd() interface{} {
	left := e.parseComparison()

	for e.pos < len(e.tokens) && e.tokens[e.pos].Value == "&&" {
		e.pos++ // skip operator
		right := e.parseComparison()
		left = toBool(left) && toBool(right)
	}

	return left
}

// Parse comparison operators
func (e *Evaluator) parseComparison() interface{} {
	left := e.parseAddSub()

	if e.pos < len(e.tokens) && e.tokens[e.pos].Type == "OP" {
		op := e.tokens[e.pos].Value
		switch op {
		case "==", "!=", "<", "<=", ">", ">=":
			e.pos++ // skip operator
			right := e.parseAddSub()
			return compare(left, op, right)
		}
	}

	return left
}

// Parse addition and subtraction
func (e *Evaluator) parseAddSub() interface{} {
	left := e.parseMulDiv()

	for e.pos < len(e.tokens) {
		if e.tokens[e.pos].Type != "OP" {
			break
		}
		op := e.tokens[e.pos].Value
		if op != "+" && op != "-" {
			break
		}
		e.pos++ // skip operator
		right := e.parseMulDiv()

		if op == "+" {
			if fmt.Sprintf("%T", left) == "string" && fmt.Sprintf("%T", right) == "string" {
				left = fmt.Sprint(left) + fmt.Sprint(right)
			} else {
				left = toFloat(left) + toFloat(right)
			}
		} else {
			left = toFloat(left) - toFloat(right)
		}
	}

	return left
}

// Parse multiplication and division
func (e *Evaluator) parseMulDiv() interface{} {
	left := e.parsePrimary()
	for e.pos < len(e.tokens) {
		if e.tokens[e.pos].Type != "OP" {
			break
		}
		op := e.tokens[e.pos].Value
		if op != "*" && op != "/" {
			break
		}
		e.pos++ // skip operator
		right := e.parsePrimary()

		if op == "*" {
			left = toFloat(left) * toFloat(right)
		} else {
			left = toFloat(left) / toFloat(right)
		}
	}

	return left
}

// Parse primary values (numbers, variables, parentheses)
func (e *Evaluator) parsePrimary() interface{} {
	if e.pos >= len(e.tokens) {
		panic("unexpected end of expression")
	}

	token := e.tokens[e.pos]

	if token.Type == "STRING" {
		e.pos++
		return fmt.Sprint(token.Value)
	}

	// Handle NOT operator
	if token.Type == "OP" && token.Value == "!" {
		e.pos++
		return !toBool(e.parsePrimary())
	}

	// Handle negative numbers
	if token.Type == "OP" && token.Value == "-" {
		e.pos++
		return -toFloat(e.parsePrimary())
	}

	// Handle parentheses
	if token.Type == "PAREN" && token.Value == "(" {
		e.pos++ // skip (
		result := e.parseExpression()
		if e.pos >= len(e.tokens) || e.tokens[e.pos].Value != ")" {
			panic("missing closing parenthesis")
		}
		e.pos++ // skip )
		return result
	}

	// Handle numbers
	if token.Type == "NUMBER" {
		e.pos++
		val, _ := strconv.ParseFloat(token.Value, 64)
		return val
	}

	// Handle booleans
	if token.Type == "BOOL" {
		e.pos++
		return token.Value == "true"
	}

	// Handle variables
	if token.Type == "VAR" {
		e.pos++
		if val, ok := e.vars[token.Value]; ok {
			return val
		}
		panic(fmt.Sprintf("undefined variable: %s", token.Value))
	}

	if token.Type == "FUNC" {
		funcName := token.Value
		e.pos++
		e.pos++

		if funcName == "read" {
			filename := e.tokens[e.pos].Value

			conts, _ := ReadFile(filename)
			e.pos++
			if e.pos+1 < len(e.tokens) {
				e.pos++
			}
			return fmt.Sprint(conts)
		}
	}

	panic(fmt.Sprintf("unexpected token: %s", token.Value))
}

// Helper function to compare values
func compare(left interface{}, op string, right interface{}) bool {
	switch op {
	case "==":
		// Handle string compare
		if lftstr, ok := left.(string); ok {
			if rtstr, ok := right.(string); ok {
				return lftstr == rtstr
			}
		}
		// Handle boolean comparison
		var ok bool
		lb, err := strconv.ParseBool(fmt.Sprint(left))
		if lb, ok = left.(bool); ok || err == nil {
			rb, err := strconv.ParseBool(fmt.Sprint(right))
			if rb, ok = right.(bool); ok || err == nil {
				return lb == rb
			}
		}
		fmt.Printf("%T %T\n", left, right)
		return toFloat(left) == toFloat(right)
	case "!=":
		// Handle string compare
		if lftstr, ok := left.(string); ok {
			if rtstr, ok := right.(string); ok {
				return lftstr != rtstr
			}
		}
		// Handle boolean comparison
		if lb, ok := left.(bool); ok {
			if rb, ok := right.(bool); ok {
				return lb != rb
			}
		}
		return toFloat(left) != toFloat(right)
	case "<":
		return toFloat(left) < toFloat(right)
	case "<=":
		return toFloat(left) <= toFloat(right)
	case ">":
		return toFloat(left) > toFloat(right)
	case ">=":
		return toFloat(left) >= toFloat(right)
	}
	return false
}

// Convert value to float
func toFloat(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case int:
		return float64(val)
	case bool:
		if val {
			return 1.0
		}
		return 0.0
	default:
		panic(fmt.Sprintf("cannot convert %v to float", v))
	}
}

// Convert value to bool
func toBool(v interface{}) bool {
	switch val := v.(type) {
	case bool:
		return val
	case float64:
		return val != 0
	case int:
		return val != 0
	default:
		return false
	}
}

// (Dev Note) How to add new features:
//
// 1. To add a new operator (like %):
//    - Add it to the Tokenize function
//    - Add handling in the appropriate parse function
//
// 2. To add strings:
//    - Add string detection in Tokenize (look for quotes)
//    - Add string comparison in the compare function
//
// 3. To add functions like max(a,b):
//    - Detect function names in Tokenize
//    - Add a parseFunction method
//
// 4. To add arrays:
//    - Detect [ ] in Tokenize
//    - Add array parsing logic
