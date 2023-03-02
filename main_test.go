package requestrate

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_NewRequestRate(t *testing.T) {
	a := NewRequestRate()
	assert.IsType(t, a, &RequestRate{})
}

func Test_Incr(t *testing.T) {
	a := NewRequestRate()
	a.Incr(1)
	assert.Equal(t, a.currentCounter.count, 1)
}

func Test_PruneHistory(t *testing.T) {
	a := NewRequestRate()
	a.Start()
	a.Incr(1)
	time.Sleep(2 * time.Second)
	assert.Equal(t, 0, a.currentCounter.count)
	assert.Equal(t, 1, a.counterHistory[0].count)
	assert.Equal(t, 2, len(a.counterHistory))
	a.PruneHistory()
	assert.Equal(t, len(a.counterHistory), 0)
}

func Test_Rate(t *testing.T) {
	a := NewRequestRate()
	a.Start()
	a.Incr(1)
	time.Sleep(1010 * time.Millisecond)
	assert.Equal(t, a.latestCounter.count, 1)
	assert.Equal(t, a.Rate(), 1)
}

func Test_Observerables(t *testing.T) {
	a := NewRequestRate()
	b := a.Observe(1)
	time.Sleep(time.Second * 1)
	a.Finish(b)

	assert.GreaterOrEqual(t, a.currentCounter.duration, int64(1000))
}

func Test_ObserverablesError(t *testing.T) {
	a := NewRequestRate()
	_, err := a.Finish("asd")
	assert.NotNil(t, err)
}
