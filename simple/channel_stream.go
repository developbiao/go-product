package main

import (
	"fmt"
	"time"
)

func main() {
	stream := pump()
	go suck(stream)
	time.Sleep(1e9)
}

func pump() chan int {
	ch := make(chan int)
	go func() {
		for i := 0; ; i++ {
			ch <- 1
		}
	}()
	return ch
}

func suck(ch chan int) {
	for v := range ch {
		fmt.Println(v)
	}
}
