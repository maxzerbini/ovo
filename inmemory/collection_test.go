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

func TestPutAndGet(t *testing.T) {
	t.Log("TestPutAndGet started")
	coll := NewCollection()
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

func TestPutLoop(t *testing.T) {
	t.Log("TestPutLoop started")
	coll := NewCollection()
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

func TestGetLoop(t *testing.T) {
	t.Log("TestGetLoop started")
	coll := NewCollection()
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

func _TestCount(t *testing.T) {
	t.Log("TestPutAndGet started")
	coll := NewCollection()
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

func _TestCountConcurrent(t *testing.T) {
	t.Log("TestPutAndGet started")
	coll := NewCollection()
	max := 100000
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
