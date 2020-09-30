package message

import (
	"errors"
	"fmt"
	"net"
	"time"
)
// program task package
type Task struct {
	Issue string
	Timestamp int64
	Task_id string
	Callback string
	Attributes int
	Source []string
}

func (t *Task) Check() error {
	if (len(t.Issue) == 0 || len(t.Task_id) == 0 || time.Now().Unix() - t.Timestamp > 30) {
		return errors.New(fmt.Sprintf("task load check err: issue:%v,timestamp:%v,task_id:%v",
			t.Issue, t.Timestamp, t.Task_id))
	}
	return nil

}

func (t *Task) String() string {
	return fmt.Sprintf("{issue:%v,timestamp:%v,task_id:%v,callback:%v,attributes:%v,source:%v",
		t.Issue, t.Timestamp, t.Task_id, t.Callback, t.Attributes, t.Source)
}
// Callback data
type Callback struct {
	Task_id string `json:"task_id"`
	Code int       `json:"code"`
	Data string    `json:"data"`
}

// submit task package
type Submit struct {
	Issue string          `json:"issue"`
	Callback string       `json:"callback"`
	Internet_dest string  `json:"internet_dest"`  // 外网地址
	Nat_dest string       `json:"nat_dest"`      // NAT 映射内网地址
}

func (s *Submit) Check() error{
	if len(s.Issue) == 0 || len(s.Internet_dest) == 0  || net.ParseIP(s.Internet_dest) == nil {
		return errors.New("Data verification error")
	}

	return nil
}

type Res struct {
	Code int16
	Result interface{}
}
