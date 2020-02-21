package metrics

import (
	"encoding/json"
	"sync"
	"time"
)

func NewTreeCounter() *treeCounter {
	return &treeCounter{
		Buckets: map[string]*yearCounter{},
	}
}

type treeCounter struct {
	sync.Mutex
	Buckets map[string]*yearCounter
}

func (self *treeCounter) Event(at time.Time, ids ...string) {
	self.Lock()
	defer self.Unlock()

	for _, id := range ids {
		if counter, ok := self.Buckets[id]; !ok {
			self.Buckets[id] = newYearCounter(at)
		} else {
			counter.Event(at)
		}
	}
}

func (self *treeCounter) Summary(at time.Time) interface{} {
	self.Lock()
	defer self.Unlock()

	return self.Buckets
}

func newYearCounter(at time.Time) *yearCounter {
	return &yearCounter{
		Buckets: map[int]*monthCounter{
			at.Year(): newMonthCounter(at),
		},
	}
}

type yearCounter struct {
	Buckets map[int]*monthCounter
}

func (self *yearCounter) Event(at time.Time) {
	if counter, ok := self.Buckets[at.Year()]; !ok {
		self.Buckets[at.Year()] = newMonthCounter(at)
	} else {
		counter.Event(at)
	}
}

func (self *yearCounter) MarshalJSON() ([]byte, error) {
	return json.Marshal(self.Buckets)
}

func newMonthCounter(at time.Time) *monthCounter {
	return &monthCounter{
		Buckets: map[int]*dayCounter{
			int(at.Month()): newDayCounter(at),
		},
	}
}

type monthCounter struct {
	Buckets map[int]*dayCounter
}

func (self *monthCounter) Event(at time.Time) {
	if counter, ok := self.Buckets[int(at.Month())]; !ok {
		self.Buckets[int(at.Month())] = newDayCounter(at)
	} else {
		counter.Event(at)
	}
}

func (self *monthCounter) MarshalJSON() ([]byte, error) {
	return json.Marshal(self.Buckets)
}

func newDayCounter(at time.Time) *dayCounter {
	return &dayCounter{
		Buckets: map[int]*hourCounter{
			at.Day(): newHourCounter(at),
		},
	}
}

type dayCounter struct {
	Buckets map[int]*hourCounter
}

func (self *dayCounter) Event(at time.Time) {
	if counter, ok := self.Buckets[at.Day()]; !ok {
		self.Buckets[at.Day()] = newHourCounter(at)
	} else {
		counter.Event(at)
	}
}

func (self *dayCounter) MarshalJSON() ([]byte, error) {
	return json.Marshal(self.Buckets)
}

func newHourCounter(at time.Time) *hourCounter {
	return &hourCounter{
		Buckets: map[int]*minuteCounter{
			at.Hour(): newMinuteCounter(at),
		},
	}
}

type hourCounter struct {
	Buckets map[int]*minuteCounter
}

func (self *hourCounter) Event(at time.Time) {
	if counter, ok := self.Buckets[at.Hour()]; !ok {
		self.Buckets[at.Hour()] = newMinuteCounter(at)
	} else {
		counter.Event(at)
	}
}

func (self *hourCounter) MarshalJSON() ([]byte, error) {
	return json.Marshal(self.Buckets)
}

func newMinuteCounter(at time.Time) *minuteCounter {
	return &minuteCounter{
		Buckets: map[int]int64{at.Minute(): 1},
	}
}

type minuteCounter struct {
	Buckets map[int]int64
}

func (self *minuteCounter) Event(at time.Time) {
	self.Buckets[at.Minute()] += 1
}

func (self *minuteCounter) MarshalJSON() ([]byte, error) {
	return json.Marshal(self.Buckets)
}
