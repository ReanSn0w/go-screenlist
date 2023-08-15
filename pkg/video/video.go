package video

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"os"

	"github.com/ReanSn0w/go-screenlist/pkg/ffmpeg"
	"github.com/ReanSn0w/go-screenlist/pkg/utils"
)

type Spec struct {
	Title      string
	Resolution [2]int
	Fps        int
	Size       int
	Duration   int
}

func Specs(file string) (*Spec, error) {
	seconds, err := ffmpeg.Seconds(file)
	if err != nil {
		return nil, err
	}

	resolution, err := ffmpeg.Resolution(file)
	if err != nil {
		return nil, err
	}

	fps, err := ffmpeg.FPS(file)
	if err != nil {
		return nil, err
	}

	size, err := ffmpeg.Filesize(file)
	if err != nil {
		return nil, err
	}

	return &Spec{
		Title:      file,
		Resolution: resolution,
		Fps:        fps,
		Size:       size,
		Duration:   seconds,
	}, nil
}

// Load return video screenshots
func Load(file string, count int, removeOriginals bool) ([]image.Image, error) {
	if count == 0 {
		return nil, errors.New("count must be greater than 0")
	}

	seconds, err := ffmpeg.Seconds(file)
	if err != nil {
		return nil, err
	}

	stack := utils.NewErrorStack()

	images, err := makeframes(file, seconds, count)
	if len(images) == 0 {
		return nil, errors.New("no frames captured")
	}
	stack.Add(err)

	if !removeOriginals {
		stack.Add(saveImages(file, images))
	}

	return images, stack.Get()
}

func frames(seconds int, count int) []int {
	seconds = seconds - seconds/5 // remove last 5% of the video, usually credits or black screen
	frame := seconds / (count + 1)

	frames := make([]int, count)
	for i := 0; i < count; i++ {
		frames[i] = frame * (i + 1)
	}

	return frames
}

func makeframes(file string, seconds int, count int) ([]image.Image, error) {
	es := utils.NewErrorStack()

	images := make([]image.Image, 0)
	frames := frames(seconds, count)

	for _, frame := range frames {
		img, err := ffmpeg.Shot(file, frame)
		if err != nil {
			es.Add(err)
		}

		images = append(images, img)
	}

	return images, es.Get()
}

func saveImages(filename string, images []image.Image) error {
	es := utils.ErrorStack{}

	for i, img := range images {
		f, err := os.Create(fmt.Sprintf("%v_%v.jpg", filename, i))
		if err != nil {
			es.Add(err)
			continue
		}
		defer f.Close()

		err = jpeg.Encode(f, img, &jpeg.Options{Quality: 75})
		es.Add(err)
	}

	return es.Get()
}
