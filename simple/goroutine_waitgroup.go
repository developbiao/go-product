package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

func main() {
	wg.Add(3)
	go Relief1()
	go Relief2()
	go Relief3()
	wg.Wait() // block main goroutine
	fmt.Println("Main end...")

}

func Relief1() {
	fmt.Println("function #01...")
	wg.Done()
}
func Relief2() {
	defer wg.Done()
	fmt.Println("function #02...")
}
func Relief3() {
	defer wg.Done()
	fmt.Println("function #03...")
}
