package utils_test

import (
	"bytes"
	"image/png"
	"os"
	"testing"

	"github.com/ReanSn0w/go-screenlist/pkg/utils"
)

func TestFFMPEG_Shot(t *testing.T) {
	image, err := utils.NewFFMPEG(t).Shot("test.mp4", 60)
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

func TestFFMPEG_ShotJPEG(t *testing.T) {
	err := utils.NewFFMPEG(t).ShotJPEG("test.mp4", "image.jpeg", 60, 100)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
		return
	}
}

func TestFFMPEG_ShotPNG(t *testing.T) {
	err := utils.NewFFMPEG(t).ShotPNG("test.mp4", "image.png", 60)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
		return
	}
}
