package main

import "fmt"

func QuickSort(list []int, low, high int) {
	if high > low {
		// 位置划分
		pivot := partition(list, low, high)
		// 左部分排序
		QuickSort(list, low, pivot-1)
		// 右部分排序
		QuickSort(list, pivot+1, high)
	}
}

func partition(list []int, low, high int) int {
	pivot := list[low]
	for low < high {
		// high 指针值 >= pivot high 指针左移
		for low < high && pivot <= list[high] {
			high--
		}
		list[low] = list[high]
		for low < high && pivot >= list[low] {
			low++
		}
		list[high] = list[low]
	}
	list[low] = pivot
	return low
}

func main() {
	list := []int{55, 90, 74, 20, 16, 46, 43, 59, 2, 99, 79, 10, 73, 1, 68, 56, 3, 87, 40, 78, 14, 18, 51, 24, 57, 89}

	QuickSort(list, 0, len(list)-1)
	fmt.Println(list)
}
