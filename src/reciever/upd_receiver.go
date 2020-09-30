package reciever

import (
	"fmt"
	"github.com/devops-salt/src/log"
	"github.com/devops-salt/src/sender"
	"golang.org/x/net/context"
	"net"
	"sync"
)

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
type UdpService struct {
	*net.UDPConn
	ctx    context.Context
	cancel context.CancelFunc
}

func New(ctx context.Context, opt *sender.Option) *UdpService {
	ctx, cancel := context.WithCancel(ctx)
	reciever := &UdpService{ctx: ctx, cancel: cancel}
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", opt.Host, opt.Port))
	if err != nil {
		log.Fatal("Reciever net.ResolveTCPAddr err=%s", err)
	}
	reciever.UDPConn, err = net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatal("Reciever net.ListenUDP err=%s", err)
	}
	log.Info("reciever udp server[%s] start successful", fmt.Sprintf("%s:%d", opt.Host, opt.Port))
	return reciever

}

func (r *UdpService) ReadBytes(data []byte) (int, *net.UDPAddr, error){
	return r.UDPConn.ReadFromUDP(data)
}

func (r *UdpService) PushBytes(data []byte, addr net.Addr) (int, error) {
	return r.UDPConn.WriteTo(data, addr)
}

func (r *UdpService) Close() error {
	r.cancel()
	err := r.UDPConn.Close()
	if err != nil {
		log.Error("Reciever.Close err=%s", err)
	}
	log.Info("reciever stoped")
	return err
}
