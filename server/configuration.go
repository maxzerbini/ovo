package server

import (
	"time"
	"github.com/maxzerbini/ovo/cluster"
	"encoding/json"
    "io/ioutil"
	"os"
	"log"
)

const CONF_PATH string = "./conf/severconf.json"

type ServerConf struct {
	ServerNode *cluster.ClusterTopologyNode
	Topology cluster.ClusterTopology
	tmpPath string
}

func ( cnf *ServerConf) Init(tmpPath string) { 
	cnf.ServerNode.StartDate = time.Now()
	cnf.ServerNode.Node.State = cluster.Active
	if cnf.ServerNode.Twins == nil { cnf.ServerNode.Twins = make([]string,0)}
	if cnf.ServerNode.Stepbrothers == nil { cnf.ServerNode.Stepbrothers = make([]string,0)}
	cnf.ServerNode.UpdateDate = time.Now()
	cluster.SetCurrentNode(cnf.ServerNode, &cnf.Topology)
	cnf.tmpPath = tmpPath
}

func ( cnf *ServerConf) WriteTmp() { 
	WriteConfiguration(cnf.tmpPath, cnf)
}

func LoadConfiguration(path string) ServerConf {
	file, e := ioutil.ReadFile(path)
    if e != nil {
		log.Fatalf("Configuration file not found at %s", path)
        os.Exit(1)
    }
    var jsontype ServerConf
    json.Unmarshal(file, &jsontype)
	return jsontype;
}

func WriteConfiguration(path string, conf *ServerConf) {
	data, _ := json.Marshal(conf)
	e := ioutil.WriteFile(path, data, 0x666)
    if e != nil {
		log.Printf("Configuration file write error at %s\r\n", path)
    }
}
