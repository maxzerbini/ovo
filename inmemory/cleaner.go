package inmemory

import (
	"github.com/maxzerbini/ovo/storage"
	"log"
	"time"
)

// Cleaner removes expired elements from the storage
type Cleaner struct {
	ks         *InMemoryStorage
	period     int64 // secs
	step       int64
	expirables map[int64][]*storage.MetaDataObj
	tickChan   <-chan time.Time
	doneChan   chan bool
	commands   chan func()
}

// Create a new Cleaner
func NewCleaner(ks *InMemoryStorage, period int64) (cl *Cleaner) {
	if period < 60 {
		period = 60
	}
	cl = &Cleaner{}
	cl.ks = ks
	cl.expirables = make(map[int64][]*storage.MetaDataObj)
	cl.step = 0
	cl.period = period
	cl.commands = make(chan func(), 100)
	go cl.execCmd()
	cl.doneChan = make(chan bool)
	cl.tickChan = time.NewTicker(time.Second * 30).C
	go cl.clean()
	return cl
}

// Execute the commands in serie.
func (cl *Cleaner) execCmd() {
	for f := range cl.commands {
		f()
	}
}

// Clean expired elements periodically
func (cl *Cleaner) clean() {
	for {
		select {
		case <-cl.tickChan:
			cl.RemoveExpiredElements()
		case <-cl.doneChan:
			return
		}
	}
}

// Stop the cleaner
func (cl *Cleaner) Stop() {
	cl.doneChan <- true
}

// Add an element to the cleaner
func (cl *Cleaner) AddElement(element *storage.MetaDataObj) {
	duration := element.CreationDate.Add(time.Duration(element.TTL) * time.Second).Sub(time.Now())
	if duration < 0 {
		duration = 0
	}
	nearestStep := cl.step + (int64(duration.Seconds()) / cl.period) + 1
	cl.commands <- func() {
		if _, ok := cl.expirables[nearestStep]; !ok {
			cl.expirables[nearestStep] = make([]*storage.MetaDataObj, 0)
			//log.Printf("Create list for step  %d\r\n",nearestStep)
		}
		cl.expirables[nearestStep] = append(cl.expirables[nearestStep], element)
		//log.Printf("Added element to the cleaner for step %d TTL %d list element %d\r\n",nearestStep, element.TTL, len(cl.expirables[nearestStep]))
	}
}

// Remove the expired elements
func (cl *Cleaner) RemoveExpiredElements() {
	if list, ok := cl.expirables[cl.step]; ok {
		for _, obj := range list {
			if obj.IsExpired() {
				//log.Printf("Remove expired element of key: %s\r\n",obj.Key)
				cl.ks.DeleteExpired(obj.Key)
			} else {
				//log.Printf("Sheduled clean for element of key: %s\r\n",obj.Key)
				cl.AddElement(obj)
			}
		}
		log.Printf("Expired elements removed: %d\r\n", len(list))
	} else {
		//log.Printf("List for step %d is void %d\r\n",cl.step, len(list))
	}
	step := cl.step // clone step before increasing its value
	cl.commands <- func() {
		delete(cl.expirables, step)
	}
	cl.step++
	//log.Printf("RemoveExpiredElements done new step is %d\r\n", cl.step)
}
