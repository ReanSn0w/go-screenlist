package utils

func NewRoutineLimiter(limit int) *RoutineLimiter {
	return &RoutineLimiter{
		ch: make(chan struct{}, limit),
	}
}

type RoutineLimiter struct {
	ch chan struct{}
}

func (r *RoutineLimiter) Run(f func()) {
	r.ch <- struct{}{}
	go func() {
		f()
		<-r.ch
	}()
}
