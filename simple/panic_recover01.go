package main

import "fmt"

// panic recover example 01
func main() {
	foo()
	fmt.Println("x")
}

func foo() {
	defer func() {
		fmt.Println("b")
		if err := recover(); err != nil {
			fmt.Println(err)
		}
		fmt.Println("d")
	}()
	fmt.Println("a")
	panic("a bug occur")
	fmt.Println("c")
}
