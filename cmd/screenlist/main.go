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

	query := make(chan []string)
	done := make(chan bool)

	// make pages
	go func() {
		files := len(params.Files)
		pages := files / params.Treads
		if files%params.Treads != 0 {
			pages++
		}

		for i := 0; i <= pages; i++ {
			from := i * params.Treads
			to := from + params.Treads
			if to > files {
				query <- params.Files[from:]
				continue
			}

			query <- params.Files[from:to]
		}

		done <- true
	}()

	wg := sync.WaitGroup{}
	for {
		wg.Wait()

		select {
		case files := <-query:
			wg.Add(len(files))

			for _, file := range files {
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
					images, err := video.Load(file, params.Screenshots, !params.Delta && !params.ScreenshotMode)
					if err != nil {
						log.Logf("[ERROR] file: %s err: %v", file, err)
						if !params.Force {
							return
						}
					}

					if params.ScreenshotMode {
						return
					}

					log.Logf("[INFO] Saving screenlist for: %s", file)
					err = list.Save(file+"_screenlist.jpg", specs, params.Grid, params.ResultWidth, images)
					if err != nil {
						log.Logf("[ERROR] %v", err)
					}
				}(file)
			}
		case <-done:
			os.Exit(0)
		}
	}
}
