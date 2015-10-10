package server

import (
	//"github.com/maxzerbini/ovo/cluster"
	"testing"
)
/*
func TestConfigurationWrite(t *testing.T) {
	t.Log("TestConfigurationWrite started")
	var conf ServerConf
	var node = &cluster.OvoNode{Name:"mizard",Host:"0.0.0.0",Port:5050,APIHost:"0.0.0.0",APIPort:5052,Debug:true}
	var topoNode = &cluster.ClusterTopologyNode{Node:*node}
	topoNode.Twins = append(topoNode.Twins, "righel") 
	conf.ServerNode = *topoNode 
	var node2 = &cluster.OvoNode{Name:"righel",Host:"127.0.0.1",Port:5060,APIHost:"127.0.0.1",APIPort:5062,Debug:true}
	var topoNode2 = &cluster.ClusterTopologyNode{Node:*node2} 
	topoNode2.Twins = append(topoNode2.Twins, "mizard")
	conf.Topology = cluster.ClusterTopology{}
	conf.Topology.Nodes = append(conf.Topology.Nodes, *topoNode2)
	WriteConfiguration("../conf/serverconf.json", conf)
}
*/

func TestConfigurationLoad(t *testing.T) {
	t.Log("TestConfigurationLoad started")
	var conf = LoadConfiguration("../conf/serverconf.json")
	t.Logf("conf = %s", conf)
	conf.Init()
}