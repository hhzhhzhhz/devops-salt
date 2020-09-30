package sender

import (
	"fmt"
	"strings"
	"testing"
)

func TestUdpSender_Close(t *testing.T) {
	name := "xxxx"
	s := strings.Split(name, " ")
	fmt.Println(s[1:])

}
