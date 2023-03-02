package requestrate

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Observables provides a type to observe requests
type Observables struct {
	items map[string]time.Time
}

// Begin starts measuring a request
func (o *Observables) Begin() string {
	uuid := uuid.New().String()
	o.items[uuid] = time.Now()
	return uuid
}

// Finish will
func (o *Observables) Finish(uuid string) (int64, error) {
	if val, ok := o.items[uuid]; ok {
		delete(o.items, uuid)
		return time.Since(val).Milliseconds(), nil
	}

	return 0, errors.New("item does not exist")
}
