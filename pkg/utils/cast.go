package utils

func ItoB(i int) bool {
	if i == 0 {
		return false
	}
	return true
}

func BtoI(b bool) int {
	if !b {
		return 0
	}
	return 1
}

func BtoS(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func StoB(s string) bool {
	if s == "true" {
		return true
	}
	return false
}
