package main

import (
	"os"

	"github.com/ReanSn0w/go-screenlist/pkg/app"
	"github.com/ReanSn0w/go-screenlist/pkg/params"
	"github.com/ReanSn0w/go-screenlist/pkg/utils"
	"github.com/go-pkgz/lgr"
)

var (
	log utils.Logger
)

func main() {
	params := params.Load(lgr.Default())
	log = params.Log()
	app := app.New(log, params)
	err := app.Run()

	if err != nil {
		os.Exit(2)
	}

	os.Exit(0)
}
