package main

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		source   string
		expected Expression
	}{
		{"7", 7},
		{"x", Symbol("x")},
		{"(sum 1 2 3)", []Expression{Symbol("sum"), 1, 2, 3}},
		{"(+ (* 2 100) (* 1 10))", []Expression{Symbol("+"), []Expression{Symbol("*"), 2, 100}, []Expression{Symbol("*"), 1, 10}}},
		{"99 100", 99}, // parse stops at the first complete expression
		{"(a)(b)", []Expression{Symbol("a")}},
		{"{if (< x 0) 0 x}",
			[]Expression{Symbol("if"),
				[]Expression{Symbol("<"), Symbol("x"), 0},
				0, Symbol("x")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.source, func(t *testing.T) {
			got, _ := parse(tt.source)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("parse(%q) = %#v, want %#v", tt.source, got, tt.expected)
			}
		})
	}
}

func TestLispstr(t *testing.T) {
	tests := []struct {
		obj      Expression
		expected string
	}{
		{0, "0"},
		{1, "1"},
		{false, "#f"},
		{true, "#t"},
		{1.5, "1.5"},
		{"sin", "sin"},
		{[]Expression{"+", 1, 2}, "(+ 1 2)"},
		{[]Expression{"if", []Expression{"<", "a", "b"}, true, false}, "(if (< a b) #t #f)"},
		{[]Expression{}, "()"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			got := lispstr(tt.obj)
			if got != tt.expected {
				t.Errorf("lispstr(%#v) = %q, want %q", tt.obj, got, tt.expected)
			}
		})
	}
}
