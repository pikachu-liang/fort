package main

import (
	"flag"
	"go.etcd.io/etcd/raft/v3/raftpb"
	"strings"
)

/*
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
*/

func main() {

	cluster := flag.String("cluster", "http://127.0.0.1:9021", "comma separated cluster peers")
	id := flag.Int("id", 1, "node ID")
	kvport := flag.Int("port", 9121, "key-value server port")
	join := flag.Bool("join", false, "join an existing cluster")
	flag.Parse()

	proposeC := make(chan string)
	defer close(proposeC)
	confChangeC := make(chan raftpb.ConfChange)
	defer close(confChangeC)

	// raft provides a commit stream for the proposals from the http api
	var kvs *kvstore
	getSnapshot := func() ([]byte, error) { return kvs.getSnapshot() }
	commitC, errorC, snapshotterReady := newRaftNode(*id, strings.Split(*cluster, ","), *join, getSnapshot, proposeC, confChangeC)

	kvs = newKVStore(<-snapshotterReady, proposeC, commitC, errorC)

	// the key-value http handler will propose updates to raft
	serveHttpKVAPI(kvs, *kvport, confChangeC, errorC)

}
