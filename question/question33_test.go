package question

import "testing"

type S struct {
}

func m(x interface{}) {
}

func g(x *interface{}) {
}

func Test33(t *testing.T) {
	s := S{}
	p := &s
	m(s) //A
	//g(s) //B
	m(p) //C
	//g(p) //D
}
