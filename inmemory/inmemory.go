// This package contains the InMemory storage and its data structures.
//
//
package inmemory
import (
	"errors"
	"time"
	"github.com/maxzerbini/ovo/storage"
)

// The InMemoryStorage struct implements the OvoStorage interface. 
type InMemoryStorage struct {
	collection *InMemoryCollection
	cleaner *Cleaner
}

// Create a InMemoryStorage.
func NewInMemoryStorage()(*InMemoryStorage){
	ks := new(InMemoryStorage)
	ks.collection = NewCollection()
	ks.cleaner = NewCleaner(ks, 60)
	return ks
}

// Add an item to the storage.
func (ks *InMemoryStorage) Put(obj *storage.MetaDataObj) (error){
	if(obj != nil){
		if len(obj.Key) == 0 {
			return errors.New("Object key is null.")
		}
		if len(obj.Collection) == 0 {
			obj.Collection = "default"
		}
		obj.CreationDate = time.Now()
		ks.collection.Put(obj)
		if obj.TTL > 0 {
			go ks.cleaner.AddElement(obj)
		}
	}
	return errors.New("Object is null.")
}

// Get an item from the storage by key.
func (ks *InMemoryStorage) Get(key string) (*storage.MetaDataObj, error){
	if obj, ok := ks.collection.Get(key); ok {
		if obj.IsExpired(){
			return nil, errors.New("Not found.")
		}
		return obj, nil
	}
	return nil, errors.New("Not found.")
}

// Remove the item of the storage
func (ks *InMemoryStorage) Delete(key string) {
	ks.collection.Delete(key)
}

// Remove the item of the storage
func (ks *InMemoryStorage) DeleteExpired(key string) {
	ks.collection.DeleteExpired(key)
}

// Get an item and remove it from the storage in a single operation.
func (ks *InMemoryStorage) GetAndRemove(key string) (*storage.MetaDataObj, error) {
	if obj, ok := ks.collection.GetAndRemove(key); ok {
		if obj.IsExpired(){
			return nil, errors.New("Not found.")
		}
		return obj, nil
	}
	return nil, errors.New("Not found.")
}

// Update an item if the value is not changed.
func (ks *InMemoryStorage) UpdateValueIfEqual(obj *storage.MetaDataUpdObj) (error){
	if(obj != nil){
		if len(obj.Key) == 0 {
			return errors.New("Object key is null.")
		}
		if len(obj.Collection) == 0 {
			obj.Collection = "default"
		}
		obj.CreationDate = time.Now()
		ks.collection.UpdateValueIfEqual(obj)
	}
	return errors.New("Object is null.")
}

// Update an item (key and value) if the value is not changed.
func (ks *InMemoryStorage) UpdateKeyAndValueIfEqual(obj *storage.MetaDataUpdObj) (error){
	if(obj != nil){
		if len(obj.Key) == 0 {
			return errors.New("Object key is null.")
		}
		if len(obj.NewKey) == 0 {
			return errors.New("Object new key is null.")
		}
		if len(obj.Collection) == 0 {
			obj.Collection = "default"
		}
		obj.CreationDate = time.Now()
		ks.collection.UpdateKeyAndValueIfEqual(obj)
	}
	return errors.New("Object is null.")
}

// Change the key of an item.
func (ks *InMemoryStorage) UpdateKey(obj *storage.MetaDataUpdObj) (error){
	if(obj != nil){
		if len(obj.Key) == 0 {
			return errors.New("Object key is null.")
		}
		if len(obj.NewKey) == 0 {
			return errors.New("Object new key is null.")
		}
		if len(obj.Collection) == 0 {
			obj.Collection = "default"
		}
		obj.CreationDate = time.Now()
		ks.collection.UpdateKey(obj)
	}
	return errors.New("Object is null.")
}

// Touch an item restarting the time to live.
func (ks *InMemoryStorage) Touch(key string) {
	ks.collection.Touch(key, time.Now())
}

// Count the keys in the storage.
func (ks *InMemoryStorage) Count() (int){
	return ks.collection.Count()
}

func (ks *InMemoryStorage) List() ([]*storage.MetaDataObj){
	return ks.collection.List()
}

func (ks *InMemoryStorage) ListExpired()(elements []*storage.MetaDataObj) {
	return ks.collection.ListExpired()
}
