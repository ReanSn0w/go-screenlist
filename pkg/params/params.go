package params

import (
	"os"

	"github.com/ReanSn0w/go-screenlist/pkg/utils"
	"github.com/go-pkgz/lgr"
	"github.com/umputun/go-flags"
)

func Load(log utils.Logger) *Parameters {
	opts := &Parameters{}

	p := flags.NewParser(opts, flags.PrintErrors|flags.PassDoubleDash|flags.HelpFlag)
	p.SubcommandsOptional = true
	if _, err := p.Parse(); err != nil {
		if err.(*flags.Error).Type != flags.ErrHelp {
			log.Logf("[ERROR] cli error: %v", err)
		}
		os.Exit(2)
	}

	log.Logf("[INFO] Verbose: %v", opts.Verbose)
	log.Logf("[INFO] Number of Screenshots: %v", opts.Screenshots)
	log.Logf("[INFO] Files: %v", opts.Files)

	return opts
}

type Parameters struct {
	Verbose     bool     `short:"v" long:"verbose" description:"verbose mode"`
	Screenshots int      `long:"count" default:"15" description:"number of screenshots"`
	ResultWidth int      `long:"width" default:"1920" description:"resulting image width"`
	Treads      int      `long:"treads" default:"4" description:"number of treads"`
	Delta       bool     `long:"delta" description:"save delta images"`
	Force       bool     `short:"f" long:"force" description:"force execution (ignore errors)"`
	Grid        int      `long:"grid" default:"3" description:"grid size"`
	Files       []string `short:"i" long:"input" description:"file destinations"`
}

func (p *Parameters) Log() utils.Logger {
	if p.Verbose {
		return lgr.New(lgr.Msec, lgr.Debug, lgr.CallerFile, lgr.CallerFunc)
	}

	return lgr.Default()
}
