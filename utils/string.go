package utils

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func WordLimiter(s string, limit int) string {

	if strings.TrimSpace(s) == "" {
		return s
	}

	// convert string to slice
	strSlice := strings.Fields(s)

	// count the number of words
	numWords := len(strSlice)

	var result string

	if numWords > limit {
		// convert slice/array back to string
		result = strings.Join(strSlice[0:limit], " ")

		// the three dots for end characters are optional
		// you can change it to something else or remove this line
		result = result + "..."
	} else {

		// the number of limit is higher than the number of words
		// return default or else will cause
		// panic: runtime error: slice bounds out of range
		result = s
	}

	return result

}
