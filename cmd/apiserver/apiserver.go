package main

import (
	"github.com/golang/glog"
	"github.com/spf13/pflag"
)

var (
	port           = pflag.Int("port", 8080, "The port to listen on.  Default 8080.")
	address        = pflag.String("address", "127.0.0.1", "The address on the local server to listen to. Default 127.0.0.1")
	apiPrefix      = pflag.String("api_prefix", "/api/alpha", "The prefix for API requests on the server. Default '/api/alpha'")
	etcdServerList = pflag.StringArray("etcd_servers", []string{"192.168.1.21:2379"}, "Servers for the etcd (http://ip:port), comma separated")
	// TODO: let node register itself to API server
	nodeList = pflag.StringArray("nodes", []string{"nodes"}, "List of nodes to place segments onto, comma separated.")
)

func main() {
	if len(*nodeList) == 0 {
		glog.Fatal("No nodesList specified!")
	}
	glog.Infof("apiserver started successfully. address: %s, port: %d, apiPrefix: %s, etcdServerList: %s",
		*address, *port, *apiPrefix, *etcdServerList)
}
