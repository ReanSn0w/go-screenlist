package list

import (
	"go-screenlist/pkg/video"
	"testing"
)

func Test_Make(t *testing.T) {
	spec, err := video.Specs("test.mp4")
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
		t.FailNow()
	}

	images, err := video.Load("test.mp4", 2, false)
	if err != nil {
		t.Errorf("Expected nil, got \n %v", err)
		t.FailNow()
	}

	err = Save("list.jpg", spec, 3, 1000, images)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
		t.FailNow()
	}
}
