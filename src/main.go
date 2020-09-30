package main

import (
	"fmt"
	"github.com/devops-salt/src/config"
	"github.com/devops-salt/src/log"
	"github.com/devops-salt/src/salt"
	"github.com/devops-salt/src/util"
	"github.com/devops-salt/src/version"
	"os"
	"runtime"
)


func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	//if err := util.WritePidFile(config.PidFile); err != nil {
	//	log.Fatal("write pidfile failed, err=%s", err)
	//}
}

// usage prints command line usage to stderr.
func usage() {
	fmt.Fprintf(os.Stderr, "Usage:\n")
	//fmt.Fprintf(os.Stderr, "  ./salt -c filename [-d]\n")
	fmt.Fprintf(os.Stderr, "  ./salt -configfile filename\n")
	fmt.Fprintf(os.Stderr, "  ./salt -server\n")
	fmt.Fprintf(os.Stderr, "%s@%s %s\n", version.Product, version.Version, version.Website)
	os.Exit(2)
}

func main() {
	if (config.GetServer() != "server") {
		client := &salt.Program_client{}
		err := util.Run(client)
		if err != nil {
			log.Error("start client failed")
		}
		return
	}
	server := &salt.Program_server{}
	err := util.Run(server)
	if err != nil {
		log.Error("start service failed")
	}


}

