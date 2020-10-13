package main

import "fmt"

func fibonacci02(num int) int {
	if num < 2 {
		return 1
	}
	return fibonacci02(num-1) + fibonacci02(num-2)
}

func main() {
	for i := 0; i < 10; i++ {
		nums := fibonacci02(i)
		fmt.Println(nums)
	}

}
