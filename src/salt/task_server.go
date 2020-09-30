package salt

import (
	"fmt"
	"github.com/devops-salt/src/app/routers"
	myhttp "github.com/devops-salt/src/http"
	"github.com/devops-salt/src/log"
	"github.com/devops-salt/src/reciever"
	"github.com/devops-salt/src/sender"
	"github.com/devops-salt/src/util"
	"golang.org/x/net/context"
	"net/http"
	"time"
)

type Server struct {
	factory *reciever.Factoty
	waitGroup util.WaitGroupWrapper
	http.Server
	ctx context.Context
}

func NewServer(Udpopt *sender.Option, HttpOpt *myhttp.Option) *Server {
	server := &Server{
		ctx: context.Background(),
	}
	receiver := reciever.New(server.ctx, Udpopt)
	factory := reciever.NewFactory(receiver)
	server.factory = factory

	router := routers.InitRouter()
	server.Server = http.Server{
		Addr:           fmt.Sprintf(":%d", HttpOpt.HTTPPort),
		Handler:        router,
		ReadTimeout:    time.Duration(HttpOpt.Timeout) * time.Second,
		WriteTimeout:   time.Duration(HttpOpt.Timeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return server
}

func (s *Server) Main() {
	s.waitGroup.Wrap(func() {
		s.factory.SyncRecvHandle()
	})
	s.waitGroup.Wrap(func() {
		log.Info("HTTP Server start")
		if err := s.ListenAndServe(); err != nil {
			log.Error("HTTP Server err:%s", err)
		}
	})
}


func (s *Server) Exit() {
	log.Info("Server exiting")
	s.factory.Close()
	if err := s.Shutdown(s.ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}
	log.Info("HTTP Server Stoped")
	s.waitGroup.Wait()
	log.Info("Server exited")
}

