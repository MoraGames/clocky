package events

import "github.com/MoraGames/clockyuwu/pkg/types"

type SetSlice []*Set
type SetJsonSlice []*SetJson
type FuncMap map[string]func(h1, h2, m1, m2 int) bool
type Set struct {
	Name     string
	Typology string
	Enabled  bool
	Verify   func(h1, h2, m1, m2 int) bool
}
type SetJson struct {
	Name     string
	Typology string
	Enabled  bool
}

var (
	SetsFunctions = FuncMap{
		"aa:aa":  aaaa,
		"xa:aa":  xaaa,
		"ab:ab":  abab,
		"ab:ba":  abba,
		"ab:cd":  abcd,
		"xa:bc":  xabc,
		"xc:ba":  xcba,
		"ac:eg":  aceg,
		"xa:ce":  xace,
		"xe:ca":  xeca,
		"n:2*n":  n2n,
		"xn:3*n": xn3n,
	}
	Sets = SetSlice{
		{"aa:aa", "standard", false, aaaa},
		{"xa:aa", "standard", false, xaaa},
		{"ab:ab", "standard", false, abab},
		{"ab:ba", "standard", false, abba},
		{"ab:cd", "standard", false, abcd},
		{"xa:bc", "standard", false, xabc},
		{"xc:ba", "standard", false, xcba},
		{"ac:eg", "standard", false, aceg},
		{"xa:ce", "standard", false, xace},
		{"xe:ca", "standard", false, xeca},
		{"n:2*n", "standard", false, n2n},
		{"xn:3*n", "standard", false, xn3n},
	}
	SetsJson = SetJsonSlice{}

	AssignSetsFromSetsJson = func(utils types.Utils) {
		Sets = SetsJson.ToSlice()
	}
	AssignSetsWithDefault = func(utils types.Utils) {
		Sets = SetSlice{
			{"aa:aa", "standard", false, aaaa},
			{"xa:aa", "standard", false, xaaa},
			{"ab:ab", "standard", false, abab},
			{"ab:ba", "standard", false, abba},
			{"ab:cd", "standard", false, abcd},
			{"xa:bc", "standard", false, xabc},
			{"dc:ba", "standard", false, dcba},
			{"xc:ba", "standard", false, xcba},
			{"ac:eg", "standard", false, aceg},
			{"xa:ce", "standard", false, xace},
			{"xe:ca", "standard", false, xeca},
			{"n:2*n", "standard", false, n2n},
			{"xn:3*n", "standard", false, xn3n},
		}
	}
)

func (s SetSlice) ToJsonSlice() SetJsonSlice {
	jsonSlice := make(SetJsonSlice, 0)
	for _, set := range s {
		jsonSlice = append(jsonSlice, &SetJson{
			Name:     set.Name,
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
			Typology: setjson.Typology,
			Enabled:  setjson.Enabled,
			Verify:   SetsFunctions[setjson.Name],
		})
	}
	return slice
}

// aa:aa
func aaaa(a, b, c, d int) bool {
	return a == b && b == c && c == d
}

// ?a:aa
func xaaa(_, b, c, d int) bool {
	return b == c && c == d
}

// ab:ab
func abab(a, b, c, d int) bool {
	return a == c && b == d && a != b
}

// ab:ba
func abba(a, b, c, d int) bool {
	return a == d && b == c && a != b
}

// ab:cd
func abcd(a, b, c, d int) bool {
	return b == a+1 && c == b+1 && d == c+1
}

// ?a:bc
func xabc(_, b, c, d int) bool {
	return c == b+1 && d == c+1
}

// dc:ba
func dcba(a, b, c, d int) bool {
	return c == d+1 && b == c+1 && a == b+1
}

// ?c:ba
func xcba(_, b, c, d int) bool {
	return c == d+1 && b == c+1
}

// ac:eg
func aceg(a, b, c, d int) bool {
	return b == a+2 && c == b+2 && d == c+2
}

// ?a:ce
func xace(_, b, c, d int) bool {
	return c == b+2 && d == c+2
}

// ?e:ca
func xeca(_, b, c, d int) bool {
	return c == d+2 && b == c+2
}

// n:2*n
func n2n(a, b, c, d int) bool {
	return 2*((a*10)+b) == (c*10)+d
}

// ?n:3*n
func xn3n(_, b, c, d int) bool {
	return 3*b == (c*10)+d
}
