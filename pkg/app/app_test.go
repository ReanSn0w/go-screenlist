package app_test

import (
	"testing"

	"github.com/ReanSn0w/go-screenlist/pkg/app"
	"github.com/ReanSn0w/go-screenlist/pkg/params"
)

var (
	p = &params.Parameters{
		Verbose:     true,
		Screenshots: 16,
		ResultWidth: 1080,
		Treads:      4,
		Delta:       true,
		Force:       false,
		Grid:        4,
		Files:       []string{"test.mp4"},
	}
)

func TestApp_Run(t *testing.T) {
	log := p.Log()
	app.New(log, p).Run()
}
