package main

import "fmt"

func main() {
	bar()
	fmt.Println("x")

	// output should be
	// a
	// b
	// a bug occur
	// d
	// x

}

func bar() {
	defer func() {
		fmt.Println("b")
		if err := recover(); err != nil {
			fmt.Println(err)
		}
		fmt.Println("d")
	}()

	// invoke foo()
	foo()
	fmt.Println("e")

}

func foo() {
	fmt.Println("a")
	panic("a bug occur")
	fmt.Println("c")
}
