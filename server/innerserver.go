package server

import(
	"github.com/maxzerbini/ovo/storage"
	"github.com/maxzerbini/ovo/processor"
	"github.com/maxzerbini/ovo/command"
	"net"
	"net/rpc"
	"net/http"
	"log"
	"errors"
)

type InnerServer struct {
	keystorage storage.OvoStorage
	incmdproc *processor.InCommandQueue
	config *ServerConf	
}

func NewInnerServer(conf *ServerConf, ks storage.OvoStorage, in *processor.InCommandQueue) *InnerServer{
	return &InnerServer{keystorage:ks, incmdproc:in, config:conf}
}

func (srv *InnerServer)Do(){
	rpc.Register(srv)
	rpc.HandleHTTP()
	listener, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("Starting RPC-server -listen error:", e)
	}
	http.Serve(listener, nil)
}
// Enqueue a remote command.
func (srv *InnerServer) ExecuteCommand(rpccmd *command.RpcCommand, reply *int) (err error) {
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

