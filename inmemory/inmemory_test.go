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

func TestKSPutAndGet(t *testing.T) {
	t.Log("TestPutAndGet started")
	ks := NewInMemoryStorage()
	var data = storage.NewMetaDataObj("test", []byte("test string"), "default", 60, 0)
	ks.Put(&data)
	res, ok := ks.Get("test")
	if ok != nil {
		t.Fail()
	} else {
		t.Log(res.Key)
	}
}

func TestKSPutAndGetLoop(t *testing.T) {
	t.Log("TestPutAndGet started")
	ks := NewInMemoryStorage()
	for i := 0; i < 1000; i++ {
		var data = storage.NewMetaDataObj("testloopkey_"+strconv.Itoa(i), []byte("test string"), "default", 60, 0)
		ks.Put(&data)
	}
	for i := 0; i < 1000; i++ {
		_, ok := ks.Get("testloopkey_" + strconv.Itoa(i))
		if ok != nil {
			t.Fatal("key " + "testloopkey_" + strconv.Itoa(i) + " not found")
		} else {
			//fmt.Println(res.Key)
		}
	}
}

func TestKSCount(t *testing.T) {
	t.Log("TestPutAndGet started")
	ks := NewInMemoryStorage()
	max := 1000
	for i := 0; i < max; i++ {
		var data = storage.NewMetaDataObj("testloopkey_"+strconv.Itoa(i), []byte("test string"), "default", 60, 0)
		ks.Put(&data)
	}
	var count = ks.Count()
	if count != max {
		t.Fatal("Incorrect count " + strconv.Itoa(count))
	} else {
		t.Log("Correct count " + strconv.Itoa(count))
	}
}

func TestKSCountConcurrent(t *testing.T) {
	t.Log("TestPutAndGet started")
	ks := NewInMemoryStorage()
	max := 10000
	for i := 0; i < max; i++ {
		go func(j int) {
			var d = []byte(` data data data ..... data asdfghjklòzxcvbnm,wertyuiop
							 data data data ..... data asdfghjklòzxcvbnm,wertyuiop
							 data data data ..... data asdfghjklòzxcvbnm,wertyuiop
							 data data data ..... data asdfghjklòzxcvbnm,wertyuiop
							 data data data ..... data asdfghjklòzxcvbnm,wertyuiop
							 data data data ..... data asdfghjklòzxcvbnm,wertyuiop
							 data data data ..... data asdfghjklòzxcvbnm,wertyuiop
							 data data data ..... data asdfghjklòzxcvbnm,wertyuiop
							fghjklòdsfasdgfdasgsfadjgklfdagjkldfajglfdjdsgdfgfdgfdg
							qwertyuiopdfghjklzxcvbnm,12345678901234567890qwertyuiobn
			`)
			data := storage.NewMetaDataObj("testloopkey_"+strconv.Itoa(j), d, "default", 60, 0)
			ks.Put(&data)
			//fmt.Println("Put data "+data.Key)
		}(i)
	}
	time.Sleep(11 * 1e8)
	var count = ks.Count()
	if count != max {
		t.Fatal("Incorrect count " + strconv.Itoa(count))
	} else {
		t.Log("Correct count " + strconv.Itoa(count))
	}
}
