package main

import "fmt"

// Dynamic programming
func fibnocci03(n int) int {
	if n == 0 || n == 1 {
		return n
	}

	dp := make([]int, n+1)
	dp[0] = 0
	dp[1] = 1
	for i := 2; i <= n; i++ {
		dp[i] = (dp[i-1] + dp[i-2]) % 1000000007
	}
	return dp[n]
}

func main() {
	for i := 1; i <= 5; i++ {
		num := fibnocci03(i)
		fmt.Println(num)
	}
	fmt.Println("ok")
}