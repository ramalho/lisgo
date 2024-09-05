package main

import (
	"fmt"
)

type Symbol string
type Atom interface{}
type Expression interface{}
type Environment map[Symbol]interface{}

func main() {
	source := "((lambda (x) (* x 6)) 7)"
	fmt.Println(source, " // source")
	exp, err := parse(source)
	if err != nil {
		fmt.Println(err, " // error")
	} else {
		fmt.Println(exp, " // AST")
		fmt.Println(lispstr(exp), " // reconstructed source")
	}
}
