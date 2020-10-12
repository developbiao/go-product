package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string)
	var answer string

	fmt.Println("Reveal romaintic feelings... ")
	go sendLove()
	go responseLove(ch)
	waitFor()
	answer = <-ch
	if answer != "" {
		fmt.Println(answer)
	} else {
		fmt.Println("Dead ☠️....")
	}
}

func waitFor() {
	for i := 0; i < 5; i++ {
		time.Sleep(1 * 1e9)
		fmt.Println("Keep Waiting...")
	}
}

func sendLove() {
	fmt.Println("Love you, 美女老师 ❤️")
}

func responseLove(ch chan string) {
	time.Sleep(6 * 1e9)
	ch <- "Love you, too"
}
