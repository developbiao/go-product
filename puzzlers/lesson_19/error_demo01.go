package main

import (
	"errors"
	"fmt"
)

func echo(request string) (response string, err error) {
	if request == "" {
		err = errors.New("Empty request")
		return
	}
	response = fmt.Sprintf("echo: %s", request)
	return
}

func main() {

	// Example01
	for _, req := range []string{"", "hello!"} {
		fmt.Printf("request: %s\n", req)
		resp, err := echo(req)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			continue
		}
		fmt.Printf("response: %s\n", resp)
	}

	// Example 02
	err1 := fmt.Errorf("invalid contents: %s", "#$%")
	err2 := errors.New(fmt.Sprintf("invlaid contents: %s", "#$%"))
	if err1.Error() == err2.Error() {
		fmt.Println("The error messages in err1 and err2 are the same")
	}

	fmt.Println("ok")
}
