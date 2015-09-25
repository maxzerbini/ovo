package processor

import (
	"github.com/maxzerbini/ovo/storage"
	"github.com/maxzerbini/ovo/command"
	"github.com/maxzerbini/ovo/cluster"
	"time"
)

type OutCommandQueue struct {
	commands chan *command.Command
	errors chan *commandError
	serverNode *cluster.ClusterTopologyNode
	topology *cluster.ClusterTopology
	caller *NodeCaller
}

// Create the outcoming command processor queue
func NewOutCommandQueue(serverNode *cluster.ClusterTopologyNode, topology *cluster.ClusterTopology) *OutCommandQueue {
	cq := new(OutCommandQueue)
	cq.commands = make(chan *command.Command, commands_buffer_size)
	cq.errors = make(chan *commandError, commands_buffer_size)
	cq.serverNode = serverNode
	cq.topology = topology
	cq.caller = NewNodeCaller(serverNode.Node.Name)
	go cq.backend()
	go cq.errorBackend()
	return cq
}

func (cq *OutCommandQueue) Enqueu(cmd *command.Command){
	cq.commands <- cmd
}

func (cq *OutCommandQueue) backend() {
	for  cmd := range cq.commands {
		if (cmd != nil){
			switch cmd.OpCode {
				case "put" : cq.execute(cmd.Obj, cmd.OpCode)
				case "delete" : cq.execute(cmd.Obj, cmd.OpCode)
				case "touch" : cq.execute(cmd.Obj, cmd.OpCode)
				case "updatevalue" : cq.execute(cmd.Obj, cmd.OpCode)
				case "updatekey" : cq.execute(cmd.Obj, cmd.OpCode)
				case "updatekeyvalue" : cq.execute(cmd.Obj, cmd.OpCode)
				case "move" : cq.move(cmd.Obj)
				default : println("usupported command: "+cmd.OpCode)
			}	
		}
	}
}

func (cq *OutCommandQueue) execute(obj *storage.MetaDataUpdObj, operation string){
	for _, node := range cq.topology.GetTwins(cq.serverNode.Twins){
		err := cq.caller.ExecuteOperation(obj, &node.Node, operation)
		if err != nil {
			cq.enqueuError(newCommandError(obj, &node.Node, operation))
		}
	}
}

func (cq *OutCommandQueue) move(obj *storage.MetaDataUpdObj){
	if node := cq.topology.GetNodeByHash(obj.Hash); node != nil {
		err := cq.caller.ExecuteOperation(obj, &node.Node, "put")
		if err != nil {
			cq.enqueuError(newCommandError(obj, &node.Node, "put"))
		}
	}
}

func (cq *OutCommandQueue) enqueuError(cmd *commandError){
	go func(){
		cmd.count++
		time.Sleep(5000000000) // 5 secs
		cq.errors <- cmd
	}()
}

type commandError struct {
	obj *storage.MetaDataUpdObj
	operation string
	destination *cluster.OvoNode
	count int
}

func newCommandError (o *storage.MetaDataUpdObj, dest *cluster.OvoNode, op string) *commandError{
	return &commandError{obj:o, operation:op, destination:dest}
}

func (cq *OutCommandQueue) errorBackend() {
	for  cmd := range cq.errors {
		if cmd != nil{
			if cmd.count < 4 {
				err := cq.caller.ExecuteOperation(cmd.obj, cmd.destination, cmd.operation)
				if err != nil {
					cq.enqueuError(cmd)
				}
			} else {
				// TODO remove node from topology
			}
		}
	}
}