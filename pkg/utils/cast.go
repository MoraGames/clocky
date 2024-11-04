package utils

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
