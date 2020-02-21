package metrics

import (
	"encoding/json"
	"net/http"
	"time"
)

type Tracker interface {
	Event(id string, value interface{})
	Summary() interface{}
}

func NewCollector() *collector {
	return &collector{NewHistogram()}
}

type collector struct {
	Tracker
}

func (self *collector) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(self.Tracker.Summary())
}

func (self *collector) Handle(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer self.track(r, time.Now())

		next.ServeHTTP(w, r)
	}
}

func (self *collector) track(r *http.Request, start time.Time) {
	self.Tracker.Event(r.URL.Path, time.Since(start).Seconds())
}
