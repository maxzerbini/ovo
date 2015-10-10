package processor

import (
	"github.com/maxzerbini/ovo/storage"
	"github.com/maxzerbini/ovo/cluster"
	"github.com/maxzerbini/ovo/command"
	"github.com/maxzerbini/ovo/util"
)

func NewPartitioner(storage storage.OvoStorage, serverNode *cluster.ClusterTopologyNode, outcomingQueue *OutCommandQueue) *Partitioner {
	return &Partitioner{storage:storage, serverNode:serverNode, outcomingQueue:outcomingQueue}
}

type Partitioner struct {
	storage storage.OvoStorage
	serverNode *cluster.ClusterTopologyNode
	outcomingQueue *OutCommandQueue
}

func (p *Partitioner) MoveData(){
	var list = p.storage.List()
	for _, obj := range list {
		if !util.Contains(p.serverNode.Node.HashRange, obj.Hash) {
			p.outcomingQueue.Enqueu(&command.Command{OpCode:"move",Obj:obj.MetaDataUpdObj()})	
		}
	}
}