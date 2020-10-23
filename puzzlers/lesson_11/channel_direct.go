package main

import (
	"fmt"
	"math/rand"
)

// Practice code from
// https://github.com/hyper0x/Golang_Puzzlers/blob/master/src/puzzlers/article11/q1/demo23.go

func main() {
	// example 01
	// 只能发不能收的通道
	var uselessChan = make(chan<- int, 1)
	// 只能收不能发的通道
	var anotherUselessChan = make(<-chan int, 1)
	// 这里打印的可以分别代表两个通道的指针的16进制表示
	fmt.Printf("The useless channels: %v, %v\n", uselessChan, anotherUselessChan)

	// example 2
	intChan1 := make(chan int, 3)
	SendInt(intChan1)

	// example 4
	intChan2 := getIntChan()
	for elem := range intChan2 {
		fmt.Printf("The element in intChan2: %v\n", elem)
	}

	// example 5
	_ = GetIntChan(getIntChan)

}

// example 02
func SendInt(ch chan<- int) {
	ch <- rand.Intn(1000)
}

// example 03
type Noitifer interface {
	SendInt(ch chan<- int)
}

// example 04
func getIntChan() <-chan int {
	num := 5
	ch := make(chan int, num)
	for i := 0; i < num; i++ {
		ch <- i
	}
	close(ch)
	return ch
}

type GetIntChan func() <-chan int
