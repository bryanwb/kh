package main

import (
	"errors"
	"github.com/bryanwb/kh"
	"github.com/dullgiulio/pingo"
)

type FingerServer struct{}

func help() string {
	return "hi, I am easy to use"
}

func (p *FingerServer) Execute(fa *kh.FingerArgs, resp *kh.Response) error {
	resp.Debug("Inside execute")
	// wish there was a more elegant way to set the logging level
	resp.SetVerbose(fa.Flags.Verbose)
	if fa.Flags.Help {
		resp.Debugf("About to execute %s", "help")
		resp.WriteStdout(help())
		return nil
	}
	args := fa.Args
	if len(args) < 1 {
		return errors.New("Hey buddy, I need a name")
	} else {
		resp.WriteStdout("Hello, " + args[0])
	}
	return nil
}

func main() {
	plugin := &FingerServer{}
	pingo.Register(plugin)
	pingo.Run()
}
