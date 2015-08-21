package inmemory
import (
	"bytes"
	"time"
	"github.com/maxzerbini/ovo/storage"
)

const collection_buffer_size = 100

// Collection (Map) of MetaDataObj. This collection is thread-safe.
type InMemoryCollection struct {
	storage map[string]*storage.MetaDataObj
	commands chan func()
}
// Execute the commands in serie.
func (coll *InMemoryCollection) execCmd() {
	for f := range coll.commands {
		f()
	}
}

// Create a InMemoryCollection.
func NewCollection()(*InMemoryCollection){
	coll := new(InMemoryCollection)
	coll.storage = make(map[string]*storage.MetaDataObj,10)	
	coll.commands = make (chan func(), collection_buffer_size)
	go coll.execCmd()
	return coll
}

// Add an item to the collection.
func (coll *InMemoryCollection) Put(obj *storage.MetaDataObj) {
	coll.commands <- func() { coll.storage[obj.Key] = obj }
}

// Get an item from the collection by key.
func (coll *InMemoryCollection) Get(key string) (*storage.MetaDataObj, bool) {
	retChan := make(chan *storage.MetaDataObj)
	defer close(retChan)
	coll.commands <- func() {
			if ret, ok := coll.storage[key]; ok {
				retChan <- ret
			} else {
				retChan <- nil
			}
		}
	var result = <- retChan
	if result == nil { 
		return nil, false 
	} else { 
		return result , true 
	}
}

// Remove the item of the collection
func (coll *InMemoryCollection) Delete(key string) {
	coll.commands <- func() { delete(coll.storage, key) }
}

// Get an item and remove it from the collection in a single operation.
func (coll *InMemoryCollection) GetAndRemove(key string) (*storage.MetaDataObj, bool)  {
	retChan := make(chan *storage.MetaDataObj)
	defer close(retChan)
	coll.commands <- func() {
			if ret, ok := coll.storage[key]; ok {
				delete(coll.storage, key)
				retChan <- ret
			} else {
				retChan <- nil
			}
		}
	var result = <- retChan
	if result == nil { 
		return nil, false 
	} else { 
		return result , true 
	}
}

// Update an item if the value is not changed.
func (coll *InMemoryCollection) UpdateValueIfEqual(obj *storage.MetaDataUpdObj) {
	coll.commands <- func() { 
		if ret, ok := coll.storage[obj.Key]; ok {
			if bytes.Equal(ret.Data, obj.Data) {
				ret.Data = obj.NewData
				ret.CreationDate = obj.CreationDate
			}
		} 
	}
}

// Update an item (key and value) if the value is not changed.
func (coll *InMemoryCollection) UpdateKeyAndValueIfEqual(obj *storage.MetaDataUpdObj) {
	coll.commands <- func() { 
		if ret, ok := coll.storage[obj.Key]; ok {
			if bytes.Equal(ret.Data, obj.Data) {
				delete(coll.storage, obj.Key)
				ret.Data = obj.NewData
				ret.Key = obj.NewKey
				ret.CreationDate = obj.CreationDate
				coll.storage[obj.NewKey]= ret
			}
		} 
	}
}

// Change the key of an item.
func (coll *InMemoryCollection) UpdateKey(obj *storage.MetaDataUpdObj) {
	coll.commands <- func() { 
		if ret, ok := coll.storage[obj.Key]; ok {
			delete(coll.storage, obj.Key)
			ret.Key = obj.NewKey
			ret.CreationDate = obj.CreationDate
			coll.storage[obj.NewKey]= ret
		} 
	}
}

// Count the items of the collection.
func (coll *InMemoryCollection) Count() (int){
	retChan := make(chan int)
	coll.commands <- func() { retChan <- len(coll.storage) }
	return <- retChan
}

// Touch an item restarting the time to live.
func (coll *InMemoryCollection) Touch(key string, updateDate time.Time ) {
	coll.commands <- func() { 
		if ret, ok := coll.storage[key]; ok {
			ret.CreationDate = updateDate
		} 
	}
}

// Remove the item of the collection
func (coll *InMemoryCollection) Iterate() ([]*storage.MetaDataObj){
	retChan := make(chan int)
	defer close(retChan)
	list := make([]*storage.MetaDataObj, 10)
	coll.commands <- func() { 
		for _, val := range coll.storage {
			list = append(list, val)
		}
		retChan <- 1
	}
	<- retChan //wait for result
	return list
}

// Remove the item of the collection
func (coll *InMemoryCollection) Expired() ([]*storage.MetaDataObj){
	retChan := make(chan int)
	defer close(retChan)
	list := make([]*storage.MetaDataObj, 10)
	coll.commands <- func() { 
		for _, val := range coll.storage {
			if val.IsExpired() {
				list = append(list, val)
			}
		}  
		retChan <- 1
	}
	<- retChan //wait for result
	return list
}