package auth

import (
	"unicode"
)

func CheckPasswordStrength(password string) (longEnough, hasSpecial, hasUpper bool) {
	if len(password) < 14 {
		longEnough = false
		return
	}

	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpper = true
		}

		if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
			hasSpecial = true
		}
		if hasUpper && hasSpecial {
			break
		}
	}

	return
}
