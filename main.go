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

var ks storage.OvoStorage
var conf server.ServerConf
var srv *server.Server
var configPath string = "./conf/serverconf.json"
var configPathTemp string = "./conf/serverconf.json.temp"

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
	conf.Init(configPathTemp)
	log.Println("Start server node.")
	ks = inmemory.NewInMemoryStorage()
	srv = server.NewServer(&conf, ks)
	srv.Do()
}

func stop(){
	log.Println("Stop server node.")
}

