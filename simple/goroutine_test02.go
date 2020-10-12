package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("In main()")
	go longWait()
	go shortWait()
	// time.Sleep(4 * 1e9) // when main goroutine is end long wait also end
	time.Sleep(10 * 1e9)
	fmt.Println("About to sleep in main()")

}

func longWait() {
	fmt.Println("Beginning longWait()")
	time.Sleep(5 * 1e9)
	fmt.Println("End of longWait()")
}

func shortWait() {
	fmt.Println("Beginning shortWait()")
	time.Sleep(2 * 1e9)
	fmt.Println("End of shortWait()")
}
