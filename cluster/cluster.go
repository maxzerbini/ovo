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
		if contains(nd.Node.HashRange, hash) {
			return &nd
		}
	}
	return nil
}

func contains(s []int, e int) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}

