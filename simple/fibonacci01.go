package main

import "fmt"

func fibonacci(n int) {
	var a, b int
	a = 0
	b = 1
	// 1, 1, 2, 3, 5
	for i := 0; i < n; i++ {
		a, b = b, a+b
		fmt.Println(a)
	}
}

func main() {
	fibonacci(5)
}
