package events

import (
	"math"

	"github.com/MoraGames/clockyuwu/pkg/types"
)

type SetSlice []*Set
type SetJsonSlice []*SetJson
type FuncMap map[string]func(h1, h2, m1, m2 int) bool
type Set struct {
	Name     string
	Pattern  string
	Typology string
	Enabled  bool
	Verify   func(h1, h2, m1, m2 int) bool
}
type SetJson struct {
	Name     string
	Pattern  string
	Typology string
	Enabled  bool
}

var (
	SetsFunctions = FuncMap{
		//"Equal":            equal,
		"Short Equal":      shortEqual,
		"Repeat":           repeat,
		"Mirror":           mirror,
		"Rise":             rise,
		"Short Rise":       shortRise,
		"Short Fall":       shortFall,
		"Rapid Rise":       rapidRise,
		"Short Rapid Rise": shortRapidRise,
		"Short Rapid Fall": shortRapidFall,
		"Double":           double,
		"Short Triple":     shortTriple,
		"Perfect Square":   perfectSquare,
		"Equal Twins":      equalTwins,
		"Half":             half,
	}
	Sets = SetSlice{
		//{"Equal", "aa:aa", "static", false, equal},
		{"Short Equal", "?a:aa", "static", false, shortEqual},
		{"Repeat", "ab:ab", "static", false, repeat},
		{"Mirror", "ab:ba", "static", false, mirror},
		{"Rise", "ab:cd", "static", false, rise},
		{"Short Rise", "?a:bc", "static", false, shortRise},
		{"Short Fall", "?c:ba", "static", false, shortFall},
		{"Rapid Rise", "ac:eg", "static", false, rapidRise},
		{"Short Rapid Rise", "?a:ce", "static", false, shortRapidRise},
		{"shortRapidFall", "?e:ca", "static", false, shortRapidFall},
		{"double", "n:2*n", "static", false, double},
		{"shortTriple", "[unnamed]", "static", false, shortTriple},
		{"Perfect Square", "[unnamed]", "static", false, perfectSquare},
		{"Equal Twins", "aa:bb", "static", false, equalTwins},
		{"Half", "2*n:n", "static", false, half},
	}
	SetsJson = SetJsonSlice{}

	AssignSetsFromSetsJson = func(utils types.Utils) {
		Sets = SetsJson.ToSlice()
	}
	AssignSetsWithDefault = func(utils types.Utils) {
		Sets = SetSlice{
			//{"Equal", "aa:aa", "static", false, equal},
			{"Short Equal", "?a:aa", "static", false, shortEqual},
			{"Repeat", "ab:ab", "static", false, repeat},
			{"Mirror", "ab:ba", "static", false, mirror},
			{"Rise", "ab:cd", "static", false, rise},
			{"Short Rise", "?a:bc", "static", false, shortRise},
			{"Short Fall", "?c:ba", "static", false, shortFall},
			{"Rapid Rise", "ac:eg", "static", false, rapidRise},
			{"Short Rapid Rise", "?a:ce", "static", false, shortRapidRise},
			{"Short Rapid Fall", "?e:ca", "static", false, shortRapidFall},
			{"Double", "n:2*n", "static", false, double},
			{"Short Triple", "[unnamed]", "static", false, shortTriple},
			{"Perfect Square", "[unnamed]", "static", false, perfectSquare},
			{"Equal Twins", "aa:bb", "static", false, equalTwins},
			{"Half", "2*n:n", "static", false, half},
		}
	}
)

func (s SetSlice) ToJsonSlice() SetJsonSlice {
	jsonSlice := make(SetJsonSlice, 0)
	for _, set := range s {
		jsonSlice = append(jsonSlice, &SetJson{
			Name:     set.Name,
			Pattern:  set.Pattern,
			Typology: set.Typology,
			Enabled:  set.Enabled,
		})
	}
	return jsonSlice
}

func (sj SetJsonSlice) ToSlice() SetSlice {
	slice := make(SetSlice, 0)
	for _, setjson := range sj {
		slice = append(slice, &Set{
			Name:     setjson.Name,
			Pattern:  setjson.Pattern,
			Typology: setjson.Typology,
			Enabled:  setjson.Enabled,
			Verify:   SetsFunctions[setjson.Name],
		})
	}
	return slice
}

// Notes: The set is replaced by edits on repeat and mirror sets and the new "equal twins" set
// aa:aa
// func equal(h1, h2, m1, m2 int) bool {
// 	return h1 == h2 && h2 == m1 && m1 == m2
// }

// ?a:aa
func shortEqual(_, h2, m1, m2 int) bool {
	return h2 == m1 && m1 == m2
}

// ab:ab
func repeat(h1, h2, m1, m2 int) bool {
	return h1 == m1 && h2 == m2
}

// ab:ba
func mirror(h1, h2, m1, m2 int) bool {
	return h1 == m2 && h2 == m1
}

// ab:cd
func rise(h1, h2, m1, m2 int) bool {
	return h2 == h1+1 && m1 == h2+1 && m2 == m1+1
}

// ?a:bc
func shortRise(_, h2, m1, m2 int) bool {
	return m1 == h2+1 && m2 == m1+1
}

// ?c:ba
func shortFall(_, h2, m1, m2 int) bool {
	return m1 == m2+1 && h2 == m1+1
}

// ac:eg
func rapidRise(h1, h2, m1, m2 int) bool {
	return h2 == h1+2 && m1 == h2+2 && m2 == m1+2
}

// ?a:ce
func shortRapidRise(_, h2, m1, m2 int) bool {
	return m1 == h2+2 && m2 == m1+2
}

// ?e:ca
func shortRapidFall(_, h2, m1, m2 int) bool {
	return m1 == m2+2 && h2 == m1+2
}

// n:2*n
func double(h1, h2, m1, m2 int) bool {
	return 2*((h1*10)+h2) == (m1*10)+m2
}

// [unnamed]
func shortTriple(_, h2, m1, m2 int) bool {
	return 3*h2 == (m1*10)+m2
}

// [unnamed]
func perfectSquare(h1, h2, m1, m2 int) bool {
	total := (h1 * 1000) + (h2 * 100) + (m1 * 10) + (m2)
	sqrt := int(math.Sqrt(float64(total)))
	return (sqrt * sqrt) == total
}

// aa:bb
func equalTwins(h1, h2, m1, m2 int) bool {
	return h1 == h2 && m1 == m2
}

// 2*n:n
func half(h1, h2, m1, m2 int) bool {
	return (h1*10 + h2) == (m1*10+m2)*2
}
