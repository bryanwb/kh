package main

import (
	"errors"
	"github.com/bryanwb/kh"
)

// This object has to be named "Finger" due to a quirk
// in the plugin system that kh uses
type Finger struct{}

func (p *Finger) Help(fa *kh.FingerArgs, resp *kh.Response) error {
	resp.SetVerbose(fa.Flags.Verbose)
	resp.WriteStdout("Hi, I am easy to use")
	return nil
}

func (p *Finger) Execute(fa *kh.FingerArgs, resp *kh.Response) error {
	// wish there was a more elegant way to set the logging level
	resp.SetVerbose(fa.Flags.Verbose)
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
