package requestrate

import (
	"sync"
	"time"
)

type RequestRate struct {
	mu             sync.Mutex
	latestCounter  RequestCounter
	counterHistory []RequestCounter
	currentCounter RequestCounter
	observables    Observables
}

func NewRequestRate() *RequestRate {
	return &RequestRate{
		mu: sync.Mutex{},
		observables: Observables{
			items: map[string]time.Time{},
		},
	}
}

func (r *RequestRate) Incr(counter int) {
	r.mu.Lock()
	r.currentCounter.Incr(counter)
	r.mu.Unlock()
}

func (r *RequestRate) Observe(counter int) string {
	r.mu.Lock()
	uuid := r.observables.Begin()
	r.currentCounter.Incr(counter)
	r.mu.Unlock()
	return uuid
}

func (r *RequestRate) Finish(uuid string) (int64, error) {
	r.mu.Lock()
	d, err := r.observables.Finish(uuid)
	if err != nil {
		return d, err
	}

	r.currentCounter.duration = r.currentCounter.duration + d
	r.mu.Unlock()
	return d, nil
}

func (r *RequestRate) Start() {
	go r.ticker()
}

func (r *RequestRate) ticker() {
	for range time.Tick(time.Second * 1) {
		r.mu.Lock()
		r.latestCounter = r.currentCounter
		r.currentCounter.count = 0
		r.currentCounter.duration = 0
		r.mu.Unlock()

		r.counterHistory = append(r.counterHistory, r.latestCounter)
	}
}

func (r *RequestRate) Rate() int {
	return r.latestCounter.count
}

func (r *RequestRate) PruneHistory() {
	r.mu.Lock()
	r.counterHistory = make([]RequestCounter, 0)
	r.mu.Unlock()
}
