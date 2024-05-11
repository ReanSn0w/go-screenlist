package main

import (
	"os"

	"github.com/ReanSn0w/go-screenlist/pkg/engine"
	"github.com/ReanSn0w/go-screenlist/pkg/utils"
	"github.com/ReanSn0w/tk4go/pkg/config"
	"github.com/go-pkgz/lgr"
)

var (
	revision = "unknown"
	log      = lgr.Default()

	opts = struct {
		Verbose bool         `short:"v" long:"verbose" description:"verbose mode"`
		Force   bool         `short:"f" long:"force" description:"force execution (ignore errors)"`
		Treads  int          `short:"t" long:"treads" default:"4" description:"number of treads"`
		Files   []utils.File `short:"i" long:"input" description:"file destinations"`

		Screenlist engine.ScreenListPreferences `group:"screenlist" namespace:"screenlist" env-namespace:"SCREENLIST"`
		Delta      engine.DeltaPreferences      `group:"delta" namespace:"delta" env-namespace:"DELTA"`
	}{}
)

func main() {
	err := config.Parse(&opts)
	if err != nil {
		log.Logf("%v", err)
		os.Exit(2)
	}

	if opts.Verbose {
		lgr.Setup(lgr.Debug, lgr.CallerFile, lgr.CallerFunc)
	}

	config.Print(log, "Screenlist", revision, opts)

	engine := engine.New(
		log, opts.Force, opts.Treads,
		&opts.Screenlist, &opts.Delta)

	engine.Add(opts.Files...)
	engine.Wait()
}
