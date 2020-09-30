package salt

import (
	"github.com/devops-salt/src/config"
	"github.com/devops-salt/src/http"
	"github.com/devops-salt/src/sender"
)

type Program_server struct {
	server *Server
}

func (c *Program_server) Init() error  {
	UdpOpt := &sender.Option{
		Host: config.GetUDPIPv4(),
		Port: config.GetUDPPort(),
	}

	HttpOpt := &http.Option{
		HTTPPort: config.GetHTTPPort(),
		Timeout: 120,

	}
	server := NewServer(UdpOpt, HttpOpt)
	c.server = server
	return nil
}

func (s *Program_server) Start() error  {
	s.server.Main()
	return nil
}

func (s *Program_server) Stop() error  {
	s.server.Exit()
	return nil
}
