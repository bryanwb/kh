package main

import (
	"errors"
	"github.com/bryanwb/kh"
)

// This object has to be named "Finger" due to a quirk
// in the plugin system that kh uses
type Finger struct{}

func (p *Finger) Help() string {
	return "hi, I am easy to use"
}

func (p *Finger) Execute(fa *kh.FingerArgs, resp *kh.Response) error {
	resp.Debug("Inside execute")
	// wish there was a more elegant way to set the logging level
	resp.SetVerbose(fa.Flags.Verbose)
	if fa.Flags.Help {
		resp.Debugf("About to execute %s", "help")
		resp.WriteStdout(p.Help())
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
	finger := &Finger{}
	kh.Register(finger)
	kh.Run()
}
