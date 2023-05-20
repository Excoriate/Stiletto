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

	switch value := target.(type) {
	case string:
		return value != ""
	case bool:
		return true
	case int:
		return value != 0
	case []interface{}:
		return len(value) > 0
	case map[string]interface{}:
		return len(value) > 0
	default:
		return false
	}
}

func RemoveDoubleQuotes(target string) string {
	return strings.Trim(target, "\"")
}
