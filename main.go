package main

import(
	"github.com/maxzerbini/ovo/storage"
	"github.com/maxzerbini/ovo/inmemory"
	"github.com/maxzerbini/ovo/processor"
	"github.com/maxzerbini/ovo/server"
	"log"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

var ks storage.OvoStorage
var incoming *processor.InCommandQueue
var conf server.ServerConf
var srv *server.Server

func main(){
	protect(start)
	protect(stop)
}

func start() {
	conf = server.LoadConfiguration("./conf/severconf.json")
	log.Println("Start server node.")
	ks = inmemory.NewInMemoryStorage()
	incoming = processor.NewCommandQueue(ks)
	srv = server.NewServer(&conf, ks, incoming, nil, )
	srv.Do()
}

func stop(){
	log.Println("Stop server node.")
}

func protect(g func()) {
	defer func() {
		// Println executes normally even if there is a panic
		if err := recover(); err != nil {
			log.Println("run time panic: %v", err)
		}
	}()
	g() // possible runtime-error
}