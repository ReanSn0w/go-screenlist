package video_test

import (
	"testing"

	"github.com/ReanSn0w/go-screenlist/pkg/video"
)

func Test_Load(t *testing.T) {
	count := 3

	images, err := video.Load("test.mp4", count, false)
	if err != nil {
		t.Errorf("Expected nil, got \n %v", err)
	}

	if len(images) != count {
		t.Errorf("Expected 16, got %v", len(images))
	}
}

func Test_Specs(t *testing.T) {
	spec, err := video.Specs("test.mp4")
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
		t.FailNow()
	}

	if spec.Resolution[0] != 1280 {
		t.Errorf("Expected 1280, got %v", spec.Resolution[0])
	}

	if spec.Resolution[1] != 720 {
		t.Errorf("Expected 720, got %v", spec.Resolution[1])
	}

	if spec.Fps != 25 {
		t.Errorf("Expected 25, got %v", spec.Fps)
	}
}
