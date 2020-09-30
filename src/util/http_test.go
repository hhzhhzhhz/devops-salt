package util

import (
	"fmt"
	"testing"
)

func TestGet(t *testing.T)  {
	res, err := LoadTask("http://127.0.0.1:8080/someJSON")
	if (err != nil) {
		fmt.Println(err)
		return
	}
	fmt.Println(res.String())
}

