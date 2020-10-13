package main

import "fmt"

func fibonacci03() func() int {
	a, b := 0, 1
	return func() int {
		a, b = b, a+b
		return a
	}
}

func main() {
	f := fibonacci03()
	for i := 0; i < 5; i++ {
		fmt.Println(f())
	}

}
