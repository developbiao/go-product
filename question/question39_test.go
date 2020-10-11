package question

import (
	"fmt"
	"testing"
)

type People interface {
	Show()
}

type Student struct{}

func (stu *Student) Show() {

}

func Test40(t *testing.T) {
	var s *Student
	if s == nil {
		fmt.Println("s is nil")
	} else {
		fmt.Println("s is not nil")
	}
	var p People = s
	if p == nil {
		fmt.Println("p is nil")
	} else {
		fmt.Println("p is not nil")
	}

	var p2 People
	if p2 == nil {
		fmt.Println("p2 is nil")
	} else {
		fmt.Println("p2 is not nil")
	}

}
