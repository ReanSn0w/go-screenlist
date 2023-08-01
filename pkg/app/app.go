package app

import (
	"github.com/ReanSn0w/go-screenlist/pkg/list"
	"github.com/ReanSn0w/go-screenlist/pkg/params"
	"github.com/ReanSn0w/go-screenlist/pkg/utils"
	"github.com/ReanSn0w/go-screenlist/pkg/video"
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

func (a *Application) Run() error {
	ch := make(chan *task, a.pref.Treads)

	go a.createTasks(ch)
	a.createWorkers(ch)

	return nil
}

func (a *Application) createTasks(ch chan<- *task) {
	for _, file := range a.pref.Files {
		task := newTask(file)
		ch <- task
	}

	// Последний элемент канала
	// задача которого завершить выполнение программы
	close(ch)
}

func (a *Application) createWorkers(in <-chan *task) {
	rl := utils.NewRoutineLimiter(a.pref.Treads)

	for task := range in {
		if len(in) == 0 {
			break
		}

		rl.Run(func() {
			a.processTask(task)
		})
	}
}

func (a *Application) processTask(task *task) {
	if task == nil {
		return
	}

	a.log.Logf("[INFO] Processing file specs for: %s", task.path)
	specs, err := video.Specs(task.path)
	if err != nil {
		a.log.Logf("[ERROR] file: %s err: %v", task.filename, err)
		if !a.pref.Force {
			return
		}
	}

	a.log.Logf("[INFO] Processing file screenshots for: %s", task.path)
	images, err := video.Load(task.path, a.pref.Screenshots, a.pref.RemoveOriginals())
	if err != nil {
		a.log.Logf("[ERROR] file: %s err: %v", task.path, err)
		if !a.pref.Force {
			return
		}
	}

	if a.pref.ScreenshotMode {
		return
	}

	a.log.Logf("[INFO] Saving screenlist for: %s", task.path)
	err = list.Save(task.path+"_screenlist.jpg", specs, a.pref.Grid, a.pref.ResultWidth, images)
	if err != nil {
		a.log.Logf("[ERROR] %v", err)
	}
}
