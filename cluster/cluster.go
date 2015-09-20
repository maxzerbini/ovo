package cluster

import(
)

type OvoNode struct {
	Name string
	HashRange []int
	Host string
	Port int
	APIHost string
	APIPort int
	Debug bool
}

type ClusterTopologyNode struct {
	Name OvoNode
	Twins []string
}

type ClusterTopology struct {
	Nodes []ClusterTopologyNode
}

type OvoCluster interface {
	GetTwins()([]OvoNode)
}
