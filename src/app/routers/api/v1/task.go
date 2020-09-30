package v1

import (
	"encoding/json"
	"fmt"
	"github.com/devops-salt/src/cache"
	"github.com/devops-salt/src/log"
	"github.com/devops-salt/src/message"
	"github.com/devops-salt/src/util"
	"github.com/gin-gonic/gin"
	"time"
)

func DownTask(c *gin.Context)  {
	task := c.Param("task_id")
	addr := c.Param("addr")
	if len(task) != 0 && len(addr) != 0 {
		msg := cache.GetCache.GetTask(addr, task)
		if msg != nil {
			c.JSON(200, msg)
			return
		}
	}
	c.JSON(404, &message.Task{})
}

func AddTask(c *gin.Context)  {
	d, err := c.GetRawData()
	if err!=nil {
		log.Error("%s", err)
	}
	sub := []*message.Submit{}

	if err := json.Unmarshal(d, &sub); err != nil {
		log.Error("[Error] ", err.Error())
		c.JSON(400, message.Res{Code: 500, Result: "Data format error!"})
		return
	}
	res := [] *TaskStatus{}
	for _, msg :=range sub {
		if err := msg.Check(); err != nil {
			log.Error("[Error] ", err.Error())
			c.JSON(200, message.Res{Code: 500, Result: err.Error()})
			return
		}
		if len(msg.Nat_dest) == 0 {
			msg.Nat_dest = msg.Internet_dest
		}
		timestamp := time.Now().Unix()
		issue:= util.Bsse64encode(msg.Issue)
		task_id := util.Md5(fmt.Sprintf("%v%v%d",msg.Nat_dest, issue, timestamp))
		task := &message.Task{Issue: issue,
			Timestamp: timestamp,
			Task_id : task_id,
			Callback: msg.Callback,
			Source: []string{msg.Nat_dest},
			Attributes: 1,
		}
		cache.GetCache.Put(msg.Internet_dest, task_id, task)
		log.Info("api.task add task ip:%s,id:%s,msg:%s", msg.Internet_dest, task_id, task.String())
		status := TaskStatus{Code: 200, Task_id: task_id, Addr: msg.Internet_dest}
		res = append(res, &status)
	}
	res_json, err := json.Marshal(res)
	if err != nil {
		log.Error("json.Marshal err:%s", err)
		c.JSON(400, message.Res{Code: 500, Result: "Unknown error"})
		return
	}
	c.JSON(200, string(res_json))

}

type TaskStatus struct {
	Code int32      `json:"code"`
	Task_id string  `json:"task_id"`
	Addr string     `json:"addr"`
}
