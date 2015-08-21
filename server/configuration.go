package server

import (
	"github.com/maxzerbini/ovo/cluster"
)



type ServerConf struct {
	Node cluster.OvoNode
	Topology cluster.ClusterTopology
}