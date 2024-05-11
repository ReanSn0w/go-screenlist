package utils_test

import (
	"testing"

	"github.com/ReanSn0w/go-screenlist/pkg/utils"
)

func TestList_Make(t *testing.T) {
	ffmpeg := utils.NewFFMPEG(t)
	filename := utils.File("test.mp4")

	spec, err := ffmpeg.Specs(string(filename))
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
		t.FailNow()
	}

	images, err := spec.Load(ffmpeg, string(filename), 2)
	if err != nil {
		t.Errorf("Expected nil, got \n %v", err)
		t.FailNow()
	}

	err = filename.SaveImagesByPattern("image_{{.Count}}.jpg", images)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
		t.FailNow()
	}
}
