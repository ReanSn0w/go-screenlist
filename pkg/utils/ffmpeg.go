package utils

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/ReanSn0w/tk4go/pkg/tools"
)

func NewFFMPEG(log tools.Logger) *FFMPEG {
	return &FFMPEG{log: log}
}

type (
	FFMPEG struct {
		log tools.Logger
	}

	Spec struct {
		Title      string
		Resolution [2]int
		Fps        int
		Size       int
		Duration   int
	}
)

func (f *FFMPEG) Specs(file string) (*Spec, error) {
	seconds, err := f.Seconds(file)
	if err != nil {
		f.log.Logf("[ERROR] get video duration falied: %v", err)
		return nil, err
	}

	resolution, err := f.Resolution(file)
	if err != nil {
		f.log.Logf("[ERROR] get video resolution falied: %v", err)
		return nil, err
	}

	fps, err := f.FPS(file)
	if err != nil {
		f.log.Logf("[ERROR] get video fps falied: %v", err)
		return nil, err
	}

	size, err := f.Filesize(file)
	if err != nil {
		f.log.Logf("[ERROR] get video size falied: %v", err)
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

// Shot takes a screenshot of a video at a given second
func (f *FFMPEG) Shot(input string, s int) (image.Image, error) {
	output, err := f.command("ffmpeg", "-ss", f.seconds(s), "-i", input, "-vframes", "1", "-q:v", "2", "-f", "image2pipe", "-vcodec", "png", "-")
	if err != nil {
		return nil, err
	}

	img, err := png.Decode(strings.NewReader(output))
	if err != nil {
		return nil, err
	}

	return img, nil
}

// ShotJpeg takes a screenshot of a video at a given second and saves it as a jpeg
func (f *FFMPEG) ShotJPEG(input string, output string, s int, quality int) error {
	return f.shotImage(input, output, s, func(w io.Writer, img image.Image) error {
		return jpeg.Encode(w, img, &jpeg.Options{Quality: quality})
	})
}

func (f *FFMPEG) ShotPNG(input string, output string, s int) error {
	return f.shotImage(input, output, s, func(w io.Writer, img image.Image) error {
		return png.Encode(w, img)
	})
}

// func ShotWEBP(input string, output string, s int, quality int) error {
// 	return shotImage(input, output, s, func(w io.Writer, img image.Image) error {
// 		return webp.Encode(w, img, &webp.Options{Lossless: false, Quality: float32(quality)})
// 	})
// }

func (f *FFMPEG) shotImage(input string, output string, s int, prepare func(io.Writer, image.Image) error) error {
	img, err := f.Shot(input, s)
	if err != nil {
		return err
	}

	file, err := os.Create(output)
	if err != nil {
		return err
	}
	defer file.Close()

	return prepare(file, img)
}

// FPS returns fps of video
func (f *FFMPEG) FPS(file string) (int, error) {
	output, err := f.command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=r_frame_rate", "-of", "default=noprint_wrappers=1:nokey=1", file)
	if err != nil {
		return 0, err
	}

	fps, err := strconv.Atoi(strings.Split(strings.Split(string(output), "/")[0], "\n")[0])
	if err != nil {
		return 0, err
	}

	return fps, nil
}

// Seconds returns the length of a video in seconds
func (f *FFMPEG) Seconds(file string) (int, error) {
	output, err := f.command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", file)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(strings.Split(string(output), ".")[0])
}

// Resolution returns the resolution of a video
func (f *FFMPEG) Resolution(file string) ([2]int, error) {
	output, err := f.command("ffprobe", "-v", "error", "-show_entries", "stream=width,height", "-of", "csv=p=0:s=x", file)
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

// Bitrate returns the bitrate of a video
func (f *FFMPEG) Bitrate(file string) (int, error) {
	output, err := f.command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=bit_rate", "-of", "default=noprint_wrappers=1:nokey=1", file)
	if err != nil {
		return 0, err
	}

	bitrate, err := strconv.Atoi(strings.Split(string(output), "\n")[0])
	if err != nil {
		return 0, err
	}

	return bitrate, nil
}

func (f *FFMPEG) Filesize(file string) (int, error) {
	info, err := os.Stat(file)
	if err != nil {
		return 0, err
	}

	return int(info.Size()), nil
}

func (f *FFMPEG) command(name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)
	out, err := cmd.Output()

	if err != nil {

	}

	return string(out), err
}

// int to format 00:00:00
func (f *FFMPEG) seconds(s int) string {
	h := s / 3600
	s -= h * 3600
	m := s / 60
	s -= m * 60
	return strconv.Itoa(h) + ":" + strconv.Itoa(m) + ":" + strconv.Itoa(s)
}

// Load return video screenshots
func (s *Spec) Load(ffmpeg *FFMPEG, file string, count int) ([]image.Image, error) {
	if count == 0 {
		return nil, errors.New("count must be greater than 0")
	}

	seconds, err := ffmpeg.Seconds(file)
	if err != nil {
		return nil, err
	}

	stack := NewErrorStack()

	images, err := s.makeframes(ffmpeg, file, seconds, count)
	if len(images) == 0 {
		return nil, errors.New("no frames captured")
	}
	stack.Add(err)

	return images, stack.Get()
}

func (s *Spec) frames(seconds int, count int) []int {
	seconds = seconds - seconds/5 // remove last 5% of the video, usually credits or black screen
	frame := seconds / (count + 1)

	frames := make([]int, count)
	for i := 0; i < count; i++ {
		frames[i] = frame * (i + 1)
	}

	return frames
}

func (s *Spec) makeframes(ffmpeg *FFMPEG, file string, seconds int, count int) ([]image.Image, error) {
	es := NewErrorStack()

	images := make([]image.Image, 0)
	frames := s.frames(seconds, count)

	for _, frame := range frames {
		img, err := ffmpeg.Shot(file, frame)
		if err != nil {
			es.Add(err)
		}

		images = append(images, img)
	}

	return images, es.Get()
}
