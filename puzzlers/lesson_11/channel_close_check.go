package main

import (
	"fmt"
	"time"
)

// Question 01: 如果在select 语句中发现某个通道已关闭，那么应该怎样屏蔽掉它所在的分支？
// Question 02: 在select语句与for语句联用时，怎样直接退出外层for语句？

func main() {
	ch1 := make(chan int)
	go sendInt(ch1)
	go suckInt(ch1)

loop:
	for {
		select {
		case _, ok := <-ch1:
			if !ok {
				ch1 = make(chan int)
				fmt.Println("Detected channel has been closed!")
				break loop
			}
		default:
			//fmt.Println("No channel selected")
		}
	}

	fmt.Println("main function is end!")

}

func sendInt(ch1 chan<- int) {
	for i := 0; i < 5; i++ {
		ch1 <- i
		time.Sleep(1e9)
		fmt.Printf("Send value:%d\n", i)
	}
	close(ch1)
}

func suckInt(ch1 <-chan int) {
	for v := range ch1 {
		fmt.Printf("Received value:%d\n", v)
	}
}
