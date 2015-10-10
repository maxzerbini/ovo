package storage

import (
	"time"
)

type MetaDataObj struct {
	Key string
	Data []byte
	Collection string
	CreationDate time.Time
	TTL int
	Hash int
}

type MetaDataUpdObj struct {
	Key string
	NewKey string
	Data []byte
	NewData []byte
	Collection string
	CreationDate time.Time
	TTL int
	Hash int
	NewHash int
}

func NewMetaDataObj(key string, data []byte, collection string, ttl int, hash int) MetaDataObj {
	return MetaDataObj{Key:key, Data:data, Collection:collection, CreationDate:time.Now(), TTL:ttl, Hash:hash}
}

func (obj *MetaDataObj)MetaDataUpdObj() *MetaDataUpdObj {
	return &MetaDataUpdObj{Key:obj.Key, Data:obj.Data, Collection:obj.Collection, CreationDate:obj.CreationDate, TTL:obj.TTL, Hash:obj.Hash, NewKey:"", NewData:make([]byte, 0)}
}

func (obj MetaDataObj) IsExpired() bool {
	if obj.TTL == 0 { return false }
	return time.Now().After(obj.CreationDate.Add(time.Duration(obj.TTL)*time.Second))
}

func (obj *MetaDataUpdObj) MetaDataObj() (*MetaDataObj) {
	item := &MetaDataObj {Key:obj.Key, Data:obj.Data, Collection:obj.Collection, TTL:obj.TTL, Hash: obj.Hash}
	return item
}

type OvoStorage interface {
	Get(key string) (obj *MetaDataObj, err error)
	Put(obj *MetaDataObj) (error)
	Delete(key string)
	GetAndRemove(key string) (obj *MetaDataObj, err error)
	UpdateValueIfEqual(obj *MetaDataUpdObj) (error)
	UpdateKeyAndValueIfEqual(obj *MetaDataUpdObj) (error)
	UpdateKey(obj *MetaDataUpdObj) (error)
	Touch(key string)
	Count() (int)
	List() ([]*MetaDataObj)
}