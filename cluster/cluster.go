package cluster

import (
	"time"
	"log"
	"sync"
	"github.com/maxzerbini/ovo/util"
)

const (
	MaxNodeNumber = 128
	Active = "ACTIVE"
	Inactive= "INACTIVE"
)

// Node configuration informations
type OvoNode struct {
	Name string
	HashRange []int
	Host string
	ExtHost string
	Port int
	APIHost string
	APIPort int
	State string
}
// The cluster node struct
type ClusterTopologyNode struct {
	Node *OvoNode
	StartDate time.Time
	Twins []string
	Stepbrothers []string
	UpdateDate time.Time
}
// The cluster topology that contains the list of nodes
type ClusterTopology struct {
	Nodes []*ClusterTopologyNode
}

var currentNode *ClusterTopologyNode
var mux *sync.RWMutex = new(sync.RWMutex)
// Set the current node
func SetCurrentNode(node *ClusterTopologyNode, ct *ClusterTopology){
	currentNode = node
	ct.AddNode(currentNode)
}
// Get the current node
func GetCurrentNode()*ClusterTopologyNode{
	return currentNode
}
// Get the active twin nodes
func (ct *ClusterTopology) GetTwins(names []string)(nodes []*ClusterTopologyNode){
	nodes = make([]*ClusterTopologyNode,0)
	mux.RLock()
	defer mux.RUnlock()
	for _,nd := range ct.Nodes {
		for _,s := range names {
			if nd.Node.Name == s && Active == nd.Node.State {
				nodes = append(nodes, nd)
			}
		}
	}
	return nodes
}
// Get the all active nodes
func (ct *ClusterTopology) GetNodes()(nodes []*ClusterTopologyNode){
	nodes = make([]*ClusterTopologyNode,0)
	mux.RLock()
	defer mux.RUnlock()
	for _,nd := range ct.Nodes {
		if Active == nd.Node.State {
			nodes = append(nodes, nd)
		}
	}
	return nodes
}
// Get the all inactive nodes
func (ct *ClusterTopology) GetInactiveNodes()(nodes []*ClusterTopologyNode){
	nodes = make([]*ClusterTopologyNode,0)
	mux.RLock()
	defer mux.RUnlock()
	for _,nd := range ct.Nodes {
		if Active != nd.Node.State {
			nodes = append(nodes, nd)
		}
	}
	return nodes
}
// Get the all active cluster nodes except current
func (ct *ClusterTopology) GetClusterNodes()(nodes []*ClusterTopologyNode){
	nodes = make([]*ClusterTopologyNode,0)
	mux.RLock()
	defer mux.RUnlock()
	for _,nd := range ct.Nodes {
		if nd.Node.Name != currentNode.Node.Name && Active == nd.Node.State {
			nodes = append(nodes, nd)
		}
	}
	return nodes
}
// Get the relative nodes
func (ct *ClusterTopology) GetRelatives()(nodes []*ClusterTopologyNode){
	nodemap := make(map[string]*ClusterTopologyNode,0)
	nodes = make([]*ClusterTopologyNode,0)
	mux.RLock()
	defer mux.RUnlock()
	for _,nd := range ct.Nodes {
		if util.ContainsString(currentNode.Twins, nd.Node.Name) || util.ContainsString(currentNode.Stepbrothers, nd.Node.Name) {
			if Active == nd.Node.State{
				nodemap[nd.Node.Name] = nd
			}
		}
	}
	for _,nd := range nodemap {
		nodes = append(nodes, nd)
	}
	return nodes
}
// Get the node that contains the hashcode
func (ct *ClusterTopology) GetNodeByHash(hash int)(node *ClusterTopologyNode){
	mux.RLock()
	defer mux.RUnlock()
	for _,nd := range ct.Nodes {
		if util.Contains(nd.Node.HashRange, hash) {
			return nd
		}
	}
	
	return nil
}
// Get node by name
func (ct *ClusterTopology) GetNodeByName(name string)(node *ClusterTopologyNode, ind int){
	mux.RLock()
	defer mux.RUnlock()
	return ct.getNodeByName(name)
}
// same func for internal use 
func (ct *ClusterTopology) getNodeByName(name string)(node *ClusterTopologyNode, ind int){
	for ind,nd := range ct.Nodes {
		if nd.Node.Name == name {
			return nd, ind
		}
	}
	return nil, 0
}
// Add or update a node in the topology ordering the topology by node's startdate and rebuilding the hahscode
func (ct *ClusterTopology) AddNode(node *ClusterTopologyNode){
	ct.addNode(node)
	ct.buildHashcode()
}
// Add or update a twin in the topology ordering the topology by node's startdate and rebuilding the hahscode
func (ct *ClusterTopology) AddTwin(node *ClusterTopologyNode){
	ct.addNode(node)
	ct.buildHashcode()
	currentNode.Twins = append(currentNode.Twins, node.Node.Name)
	currentNode.UpdateDate = time.Now()
}
// Add or update a stepbrother in the topology ordering the topology by node's startdate and rebuilding the hahscode
func (ct *ClusterTopology) AddStepbrother(node *ClusterTopologyNode){
	ct.addNode(node)
	ct.buildHashcode()
	currentNode.Stepbrothers = append(currentNode.Stepbrothers, node.Node.Name)
	currentNode.UpdateDate = time.Now()
}
// Add or update a twin and stepbrother in the topology ordering the topology by node's startdate and rebuilding the hahscode
func (ct *ClusterTopology) AddTwinAndStepbrother(node *ClusterTopologyNode){
	ct.addNode(node)
	ct.buildHashcode()
	currentNode.Twins = append(currentNode.Twins, node.Node.Name)
	currentNode.Stepbrothers = append(currentNode.Stepbrothers, node.Node.Name)
	currentNode.UpdateDate = time.Now()
}
// Remove a node from the topology
func (ct *ClusterTopology) RemoveNode(nodeName string){
	if res := ct.removeNode(nodeName); res {
		log.Printf("Node %s REMOVED.\r\n", nodeName)
		//remove twins & stepbrother
		for _,nd := range ct.GetNodes(){
			nd.Stepbrothers = util.RemoveElement(nd.Stepbrothers, nodeName)
			nd.Twins = util.RemoveElement(nd.Twins, nodeName)
		}
		ct.buildHashcode()
	}
}
// internal remove
func (ct *ClusterTopology) removeNode(nodeName string) bool {
	if nd, ind := ct.GetNodeByName(nodeName); nd != nil {
		mux.Lock()
		defer mux.Unlock()
		ct.Nodes = append(ct.Nodes[:ind], ct.Nodes[ind+1:]...) //remove node if already present
		return true
	} else {
		return false
	}
}
// Merge the topology nodes with this topology
func (ct *ClusterTopology) Merge(topology *ClusterTopology){
	log.Println("Merging topologies ...")
	for _, node := range topology.Nodes {
		log.Printf("Evaluating node %s ...\r\n", node.Node.Name)
		ct.addNode(node)
	}
	ct.AddNode(currentNode) // restore current node if needed
	ct.buildHashcode()
}
// Add or update a node in the topology ordering the topology by node's startdate
func (ct *ClusterTopology) addNode(node *ClusterTopologyNode){ 
	mux.Lock()
	defer mux.Unlock()
	if nd, ind := ct.getNodeByName(node.Node.Name); nd != nil {
		if node.UpdateDate.After(nd.UpdateDate) {
			ct.Nodes = append(ct.Nodes[:ind], ct.Nodes[ind+1:]...) //remove node if already present
		} else {
			// do nothing
			return
		}
	}
	//add node
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
}
// Generate and assign the hashcode range to all nodes.
func (ct *ClusterTopology) buildHashcode(){
	mux.RLock()
	defer mux.RUnlock()
	count := 0
	for _, node := range ct.Nodes {
		if Active == node.Node.State {
			count++
		}
	}
	log.Println("Partitioning hashcode...")
	var r int = MaxNodeNumber
	var q int = 0
	if count > 0 {
		r = MaxNodeNumber / count
		q = MaxNodeNumber % count
		log.Printf("Range len = %d - Remnant = %d\r\n",r,q)
	}
	//log.Printf(" r = %d \r\n",r)
	var hashrange []int = make([]int, MaxNodeNumber)
	for i :=0; i<MaxNodeNumber; i++ {
		hashrange[i] = i
	}
	var start int = 0
	var end int = r - 1
	for _, node := range ct.Nodes {
		if Active == node.Node.State {
			
			if end < (len(hashrange) - 1) {
				if q > 0 { 
					end++
					q--
				}
				node.Node.HashRange = hashrange[start:end+1]
				log.Printf(" node = %s : start = %d - end = %d\r\n", node.Node.Name, start, end)
				start = end + 1
				end += r
			} else {
				node.Node.HashRange = hashrange[start:]
				log.Printf(" node = %s : start = %d - end = %d\r\n", node.Node.Name, start, end)
			}
		}
	}
	/* *
	for ind, node := range ct.Nodes {
		if Active == node.Node.State {
			if ind < (len(ct.Nodes) -1) {
				log.Printf(" node = %s : start = %d - end = %d\r\n", node.Node.Name, ind*r, (ind*(r)+r)-1)
				node.Node.HashRange = hashrange[ind*r:(ind*(r)+r)-1]
			} else {
				log.Printf(" node = %s : start = %d - end = %d\r\n", node.Node.Name, ind*r, len(hashrange)-1)
				node.Node.HashRange = hashrange[ind*r:]
			}
		}
	}
	* */
}

