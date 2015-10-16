package cluster

import (
	"time"
	"log"
)

import(
	"github.com/maxzerbini/ovo/util"
)

const MaxNodeNumber = 128

// Node configuration informations
type OvoNode struct {
	Name string
	HashRange []int
	Host string
	Port int
	APIHost string
	APIPort int
	Debug bool
}
// The cluster node struct
type ClusterTopologyNode struct {
	Node OvoNode
	StartDate time.Time
	Twins []string
}
// The cluster topology that contains the list of nodes
type ClusterTopology struct {
	Nodes []*ClusterTopologyNode
}
// Get the twins node
func (ct *ClusterTopology) GetTwins(names []string)(nodes []*ClusterTopologyNode){
	nodes = make([]*ClusterTopologyNode,0)
	for _,nd := range ct.Nodes {
		for _,s := range names {
			if nd.Node.Name == s {
				nodes = append(nodes, nd)
			}
		}
	}
	return nodes
}
// Get the node that contains the hashcode
func (ct *ClusterTopology) GetNodeByHash(hash int)(node *ClusterTopologyNode){
	for _,nd := range ct.Nodes {
		if util.Contains(nd.Node.HashRange, hash) {
			return nd
		}
	}
	return nil
}
// Get node by name
func (ct *ClusterTopology) GetNodeByName(name string)(node *ClusterTopologyNode, ind int){
	for ind,nd := range ct.Nodes {
		if nd.Node.Name == name {
			return nd, ind
		}
	}
	return nil, 0
}
// Add or update a node in the topology ordering the topology by node's startdate
func (ct *ClusterTopology) AddNode(node *ClusterTopologyNode){
	if nd, ind := ct.GetNodeByName(node.Node.Name); nd != nil {
		ct.Nodes = append(ct.Nodes[:ind], ct.Nodes[ind+1:]...) //remove node if already present
	}
	if len(ct.Nodes)==0 {
		ct.Nodes = append(ct.Nodes, node)
	} else {
		var ind int = 0
		var nd *ClusterTopologyNode
		for _,nd = range ct.Nodes {
			if nd.StartDate.After(node.StartDate) {
				break
			}
			ind++
		}
		ct.Nodes = append(ct.Nodes[:ind], append([]*ClusterTopologyNode{node},ct.Nodes[ind:]...)...)
	}
	ct.buildHashcode()
}
// Remove a node from the topology
func (ct *ClusterTopology) RemoveNode(nodeName string){
	if nd, ind := ct.GetNodeByName(nodeName); nd != nil {
		ct.Nodes = append(ct.Nodes[:ind], ct.Nodes[ind+1:]...) //remove node if already present
		ct.buildHashcode()
	}
}
// Merge the topology nodes with this topology
func (ct *ClusterTopology) Merge(topology *ClusterTopology){
	log.Println("Merging topologies ...")
	for _, node := range topology.Nodes {
		log.Printf("Evaluating node %s ...\r\n", node.Node.Name)
		ct.AddNode(node)
	}
}
// Generate and assign the hashcode range to all nodes.
func (ct *ClusterTopology) buildHashcode(){
	log.Println("Partitioning hashcode...")
	var r int = MaxNodeNumber
	if len(ct.Nodes) > 0 {
		r = MaxNodeNumber / len(ct.Nodes)
	}
	//log.Printf(" r = %d \r\n",r)
	var hashrange []int = make([]int, MaxNodeNumber)
	for i :=0; i<MaxNodeNumber; i++ {
		hashrange[i] = i
	}
	for ind, node := range ct.Nodes {
		if ind < (len(ct.Nodes) -1) {
			log.Printf(" node = %s : start = %d - end = %d\r\n", node.Node.Name, ind*r, (ind*(r)+r)-1)
			node.Node.HashRange = hashrange[ind*r:(ind*(r)+r)-1]
		} else {
			log.Printf(" node = %s : start = %d - end = %d\r\n", node.Node.Name, ind*r, len(hashrange)-1)
			node.Node.HashRange = hashrange[ind*r:]
		}
	}
}

