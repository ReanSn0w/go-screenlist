package video

import (
	"errors"
	"go-screenlist/pkg/utils"
	"image"
	"image/jpeg"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/mowshon/moviego"
)

// Load return video screenshots
func Load(file string, count int, removeOriginals bool) ([]image.Image, error) {
	if count == 0 {
		return nil, errors.New("count must be greater than 0")
	}

	seconds, err := Seconds(file)
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

// Seconds returns the length of a video in seconds
func Seconds(file string) (int, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", file)
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(strings.Split(string(output), ".")[0])
}

// Resolution returns the resolution of a video
func Resolution(file string) ([2]int, error) {
	// ffprobe -v error -show_entries stream=width,height -of csv=p=0:s=x test.mp4
	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "stream=width,height", "-of", "csv=p=0:s=x", file)
	output, err := cmd.Output()
	if err != nil {
		return [2]int{}, err
	}

	vals := strings.Split(string(output), "x")
	if len(vals) != 2 {
		return [2]int{}, errors.New("invalid resolution")
	}

	width, err := strconv.Atoi(vals[0])
	if err != nil {
		return [2]int{}, err
	}

	height, err := strconv.Atoi(strings.Split(vals[1], "\n")[0])
	if err != nil {
		return [2]int{}, err
	}

	return [2]int{width, height}, nil
}

// Fps returns the fps of a video
func Fps(file string) (int, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=r_frame_rate", "-of", "default=noprint_wrappers=1:nokey=1", file)
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	fps, err := strconv.Atoi(strings.Split(strings.Split(string(output), "/")[0], "\n")[0])
	if err != nil {
		return 0, err
	}

	return fps, nil
}

// Bitrate returns the bitrate of a video
func Bitrate(file string) (int, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=bit_rate", "-of", "default=noprint_wrappers=1:nokey=1", file)
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	bitrate, err := strconv.Atoi(strings.Split(string(output), "\n")[0])
	if err != nil {
		return 0, err
	}

	return bitrate, nil
}

func Filesize(file string) (int, error) {
	info, err := os.Stat(file)
	if err != nil {
		return 0, err
	}

	return int(info.Size()), nil
}

type Spec struct {
	Title      string
	Resolution [2]int
	Fps        int
	Size       int
	Duration   int
}

func Specs(file string) (*Spec, error) {
	seconds, err := Seconds(file)
	if err != nil {
		return nil, err
	}

	resolution, err := Resolution(file)
	if err != nil {
		return nil, err
	}

	fps, err := Fps(file)
	if err != nil {
		return nil, err
	}

	size, err := Filesize(file)
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
