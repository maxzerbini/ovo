package model

import (
	"github.com/maxzerbini/ovo/storage"
)

type Any interface { }

type OvoResponse struct {
	Status string
	Code string
	Data Any
}

type OvoKVRequest struct {
	Key string
	Data []byte
	Collection string
	TTL int
}

type OvoKVUpdateRequest struct {
	Key string
	NewKey string
	Data []byte
	NewData []byte
}

type OvoKVResponse struct {
	Key string
	Data []byte
}

func NewOvoResponse(status string, code string, data Any) *OvoResponse {
	return &OvoResponse{Status:status, Code:code, Data: data}
}

func NewOvoKVResponse(obj *storage.MetaDataObj) *OvoKVResponse {
	var rsp = &OvoKVResponse{Key:obj.Key, Data:obj.Data}
	return rsp
}

func NewMetaDataObj(req *OvoKVRequest) *storage.MetaDataObj {
	var obj = new(storage.MetaDataObj)
	obj.Key = req.Key
	obj.Data = req.Data
	obj.Collection = req.Collection
	obj.TTL = req.TTL
	return obj
}

func NewMetaDataUpdObj(req *OvoKVUpdateRequest) *storage.MetaDataUpdObj {
	var obj = new(storage.MetaDataUpdObj)
	obj.Key = req.Key
	obj.NewKey = req.NewKey
	obj.Data = req.Data
	obj.NewData = req.NewData
	return obj
}