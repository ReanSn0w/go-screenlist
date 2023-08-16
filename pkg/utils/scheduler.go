package utils

import "sync"

func NewScheduler[I any, O any](parallel int, task func(I) O) *Scheduler[I, O] {
	return &Scheduler[I, O]{
		wg:   &sync.WaitGroup{},
		rl:   NewRoutineLimiter(parallel),
		fifo: NewFifo[I](),
		out:  make(chan O),
		task: task,
	}
}

type Scheduler[I any, O any] struct {
	wg   *sync.WaitGroup
	rl   *RoutineLimiter
	fifo *Fifo[I]
	out  chan O

	task func(I) O
}

func (s *Scheduler[I, O]) Push(items ...I) {
	s.wg.Add(len(items))
	s.fifo.Push(items...)

	go func() {
		for {
			item := s.fifo.Pop()
			if item == nil {
				break
			}

			s.rl.Run(func() {
				s.out <- s.task(*item)
				s.wg.Done()
			})
		}
	}()
}

func (s *Scheduler[I, O]) Out() <-chan O {
	return s.out
}

func (s *Scheduler[I, O]) Wait() {
	s.wg.Wait()
}
