package model

import (
	"github.com/maxzerbini/ovo/cluster"
	"github.com/maxzerbini/ovo/storage"
)

type Any interface{}

type OvoResponse struct {
	Status string
	Code   string
	Data   Any
}

type OvoKVRequest struct {
	Key        string
	Data       []byte
	Collection string
	TTL        int
	Hash       int
}

type OvoKVUpdateRequest struct {
	Key     string
	NewKey  string
	Data    []byte
	NewData []byte
	Hash    int
	NewHash int
}

type OvoKVResponse struct {
	Key  string
	Data []byte
}

type OvoKVKeys struct {
	Keys []string
}

type OvoTopologyNode struct {
	Name      string
	HashRange []int
	Host      string
	Port      int
	State     string
	Twins     []string
}

type OvoTopology struct {
	Nodes []*OvoTopologyNode
}

type OvoCounter struct {
	Key   string
	Value int64
	TTL   int
	Hash  int
}

type OvoCounterResponse struct {
	Key   string
	Value int64
}

func NewOvoResponse(status string, code string, data Any) *OvoResponse {
	return &OvoResponse{Status: status, Code: code, Data: data}
}

func NewOvoKVResponse(obj *storage.MetaDataObj) *OvoKVResponse {
	var rsp = &OvoKVResponse{Key: obj.Key, Data: obj.Data}
	return rsp
}

func NewMetaDataObj(req *OvoKVRequest) *storage.MetaDataObj {
	var obj = new(storage.MetaDataObj)
	obj.Key = req.Key
	obj.Data = req.Data
	obj.Collection = req.Collection
	obj.TTL = req.TTL
	obj.Hash = req.Hash
	return obj
}

func NewMetaDataUpdObj(req *OvoKVUpdateRequest) *storage.MetaDataUpdObj {
	var obj = new(storage.MetaDataUpdObj)
	obj.Key = req.Key
	obj.NewKey = req.NewKey
	obj.Data = req.Data
	obj.NewData = req.NewData
	obj.Hash = req.Hash
	obj.NewHash = req.NewHash
	return obj
}

func NewOvoTopologyNode(node *cluster.ClusterTopologyNode) *OvoTopologyNode {
	return &OvoTopologyNode{Name: node.Node.Name, HashRange: node.Node.HashRange, Host: node.Node.Host, Port: node.Node.Port, State: node.Node.State, Twins: node.Twins}
}

func NewOvoTopology(topology *cluster.ClusterTopology) *OvoTopology {
	ret := &OvoTopology{Nodes: make([]*OvoTopologyNode, 0)}
	for _, node := range topology.Nodes {
		ret.Nodes = append(ret.Nodes, NewOvoTopologyNode(node))
	}
	return ret
}

func NewMetaDataCounter(counter *OvoCounter) *storage.MetaDataCounter {
	return &storage.MetaDataCounter{Key: counter.Key, Value: counter.Value, TTL: counter.TTL, Hash: counter.Hash}
}

func NewOvoCounterResponse(counter *storage.MetaDataCounter) *OvoCounterResponse {
	return &OvoCounterResponse{Key: counter.Key, Value: counter.Value}
}
