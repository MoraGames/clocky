package model

import "github.com/expr-lang/expr/vm"

type Set struct {
	ID          int64
	Name        string
	Description string
	Type        string
	Expression  string
	Program     *vm.Program
	Events      []*Event
}
