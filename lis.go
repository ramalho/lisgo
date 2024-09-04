package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Expression interface{}
type Atom interface{}

type Symbol string

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

func main() {
	// Example usage
	program := `
	(define double
    	(lambda (n)
      	(* n 2)))`
	exp, err := Parse(program)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Parsed expression:", exp)
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
