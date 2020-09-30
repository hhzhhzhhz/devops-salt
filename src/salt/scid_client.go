package salt

import (
	"github.com/devops-salt/src/config"
	"github.com/devops-salt/src/sender"
	"github.com/devops-salt/src/udp"
)

type Program_client struct {
	client *udp.UdpClient
}


func (c *Program_client) Init() error  {
	opt := &sender.Option{Host: config.GetUDPIPv4(), Port: config.GetUDPPort(), Source: config.GetSource()}
	client := udp.NewClient(opt)
	c.client = client
	return nil
}

func (c *Program_client) Start() error  {
	c.client.Main()
	return nil
}

func (c *Program_client) Stop() error  {
	c.client.Exit()
	return nil
}
