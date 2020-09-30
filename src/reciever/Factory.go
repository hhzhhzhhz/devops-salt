package reciever

import (
	"github.com/devops-salt/src/cache"
	"github.com/devops-salt/src/log"
	"github.com/devops-salt/src/message"
	"github.com/devops-salt/src/util"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"io"
	"net"
	"sync"
)

const (
	jobQueueSize int = 1 // 心跳队列 1
)


type Factoty struct {
	//jobQueue chan *Message
	worker   map[string] int64
	lock     sync.RWMutex
	receiver   *UdpService
	ctx      context.Context
	cancel   context.CancelFunc
}

type Message struct {
	message *message.Package
	addr *net.UDPAddr
}

func NewFactory(receiver *UdpService) *Factoty {
	ctx, cancel := context.WithCancel(receiver.ctx)
	return &Factoty{
		//jobQueue: make(chan *Message, jobQueueSize),
		worker: make(map[string]int64, 100),
		ctx:      ctx,
		cancel:   cancel,
		receiver: receiver,
	}
}

func (f *Factoty) Employ(id string) {
	f.lock.Lock()
	f.worker[id] = 0
	f.lock.Unlock()
	//log.Info("Factoty.Employ worker[id=%s]", id)
}

func (f *Factoty) Dismiss(id string) {
	f.lock.Lock()
	delete(f.worker, id)
	f.lock.Unlock()
	//log.Info("Factoty.Dismiss worker[id=%s]", id)
}


// Heartbeat processing
func (f *Factoty) SyncRecvHandle() {
	for {
		temp := byteSlicePool.Get().([]byte)
		n, addr, err := f.receiver.ReadBytes(temp)
		if err != nil {
			if err != io.EOF {
				log.Error("Reciever.SyncRecvHealth conn err=%s", err)
			}
			log.Error("Reciever.SyncRecvHealth conn.Read err=%s", err)
			break
		}
		msg := &message.Package{}
		if err := proto.Unmarshal(temp[:n], msg); err != nil {
			// 阿里聚石塔 会往udp 端口发送21字节数据包，打印解析失败数据包 日志会量大
			//log.Error("receiver.Factory Unmarshal failed len:%d err:%s",n, err)
			byteSlicePool.Put(temp)
			continue
		}
		if (util.CheckHeartbeat(msg)) {
			Msg := &Message{message: msg, addr: addr}
			go func(Msg *Message) {
				addr := Msg.addr.IP.String()
				f.lock.Lock()
				_, ok := f.worker[addr]
				f.lock.Unlock()
				if !ok {
					// Only process one heartbeat packet of the same address at a time
					f.Employ(addr)
					defer f.Dismiss(addr)
					// 查询任务
					v := cache.GetCache.Get(addr)
					if v == nil {
						return
					}
					for k, t := range v {
						// The order is internal network address, external network address
						addrs := []string{t.Source[0], addr }
						task := &message.Package{
							Timestamp: proto.Int64(t.Timestamp),
							Attributes: message.Package_ISSUE_TASK.Enum(),
							TaskId: proto.String(k),
							Issue: proto.String(t.Issue),
							Callback: proto.String(t.Callback),
							Source: addrs,
						}
						data, err := proto.Marshal(task)
						if err != nil {
							log.Error("recever factory SyncRecvHandle send Marshal  err:%s", err)
							return
						}
						f.receiver.PushBytes(data, Msg.addr)
					}
				}
			}(Msg)
		}
		byteSlicePool.Put(temp)
	}
	f.receiver.Close()
}

func (f *Factoty) Close() {
	log.Info("udp recever factory ready to stop")
	f.cancel()
	if (f.receiver != nil) {
		f.receiver.Close()
	}
	log.Info("udp recever factory stoped")
}
