package main

import (
	"fmt"
	"regexp"
)

var emailRegex = regexp.MustCompile("^[a-z0-9._%+\\-]+@[a-z0-9.\\-]+\\.[a-z]{2,4}$")

func main() {
	e := "developbiao@gmail.com"
	if isEmailValid(e) {
		fmt.Println(e + " is a valid email")
	} else {
		fmt.Println(e + " not valid email")
	}

}

func isEmailValid(e string) bool {
	if len(e) < 3 || len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}
