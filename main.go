package main

import(
	"github.com/maxzerbini/ovo/util"
	"github.com/maxzerbini/ovo/storage"
	"github.com/maxzerbini/ovo/inmemory"
	"github.com/maxzerbini/ovo/server"
	"log"
	"runtime"
	"runtime/debug"
	"flag"
)

const (
	Version = "1.0"
)

var (
	ks storage.OvoStorage
	conf server.ServerConf
	srv *server.Server
	configPath string = "./conf/serverconf.json"
	configPathTemp string = "./conf/serverconf.json.temp"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	debug.SetGCPercent(700)
	flag.StringVar(&configPath, "conf", "./conf/serverconf.json", "path of the file severconf.json")
}

// Start the Ovo Key/Value Storage
func main(){
	flag.Parse()
	configPathTemp = configPath + ".temp"
	util.Protect(start)
	util.Protect(stop)
}

// Start the server node
func start() {
	conf = server.LoadConfiguration(configPath)
	log.Printf("Start server node OVO Engine v.%s .\r\n",Version)
	conf.Init(configPathTemp)
	ks = inmemory.NewInMemoryStorage()
	srv = server.NewServer(&conf, ks)
	srv.Do()
}

func stop(){
	log.Println("Stop server node.")
}

