package ffmpeg_test

import (
	"bytes"
	"image/png"
	"os"
	"testing"

	"github.com/ReanSn0w/go-screenlist/pkg/ffmpeg"
)

func Test_Shot(t *testing.T) {
	image, err := ffmpeg.Shot("test.mp4", 60)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
		return
	}

	buffer := new(bytes.Buffer)
	err = png.Encode(buffer, image)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
		return
	}

	err = os.WriteFile("image.png", buffer.Bytes(), 0644)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
		return
	}
}

func Test_ShotJPEG(t *testing.T) {
	err := ffmpeg.ShotJPEG("test.mp4", "image.jpeg", 60, 100)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
		return
	}
}

func Test_ShotPNG(t *testing.T) {
	err := ffmpeg.ShotPNG("test.mp4", "image.png", 60)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
		return
	}
}
