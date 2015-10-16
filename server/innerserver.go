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
	config *ServerConf	
	partitioner *processor.Partitioner
}

// Creata a new inner server.
func NewInnerServer(conf *ServerConf, ks storage.OvoStorage, in *processor.InCommandQueue, out *processor.OutCommandQueue, partitioner *processor.Partitioner) *InnerServer{
	return &InnerServer{keystorage:ks, incmdproc:in, config:conf, partitioner:partitioner}
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
	srv.incmdproc.Enqueu(rpccmd.Command())
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
	srv.config.Topology.AddNode(node)
	srv.config.WriteTmp()
	reply.Nodes = srv.config.Topology.Nodes
	// start data partitioner
	go srv.partitioner.MoveData() 
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
	for _,node:= range topology.Nodes {
		srv.config.Topology.AddNode(node)
	}
	*reply = srv.config.Topology
	// start data partitioner
	go srv.partitioner.MoveData()
	return nil
}
