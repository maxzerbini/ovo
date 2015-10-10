package main

import(
	"github.com/maxzerbini/ovo/util"
	"github.com/maxzerbini/ovo/storage"
	"github.com/maxzerbini/ovo/inmemory"
	"github.com/maxzerbini/ovo/processor"
	"github.com/maxzerbini/ovo/server"
	"log"
	"runtime"
	"flag"
)

var ks storage.OvoStorage
var incoming *processor.InCommandQueue
var outcmdproc *processor.OutCommandQueue
var conf server.ServerConf
var srv *server.Server
var configPath string = "./conf/serverconf.json"
var configPathTemp string = ".temp"

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.StringVar(&configPath, "conf", "./conf/serverconf.json", "path of the file severconf.json")
}

// Start the Ovo Key/Value Storage
func main(){
	flag.Parse()
	util.Protect(start)
	util.Protect(stop)
}

// Start the server node
func start() {
	conf = server.LoadConfiguration(configPath)
	conf.Init()
	server.WriteConfiguration(configPath+configPathTemp, conf)
	log.Println("Start server node.")
	ks = inmemory.NewInMemoryStorage()
	incoming = processor.NewCommandQueue(ks)
	outcmdproc = processor.NewOutCommandQueue(&conf.ServerNode, &conf.Topology, incoming)
	srv = server.NewServer(&conf, ks, incoming, outcmdproc )
	srv.Do()
}

func stop(){
	log.Println("Stop server node.")
}

