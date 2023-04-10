package common

import (
	"math/rand"
	"strings"
	"time"
)

const (
	lowercaseLetters = "abcdefghijklmnopqrstuvwxyz"
	uppercaseLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	specialChars     = "!@#$%^&*()-_=+[]{}|;:,.<>?/"
)

func GenerateRandomStringWithPrefix(limit int, upper, lower, withSpecialChars bool, prefix string) string {
	var characters string

	if upper {
		characters += uppercaseLetters
	}
	if lower {
		characters += lowercaseLetters
	}
	if withSpecialChars {
		characters += specialChars
	}

	if characters == "" {
		return prefix
	}

	remainingLimit := limit - len(prefix)

	if remainingLimit <= 0 {
		return prefix
	}

	rand.NewSource(time.Now().UnixNano())
	builder := strings.Builder{}
	builder.WriteString(prefix)
	for i := 0; i < remainingLimit; i++ {
		randomIndex := rand.Intn(len(characters))
		builder.WriteByte(characters[randomIndex])
	}

	return builder.String()
}
