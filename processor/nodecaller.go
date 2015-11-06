package processor

import(
	"github.com/maxzerbini/ovo/storage"
	"github.com/maxzerbini/ovo/cluster"
	"github.com/maxzerbini/ovo/command"
	"log"
	"net/rpc"
	"strconv"
	"sync"
)

type NodeCaller struct {
	Source string
	Name string
	clients map[string]*rpc.Client
	mux *sync.RWMutex 
}

// Create the node caller
func NewNodeCaller(source string) *NodeCaller {
	nc := new(NodeCaller)
	nc.Source = source
	nc.clients = make(map[string]*rpc.Client)
	nc.mux = new (sync.RWMutex)
	return nc
}
// Add a client if it is not already in the map
func (nc *NodeCaller) addCaller(name string, client *rpc.Client) {
	nc.mux.Lock()
	defer nc.mux.Unlock()
	if _,ok:=nc.clients[name]; !ok{
		log.Printf("Create client for %s\r\n", name)
		nc.clients[name] = client
	} 
}
// Delete a client
func (nc *NodeCaller) deleteCaller(name string) {
	nc.mux.Lock()
	defer nc.mux.Unlock()
	delete(nc.clients, name) 
}
// Get the caller
func (nc *NodeCaller) getCaller(name string) (*rpc.Client, bool) {
	nc.mux.RLock()
	defer nc.mux.RUnlock()
	if val, ok := nc.clients[name]; ok {
		return val, true
	} else {
		return nil, false
	}
}

// Execute remote operation on destination server
func (nc *NodeCaller) ExecuteOperation(obj *storage.MetaDataUpdObj, destination *cluster.OvoNode, operation string) error {
	defer func() {
		// executes normally even if there is a panic
		if err2 := recover(); err2 != nil {
			//remove the client
			nc.deleteCaller(destination.Name)
		}
	}()
	var client *rpc.Client
	var ok bool
	if client, ok = nc.getCaller(destination.Name); !ok {
		client = nc.createClient(destination)
	}
	rpccmd := &command.RpcCommand{Source:"test", OpCode:operation, Obj:obj}
	var reply int = 0
	err := client.Call("InnerServer.ExecuteCommand", rpccmd, &reply)
	if err != nil {
		log.Println("InnerServer.ExecuteCommand error: ", err)
	}
	return err
}
// Ask the destination to register the node
func (nc *NodeCaller) RegisterNode(node *cluster.ClusterTopologyNode, destination *cluster.OvoNode) ( *cluster.ClusterTopology, error) {
	defer func() {
		// executes normally even if there is a panic
		if err2 := recover(); err2 != nil {
			//remove the client
			nc.deleteCaller(destination.Name)
		}
	}()
	var client *rpc.Client
	var ok bool
	if client, ok = nc.getCaller(destination.Name); !ok {
		client = nc.createClient(destination)
	}
	var topology = new (cluster.ClusterTopology)
	var err = client.Call("InnerServer.RegisterNode", node, topology)
	if err != nil {
		topology = nil
		log.Println("InnerServer.RegisterNode error: ", err)
	}
	return topology, err
}
// Ask the destination to register the node as a twin
func (nc *NodeCaller) RegisterTwin(node *cluster.ClusterTopologyNode, destination *cluster.OvoNode) ( *cluster.ClusterTopology, error) {
	defer func() {
		// executes normally even if there is a panic
		if err2 := recover(); err2 != nil {
			//remove the client
			nc.deleteCaller(destination.Name)
		}
	}()
	var client *rpc.Client
	var ok bool
	if client, ok = nc.getCaller(destination.Name); !ok {
		client = nc.createClient(destination)
	}
	var topology = new (cluster.ClusterTopology)
	var err = client.Call("InnerServer.RegisterTwin", node, topology)
	if err != nil {
		topology = nil
		log.Println("InnerServer.RegisterTwin error: ", err)
	}
	return topology, err
}
// Ask the destination to register the node as a stepbrother
func (nc *NodeCaller) RegisterStepbrother(node *cluster.ClusterTopologyNode, destination *cluster.OvoNode) ( *cluster.ClusterTopology, error) {
	defer func() {
		// executes normally even if there is a panic
		if err2 := recover(); err2 != nil {
			//remove the client
			nc.deleteCaller(destination.Name)
		}
	}()
	var client *rpc.Client
	var ok bool
	if client, ok = nc.getCaller(destination.Name); !ok {
		client = nc.createClient(destination)
	}
	var topology = new (cluster.ClusterTopology)
	var err = client.Call("InnerServer.RegisterStepbrother", node, topology)
	if err != nil {
		topology = nil
		log.Println("InnerServer.RegisterStepbrother error: ", err)
	}
	return topology, err
}
// Ask the destination to register the node as a twin and a stepbrother
func (nc *NodeCaller) RegisterTwinAndStepbrother(node *cluster.ClusterTopologyNode, destination *cluster.OvoNode) ( *cluster.ClusterTopology, error) {
	defer func() {
		// executes normally even if there is a panic
		if err2 := recover(); err2 != nil {
			//remove the client
			nc.deleteCaller(destination.Name)
		}
	}()
	var client *rpc.Client
	var ok bool
	if client, ok = nc.getCaller(destination.Name); !ok {
		client = nc.createClient(destination)
	}
	var topology = new (cluster.ClusterTopology)
	var err = client.Call("InnerServer.RegisterTwinAndStepbrother", node, topology)
	if err != nil {
		topology = nil
		log.Println("InnerServer.RegisterTwinAndStepbrother error: ", err)
	}
	return topology, err
}
// Ask the destination to give the topology
func (nc *NodeCaller) GetTopology(currentNode string, destination *cluster.OvoNode) ( *cluster.ClusterTopology, error) {
	defer func() {
		// executes normally even if there is a panic
		if err2 := recover(); err2 != nil {
			//remove the client
			nc.deleteCaller(destination.Name)
		}
	}()
	var client *rpc.Client
	var ok bool
	if client, ok = nc.getCaller(destination.Name); !ok {
		client = nc.createClient(destination)
	}
	var topology = new (cluster.ClusterTopology)
	var err = client.Call("InnerServer.GetTopology", currentNode, topology)
	if err != nil {
		topology = nil
		log.Println("InnerServer.GetTopology error: ", err)
	}
	return topology, err
}
// Ask the destination to update the topology
func (nc *NodeCaller) UpdateTopology(topology *cluster.ClusterTopology, destination *cluster.OvoNode)( *cluster.ClusterTopology, error) {
	defer func() {
		// executes normally even if there is a panic
		if err2 := recover(); err2 != nil {
			//remove the client
			nc.deleteCaller(destination.Name)
		}
	}()
	var client *rpc.Client
	var ok bool
	if client, ok = nc.getCaller(destination.Name); !ok {
		client = nc.createClient(destination)
	}
	var mergedtopology = new (cluster.ClusterTopology)
	var err = client.Call("InnerServer.UpdateTopology", topology, mergedtopology)
	if err != nil {
		mergedtopology = nil
		log.Println("InnerServer.UpdateTopology error: ", err)
	}
	return mergedtopology, err
}
// Ask the destination to update the node
func (nc *NodeCaller) UpdateNode(node *cluster.ClusterTopologyNode, destination *cluster.OvoNode) ( *cluster.ClusterTopology, error) {
	defer func() {
		// executes normally even if there is a panic
		if err2 := recover(); err2 != nil {
			//remove the client
			nc.deleteCaller(destination.Name)
		}
	}()
	var client *rpc.Client
	var ok bool
	if client, ok = nc.getCaller(destination.Name); !ok {
		client = nc.createClient(destination)
	}
	var topology = new (cluster.ClusterTopology)
	var err = client.Call("InnerServer.UpdateNode", node, topology)
	if err != nil {
		topology = nil
		log.Println("InnerServer.UpdateNode error: ", err)
	}
	return topology, err
}
// Ask the destination to give the topology
func (nc *NodeCaller) Ping(currentNode string, destination *cluster.OvoNode) error {
	defer func() {
		// executes normally even if there is a panic
		if err2 := recover(); err2 != nil {
			//remove the client
			nc.deleteCaller(destination.Name)
		}
	}()
	var client *rpc.Client
	var ok bool
	if client, ok = nc.getCaller(destination.Name); !ok {
		client = nc.createClient(destination)
	}
	var reply int = 0
	var err = client.Call("InnerServer.Ping", currentNode, &reply)
	if err != nil {
		log.Println("InnerServer.Ping error: ", err)
	}
	return err
}
// Remove a client by name
func (nc *NodeCaller) RemoveClient(name string){
	delete (nc.clients, name)
}
// Create the client.
func (nc *NodeCaller) createClient(destination *cluster.OvoNode) *rpc.Client{
	defer func() {
		// Println executes normally even if there is a panic
		if err := recover(); err != nil {
			log.Println("run time panic: %v", err)
		}
	}()
	client, err := rpc.DialHTTP("tcp", destination.APIHost + ":"+strconv.Itoa(destination.APIPort))
	if err != nil {
		log.Printf("dialing: %v \r\n", err)
		return nil
	} else {
		nc.addCaller(destination.Name, client)
		return client
	}
}