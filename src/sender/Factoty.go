package sender

import (
	"encoding/json"
	"fmt"
	"github.com/devops-salt/src/config"
	"github.com/devops-salt/src/log"
	"github.com/devops-salt/src/message"
	"github.com/devops-salt/src/util"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"net/http"
	"runtime"
	"sync"
	"time"
)

const (
	jobQueueSize int = 1024
	Sleep int = 2
	Timeout int = 3
)

type Factoty struct {
	jobQueue chan *message.Package
	lock     sync.RWMutex
	worker   map[string] int64
	sender   *UdpSender
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewFactory(sender *UdpSender) *Factoty {
	ctx, cancel := context.WithCancel(sender.ctx)
	return &Factoty{
		jobQueue: make(chan *message.Package, jobQueueSize),
		worker: map[string]int64{},
		ctx:      ctx,
		cancel:   cancel,
		sender: sender,
	}
}

func (f *Factoty) Employ(id string) {
	f.lock.Lock()
	f.worker[id] = time.Now().Unix()
	f.lock.Unlock()
	log.Info("Factoty.Employ worker[id=%s]", id)
}

func (f *Factoty) Dismiss(id string) {
	f.lock.Lock()
	delete(f.worker, id)
	f.lock.Unlock()
	log.Info("Factoty.Dismiss worker[id=%s]", id)
}
// Send heartbeat regularly
func (f *Factoty)  SyncSendHealth() {
	msg := &message.Package{
		Source: config.GetSource(),
		Attributes: message.Package_HEARTBEAT.Enum(),
	}
	ticker := time.NewTicker(time.Second * 2)
	for {
		msg.Timestamp = proto.Int64(time.Now().Unix())
		data, err := proto.Marshal(msg)
		if err != nil {
			log.Error("Resolve package failed: err=%s", err)
			break
		}
		if err := f.sender.PushBytes(data); err != nil {
			break
		}
		<- ticker.C
	}
	defer f.sender.Close()
}
// Receive heartbeat return, generate task and join queue
func (f *Factoty) SyncRecv() {
	for {
		temp := byteSlicePool.Get().([]byte)
		n, err := f.sender.ReadBytes(temp)
		if err != nil {
			log.Error("Factoty.JobTask read failed: %v", err)
			select {
			case <-time.After(2 * time.Second):
				f.sender.Reconnect()
			case <-f.ctx.Done():
				return
			}
			runtime.Gosched()
			continue
		}
		msg := &message.Package{}
		if err = proto.Unmarshal(temp[:n],msg); err != nil {
			log.Error("Factoty.JobTask Unmarshal failed: %v", err)
			continue
		}
		go func() {
			if (!util.CheckTask(msg, f.sender.opt.Source)) {
				log.Error("Factoty.JobTask CheckTask failed: %v", msg.String())
				return
			}
			log.Info("Factoty.JobTask push task_id :%s", *msg.TaskId)
			f.jobQueue <- msg
		}()
		byteSlicePool.Put(temp)
	}
	defer f.sender.Close()
}
// Processing task queue
func (f *Factoty) SyncProcess()  {
	for {
		select {
		case <-f.ctx.Done():
			return
		case job := <-f.jobQueue:
			go func(job *message.Package) {
				f.lock.Lock()
				_, ok := f.worker[*job.TaskId]
				f.lock.Unlock()
				if !ok {
					f.Employ(*job.TaskId)
					defer f.Dismiss(*job.TaskId)
					log.Info("Factoty.SyncProcess task_id :%s", *job.TaskId)
					// Download the task through task_id, the task download is complete, the platform task disappears, carry the second address to download
					task_url := fmt.Sprintf("%s/%s/%s", config.GetDownloadUrl(), job.Source[1], *job.TaskId)
					log.Info("Factoty.SyncProcess load task_id=%s", *job.TaskId)
					job_task, err:= util.LoadTask(task_url)
					if err != nil {
						log.Error("Factoty.SyncProcess load task_id=%s err:%s", *job.TaskId, err)
						return
					}
					if err := job_task.Check(); err != nil {
						log.Error("Factoty.SyncProcess load task_id=%s err:%s", *job.TaskId, err)
						return
					}
					log.Info("task_id:%s, issue:%s, source:%s", job_task.Task_id, job_task.Issue, job_task.Source)
					util.Command(job_task, f.CallbackData)
					log.Info("Factoty.SyncProcess task_id=%s finish", *job.TaskId)
				}
			}(job)
		}
	}
}

func (f *Factoty) CallbackData(msg *message.Callback, host string) {
	times := 0
	addr := []string{config.GetCallbackUrl(), host}
	msg_json, err := json.Marshal(msg)
	if err != nil {
		log.Error("Factoty.CallbackData json.Marshal failed err:%s ", err.Error())
		return
	}
	for _, value := range addr {
		if len(value) > 0 {
			for ; times < Timeout ; {
				times++
				code, err := util.Post(value, string(msg_json))
				if err != nil {
					continue
				}
				if (code == http.StatusOK) {
					break
				}
			}
			times = 0
		}
	}

}

func (f *Factoty) Close()  {
	log.Info("sender factory ready to stop")
	if (f.sender != nil) {
		f.sender.Close()
	}
	close(f.jobQueue)
	log.Info("sender factory stoped")

}
