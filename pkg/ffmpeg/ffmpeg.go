package ffmpeg

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
)

// Shot takes a screenshot of a video at a given second
func Shot(input string, s int) (image.Image, error) {
	output, err := command("ffmpeg", "-ss", seconds(s), "-i", input, "-vframes", "1", "-q:v", "2", "-f", "image2pipe", "-vcodec", "png", "-")
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
func ShotJPEG(input string, output string, s int, quality int) error {
	return shotImage(input, output, s, func(w io.Writer, img image.Image) error {
		return jpeg.Encode(w, img, &jpeg.Options{Quality: quality})
	})
}

func ShotPNG(input string, output string, s int) error {
	return shotImage(input, output, s, func(w io.Writer, img image.Image) error {
		return png.Encode(w, img)
	})
}

// func ShotWEBP(input string, output string, s int, quality int) error {
// 	return shotImage(input, output, s, func(w io.Writer, img image.Image) error {
// 		return webp.Encode(w, img, &webp.Options{Lossless: false, Quality: float32(quality)})
// 	})
// }

func shotImage(input string, output string, s int, prepare func(io.Writer, image.Image) error) error {
	img, err := Shot(input, s)
	if err != nil {
		return err
	}

	f, err := os.Create(output)
	if err != nil {
		return err
	}
	defer f.Close()

	return prepare(f, img)
}

// FPS returns fps of video
func FPS(file string) (int, error) {
	output, err := command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=r_frame_rate", "-of", "default=noprint_wrappers=1:nokey=1", file)
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
func Seconds(file string) (int, error) {
	output, err := command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", file)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(strings.Split(string(output), ".")[0])
}

// Resolution returns the resolution of a video
func Resolution(file string) ([2]int, error) {
	output, err := command("ffprobe", "-v", "error", "-show_entries", "stream=width,height", "-of", "csv=p=0:s=x", file)
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
	output, err := command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=r_frame_rate", "-of", "default=noprint_wrappers=1:nokey=1", file)
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
	output, err := command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=bit_rate", "-of", "default=noprint_wrappers=1:nokey=1", file)
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

func command(name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)
	out, err := cmd.Output()
	return string(out), err
}

// int to format 00:00:00
func seconds(s int) string {
	h := s / 3600
	s -= h * 3600
	m := s / 60
	s -= m * 60
	return strconv.Itoa(h) + ":" + strconv.Itoa(m) + ":" + strconv.Itoa(s)
}
