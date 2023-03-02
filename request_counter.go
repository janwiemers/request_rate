package requestrate

type RequestCounter struct {
	count    int
	duration int64
}

func (r *RequestCounter) Incr(counter int) {
	r.count = r.count + counter
}

// func (r *RequestCounter)
