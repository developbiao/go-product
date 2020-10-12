package main

import (
	"fmt"
	"strconv"
	"time"
)

func main() {
	// producer
	ch := producer()
	// consumer in the main conroutine
	consumer(ch)
	//time.Sleep(3 * 1e9)
	fmt.Println("Main func is over!")
}

func producer() chan string {
	ch := make(chan string)
	go func() {
		for i := 0; i <= 10; i++ {
			msg := "Product #" + strconv.Itoa(i)
			fmt.Println(msg)
			ch <- msg
			time.Sleep(1e9)
		}
		close(ch)
	}()
	return ch
}

func consumer(ch chan string) {
	for msg := range ch {
		fmt.Println("Consumer<< " + msg)
	}
	fmt.Println("ch closed.")
}
