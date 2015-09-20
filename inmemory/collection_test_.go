package inmemory
import (
	"runtime"
	"testing"
	"strconv"
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

func TestPutAndGetLoop(t *testing.T) {
	t.Log("TestPutAndGet started")
	coll := NewCollection()
	for i:=0; i< 1000; i++ {
		var d = []byte("test string")
		var data = storage.NewMetaDataObj("testloopkey_"+strconv.Itoa(i), d, "default", 60, 1)
		coll.Put(&data)
	}
	for i:=0; i< 1000; i++ {
		_, ok := coll.Get("testloopkey_"+strconv.Itoa(i))
		if !ok {
			t.Fatal("key "+"testloopkey_"+strconv.Itoa(i)+" not found")
		} else {
			//fmt.Println(res.Key)
		}
	}
}

func TestCount(t *testing.T) {
	t.Log("TestPutAndGet started")
	coll := NewCollection()
	max := 1000
	for i:=0; i< max; i++ {
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

func TestCountConcurrent(t *testing.T) {
	t.Log("TestPutAndGet started")
	coll := NewCollection()
	max := 1000000
	for i:=0; i< max; i++ {
		go func (j int){
			var d = []byte("test string")
			data := storage.NewMetaDataObj("testloopkey_"+strconv.Itoa(j), d, "default", 60, 1)
			coll.Put(&data)
			//fmt.Println("Put data "+data.Key)
		} (i)
	}
	time.Sleep(1 * 1e9)
	var count = coll.Count()
	if count != max {
		t.Fatal("Incorrect count " + strconv.Itoa(count))
	} else {
		t.Log("Correct count " + strconv.Itoa(count))
	}
}