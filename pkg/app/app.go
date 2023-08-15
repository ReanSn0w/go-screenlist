package app

import (
	"sync"

	"github.com/ReanSn0w/go-screenlist/pkg/params"
	"github.com/ReanSn0w/go-screenlist/pkg/utils"
)

func New(log utils.Logger, pref *params.Parameters) *Application {
	return &Application{
		log:  log,
		pref: pref,
		wg:   &sync.WaitGroup{},
		rl:   utils.NewRoutineLimiter(pref.Treads),
	}
}

type Application struct {
	log  utils.Logger
	pref *params.Parameters
	wg   *sync.WaitGroup
	rl   *utils.RoutineLimiter
}

func (a *Application) Run() error {
	es := utils.ErrorStack{}
	a.wg.Add(1)

	ch := make(chan *task, a.pref.Treads)

	go a.createTasks(ch)
	go a.createWorkers(ch)

	a.wg.Wait()
	return es.Get()
}

func (a *Application) createTasks(ch chan<- *task) {
	for _, file := range a.pref.Files {
		task := newTask(file)
		ch <- task
	}

	// Последний элемент канала
	// задача которого завершить выполнение программы
	ch <- nil
}

func (a *Application) createWorkers(in <-chan *task) {
	for task := range in {
		a.wg.Add(1)

		if task == nil {
			a.wg.Done()
			a.wg.Done()
			break
		}

		a.rl.Run(func() {
			task.process(a.pref.Log(), a.pref)
			a.wg.Done()
		})
	}
}
