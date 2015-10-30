package server

import(
	"github.com/maxzerbini/ovo/storage"
	"github.com/maxzerbini/ovo/processor"
	"github.com/maxzerbini/ovo/command"
	"github.com/maxzerbini/ovo/cluster"
	"net"
	"net/rpc"
	"net/http"
	"log"
	"errors"
	"strconv"
)

// The innser server implementation. Listen for incoming commands.
type InnerServer struct {
	keystorage storage.OvoStorage
	incmdproc *processor.InCommandQueue
	outcmdproc *processor.OutCommandQueue
	config *ServerConf	
	partitioner *processor.Partitioner
}

// Creata a new inner server.
func NewInnerServer(conf *ServerConf, ks storage.OvoStorage, in *processor.InCommandQueue, out *processor.OutCommandQueue, partitioner *processor.Partitioner) *InnerServer{
	return &InnerServer{keystorage:ks, incmdproc:in, config:conf, partitioner:partitioner, outcmdproc:out}
}

// Start listening commands.
func (srv *InnerServer)Do(){
	rpc.Register(srv)
	rpc.HandleHTTP()
	listener, e := net.Listen("tcp", srv.config.ServerNode.Node.APIHost+":"+strconv.Itoa(srv.config.ServerNode.Node.APIPort))
	if e != nil {
		log.Fatal("Starting RPC-server -listen error:", e)
	}
	http.Serve(listener, nil)
}

// Enqueue a remote command.
func (srv *InnerServer) ExecuteCommand(rpccmd command.RpcCommand, reply *int) (err error) {
	defer func() {
		// Executes normally even if there is a panic
		if e:= recover(); e != nil {
			log.Println("Run time panic: %v", e)
			*reply = -1
			err = errors.New("Runtime error.")
		}
	}()
	cmd := rpccmd.Command()
	if rpccmd.OpCode == "move" {
		cmd.OpCode = "put"
		srv.incmdproc.Enqueu(cmd)
		// replicate data on twins
		srv.outcmdproc.Enqueu(cmd)
	} else {
		srv.incmdproc.Enqueu(cmd)
	}
	*reply = 0
	return nil
}

// Register a new node in the cluster.
func (srv *InnerServer) RegisterNode(node *cluster.ClusterTopologyNode, reply *cluster.ClusterTopology) (err error) {
	defer func() {
		// Executes normally even if there is a panic
		if e:= recover(); e != nil {
			log.Println("Run time panic: %v", e)
			*reply = srv.config.Topology
			err = errors.New("Runtime error.")
		}
	}()
	log.Printf("Node %s registration or update state %s\r\n", node.Node.Name, node.Node.State)
	srv.config.Topology.AddNode(node)
	srv.config.WriteTmp()
	reply.Nodes = srv.config.Topology.Nodes
	if cluster.Inactive == node.Node.State {
		srv.outcmdproc.Caller.RemoveClient(node.Node.Name)
	} 
	// start data partitioner
	go srv.partitioner.MoveData() 
	return nil
}

// Register the new node as a twin.
func (srv *InnerServer) RegisterTwin(node *cluster.ClusterTopologyNode, reply *cluster.ClusterTopology) (err error) {
	defer func() {
		// Executes normally even if there is a panic
		if e:= recover(); e != nil {
			log.Println("Run time panic: %v", e)
			*reply = srv.config.Topology
			err = errors.New("Runtime error.")
		}
	}()
	log.Printf("Node %s ask registration as twin\r\n", node.Node.Name)
	srv.config.Topology.AddTwin(node)
	srv.config.WriteTmp()
	reply.Nodes = srv.config.Topology.Nodes
	// start data partitioner
	go srv.partitioner.MoveData() 
	go srv.updateAllClusterNodes() // update the state on the other nodes
	return nil
}

// Register the new node as a stepbrother.
func (srv *InnerServer) RegisterStepbrother(node *cluster.ClusterTopologyNode, reply *cluster.ClusterTopology) (err error) {
	defer func() {
		// Executes normally even if there is a panic
		if e:= recover(); e != nil {
			log.Println("Run time panic: %v", e)
			*reply = srv.config.Topology
			err = errors.New("Runtime error.")
		}
	}()
	log.Printf("Node %s ask registration as stepbrother\r\n", node.Node.Name)
	srv.config.Topology.AddStepbrother(node)
	srv.config.WriteTmp()
	reply.Nodes = srv.config.Topology.Nodes
	// start data partitioner
	go srv.partitioner.MoveData()
	go srv.updateAllClusterNodes() // update the state on the other nodes
	return nil
}
// Register the new node as a twin and a stepbrother.
func (srv *InnerServer) RegisterTwinAndStepbrother(node *cluster.ClusterTopologyNode, reply *cluster.ClusterTopology) (err error) {
	defer func() {
		// Executes normally even if there is a panic
		if e:= recover(); e != nil {
			log.Println("Run time panic: %v", e)
			*reply = srv.config.Topology
			err = errors.New("Runtime error.")
		}
	}()
	log.Printf("Node %s ask registration as stepbrother\r\n", node.Node.Name)
	srv.config.Topology.AddTwinAndStepbrother(node)
	srv.config.WriteTmp()
	reply.Nodes = srv.config.Topology.Nodes
	// start data partitioner
	go srv.partitioner.MoveData()
	go srv.updateAllClusterNodes() // update the state on the other nodes
	return nil
}
// Merge the cluster topology configuration.
func (srv *InnerServer) UpdateTopology(topology *cluster.ClusterTopology, reply *cluster.ClusterTopology) (err error) {
	defer func() {
		// Executes normally even if there is a panic
		if e:= recover(); e != nil {
			log.Println("Run time panic: %v", e)
			*reply = srv.config.Topology
			err = errors.New("Runtime error.")
		}
	}()
	srv.config.Topology.Merge(topology)
	srv.config.WriteTmp()
	reply.Nodes = srv.config.Topology.Nodes
	// start data partitioner
	go srv.partitioner.MoveData()
	return nil
}
// Update the node state without moving data.
func (srv *InnerServer) UpdateNode(node *cluster.ClusterTopologyNode, reply *cluster.ClusterTopology) (err error) {
	defer func() {
		// Executes normally even if there is a panic
		if e:= recover(); e != nil {
			log.Println("Run time panic: %v", e)
			*reply = srv.config.Topology
			err = errors.New("Runtime error.")
		}
	}()
	log.Printf("Node %s update state %s\r\n", node.Node.Name, node.Node.State)
	srv.config.Topology.AddNode(node)
	srv.config.WriteTmp()
	reply.Nodes = srv.config.Topology.Nodes
	return nil
}
// Update all cluster node
func (srv *InnerServer) updateAllClusterNodes(){
	for _, nd := range srv.config.Topology.GetClusterNodes(){
		log.Printf("Notifies the status change to the node %s ...\r\n", nd.Node.Name)
		srv.outcmdproc.Caller.UpdateNode(srv.config.ServerNode, nd.Node)
	}
}
// Get the topology
func (srv *InnerServer) GetTopology(name *string, reply *cluster.ClusterTopology) (err error){
	log.Printf("Node %s asked topology\r\n", *name)
	reply.Nodes = srv.config.Topology.GetNodes()
	return nil
}
// Ping
func (srv *InnerServer) Ping(name *string, reply *int) (err error){
	log.Printf("Node %s calls ping\r\n", *name)
	*reply = 1
	return nil
}