package metrics

import (
	"encoding/json"
	"sync"
	"time"
)

func NewLinearCounter() *linearCounter {
	return &linearCounter{
		Buckets: map[string]*intCounter{},
	}
}

type linearCounter struct {
	sync.Mutex
	Buckets map[string]*intCounter
}

func (self *linearCounter) Event(at time.Time, ids ...string) {
	self.Lock()
	defer self.Unlock()

	for _, id := range ids {
		if counter, ok := self.Buckets[id]; !ok {
			self.Buckets[id] = newIntCounter(at)
		} else {
			counter.Event(at)
		}
	}
}

func (self *linearCounter) Summary(at time.Time) interface{} {
	self.Lock()
	defer self.Unlock()

	results := map[string]interface{}{}

	for topic, counter := range self.Buckets {
		results[topic] = counter.Summary(at)
	}

	return results
}

func newIntCounter(at time.Time) *intCounter {
	return &intCounter{
		Start:   at.Truncate(time.Minute),
		Minutes: []uint32{1},
	}
}

type intCounter struct {
	Start   time.Time
	Minutes []uint32
}

func (self *intCounter) Event(at time.Time) {
	bucket := self.fillForward(at)

	self.Minutes[bucket] += 1
}

func (self *intCounter) Summary(at time.Time) interface{} {
	self.fillForward(at)

	return self
}

func (self *intCounter) fillForward(at time.Time) int {
	bucket := int(at.Truncate(time.Minute).Sub(self.Start).Minutes())

	for bucket >= len(self.Minutes) {
		self.Minutes = append(self.Minutes, 0)
	}

	return bucket
}

func (self *intCounter) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"start":   self.Start,
		"minutes": self.Minutes,
	})
}
