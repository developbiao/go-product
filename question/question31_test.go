package question

import (
	"fmt"
	"testing"
)

type Person struct {
	age int
}

func Test31(t *testing.T) {
	person := &Person{28}
	// 1
	defer fmt.Println(person.age)
	// 2
	defer func(p *Person) {
		fmt.Println(p.age)
	}(person)
	// 3
	defer func() {
		fmt.Println(person.age)
	}()

	person.age = 29
}
