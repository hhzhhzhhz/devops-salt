package cache

import (
	"github.com/devops-salt/src/log"
	"github.com/devops-salt/src/message"
	"sync"
	"time"
)
var (
	flushTime int64 = 5*60
	GetCache Cache
)

type Cache struct {
	lock      sync.RWMutex
	data map[string] *Node // ip task
}

type Node struct {
	task      map[string] *message.Task // task_id message
	timestamp int64
	lock      sync.RWMutex
}

func (c *Cache) Put(addr string, task_id string, msg *message.Task) {
	c.lock.Lock()
	t, ok := c.data[addr]
	c.lock.Unlock()
	if ok {
		t.lock.Lock()
		defer t.lock.Unlock()
		t.timestamp = time.Now().Unix()
		t.task[task_id] = msg
	} else {
		c.lock.Lock()
		defer c.lock.Unlock()
		task := make(map[string] *message.Task, 10)
		task[task_id] = msg
		node := &Node{task: task,
			timestamp: time.Now().Unix(),
		}
		c.data[addr] = node
	}
}

func (c *Cache)  Get(addr string) map[string] *message.Task {
	c.lock.Lock()
	t, ok := c.data[addr]
	c.lock.Unlock()
	if ok {
		log.Info("Cache Get addr:%s",addr)
		return t.task
	}

	return nil
}

func (c *Cache) GetTask(addr string, task_id string) *message.Task {
	c.lock.Lock()
	t, ok := c.data[addr]
	c.lock.Unlock()
	if ok {
		t.lock.Lock()
		m, ok := t.task[task_id]
		defer t.lock.Unlock()
		if ok {
			log.Info("Cache GetTask delete addr:%s tak_id=%s",addr, task_id)
			delete(t.task, task_id)
			return m
		}
	}
	return nil
}

func (c *Cache)  Flush() {
	c.lock.Lock()
	defer c.lock.Unlock()
	time := time.Now().Unix()
	for addr, task := range c.data {
		if (time - task.timestamp > flushTime) {
			log.Info("Cache Flush delete addr:%s",addr)
			delete(c.data, addr)
		}
	}
}

func init()  {
	GetCache = Cache{
		data: make(map[string] *Node, 100),
	}
	go func(GetCache Cache) {
		ticker := time.NewTicker(time.Second * 30)
		for {
			GetCache.Flush()
			<- ticker.C
		}
	}(GetCache)


}


