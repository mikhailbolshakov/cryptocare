package kit

import (
	"regexp"
)

// IsEmailValid checks email format
func IsEmailValid(email string) bool {
	if len(email) < 3 || len(email) > 254 {
		return false
	}
	match, err := regexp.MatchString("^(((\\\\.)|[^\\s[:cntrl:]\\(\\)<>@,;:'\\\\\\\"\\.\\[\\]]|')+|(\"(\\\\\"|[^\"])*\"))(\\.(((\\\\.)|[^\\s[:cntrl:]\\(\\)<>@,;:'\\\\\\\"\\.\\[\\]]|')+|(\"(\\\\\"|[^\"])*\")))*@[a-zA-Z0-9а-яА-Я](?:[a-zA-Z0-9а-яА-Я-]{0,61}[a-zA-Z0-9а-яА-Я])?(?:\\.[a-zA-Z0-9а-яА-Я](?:[a-zA-Z0-9а-яА-Я-]{0,61}[a-zA-Z0-9а-яА-Я])?)*$", email)
	return match && err == nil
}
