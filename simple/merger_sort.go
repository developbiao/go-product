package main

import "fmt"

func mergeSort(data []int) []int {
	if len(data) <= 1 {
		return data
	}
	middle := len(data) / 2
	left := mergeSort(data[:middle])
	right := mergeSort(data[middle:])
	return merge(left, right)
}

func merge(left, right []int) (result []int) {
	l, r := 0, 0
	if l < len(left) && r < len(right) {
		if left[l] > right[r] {
			result = append(result, right[r])
			r++
		} else {
			result = append(result, left[l])
			l++
		}
	}
	result = append(result, left[l:]...)
	result = append(result, right[r:]...)
	return result
}

func main() {
	var data = []int{10, 7, 3, 50, 88, 17, 69, 52}
	data = mergeSort(data)
	fmt.Println(data)

}
