package main

import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		for {
			select {
			case <-ch1:
				fmt.Println("成功收到ch1的数据: ", <-ch1)
			case ch2 <- 1:
				fmt.Println("成功向通道ch2中写入数据")
			case <-time.After(time.Second * 2):
				fmt.Println("Time out!!")
			}
		}
	}()

	for i := 0; i < 10; i++ {
		ch1 <- i
		fmt.Println("ch1写入数据: ", i)
		time.Sleep(1000 * time.Millisecond)
	}

	for i := 0; i < 10; i++ {
		str := <-ch2
		fmt.Println("获取到ch2的数据: ", str)
	}

	fmt.Println("Main function is over!")

}
