// Package main is the entry point to directly run an Ios server as an executable
package main

import (
	"flag"
	"github.com/golang/glog"
	"github.com/heidi-ann/ios/config"
	"github.com/heidi-ann/ios/ios/server"
	"os"
	"os/signal"
	"syscall"
)

// command line flags
var id = flag.Int("id", -1, "server ID [REQUIRED]")                                                                                          // required flag
var configFile = flag.String("config", os.Getenv("GOPATH")+"/src/github.com/heidi-ann/ios/ios/example.conf", "Server configuration file") // optional flag
var diskPath = flag.String("disk", ".", "Path to directory to store persistent storage")                                                     // optional flag

// entry point of server executable
func main() {
	// set up logging
	flag.Parse()
	defer glog.Flush()

	// parse configuration file
	conf := config.ParseServerConfig(*configFile)
	if *id == -1 {
		glog.Fatal("ID is required")
	}
	if *id >= len(conf.Peers.Address) {
		glog.Fatal("Node ID is ", *id, " but is configured with a ", len(conf.Peers.Address), " node cluster")
	}

	// logging
	glog.Info("Starting server ", *id)
	defer glog.Warning("Shutting down server ", *id)

	// start Ios server
	go server.RunIos(*id, conf, *diskPath)

	// waiting for exit
	// always flush (whatever happens)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs
	glog.Flush()
	glog.Warning("Shutting down due to ", sig)
}