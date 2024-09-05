package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func parse(source string) (Expression, error) {
	tokens := tokenize(source)
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

// lispstr converts a Go object back into a Lisp-readable string.
func lispstr(obj Expression) string {
	switch obj := obj.(type) {
	case bool:
		if obj {
			return "#t"
		}
		return "#f"
	case []Expression:
		items := make([]string, len(obj))
		for i, x := range obj {
			items[i] = lispstr(x)
		}
		return "(" + strings.Join(items, " ") + ")"
	case Symbol:
		return string(obj)
	default:
		return fmt.Sprintf("%v", obj)
	}
}
