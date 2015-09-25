package processor

import (
	"github.com/maxzerbini/ovo/storage"
	"github.com/maxzerbini/ovo/command"
	"github.com/maxzerbini/ovo/cluster"
)

type OutCommandQueue struct {
	commands chan *command.Command
	serverNode *cluster.ClusterTopologyNode
	topology *cluster.ClusterTopology
	caller *NodeCaller
}

// Create the outcoming command processor queue
func NewOutCommandQueue(serverNode *cluster.ClusterTopologyNode, topology *cluster.ClusterTopology) *OutCommandQueue {
	cq := new(OutCommandQueue)
	cq.commands = make(chan *command.Command, commands_buffer_size)
	cq.serverNode = serverNode
	cq.topology = topology
	cq.caller = NewNodeCaller(serverNode.Node.Name)
	go cq.backend()
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
		cq.caller.ExecuteOperation(obj, &node.Node, operation)
	}
}

func (cq *OutCommandQueue) move(obj *storage.MetaDataUpdObj){
	if node := cq.topology.GetNodeByHash(obj.Hash); node != nil {
		cq.caller.ExecuteOperation(obj, &node.Node, "put")
	}
}

