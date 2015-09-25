package processor
import(
	"github.com/maxzerbini/ovo/storage"
	"github.com/maxzerbini/ovo/cluster"
	"github.com/maxzerbini/ovo/command"
	"log"
	"net/rpc"
	"strconv"
)

type NodeCaller struct {
	Name string
	clients map[string]*rpc.Client
}

// Create the node caller
func NewNodeCaller() *NodeCaller {
	nc := new(NodeCaller)
	nc.clients = make(map[string]*rpc.Client)
	return nc
}

// Execute remote operation on destination server
func (nc *NodeCaller) ExecuteOperation(obj *storage.MetaDataUpdObj, destination *cluster.OvoNode, operation string){
	defer func() {
		// Println executes normally even if there is a panic
		if err := recover(); err != nil {
			log.Println("run time panic: %v", err)
		}
	}()
	var client *rpc.Client
	var ok bool
	if client, ok = nc.clients[destination.Name]; !ok {
		client = nc.createClient(destination)
	}
	rpccmd := &command.RpcCommand{Source:"test", OpCode:operation, Obj:obj}
	var reply int = 0
	err := client.Call("InnerServer.ExecuteCommand", rpccmd, &reply)
	if err != nil {
		log.Println("InnerServer.ExecuteCommand error: ", err)
	}
}

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
		nc.clients[destination.Name] = client
		return client
	}
}