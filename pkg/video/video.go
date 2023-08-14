package video

import (
	"errors"
	"image"
	"image/jpeg"
	"os"
	"strconv"

	"github.com/ReanSn0w/go-screenlist/pkg/ffmpeg"
	"github.com/ReanSn0w/go-screenlist/pkg/utils"
	"github.com/mowshon/moviego"
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

	framePaths, err := makeframes(file, seconds, count)
	if len(framePaths) == 0 {
		return nil, errors.New("no frames captured")
	}
	stack.Add(err)

	images, err := loadImages(framePaths, removeOriginals)
	stack.Add(err)

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

func makeframes(file string, seconds int, count int) ([]string, error) {
	video, err := moviego.Load(file)
	if err != nil {
		return nil, err
	}

	es := utils.NewErrorStack()

	paths := make([]string, count) // generate frames
	frames := frames(seconds, count)

	for i, frame := range frames {
		abs, err := video.Screenshot(float64(frame), file+"-screenshot-"+strconv.Itoa(i)+".jpg")
		if err != nil {
			es.Add(err)
			continue
		}

		paths[i] = abs
	}

	return paths, es.Get()
}

func loadImages(paths []string, deleteOriginals bool) ([]image.Image, error) {
	result := []image.Image{}
	es := utils.NewErrorStack()

	for _, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			es.Add(err)
			continue
		}

		image, err := jpeg.Decode(file)
		if err != nil {
			es.Add(err)
			continue
		}

		result = append(result, image)

		if deleteOriginals {
			err = os.Remove(path)
			es.Add(err)
		}
	}

	return result, es.Get()
}
