package events

type SetSlice []*Set
type Set struct {
	Name     string
	Typology string
	Enabled  bool
	Verify   func(h1, h2, m1, m2 int) bool
}

var Sets = SetSlice{
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

/* func points(time time.Time) int {
 hour := time.Hour()
 hour1 := hour / 10
 hour2 := hour % 10
 minute := time.Minute()
 minute1 := minute / 10
 minute2 := minute % 10
 points := 0
 if aaaa(hour1, hour2, minute1, minute2) {
  points += 1
 }
 if xaaa(hour1, hour2, minute1, minute2) {
  points += 1
 }
 if abab(hour1, hour2, minute1, minute2) {
  points += 1
 }
 if abba(hour1, hour2, minute1, minute2) {
  points += 1
 }
 if abcd(hour1, hour2, minute1, minute2) {
  points += 1
 }
 if xabc(hour1, hour2, minute1, minute2) {
  points += 1
 }
 if dcba(hour1, hour2, minute1, minute2) {
  points += 1
 }
 if xcba(hour1, hour2, minute1, minute2) {
  points += 1
 }
 if aceg(hour1, hour2, minute1, minute2) {
  points += 1
 }
 if xace(hour1, hour2, minute1, minute2) {
  points += 1
 }
 if xeca(hour1, hour2, minute1, minute2) {
  points += 1
 }
 if ab2ab(hour1, hour2, minute1, minute2) {
  points += 1
 }
 return points
} */
