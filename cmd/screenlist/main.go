package main

import (
	"os"
	"sync"

	"github.com/ReanSn0w/go-screenlist/pkg/list"
	"github.com/ReanSn0w/go-screenlist/pkg/params"
	"github.com/ReanSn0w/go-screenlist/pkg/utils"
	"github.com/ReanSn0w/go-screenlist/pkg/video"
	"github.com/go-pkgz/lgr"
)

var (
	log utils.Logger
)

func main() {
	params := params.Load(lgr.Default())
	log = params.Log()

	wg := sync.WaitGroup{}
	wg.Add(len(params.Files))

	for _, file := range params.Files {
		go func(file string) {
			defer wg.Done()

			log.Logf("[INFO] Processing file specs for: %s", file)
			specs, err := video.Specs(file)
			if err != nil {
				log.Logf("[ERROR] file: %s err: %v", file, err)
				if !params.Force {
					return
				}
			}

			log.Logf("[INFO] Processing file screenshots for: %s", file)
			images, err := video.Load(file, params.Screenshots, !params.Delta)
			if err != nil {
				log.Logf("[ERROR] file: %s err: %v", file, err)
				if !params.Force {
					return
				}
			}

			log.Logf("[INFO] Saving screenlist for: %s", file)
			err = list.Save(file+"_screenlist.jpg", specs, params.Grid, params.ResultWidth, images)
			if err != nil {
				log.Logf("[ERROR] %v", err)
			}
		}(file)
	}

	wg.Wait()
	os.Exit(0)
}
