package udp

import (
	"github.com/devops-salt/src/log"
	"github.com/devops-salt/src/sender"
	"github.com/devops-salt/src/util"
	"golang.org/x/net/context"
)


type UdpClient struct {
	factory *sender.Factoty
	waitGroup util.WaitGroupWrapper
}

func NewClient(opt *sender.Option) *UdpClient {
	client := &UdpClient{}
	udp_sender := sender.New(context.Background(), opt)
	factory := sender.NewFactory(udp_sender)
	client.factory = factory
	return client

}

func (c *UdpClient) Main() {
	c.waitGroup.Wrap(func() {
		c.factory.SyncSendHealth()
	})

	c.waitGroup.Wrap(func() {
		c.factory.SyncRecv()
	})

	c.waitGroup.Wrap(func() {
		c.factory.SyncProcess()
	})
}

func (c *UdpClient) Exit() {
	log.Info("UdpClient exiting")
	c.factory.Close()
	c.waitGroup.Wait()
	log.Info("UdpClient exited")
}
