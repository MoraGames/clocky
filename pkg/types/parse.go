package types

import (
	"fmt"
	"regexp"
	"strings"
)

func ParseSlice(s string) ([]string, error) {
	//s must match the regexp
	exp := regexp.MustCompile("^\\[((\"[a-zA-Z0-9_+-]+\")(,\"[a-zA-Z0-9_+-]+\")*)?\\]$")
	match := exp.MatchString(s)
	if match {
		slice := make([]string, 0)
		if s != "[]" {
			for _, v := range strings.Split(strings.Trim(s, "[]"), ",") {
				slice = append(slice, strings.ReplaceAll(strings.Trim(v, "\""), "_", " "))
			}
		}
		return slice, nil
	} else {
		return nil, fmt.Errorf("string %q doesn't match the regexp", s)
	}
}
