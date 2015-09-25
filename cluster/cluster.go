package cluster

import(
	"github.com/maxzerbini/ovo/util"
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
	Node OvoNode
	Twins []string
}

type ClusterTopology struct {
	Nodes []ClusterTopologyNode
}


func (ct *ClusterTopology) GetTwins(names []string)(nodes []ClusterTopologyNode){
	nodes = make([]ClusterTopologyNode,0)
	for _,nd := range ct.Nodes {
		for _,s := range names {
			if nd.Node.Name == s {
				nodes = append(nodes, nd)
			}
		}
	}
	return nodes
}

func (ct *ClusterTopology) GetNodeByHash(hash int)(node *ClusterTopologyNode){
	for _,nd := range ct.Nodes {
		if util.Contains(nd.Node.HashRange, hash) {
			return &nd
		}
	}
	return nil
}



