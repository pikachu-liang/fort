package main

import (
	"flag"
	"github.com/fort-io/fort/pkg/store"
	"github.com/golang/glog"
	"github.com/spf13/pflag"
	"net/http"
)

var (
	port    = pflag.Int("port", 8080, "The port to listen on.  Default 8080.")
	address = pflag.String("address", "127.0.0.1", "The address on the local server to listen to. Default 127.0.0.1")
	dataDir = pflag.String("data_dir", "~/fort", "The data dir to put all data.  Default ~/fort.")
)

func queryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		_, _ = w.Write([]byte("hello"))
		return
	}

}

func main() {
	flag.Parse()
	kvstore, _ := store.NewBBoltStore(store.BBoltStoreDefaultOptions)
	_ = kvstore.Set("hello", []byte("world"))
	_ = dataDir

	http.HandleFunc("/", queryHandler)
	if err := http.ListenAndServe("127.0.0.1:8080", nil); err != nil {
		glog.Fatalf("ListenAndServe failed. address: %s, port: %d", *address, *port)
	}

}
