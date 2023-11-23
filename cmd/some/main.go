package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ReanSn0w/go-screenlist/pkg/ffmpeg"
	"github.com/go-pkgz/lgr"
)

const (
	filename = "video.mp4"
	total    = 1440
)

func main() {
	lgr.Setup(lgr.Debug)

	seconds, err := ffmpeg.Seconds(filename)
	if err != nil {
		lgr.Default().Logf("[ERROR] %v", err)
	}

	seconds *= 10

	tasks(context.Background(), seconds, 4, func(i int) {
		err := ffmpeg.ShotPNG(filename, fmt.Sprintf("image-test/%v.png", i), time.Second/10*time.Duration(i))
		if err != nil {
			lgr.Default().Logf("[WARN] failed to set screenshot #%v. err: %v", i, err)
		} else {
			lgr.Default().Logf("[DEBUG] shot done for %v", i)
		}
	})
}

func tasks(ctx context.Context, count int, quenue int, run func(i int)) {
	wg := &sync.WaitGroup{}
	wg.Add(count)

	ch := make(chan func(), quenue)
	defer close(ch)

	// Context close
	force := false
	go func() {
		<-ctx.Done()
		force = true
	}()

	go func() {
		// task builder

		for t := 0; t < count; t++ {
			val := t

			if force {
				lgr.Default().Logf("[INFO] task %v force closed by context", t)
				wg.Done()
				continue
			}

			ch <- func() {
				run(val)
				wg.Done()
			}
		}
	}()

	go func() {
		// runner

		for val := range ch {
			val()
		}
	}()

	wg.Wait()
}
