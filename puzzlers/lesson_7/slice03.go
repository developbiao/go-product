package main

import "fmt"

func main() {

	// 一个切片的底层数据永远不会被替换
	// Go语言在扩容时候一定会生成新的底层数组，但是它也同时生成了新的切片。
	// 它是把新的切片作为了新底层数组的窗口，没有对原切片其底层数组做任何改动。
	// 示例1。
	a1 := [7]int{1, 2, 3, 4, 5, 6, 7}
	fmt.Printf("a1: %v (len: %d, cap: %d)\n",
		a1, len(a1), cap(a1))
	s9 := a1[1:4]
	s9[0] = 97
	fmt.Printf("s9: %v (len: %d, cap: %d)\n",
		s9, len(s9), cap(s9))
	for i := 1; i <= 5; i++ {
		s9 = append(s9, i)
		fmt.Printf("s9(%d): %v (len: %d, cap: %d)\n",
			i, s9, len(s9), cap(s9))
	}
	fmt.Printf("a1: %v (len: %d, cap: %d)\n",
		a1, len(a1), cap(a1))
	fmt.Println()

}
