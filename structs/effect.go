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
	QuintupleNegativePoints = &Effect{"Mul -5", "Event", "*", -5}
	QuadrupleNegativePoints = &Effect{"Mul -4", "Event", "*", -4}
	TripleNegativePoints    = &Effect{"Mul -3", "Event", "*", -3}
	DoubleNegativePoints    = &Effect{"Mul -2", "Event", "*", -2}
	SingleNegativePoints    = &Effect{"Mul -1", "Event", "*", -1}
	DoublePositivePoints    = &Effect{"Mul +2", "Event", "*", 2}
	TriplePositivePoints    = &Effect{"Mul +3", "Event", "*", 3}
	QuadruplePositivePoints = &Effect{"Mul +4", "Event", "*", 4}
	QuintuplePositivePoints = &Effect{"Mul +5", "Event", "*", 5}
	SixtuplePositivePoints  = &Effect{"Mul +6", "Event", "*", 6}

	// Additive
	SubFourPoints  = &Effect{"Sub 4", "Event", "-", 4}
	SubThreePoints = &Effect{"Sub 3", "Event", "-", 3}
	SubTwoPoints   = &Effect{"Sub 2", "Event", "-", 2}
	SubOnePoint    = &Effect{"Sub 1", "Event", "-", 1}
	AddOnePoint    = &Effect{"Add 1", "Event", "+", 1}
	AddTwoPoints   = &Effect{"Add 2", "Event", "+", 2}
	AddThreePoints = &Effect{"Add 3", "Event", "+", 3}
	AddFourPoints  = &Effect{"Add 4", "Event", "+", 4}
	AddFivePoints  = &Effect{"Add 5", "Event", "+", 5}

	// Special Effects
	ComebackBonus1         = &Effect{"Comeback 1", "User", "+", 1}
	ComebackBonus2         = &Effect{"Comeback 2", "User", "+", 2}
	ComebackBonus3         = &Effect{"Comeback 3", "User", "+", 3}
	ComebackBonus4         = &Effect{"Comeback 4", "User", "+", 4}
	ComebackBonus5         = &Effect{"Comeback 5", "User", "+", 5}
	LastChanceBonus        = &Effect{"Last Chance 1", "User", "+", 2}
	LastChanceBonus2       = &Effect{"Last Chance 2", "User", "+", 5}
	ReigningLeader         = &Effect{"Reigning Leader", "User", "+", 1}
	NoNegative             = &Effect{"No Negative", "User", "+", 0}
	ConsistentParticipant1 = &Effect{"Consistent Participant 1", "User", "+", 3}
	ConsistentParticipant2 = &Effect{"Consistent Participant 2", "User", "+", 5}

	// Map of all the effects
	Effects = map[string]*Effect{
		"Mul -5":                   QuintupleNegativePoints,
		"Mul -3":                   TripleNegativePoints,
		"Mul -2":                   DoubleNegativePoints,
		"Mul -1":                   SingleNegativePoints,
		"Mul +2":                   DoublePositivePoints,
		"Mul +3":                   TriplePositivePoints,
		"Mul +5":                   QuintuplePositivePoints,
		"Mul +6":                   SixtuplePositivePoints,
		"Sub 4":                    SubFourPoints,
		"Sub 3":                    SubThreePoints,
		"Sub 2":                    SubTwoPoints,
		"Sub 1":                    SubOnePoint,
		"Add 1":                    AddOnePoint,
		"Add 2":                    AddTwoPoints,
		"Add 3":                    AddThreePoints,
		"Add 4":                    AddFourPoints,
		"Add 5":                    AddFivePoints,
		"Comeback 1":               ComebackBonus1,
		"Comeback 2":               ComebackBonus2,
		"Comeback 3":               ComebackBonus3,
		"Comeback 4":               ComebackBonus4,
		"Comeback 5":               ComebackBonus5,
		"Last Chance":              LastChanceBonus,
		"Last Chance 2":            LastChanceBonus2,
		"Reigning Leader":          ReigningLeader,
		"No Negative":              NoNegative,
		"Consistent Participant 1": ConsistentParticipant1,
		"Consistent Participant 2": ConsistentParticipant2,
	}
)
