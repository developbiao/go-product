package main

import (
	"fmt"
	"time"
)

func testgo1() {
	for i := 0; i <= 10; i++ {
		fmt.Println("Test goroutine01")
	}
}

func testgo2() {
	for i := 0; i <= 10; i++ {
		fmt.Println("Test goroutine02")
	}
}

func main() {
	go testgo1()
	go testgo2()

	for i := 0; i <= 10; i++ {
		fmt.Println("Main function execution", i)
	}
	time.Sleep(3000 * time.Millisecond)
}
