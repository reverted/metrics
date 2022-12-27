package metrics

import (
	"sync"
	"time"
)

func NewHistogram() *histogram {
	return &histogram{
		Start:   time.Now().UTC(),
		Buckets: map[string]*eventCounter{},
	}
}

type histogram struct {
	sync.RWMutex
	Start   time.Time
	Buckets map[string]*eventCounter
}

func (self *histogram) Event(id string, value interface{}) {
	self.Lock()
	defer self.Unlock()

	if counter, ok := self.Buckets[id]; !ok {
		self.Buckets[id] = newEventCounter(value)
	} else {
		counter.Event(value)
	}
}

func (self *histogram) Summary() interface{} {
	self.RLock()
	defer self.RUnlock()

	var results []map[string]interface{}
	for k, v := range self.Buckets {
		results = append(results, map[string]interface{}{
			"path":   k,
			"events": v.Events,
		})
	}
	return results
}

type event struct {
	At    time.Time   `json:"at"`
	Value interface{} `json:"value"`
}

func newEventCounter(value interface{}) *eventCounter {
	counter := &eventCounter{
		Events: []event{},
	}

	counter.Event(value)

	return counter
}

type eventCounter struct {
	sync.RWMutex
	Events []event
}

func (self *eventCounter) Event(value interface{}) {
	self.Lock()
	defer self.Unlock()

	self.Events = append(self.Events, event{
		At:    time.Now().UTC(),
		Value: value,
	})
}
