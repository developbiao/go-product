package main

import "fmt"

func BubblingASC(values []int) {
	for i := 0; i < len(values)-1; i++ {
		for j := i + 1; j < len(values); j++ {
			if values[i] > values[j] {
				values[i], values[j] = values[j], values[i]
			}
		}
	}
	fmt.Println(values)
}

func BubblingDESC(values []int) {
	for i := 0; i < len(values)-1; i++ {
		for j := i + 1; j < len(values); j++ {
			if values[i] < values[j] {
				values[i], values[j] = values[j], values[i]
			}
		}
	}
	fmt.Println(values)
}

func main() {
	values := []int{3, 8, 10, 90, 20, 100, 52, 33, 77}
	fmt.Println(values)
	BubblingASC(values)
	BubblingDESC(values)
	fmt.Println("Programing is ok!")

}
