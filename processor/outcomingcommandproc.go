package processor

import (
	"github.com/maxzerbini/ovo/cluster"
	"github.com/maxzerbini/ovo/command"
	"github.com/maxzerbini/ovo/storage"
	"github.com/maxzerbini/ovo/util"
	"time"
)

type OutCommandQueue struct {
	commands      chan *command.Command
	errors        chan *commandError
	serverNode    *cluster.ClusterTopologyNode
	topology      *cluster.ClusterTopology
	Caller        *NodeCaller
	incomingQueue *InCommandQueue
}

// Create the outcoming command processor queue
func NewOutCommandQueue(serverNode *cluster.ClusterTopologyNode, topology *cluster.ClusterTopology, incomingQueue *InCommandQueue) *OutCommandQueue {
	cq := new(OutCommandQueue)
	cq.commands = make(chan *command.Command, commands_buffer_size)
	cq.errors = make(chan *commandError, commands_buffer_size)
	cq.serverNode = serverNode
	cq.topology = topology
	cq.incomingQueue = incomingQueue
	cq.Caller = NewNodeCaller(serverNode.Node.Name)
	go cq.backend()
	go cq.errorBackend()
	return cq
}

func (cq *OutCommandQueue) Enqueu(cmd *command.Command) {
	cq.commands <- cmd
}

func (cq *OutCommandQueue) backend() {
	for cmd := range cq.commands {
		if cmd != nil {
			switch cmd.OpCode {
			case "put":
				cq.execute(cmd.Obj, cmd.OpCode)
			case "delete":
				cq.execute(cmd.Obj, cmd.OpCode)
			case "touch":
				cq.execute(cmd.Obj, cmd.OpCode)
			case "updatevalue":
				cq.execute(cmd.Obj, cmd.OpCode)
			case "updatekey":
				cq.executeUpdateKey(cmd.Obj, cmd.OpCode)
			case "updatekeyvalue":
				cq.executeUpdateKey(cmd.Obj, cmd.OpCode)
			case "move":
				cq.move(cmd.Obj)
			case "setcounter":
				cq.execute(cmd.Obj, cmd.OpCode)
			case "deletecounter":
				cq.execute(cmd.Obj, cmd.OpCode)
			case "movecounter":
				cq.moveCounter(cmd.Obj)
			default:
				println("usupported command: " + cmd.OpCode)
			}
		}
	}
}

func (cq *OutCommandQueue) execute(obj *storage.MetaDataUpdObj, operation string) {
	for _, node := range cq.topology.GetTwins(cq.serverNode.Twins) {
		err := cq.Caller.ExecuteOperation(obj, node.Node, operation)
		if err != nil {
			cq.enqueuError(newCommandError(obj, node.Node, operation))
		}
	}
}

func (cq *OutCommandQueue) executeUpdateKey(obj *storage.MetaDataUpdObj, operation string) {
	if !util.Contains(cq.serverNode.Node.HashRange, obj.NewHash) {
		// delete the data on the twins
		for _, node := range cq.topology.GetTwins(cq.serverNode.Twins) {
			err := cq.Caller.ExecuteOperation(obj, node.Node, "delete")
			if err != nil {
				cq.enqueuError(newCommandError(obj, node.Node, "delete"))
			}
		}
		// move the data because the new hashcode does not belong to this node
		cq.move(obj)
	} else {
		// update data on the twins
		for _, node := range cq.topology.GetTwins(cq.serverNode.Twins) {
			err := cq.Caller.ExecuteOperation(obj, node.Node, operation)
			if err != nil {
				cq.enqueuError(newCommandError(obj, node.Node, operation))
			}
		}
	}
}

func (cq *OutCommandQueue) move(obj *storage.MetaDataUpdObj) {
	if node := cq.topology.GetNodeByHash(obj.Hash); node != nil {
		err := cq.Caller.ExecuteOperation(obj, node.Node, "move")
		if err != nil {
			cq.enqueuError(newCommandError(obj, node.Node, "move"))
		} else {
			cq.incomingQueue.Enqueu(&command.Command{OpCode: "delete", Obj: obj})
		}
	}
}

func (cq *OutCommandQueue) moveCounter(obj *storage.MetaDataUpdObj) {
	if node := cq.topology.GetNodeByHash(obj.Hash); node != nil {
		err := cq.Caller.ExecuteOperation(obj, node.Node, "setcounter")
		if err != nil {
			cq.enqueuError(newCommandError(obj, node.Node, "setcounter"))
		} else {
			cq.incomingQueue.Enqueu(&command.Command{OpCode: "deletecounter", Obj: obj})
		}
	}
}

func (cq *OutCommandQueue) enqueuError(cmd *commandError) {
	go func() {
		cmd.count++
		time.Sleep(5000000000) // 5 secs
		cq.errors <- cmd
	}()
}

type commandError struct {
	obj         *storage.MetaDataUpdObj
	operation   string
	destination *cluster.OvoNode
	count       int
}

func newCommandError(o *storage.MetaDataUpdObj, dest *cluster.OvoNode, op string) *commandError {
	return &commandError{obj: o, operation: op, destination: dest}
}

func (cq *OutCommandQueue) errorBackend() {
	for cmd := range cq.errors {
		if cmd != nil {
			if cmd.count < 4 {
				err := cq.Caller.ExecuteOperation(cmd.obj, cmd.destination, cmd.operation)
				if err != nil {
					cq.enqueuError(cmd)
				}
			} else {
				// TODO remove node from topology
			}
		}
	}
}
