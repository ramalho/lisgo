package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Symbol string
type Atom interface{}
type Expression interface{}
type Environment map[Symbol]interface{}

func parse(program string) (Expression, error) {
	tokens := tokenize(program)
	return readFromTokens(&tokens)
}

func tokenize(s string) []string {
	return strings.Fields(strings.ReplaceAll(strings.ReplaceAll(s, "(", " ( "), ")", " ) "))
}

func readFromTokens(tokens *[]string) (Expression, error) {
	if len(*tokens) == 0 {
		return nil, errors.New("unexpected EOF while reading")
	}
	token := (*tokens)[0]
	*tokens = (*tokens)[1:]
	if token == "(" {
		var L []Expression
		for (*tokens)[0] != ")" {
			exp, err := readFromTokens(tokens)
			if err != nil {
				return nil, err
			}
			L = append(L, exp)
		}
		*tokens = (*tokens)[1:] // pop off ')'
		return L, nil
	} else if token == ")" {
		return nil, errors.New("unexpected )")
	} else {
		return atom(token), nil
	}
}

func atom(token string) Atom {
	if i, err := strconv.Atoi(token); err == nil {
		return i
	} else if f, err := strconv.ParseFloat(token, 64); err == nil {
		return f
	} else {
		return Symbol(token)
	}
}

func main() {
	program := "((lambda (x) (* x 6)) 7)"
	exp, err := parse(program)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Parsed Expression:", exp)
	}
}
