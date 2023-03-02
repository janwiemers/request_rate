package requestrate

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Observables struct {
	items map[string]time.Time
}

func (o *Observables) Begin() string {
	uuid := uuid.New().String()
	o.items[uuid] = time.Now()
	return uuid
}

func (o *Observables) Finish(uuid string) (int64, error) {
	if val, ok := o.items[uuid]; ok {
		delete(o.items, uuid)
		return time.Since(val).Milliseconds(), nil
	}

	return 0, errors.New("item does not exist")
}
