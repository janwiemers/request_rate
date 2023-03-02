package requestrate

import (
	"sync"
	"time"
)

// RequestRate provides the basic structure for the request rate class.
// Main entry point should alway be "NewRequestRate()" to ensure all properties
// are set up as intented.
// To start the counting you need to call the Start() API.
type RequestRate struct {
	mu             sync.Mutex
	latestCounter  RequestCounter
	counterHistory []RequestCounter
	currentCounter RequestCounter
	observables    Observables
}

// NewRequestRate returns a new RequestRate object. Advice is to only have
// one Object per application instance. The solution is thread safe and relies on
// sync.Mutex.
//
// To track counts use Incr()
// To track counts and measure durations use Observe()
func NewRequestRate() *RequestRate {
	return &RequestRate{
		mu: sync.Mutex{},
		observables: Observables{
			items: map[string]time.Time{},
		},
	}
}

// Incr increases the counter. It is an integer which can be 1 or more in
// case you want to batch or sample.
func (r *RequestRate) Incr(counter int) {
	r.mu.Lock()
	r.currentCounter.Incr(counter)
	r.mu.Unlock()
}

// Observe increments the counter but also measures the duration. It does so by
// storing and returning a UUIDv4. The call to the "Finish()" API takes the UUIDv4
// to calculate the Duration
func (r *RequestRate) Observe(counter int) string {
	r.mu.Lock()
	uuid := r.observables.Begin()
	r.currentCounter.Incr(counter)
	r.mu.Unlock()
	return uuid
}

// Finish takes the UUIDv4 returned from the "Observe()" API and returns the duration
// in Miiliseconds as well as an error.
// Error will be in cases where the UUIDv4 cannot be found.
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

// Starts the ticker which does the computation and history generation
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

// Rate, returns the rate per second
func (r *RequestRate) Rate() int {
	return r.latestCounter.count
}

// PruneHistory deleted the History
func (r *RequestRate) PruneHistory() {
	r.mu.Lock()
	r.counterHistory = make([]RequestCounter, 0)
	r.mu.Unlock()
}
