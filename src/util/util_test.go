package util

import (
	"fmt"
	"github.com/devops-salt/src/message"
	"github.com/golang/protobuf/proto"
	"testing"
)

func TestCheckTask(t *testing.T) {
	msg := &message.Package{
		Source: []string{"127.0.0.1", "192.168.1.1"},
		Timestamp: proto.Int64(1600743387),
		Issue: proto.String("dir"),
		TaskId: proto.String("d2ae9f379b904ea184efaf343c8ec188"),
	}
	fmt.Println(CheckTask(msg, []string{"192.168.1.1"}))
	fmt.Println("xxx")


}
