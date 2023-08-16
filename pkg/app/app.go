package app

import (
	"github.com/ReanSn0w/go-screenlist/pkg/params"
	"github.com/ReanSn0w/go-screenlist/pkg/utils"
)

func New(log utils.Logger, pref *params.Parameters) *Application {
	return &Application{
		log:  log,
		pref: pref,
	}
}

type Application struct {
	log  utils.Logger
	pref *params.Parameters
}

func (a *Application) Run() {
	s := utils.NewScheduler[task, task](a.pref.Treads, a.processTasks)
	s.Push(a.createTasks()...)

	go func() {
		for task := range s.Out() {
			a.log.Logf("[INFO] Processed file: %s", task.filename)
		}
	}()

	s.Wait()
}

func (a *Application) createTasks() []task {
	res := []task{}

	for _, file := range a.pref.Files {
		task := newTask(file)
		res = append(res, task)
	}

	return res
}

func (a *Application) processTasks(t task) task {
	t.process(a.log, a.pref)
	return t
}
