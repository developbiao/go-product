package main

import (
	"fmt"
	"math"
)

func main() {
	var arr = []int{4, 5, 6, 5, 6, 7, 8, 9, 10, 9}
	fmt.Println(arr)
	index := findIndex(arr, 8)
	fmt.Println("find index is :", index)

}

func findIndex(data []int, number int) int {
	// Because rules abs is between 1 so
	// nextIndex = abs(number - data[0]) < number
	nextIndex := int(math.Abs(float64(number - data[0])))
	for nextIndex < len(data) {
		if data[nextIndex] == number {
			// found index
			return nextIndex
		}
		// jump index
		nextIndex += int(math.Abs(float64(number - data[nextIndex])))
	}
	return -1
}
