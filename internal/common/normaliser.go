package common

import "strings"

func NormaliseStringUpper(target string) string {
	return strings.TrimSpace(strings.ToUpper(target))
}

func NormaliseStringLower(target string) string {
	return strings.TrimSpace(strings.ToLower(target))
}

func NormaliseNoSpaces(target string) string {
	return strings.TrimSpace(target)
}

func IsNotNilAndNotEmpty(target interface{}) bool {
	if target == nil {
		return false
	}

	targetNormalised := NormaliseNoSpaces(target.(string))
	return targetNormalised != ""
}
