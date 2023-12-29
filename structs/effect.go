package structs

import "github.com/MoraGames/clockyuwu/pkg/types"

type Effect struct {
	Name  string
	Scope string
	Key   string
	Value int
}

type EffectPresence struct {
	Effect   *Effect
	Possible float64
	Amount   types.Interval
}

var (
	// Multiplier
	TripleNegativePoints    = &Effect{"Mul -3", "Event", "*", -3}
	DoubleNegativePoints    = &Effect{"Mul -2", "Event", "*", -2}
	SingleNegativePoints    = &Effect{"Mul -1", "Event", "*", -1}
	DoublePositivePoints    = &Effect{"Mul +2", "Event", "*", 2}
	TriplePositivePoints    = &Effect{"Mul +3", "Event", "*", 3}
	QuintuplePositivePoints = &Effect{"Mul +5", "Event", "*", 5}

	// Additive
	SubTwoPoints   = &Effect{"Sub 2", "Event", "-", 2}
	SubOnePoint    = &Effect{"Sub 1", "Event", "-", 1}
	AddOnePoint    = &Effect{"Add 1", "Event", "+", 1}
	AddTwoPoints   = &Effect{"Add 2", "Event", "+", 2}
	AddThreePoints = &Effect{"Add 3", "Event", "+", 3}

	// Special Effects
	ComebackBonus1  = &Effect{"Comeback 1", "User", "+", 1}
	ComebackBonus2  = &Effect{"Comeback 2", "User", "+", 2}
	ComebackBonus3  = &Effect{"Comeback 3", "User", "+", 3}
	LastChanceBonus = &Effect{"Last Chance", "User", "+", 2}

	// Map of all the effects
	Effects = map[string]*Effect{
		"Mul -3":      TripleNegativePoints,
		"Mul -2":      DoubleNegativePoints,
		"Mul -1":      SingleNegativePoints,
		"Mul +2":      DoublePositivePoints,
		"Mul +3":      TriplePositivePoints,
		"Mul +5":      QuintuplePositivePoints,
		"Sub 2":       SubTwoPoints,
		"Sub 1":       SubOnePoint,
		"Add 1":       AddOnePoint,
		"Add 2":       AddTwoPoints,
		"Add 3":       AddThreePoints,
		"Comeback 1":  ComebackBonus1,
		"Comeback 2":  ComebackBonus2,
		"Comeback 3":  ComebackBonus3,
		"Last Chance": LastChanceBonus,
	}
)
