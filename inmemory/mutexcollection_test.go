package inmemory

import (
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/maxzerbini/ovo/storage"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func TestPutAndGetMutex(t *testing.T) {
	t.Log("TestPutAndGet started")
	coll := NewMutexCollection()
	var d = []byte("test string")
	var data = storage.NewMetaDataObj("test", d, "default", 60, 1)
	coll.Put(&data)
	res, ok := coll.Get("test")
	if !ok {
		t.Fail()
	} else {
		t.Log(res.Key)
	}
}

func TestPutAndGetLoopMutex(t *testing.T) {
	t.Log("TestPutAndGet started")
	coll := NewMutexCollection()
	for i := 0; i < 1000; i++ {
		var d = []byte("test string")
		var data = storage.NewMetaDataObj("testloopkey_"+strconv.Itoa(i), d, "default", 60, 1)
		coll.Put(&data)
	}
	for i := 0; i < 1000; i++ {
		_, ok := coll.Get("testloopkey_" + strconv.Itoa(i))
		if !ok {
			t.Fatal("key " + "testloopkey_" + strconv.Itoa(i) + " not found")
		} else {
			//fmt.Println(res.Key)
		}
	}
}

func TestPutLoopMutex(t *testing.T) {
	t.Log("TestPutLoopMutex started")
	coll := NewMutexCollection()
	operationNumber := 100000
	resultChan := make(chan (bool), operationNumber)
	startTime := time.Now()
	for i := 0; i < operationNumber; i++ {
		go func(j int) {
			var d = []byte("test string")
			data := storage.NewMetaDataObj("testloopkey_"+strconv.Itoa(j), d, "default", 60, 1)
			coll.Put(&data)
			resultChan <- true
		}(i)
	}
	for i := 0; i < operationNumber; i++ {
		<-resultChan
	}
	t.Logf("All done: %d - elapsed time = %s", len(resultChan), time.Since(startTime))
}

func TestGetLoopMutex(t *testing.T) {
	t.Log("TestGetLoopMutex started")
	coll := NewMutexCollection()
	operationNumber := 100000
	for i := 0; i < operationNumber; i++ {
		var d = []byte("test string")
		data := storage.NewMetaDataObj("testloopkey_"+strconv.Itoa(i), d, "default", 60, 1)
		coll.Put(&data)
	}
	t.Logf("Number of object: %d", coll.Count())
	startTime := time.Now()
	resultChan := make(chan (bool), operationNumber)
	for i := 0; i < operationNumber; i++ {
		go func(j int) {
			_, _ = coll.Get("testloopkey_" + strconv.Itoa(i))
			resultChan <- true
		}(i)
	}
	for i := 0; i < operationNumber; i++ {
		<-resultChan
	}
	t.Logf("All done: %d - elapsed time = %s", len(resultChan), time.Since(startTime))
}

func TestCountMutex(t *testing.T) {
	t.Log("TestPutAndGet started")
	coll := NewMutexCollection()
	max := 1000
	for i := 0; i < max; i++ {
		var d = []byte("test string")
		var data = storage.NewMetaDataObj("testloopkey_"+strconv.Itoa(i), d, "default", 60, 1)
		coll.Put(&data)
	}
	var count = coll.Count()
	if count != max {
		t.Fatal("Incorrect count " + strconv.Itoa(count))
	} else {
		t.Log("Correct count " + strconv.Itoa(count))
	}
}

func _TestCountConcurrentMutex(t *testing.T) {
	t.Log("TestPutAndGet started")
	coll := NewMutexCollection()
	max := 1000000
	for i := 0; i < max; i++ {
		go func(j int) {
			var d = []byte("test string")
			data := storage.NewMetaDataObj("testloopkey_"+strconv.Itoa(j), d, "default", 60, 1)
			coll.Put(&data)
			//fmt.Println("Put data "+data.Key)
		}(i)
	}
	time.Sleep(20 * 1e9)
	var count = coll.Count()
	if count != max {
		t.Fatal("Incorrect count " + strconv.Itoa(count))
	} else {
		t.Log("Correct count " + strconv.Itoa(count))
	}
}
