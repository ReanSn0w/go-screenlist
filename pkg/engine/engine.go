package engine

import (
	"sync"

	"github.com/ReanSn0w/go-screenlist/pkg/utils"
	"github.com/ReanSn0w/tk4go/pkg/tools"
	"github.com/nfnt/resize"
)

func New(log tools.Logger, force bool, threads int, screenlist *ScreenListPreferences, delta *DeltaPreferences) *Engine {
	ffmpeg := utils.NewFFMPEG(log)
	list := utils.NewList(log, ffmpeg)

	return &Engine{
		log:    log,
		force:  force,
		rl:     tools.NewRoutineLimiter(threads),
		wg:     &sync.WaitGroup{},
		ffmpeg: ffmpeg,
		list:   list,
		preferences: struct {
			screenlist *ScreenListPreferences
			delta      *DeltaPreferences
		}{
			screenlist: screenlist,
			delta:      delta,
		},
	}
}

type (
	Engine struct {
		force       bool
		log         tools.Logger
		rl          *tools.RoutineLimiter
		wg          *sync.WaitGroup
		ffmpeg      *utils.FFMPEG
		list        *utils.List
		preferences struct {
			screenlist *ScreenListPreferences
			delta      *DeltaPreferences
		}
	}

	ScreenListPreferences struct {
		Enabled bool   `long:"enabled" env:"ENABLED" description:"enable screenlist"`
		Info    bool   `long:"info" env:"INFO" description:"enable info"`
		Images  int    `long:"images" env:"IMAGES" default:"15" description:"images directory"`
		Grid    int    `long:"grid" env:"GRID" default:"3" description:"grid size"`
		Width   int    `long:"width" env:"WIDTH" default:"1200" description:"resulting image width"`
		Result  string `long:"result" env:"RESULT" default:"{{.Name}}_screenlist.jpg" description:"resulting image name. Use {{.Name}} for filename prefix"`
	}

	DeltaPreferences struct {
		Enabled bool   `long:"enabled" env:"ENABLED" description:"enable delta saving"`
		Images  int    `long:"images" env:"IMAGES" default:"15" description:"images directory"`
		Width   int    `long:"width" env:"WIDTH" default:"1200" description:"resulting image width"`
		Result  string `long:"result" env:"RESULT" default:"{{.Name}}_screenshot_{{.Counter}}.jpg" description:"resulting image name with counter. Use {{.Counter}} for counter value and {{.Name}} for filename prefix"`
	}
)

func (e *Engine) Add(files ...utils.File) {
	e.wg.Add(len(files))

	for _, file := range files {
		e.rl.Run(func() {
			defer e.wg.Done()
			e.process(file)
		})
	}
}

func (e *Engine) Wait() {
	e.wg.Wait()
}

func (e *Engine) process(file utils.File) error {
	e.log.Logf("[INFO] processing file %v", file)

	spec, err := e.ffmpeg.Specs(string(file))
	if err != nil {
		e.log.Logf("[ERROR] get spec for %s: %v", file, err)
		return err
	}

	if e.preferences.screenlist.Enabled {
		images, err := spec.Load(e.ffmpeg, string(file), e.preferences.screenlist.Images)
		if err != nil {
			e.log.Logf("[ERROR] get images for %s: %v", file, err)
			return err
		}

		spec := spec
		if !e.preferences.screenlist.Info {
			spec = nil
		}

		img, err := e.list.Make(spec, e.preferences.screenlist.Grid, e.preferences.screenlist.Width, images)
		if err != nil {
			e.log.Logf("[ERROR] save screenlist for %s: %v", file, err)
			return err
		}

		err = file.SaveImagesByPattern(e.preferences.screenlist.Result, img)
		if err != nil {
			e.log.Logf("[ERROR] save screenlist for %s: %v", file, err)
			return err
		}
	}

	if e.preferences.delta.Enabled {
		images, err := spec.Load(e.ffmpeg, string(file), e.preferences.screenlist.Images)
		if err != nil {
			e.log.Logf("[ERROR] get images for %s: %v", file, err)
			return err
		}

		for i := range images {
			images[i] = resize.Resize(uint(e.preferences.delta.Width), 0, images[i], resize.Lanczos3)
		}

		err = file.SaveImagesByPattern(e.preferences.delta.Result, images...)
		if err != nil {
			e.log.Logf("[ERROR] save screenlist for %s: %v", file, err)
			return err
		}
	}

	return nil
}
