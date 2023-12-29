package types

import (
	"fmt"
	"regexp"
	"strings"
)

func ParseSlice(s string) ([]string, error) {
	//s must match the regexp
	exp := regexp.MustCompile("^\\[(\"[a-zA-Z0-9 ]+\")(, \"[a-zA-Z0-9 ]+\")*\\]$")
	match := exp.MatchString(s)
	if match {
		slice := strings.Split(strings.Trim(s, "[]"), ", ")
		for i, v := range slice {
			slice[i] = strings.Trim(v, "\"")
		}
		return slice, nil
	} else {
		return nil, fmt.Errorf("string %q doesn't match the regexp", s)
	}
}
