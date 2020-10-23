package main

func main() {
	// example 01
	ch1 := make(chan int, 1)
	ch1 <- 1
	//ch1 <-2 // 通道已满，因些这里会造成阻塞

	// example 02
	ch2 := make(chan int, 1)
	//elem, ok := <-ch2 // 通道已空，因些这里会阻塞
	//_, _ = elem, ok
	ch2 <- 1

	// example 03
	var ch3 chan int
	ch3 <- 1 // 通道的值为nil, 因为这里会造成永久的阻塞!
	<-ch3    // 通道的值为nil, 因些这里会造成永久的阻塞！
	_ = ch3

}
