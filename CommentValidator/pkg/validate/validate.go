package validate

import "regexp"

var (
	valid = regexp.MustCompile("(qwerty|йцукен|zxvbnm)")
)

// false if message contains stopwords, else true
func IsValid(text string) bool {
	if valid.MatchString(text) {
		return false
	}
	return true
}
