package processor

import(
	"log"
	"time"
	"github.com/maxzerbini/ovo/cluster"
)

const(
	CheckerPeriod = 5 // 5 secs
	ErrorThreshold = 3
)

type Checker struct {
	topology *cluster.ClusterTopology
	outcomingQueue *OutCommandQueue
	doneChan chan(bool)
	nodeError map[string]int
	partitioner *Partitioner
}

func NewChecker(topology *cluster.ClusterTopology, outcomingQueue *OutCommandQueue,partitioner *Partitioner) *Checker {
	return &Checker{topology:topology, outcomingQueue:outcomingQueue,doneChan:make(chan bool),nodeError:make(map[string]int,0),partitioner:partitioner}
}

func (ckr *Checker) Stop(){
	ckr.doneChan <- true
}

func (ckr *Checker) Do(){
	tickChan := time.NewTicker(time.Second * CheckerPeriod).C
	log.Printf("Start checking cluster node (relatives)...\r\n")
	for {
        select {
        	case <- tickChan: ckr.checkNodes()
			case <- ckr.doneChan: return
		}
	}
}

func (ckr *Checker) checkNodes(){
	nodes := ckr.topology.GetRelatives()
	for _,nd := range nodes {
		if _, ok := ckr.nodeError[nd.Node.Name]; !ok {
			    ckr.nodeError[nd.Node.Name] = 0
		}
		if err := ckr.outcomingQueue.Caller.Ping(cluster.GetCurrentNode().Node.Name, nd.Node); err != nil {
			ckr.nodeError[nd.Node.Name] = ckr.nodeError[nd.Node.Name] + 1
		} else {
			ckr.nodeError[nd.Node.Name] = 0
		}
	}
	for name,count := range ckr.nodeError {
		if count >= ErrorThreshold {
			ckr.notifyFaultNotification(name)
		}
	}
}

func (ckr *Checker) notifyFaultNotification(name string){
	node,_ := ckr.topology.GetNodeByName(name)
	node.Node.State = cluster.Inactive
	node.UpdateDate = time.Now()
	node.Node.HashRange = make([]int,0)
	nodes := ckr.topology.GetClusterNodes()
	ckr.topology.AddNode(node) // repartition index
	ckr.outcomingQueue.Caller.RemoveClient(name)
	for _,nd := range nodes {
		if nd.Node.Name != name {
			if _,err := ckr.outcomingQueue.Caller.RegisterNode(node, nd.Node); err != nil {
				log.Printf("Call node %s fail\r\n", nd.Node.Name)
			} else {
				log.Printf("Notify node %s of fault of node %s\r\n", nd.Node.Name, name)
			}
		} 
	}
	ckr.nodeError[name] = 0
	go ckr.partitioner.MoveData()
}