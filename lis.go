package main

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

type Symbol string
type Expression interface{}
type Atom interface{}

type Environment map[Symbol]interface{}

// Parse reads a Scheme expression from a string
func Parse(program string) (Expression, error) {
	tokens := Tokenize(program)
	return ReadFromTokens(&tokens)
}

// Tokenize converts a string into a list of tokens
func Tokenize(s string) []string {
	s = strings.ReplaceAll(s, "(", " ( ")
	s = strings.ReplaceAll(s, ")", " ) ")
	return strings.Fields(s)
}

// ReadFromTokens reads an expression from a sequence of tokens
func ReadFromTokens(tokens *[]string) (Expression, error) {
	if len(*tokens) == 0 {
		return nil, errors.New("unexpected EOF while reading")
	}

	token := (*tokens)[0]
	*tokens = (*tokens)[1:]

	if token == "(" {
		var exp []Expression
		for (*tokens)[0] != ")" {
			subExp, err := ReadFromTokens(tokens)
			if err != nil {
				return nil, err
			}
			exp = append(exp, subExp)
		}
		*tokens = (*tokens)[1:] // discard ')'
		return exp, nil
	} else if token == ")" {
		return nil, errors.New("unexpected )")
	} else {
		return ParseAtom(token), nil
	}
}

// ParseAtom converts a token into an Atom: int, float or symbol
func ParseAtom(token string) Atom {
	if i, err := strconv.Atoi(token); err == nil {
		return i
	} else if f, err := strconv.ParseFloat(token, 64); err == nil {
		return f
	} else {
		return Symbol(token)
	}
}

/*
Syntax tree for:

	(define double
    	(lambda (n)
      	(* n 2)))


							  '*'  'n'   2
                        'n'    └────┼────┘
                         │          │
           'lambda'     [ ]        [ ]
               └─────────┼──────────┘
                         │
'define'   'double'     [ ]
    └─────────┼──────────┘
              │
             [ ]
*/

func standardEnv() Environment {
	env := Environment{
		"+":        func(a, b float64) float64 { return a + b },
		"-":        func(a, b float64) float64 { return a - b },
		"*":        func(a, b float64) float64 { return a * b },
		"/":        func(a, b float64) float64 { return a / b },
		"quotient": func(a, b float64) float64 { return math.Floor(a / b) },
		">":        func(a, b float64) bool { return a > b },
		"<":        func(a, b float64) bool { return a < b },
		">=":       func(a, b float64) bool { return a >= b },
		"<=":       func(a, b float64) bool { return a <= b },
		"=":        func(a, b float64) bool { return a == b },
		"abs":      math.Abs,
		"begin":    func(args ...interface{}) interface{} { return args[len(args)-1] },
		"eq?":      func(a, b interface{}) bool { return reflect.DeepEqual(a, b) },
		"equal?":   func(a, b interface{}) bool { return reflect.DeepEqual(a, b) },
		"max":      math.Max,
		"min":      math.Min,
		"not":      func(a bool) bool { return !a },
		"number?":  func(a interface{}) bool { _, ok := a.(float64); return ok },
		"procedure?": func(a interface{}) bool {
			return reflect.TypeOf(a).Kind() == reflect.Func
		},
		"modulo": func(a, b float64) float64 { return math.Mod(a, b) },
		"round":  math.Round,
		"symbol?": func(a interface{}) bool {
			_, ok := a.(Symbol)
			return ok
		},
	}
	return env
}

func evaluate(exp Expression, env Environment) (interface{}, error) {
	switch exp := exp.(type) {
	case Symbol:
		// Variable reference
		if val, ok := env[exp]; ok {
			return val, nil
		}
		return nil, fmt.Errorf("undefined symbol: %s", exp)
	case int, float64:
		// Number literal
		return exp, nil
	case []interface{}:
		if len(exp) == 0 {
			return nil, errors.New("empty expression")
		}
		switch exp[0] {
		case "define":
			// (define var exp)
			if len(exp) != 3 {
				return nil, errors.New("invalid define syntax")
			}
			varName, ok := exp[1].(Symbol)
			if !ok {
				return nil, errors.New("invalid variable name")
			}
			value, err := evaluate(exp[2], env)
			if err != nil {
				return nil, err
			}
			env[varName] = value
			return nil, nil
		default:
			// (proc arg...)
			procVal, err := evaluate(exp[0], env)
			if err != nil {
				return nil, errors.New("invalid procVal")
			}
			proc, ok := procVal.(func(...interface{}) interface{})
			if !ok {
				return nil, errors.New("invalid proc")
			}
			args := make([]interface{}, len(exp)-1)
			for i, arg := range exp[1:] {
				argVal, err := evaluate(arg, env)
				if err != nil {
					return nil, err
				}
				args[i] = argVal
			}
			return proc(args...), nil
		}
	default:
		return nil, fmt.Errorf("unexpected expression type: %T", exp)
	}
}

// func main() {
// 	env := standardEnv()
// 	// Example usage
// 	exp := []interface{}{Symbol("define"), Symbol("x"), 42}
// 	_, err := evaluate(exp, env)
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 	} else {
// 		fmt.Println("x:", env["x"])
// 	}
// }

func main() {
	// Example usage
	env := standardEnv()
	program := `
	((lambda (n) (* n 2)) 7)`
	exp, err := Parse(program)
	if err != nil {
		fmt.Println("Parsing error:", err)
	} else {
		fmt.Println("Parsed expression:", exp)
		res, err := evaluate(exp, env)
		if err != nil {
			fmt.Println("Evaluation error:", err)
		} else {
			fmt.Print("Result: ", res)
		}
	}
}
