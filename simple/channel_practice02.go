package main

import "fmt"

func main() {
	ch1 := make(chan int)

	go func() {
		fmt.Println("=====子协程执行====")
		//for {
		//	v, ok := <- ch1
		//	if !ok {
		//		fmt.Println("通道已关闭")
		//		break
		//}
		//fmt.Println("读取到通道中的数据是:", v)
		//}

		for v := range ch1 {
			fmt.Println("读取到通道中的数据是:", v)
		}
	}()

	for i := 0; i < 5; i++ {
		ch1 <- i
	}
	close(ch1)

	fmt.Println("Main goroutine over!")
}
