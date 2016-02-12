package inmemory

import (
	"bytes"
	"sync"
	"time"

	"github.com/maxzerbini/ovo/storage"
)

const collection_buffer_size = 100

// Collection (Map) of MetaDataObj. This collection is thread-safe.
type InMemoryMutexCollection struct {
	storage  map[string]*storage.MetaDataObj
	counters map[string]*storage.MetaDataCounter
	sync.RWMutex
}

// Create a InMemoryMutexCollection.
func NewMutexCollection() *InMemoryMutexCollection {
	coll := new(InMemoryMutexCollection)
	coll.storage = make(map[string]*storage.MetaDataObj, 10)
	coll.counters = make(map[string]*storage.MetaDataCounter, 10)
	return coll
}

// Add an item to the collection.
func (coll *InMemoryMutexCollection) Put(obj *storage.MetaDataObj) {
	coll.Lock()
	defer coll.Unlock()
	coll.storage[obj.Key] = obj
}

// Get an item from the collection by key.
func (coll *InMemoryMutexCollection) Get(key string) (*storage.MetaDataObj, bool) {
	coll.RLock()
	defer coll.RUnlock()
	if ret, ok := coll.storage[key]; ok {
		return ret, true
	} else {
		return nil, false
	}
}

// Remove the item of the collection
func (coll *InMemoryMutexCollection) Delete(key string) {
	coll.Lock()
	defer coll.Unlock()
	delete(coll.storage, key)
}

// Remove the item of the collection
func (coll *InMemoryMutexCollection) DeleteExpired(key string) {
	coll.Lock()
	defer coll.Unlock()
	if ret, ok := coll.storage[key]; ok {
		if ret.IsExpired() {
			delete(coll.storage, key)
		}
	}
}

// Get an item and remove it from the collection in a single operation.
func (coll *InMemoryMutexCollection) GetAndRemove(key string) (*storage.MetaDataObj, bool) {
	coll.Lock()
	defer coll.Unlock()
	if ret, ok := coll.storage[key]; ok {
		delete(coll.storage, key)
		return ret, true
	} else {
		return nil, false
	}
}

// Update an item if the value is not changed.
func (coll *InMemoryMutexCollection) UpdateValueIfEqual(obj *storage.MetaDataUpdObj) bool {
	coll.Lock()
	defer coll.Unlock()
	if ret, ok := coll.storage[obj.Key]; ok {
		if bytes.Equal(ret.Data, obj.Data) {
			ret.Data = obj.NewData
			ret.CreationDate = obj.CreationDate
			return true
		} else {
			return false // not equal
		}
	} else {
		return false // not found
	}
}

// Update an item (key and value) if the value is not changed.
func (coll *InMemoryMutexCollection) UpdateKeyAndValueIfEqual(obj *storage.MetaDataUpdObj) bool {
	coll.Lock()
	defer coll.Unlock()
	if ret, ok := coll.storage[obj.Key]; ok {
		if bytes.Equal(ret.Data, obj.Data) {
			delete(coll.storage, obj.Key)
			ret.Data = obj.NewData
			ret.Key = obj.NewKey
			ret.Hash = obj.NewHash
			ret.CreationDate = obj.CreationDate
			coll.storage[obj.NewKey] = ret
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

// Change the key of an item.
func (coll *InMemoryMutexCollection) UpdateKey(obj *storage.MetaDataUpdObj) {
	coll.Lock()
	defer coll.Unlock()
	if ret, ok := coll.storage[obj.Key]; ok {
		delete(coll.storage, obj.Key)
		ret.Key = obj.NewKey
		ret.CreationDate = obj.CreationDate
		ret.Hash = obj.NewHash
		coll.storage[obj.NewKey] = ret
	}
}

// Count the items of the collection.
func (coll *InMemoryMutexCollection) Count() int {
	coll.RLock()
	defer coll.RUnlock()
	return len(coll.storage)
}

// Touch an item restarting the time to live.
func (coll *InMemoryMutexCollection) Touch(key string, updateDate time.Time) {
	coll.Lock()
	defer coll.Unlock()
	if ret, ok := coll.storage[key]; ok {
		ret.CreationDate = updateDate
	}
}

// List the keys of the items in the collection
func (coll *InMemoryMutexCollection) Keys() []string {
	coll.RLock()
	defer coll.RUnlock()
	list := make([]string, 0)
	for _, val := range coll.storage {
		if !val.IsExpired() {
			list = append(list, val.Key)
		}
	}
	return list
}

// List the items in the collection
func (coll *InMemoryMutexCollection) List() []*storage.MetaDataObj {
	coll.RLock()
	defer coll.RUnlock()
	list := make([]*storage.MetaDataObj, 0)
	for _, val := range coll.storage {
		if !val.IsExpired() {
			list = append(list, val)
		}
	}
	return list
}

// List the expired items of the collection
func (coll *InMemoryMutexCollection) ListExpired() []*storage.MetaDataObj {
	coll.RLock()
	defer coll.RUnlock()
	list := make([]*storage.MetaDataObj, 0)
	for _, val := range coll.storage {
		if val.IsExpired() {
			list = append(list, val)
		}
	}
	return list
}

// Increment a counter.
func (coll *InMemoryMutexCollection) Increment(c *storage.MetaDataCounter) *storage.MetaDataCounter {
	coll.Lock()
	defer coll.Unlock()
	if ret, ok := coll.counters[c.Key]; ok {
		if ret.IsExpired() {
			ret.CreationDate = time.Now()
			ret.Value = c.Value
		} else {
			ret.Value += c.Value
		}
		return ret
	} else {
		c.CreationDate = time.Now()
		coll.counters[c.Key] = c
		return c
	}
}

// Set the value of a counter.
func (coll *InMemoryMutexCollection) SetCounter(c *storage.MetaDataCounter) *storage.MetaDataCounter {
	coll.Lock()
	defer coll.Unlock()
	if ret, ok := coll.counters[c.Key]; ok {
		if ret.IsExpired() {
			ret.CreationDate = time.Now()
		}
		ret.Value = c.Value
		return ret
	} else {
		c.CreationDate = time.Now()
		coll.counters[c.Key] = c
		return c
	}
}

// Get a counter by key.
func (coll *InMemoryMutexCollection) GetCounter(key string) (*storage.MetaDataCounter, bool) {
	coll.RLock()
	defer coll.RUnlock()
	if ret, ok := coll.counters[key]; ok {
		return ret, true
	} else {
		return nil, false
	}
}

// Remove the item of the collection
func (coll *InMemoryMutexCollection) DeleteCounter(key string) {
	coll.Lock()
	defer coll.Unlock()
	delete(coll.counters, key)
}

// List the items in the collection
func (coll *InMemoryMutexCollection) ListCounters() []*storage.MetaDataCounter {
	coll.RLock()
	defer coll.RUnlock()
	list := make([]*storage.MetaDataCounter, 0)
	for _, val := range coll.counters {
		if !val.IsExpired() {
			list = append(list, val)
		}
	}
	return list
}

// Delete an item if the value is not changed.
func (coll *InMemoryMutexCollection) DeleteValueIfEqual(obj *storage.MetaDataObj) bool {
	coll.Lock()
	defer coll.Unlock()
	if ret, ok := coll.storage[obj.Key]; ok {
		if bytes.Equal(ret.Data, obj.Data) {
			delete(coll.storage, obj.Key)
			return true
		} else {
			return false // values are not equal
		}
	} else {
		return true
	}

}
