package app

import (
	"strings"

	"github.com/ReanSn0w/go-screenlist/pkg/list"
	"github.com/ReanSn0w/go-screenlist/pkg/params"
	"github.com/ReanSn0w/go-screenlist/pkg/utils"
	"github.com/ReanSn0w/go-screenlist/pkg/video"
)

func newTask(file string) *task {
	fileparts := strings.Split(file, "/")

	return &task{
		filename: fileparts[len(fileparts)-1],
		path:     file,
	}
}

type task struct {
	filename string
	path     string
}

func (t *task) process(log utils.Logger, pref *params.Parameters) {
	log.Logf("[INFO] Processing file specs for: %s", t.path)
	specs, err := video.Specs(t.path)
	if err != nil {
		log.Logf("[ERROR] file: %s err: %v", t.filename, err)
		if !pref.Force {
			return
		}
	}

	log.Logf("[INFO] Processing file screenshots for: %s", t.path)
	images, err := video.Load(t.path, pref.Screenshots, pref.RemoveOriginals())
	if err != nil {
		log.Logf("[ERROR] file: %s err: %v", t.path, err)
		if !pref.Force {
			return
		}
	}

	if pref.ScreenshotMode {
		return
	}

	log.Logf("[INFO] Saving screenlist for: %s", t.path)
	err = list.Save(t.path+"_screenlist.jpg", specs, pref.Grid, pref.ResultWidth, images)
	if err != nil {
		log.Logf("[ERROR] %v", err)
	}
}
