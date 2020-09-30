package sender

import (
	"fmt"
	"github.com/devops-salt/src/log"
	"golang.org/x/net/context"
	"net"
	"sync"
)

type Option struct {
	Host    string
	Port    int
	Source  []string
}

const (
	defaultMaxBufferSize     = 4096
	CONST_MAXRECVBUFFER  int = 32 * 1024 * 1024
	MaxSectionSize       int = 4 * 1024
)
var byteSlicePool = sync.Pool{
	New: func() interface{} {
		return make([]byte, defaultMaxBufferSize)
	},
}



type UdpSender struct {
	//url string
	net.Conn
	ctx    context.Context
	cancel context.CancelFunc
	//source []string
	opt *Option
}

func New(ctx context.Context, opt *Option) *UdpSender  {
	ctx, cancel := context.WithCancel(ctx)
	pc := &UdpSender{
		ctx:  ctx,
		cancel: cancel,
		opt: opt,
		//url: opt.Host,
		//source: opt.Source,
	}
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", opt.Host, opt.Port))
	if err != nil {
		log.Error("Resolve udp addr failed: err=%s", err)
	}
	pc.Conn, err = net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Error("Dial to server failed: %v\n", err)
	}
	log.Info("send udp client[%s] start successful", fmt.Sprintf("%s:%d", opt.Host, opt.Port))
	return pc
}

func (s *UdpSender) ReadBytes(data []byte) (int, error) {
	return s.Conn.Read(data)
}


func (s *UdpSender) PushBytes(data []byte) error {
	_, err := s.Conn.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (s *UdpSender) Reconnect() error{
	s.Conn.Close()
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", s.opt.Host, s.opt.Port))
	if err != nil {
		log.Error("Reconnect Resolve udp addr failed: err=%s", err)
	}
	s.Conn, err = net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Error("Reconnect Dial to server failed: %v\n", err)
	}
	log.Info("Reconnect send udp client[%s] start successful", fmt.Sprintf("%s:%d", s.opt.Host, s.opt.Port))
	return nil
}

func (s *UdpSender)  Close(){
	log.Info("Sender ready close")
	s.Conn.Close()
	s.cancel()
	log.Info("Sender closed")

}

