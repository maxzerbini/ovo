package processor

import (
	"log"
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
	log.Printf("Partitioner is moving data (storage size = %d)\r\n", len(list))
	for _, obj := range list {
		if (obj != nil) {
			if !util.Contains(p.serverNode.Node.HashRange, obj.Hash) {
				log.Printf("Moving key = %s\r\n", obj.Key)
				p.outcomingQueue.Enqueu(&command.Command{OpCode:"move",Obj:obj.MetaDataUpdObj()})	
			}
		}
	}
}

func (p *Partitioner) MoveObject(obj *storage.MetaDataObj){
	if (obj != nil) {
		if !util.Contains(p.serverNode.Node.HashRange, obj.Hash) {
			log.Printf("Moving key = %s\r\n", obj.Key)
			p.outcomingQueue.Enqueu(&command.Command{OpCode:"move",Obj:obj.MetaDataUpdObj()})	
		}
	}
}