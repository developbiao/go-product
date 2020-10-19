package main

import "fmt"

func main() {
	a := []bool{true, true, true}
	b := []bool{true, true, true}
	fmt.Println(BoolArrayEquals(a, b))

}

func BoolArrayEquals(a []bool, b []bool) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
