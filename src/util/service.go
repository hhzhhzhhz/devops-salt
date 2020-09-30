package util

import (
	"os"
	"os/signal"
	"syscall"
)

// Service interface
type Service interface {
	Init() error
	Start() error
	Stop() error
}

// Run 启动服务，并阻塞等待终端信号，接收到指定信号后退出，默认为（SIGINT, SIGTERM）
func Run(service Service, sigs ...os.Signal) error {
	if err := service.Init(); err != nil {
		return err
	}

	if err := service.Start(); err != nil {
		return err
	}

	if len(sigs) == 0 {
		sigs = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, sigs...)
	<-sigChan

	return service.Stop()
}

