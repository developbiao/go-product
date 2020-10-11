package main

import "fmt"

func main() {
	var channel chan int
	fmt.Printf("通道的数据类型：%T, 通道的值:%v\n", channel, channel)

	if channel == nil {
		channel = make(chan int)
		fmt.Printf("通过make创建的数据类型%T, 通道的值:%v, \n", channel, channel)
	}

}
