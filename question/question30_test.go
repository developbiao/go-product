package question

import (
	"fmt"
	"testing"
)

func f1() (r int) {
	defer func() {
		r++
	}()
	return 0
}

func f2() (r int) {
	t := 5
	defer func() {
		t = t + 5
	}()
	return t
}

func f3() (r int) {
	defer func(r int) {
		r = r + 5
	}(r)
	return 1
}

func Test30(t *testing.T) {
	fmt.Println(f1()) // 1
	fmt.Println(f2()) // 5
	fmt.Println(f3()) // 1
}
