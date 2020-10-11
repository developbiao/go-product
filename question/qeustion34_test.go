package question

import (
	"fmt"
	"testing"
)

func Test34(t *testing.T) {
	s1 := []int{1, 2, 3}
	s2 := s1[1:]             //[2, 3]
	s2[1] = 4                //[2, 4]
	fmt.Println(s1)          // [1, 2, 4] // 切片底层数据结构是数组
	s2 = append(s2, 5, 6, 7) // append 操作会导致生成新的数组，内存重新分配
	fmt.Println(s1)
}
