package main

import "fmt"

func main() {
	s2 := []int{1, 2, 3, 4}
	s3 := []int{7, 8, 9}
	copy(s2, s3)
	fmt.Println(s2)
	fmt.Println(s3)

	copy(s3, s2[2:])
	fmt.Println(s2)
	fmt.Println(s3)

}
