package server

import (
	"log"
)

import (
	"github.com/gin-gonic/gin"
	"github.com/maxzerbini/ovo/cluster"
	"github.com/maxzerbini/ovo/command"
	"github.com/maxzerbini/ovo/processor"
	"github.com/maxzerbini/ovo/server/model"
	"github.com/maxzerbini/ovo/storage"
	"github.com/maxzerbini/ovo/util"
	"net/http"
	"strconv"
)

type Server struct {
	keystorage  storage.OvoStorage
	incmdproc   *processor.InCommandQueue
	outcmdproc  *processor.OutCommandQueue
	config      *ServerConf
	partitioner *processor.Partitioner
	innerServer *InnerServer
	nodeChecker *Checker
}

func NewServer(conf *ServerConf, ks storage.OvoStorage) *Server {
	srv := &Server{keystorage: ks, config: conf}
	srv.incmdproc = processor.NewCommandQueue(ks)
	srv.outcmdproc = processor.NewOutCommandQueue(conf.ServerNode, &conf.Topology, srv.incmdproc)
	srv.partitioner = processor.NewPartitioner(ks, conf.ServerNode, srv.outcmdproc)
	srv.innerServer = NewInnerServer(conf, ks, srv.incmdproc, srv.outcmdproc, srv.partitioner)
	srv.nodeChecker = NewChecker(conf, srv.outcmdproc, srv.partitioner)
	return srv
}

func (srv *Server) Do() {
	log.Printf("Staring node %s ...\r\n", srv.config.ServerNode.Node.Name)
	go srv.innerServer.Do()
	// Creates a router without any middleware by default
	router := gin.New()
	// Global middleware
	if srv.config.Debug {
		router.Use(gin.Logger())
	}
	router.Use(gin.Recovery())
	router.GET("/ovo/keystorage", srv.count)
	router.GET("/ovo/keys", srv.keys)
	router.GET("/ovo/keystorage/:key", srv.get)
	router.POST("/ovo/keystorage", srv.post)
	router.PUT("/ovo/keystorage", srv.post)
	router.DELETE("/ovo/keystorage/:key", srv.delete)
	router.GET("/ovo/keystorage/:key/getandremove", srv.getAndRemove)
	router.POST("/ovo/keystorage/:key/updatevalueifequal", srv.updateValueIfEqual)
	router.PUT("/ovo/keystorage/:key/updatevalueifequal", srv.updateValueIfEqual)
	router.POST("/ovo/keystorage/:key/updatekeyvalueifequal", srv.updateKeyAndValueIfEqual)
	router.PUT("/ovo/keystorage/:key/updatekeyvalueifequal", srv.updateKeyAndValueIfEqual)
	router.POST("/ovo/keystorage/:key/updatekey", srv.updateKey)
	router.PUT("/ovo/keystorage/:key/updatekey", srv.updateKey)
	router.GET("/ovo/cluster", srv.getTopology)
	router.GET("/ovo/cluster/me", srv.getCurrentNode)
	router.POST("/ovo/counters", srv.increment)
	router.GET("/ovo/counters/:key", srv.getcounter)
	if srv.config.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	// register this node in the cluster
	srv.registerServer()
	// start node checker
	go srv.nodeChecker.Do()
	log.Printf("Node %s started\r\n", srv.config.ServerNode.Node.Name)
	// Listen and server on Host:Port
	router.Run(srv.config.ServerNode.Node.Host + ":" + strconv.Itoa(srv.config.ServerNode.Node.Port))
}

func (srv *Server) registerServer() {
	// update topology
	if len(srv.config.Topology.Nodes) > 0 {
		// connect first node
		for _, node := range srv.config.Topology.Nodes {
			if node.Node.Name != srv.config.ServerNode.Node.Name {
				if topology, err := srv.outcmdproc.Caller.GetTopology(srv.config.ServerNode.Node.Name, node.Node); err == nil {
					srv.config.Topology.Merge(topology)
					break
				}
			}
		}
		topologies := make([]*cluster.ClusterTopology, 0)
		failedNodes := make([]*cluster.ClusterTopologyNode, 0)
		for _, node := range srv.config.Topology.Nodes {
			if node.Node.Name != srv.config.ServerNode.Node.Name {
				if util.ContainsString(srv.config.ServerNode.Stepbrothers, node.Node.Name) && util.ContainsString(srv.config.ServerNode.Twins, node.Node.Name) {
					// register this node as twin and stepbrother
					if topology, err := srv.outcmdproc.Caller.RegisterTwinAndStepbrother(srv.config.ServerNode, node.Node); err == nil && topology != nil {
						log.Printf("Registration was successful on twin and stepbrother node %s\r\n", node.Node.Name)
						topologies = append(topologies, topology)
					} else {
						log.Printf("Registration failed on twin and stepbrother node %s\r\n", node.Node.Name)
						failedNodes = append(failedNodes, node)
					}
				} else if util.ContainsString(srv.config.ServerNode.Stepbrothers, node.Node.Name) {
					// register this node on a stepbrother node: the current node became the twin of the stepbrother node
					if topology, err := srv.outcmdproc.Caller.RegisterTwin(srv.config.ServerNode, node.Node); err == nil && topology != nil {
						log.Printf("Registration was successful on stepbrother node %s\r\n", node.Node.Name)
						topologies = append(topologies, topology)
					} else {
						log.Printf("Registration failed on stepbrother node %s\r\n", node.Node.Name)
						failedNodes = append(failedNodes, node)
					}
				} else if util.ContainsString(srv.config.ServerNode.Twins, node.Node.Name) {
					// register this node as a twin node: the node became the stepbrother of the current node
					if topology, err := srv.outcmdproc.Caller.RegisterStepbrother(srv.config.ServerNode, node.Node); err == nil && topology != nil {
						log.Printf("Registration was successful on twin node %s\r\n", node.Node.Name)
						topologies = append(topologies, topology)
					} else {
						log.Printf("Registration failed on twin node %s\r\n", node.Node.Name)
						failedNodes = append(failedNodes, node)
					}
				} else {
					// register node on the cluster
					if topology, err := srv.outcmdproc.Caller.RegisterNode(srv.config.ServerNode, node.Node); err == nil && topology != nil {
						log.Printf("Registration was successful on node %s\r\n", node.Node.Name)
						topologies = append(topologies, topology)
					} else {
						log.Printf("Registration failed on node %s\r\n", node.Node.Name)
						failedNodes = append(failedNodes, node)
					}
				}
			}
		}
		// remove failed nodes
		for _, node := range failedNodes {
			srv.config.Topology.RemoveNode(node.Node.Name)
		}
		// merge
		for _, topology := range topologies {
			srv.config.Topology.Merge(topology)
		}
	}
	srv.config.WriteTmp()
}

func (srv *Server) count(c *gin.Context) {
	res := srv.keystorage.Count()
	result := model.NewOvoResponse("done", "0", res)
	c.JSON(http.StatusOK, result)
}

func (srv *Server) keys(c *gin.Context) {
	keys := srv.keystorage.Keys()
	res := &model.OvoKVKeys{Keys: keys}
	result := model.NewOvoResponse("done", "0", res)
	c.JSON(http.StatusOK, result)
}

func (srv *Server) get(c *gin.Context) {
	key := c.Param("key")
	if res, err := srv.keystorage.Get(key); err == nil {
		obj := model.NewOvoKVResponse(res)
		result := model.NewOvoResponse("done", "0", obj)
		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusNotFound, model.NewOvoResponse("error", "101", nil))
	}
}

func (srv *Server) post(c *gin.Context) {
	var kv model.OvoKVRequest
	if c.BindJSON(&kv) == nil {
		obj := model.NewMetaDataObj(&kv)
		srv.keystorage.Put(obj)
		srv.outcmdproc.Enqueu(&command.Command{OpCode: "put", Obj: obj.MetaDataUpdObj()})
		c.JSON(http.StatusOK, model.NewOvoResponse("done", "0", nil))
	} else {
		c.JSON(http.StatusBadRequest, model.NewOvoResponse("error", "10", nil))
	}
}

func (srv *Server) delete(c *gin.Context) {
	key := c.Param("key")
	srv.keystorage.Delete(key)
	srv.outcmdproc.Enqueu(&command.Command{OpCode: "delete", Obj: &storage.MetaDataUpdObj{Key: key}})
	c.JSON(http.StatusOK, model.NewOvoResponse("done", "0", nil))
}

func (srv *Server) getAndRemove(c *gin.Context) {
	key := c.Param("key")
	if res, err := srv.keystorage.GetAndRemove(key); err == nil {
		obj := model.NewOvoKVResponse(res)
		srv.outcmdproc.Enqueu(&command.Command{OpCode: "delete", Obj: &storage.MetaDataUpdObj{Key: key}})
		result := model.NewOvoResponse("done", "0", obj)
		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusNotFound, model.NewOvoResponse("error", "101", nil))
	}
}

func (srv *Server) updateValueIfEqual(c *gin.Context) {
	key := c.Param("key")
	var kv model.OvoKVUpdateRequest
	if c.BindJSON(&kv) == nil {
		obj := model.NewMetaDataUpdObj(&kv)
		obj.Key = key
		err := srv.keystorage.UpdateValueIfEqual(obj)
		if err == nil {
			srv.outcmdproc.Enqueu(&command.Command{OpCode: "updatevalue", Obj: obj})
			c.JSON(http.StatusOK, model.NewOvoResponse("done", "0", nil))
		} else {
			c.JSON(http.StatusForbidden, model.NewOvoResponse("error", "103", nil))
		}
	} else {
		c.JSON(http.StatusBadRequest, model.NewOvoResponse("error", "10", nil))
	}
}

func (srv *Server) updateKeyAndValueIfEqual(c *gin.Context) {
	key := c.Param("key")
	var kv model.OvoKVUpdateRequest
	if c.BindJSON(&kv) == nil {
		obj := model.NewMetaDataUpdObj(&kv)
		obj.Key = key
		err := srv.keystorage.UpdateKeyAndValueIfEqual(obj)
		if err == nil {
			srv.outcmdproc.Enqueu(&command.Command{OpCode: "updatekeyvalue", Obj: obj})
			c.JSON(http.StatusOK, model.NewOvoResponse("done", "0", nil))
		} else {
			c.JSON(http.StatusForbidden, model.NewOvoResponse("error", "104", nil))
		}
	} else {
		c.JSON(http.StatusBadRequest, model.NewOvoResponse("error", "10", nil))
	}
}

func (srv *Server) updateKey(c *gin.Context) {
	key := c.Param("key")
	var kv model.OvoKVUpdateRequest
	if c.BindJSON(&kv) == nil {
		obj := model.NewMetaDataUpdObj(&kv)
		obj.Key = key
		err := srv.keystorage.UpdateKey(obj)
		if err == nil {
			srv.outcmdproc.Enqueu(&command.Command{OpCode: "updatekey", Obj: obj})
			c.JSON(http.StatusOK, model.NewOvoResponse("done", "0", nil))
		} else {
			c.JSON(http.StatusForbidden, model.NewOvoResponse("error", "105", nil))
		}
	} else {
		c.JSON(http.StatusBadRequest, model.NewOvoResponse("error", "10", nil))
	}
}

func (srv *Server) getTopology(c *gin.Context) {
	res := model.NewOvoTopology(&srv.config.Topology)
	result := model.NewOvoResponse("done", "0", res)
	c.JSON(http.StatusOK, result)
}

func (srv *Server) getCurrentNode(c *gin.Context) {
	res := model.NewOvoTopologyNode(srv.config.ServerNode)
	result := model.NewOvoResponse("done", "0", res)
	c.JSON(http.StatusOK, result)
}

func (srv *Server) increment(c *gin.Context) {
	var counter model.OvoCounter
	if c.BindJSON(&counter) == nil {
		obj := model.NewMetaDataCounter(&counter)
		cnt := srv.keystorage.Increment(obj)
		srv.outcmdproc.Enqueu(&command.Command{OpCode: "setcounter", Obj: cnt.MetaDataUpdObj()})
		c.JSON(http.StatusOK, model.NewOvoResponse("done", "0", model.NewOvoCounterResponse(cnt)))
	} else {
		c.JSON(http.StatusBadRequest, model.NewOvoResponse("error", "10", nil))
	}
}

func (srv *Server) getcounter(c *gin.Context) {
	key := c.Param("key")
	if res, err := srv.keystorage.GetCounter(key); err == nil {
		obj := model.NewOvoCounterResponse(res)
		result := model.NewOvoResponse("done", "0", obj)
		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusNotFound, model.NewOvoResponse("error", "101", nil))
	}
}
