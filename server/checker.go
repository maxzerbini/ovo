package server

import(
	"log"
	"time"
	"github.com/maxzerbini/ovo/cluster"
	"github.com/maxzerbini/ovo/processor"
)

const(
	CheckerPeriod = 5 // 5 secs
	ErrorThreshold = 3
)

type Checker struct {
	cnf *ServerConf
	topology *cluster.ClusterTopology
	outcomingQueue *processor.OutCommandQueue
	doneChan chan(bool)
	nodeError map[string]int
	partitioner *processor.Partitioner
}

func NewChecker(cnf *ServerConf, outcomingQueue *processor.OutCommandQueue,partitioner *processor.Partitioner) *Checker {
	return &Checker{cnf:cnf,topology:&cnf.Topology, outcomingQueue:outcomingQueue,doneChan:make(chan bool),nodeError:make(map[string]int,0),partitioner:partitioner}
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
	somethingheppens := false
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
	// notify faults to other nodes
	for name,count := range ckr.nodeError {
		if count >= ErrorThreshold {
			ckr.notifyFaultNotification(name)
			somethingheppens = true
		}
	}
	// clean inactive nodes
	for _,nd := range ckr.topology.GetInactiveNodes() {
		if nd.UpdateDate.Before(time.Now().Add(-1 * time.Minute)){
			ckr.topology.RemoveNode(nd.Node.Name)
			somethingheppens = true
		}
	}
	if somethingheppens {
		// write conf
		ckr.cnf.WriteTmp()
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