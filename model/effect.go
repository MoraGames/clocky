package model

import "github.com/expr-lang/expr/vm"

type (
	Effect struct {
		ID          int64
		Name        string
		Description string
		Type        string
		Order       int
		Expression  string
		Program     *vm.Program
	}

	EffectsStack struct {
		Effect    *Effect
		Amount    int
		MaxAmount int
	}
)
