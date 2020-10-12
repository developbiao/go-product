package main

import (
	"fmt"
	"time"
)

var done = make(chan bool)
var msgs = make(chan int)

func producer() {
	for i := 0; i < 10; i++ {
		msgs <- i
		time.Sleep(1000 * time.Millisecond)
	}
	done <- true
}

func consume() {
	for {
		msg := <-msgs
		fmt.Println(msg)
	}
}

func main() {
	go producer()
	go consume()
	//time.Sleep(3 * 1e9)
	<-done // block main gogroutine is over until channel done is false
	fmt.Println("main func is over!")
}
