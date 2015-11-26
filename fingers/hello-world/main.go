package main

import (
	"github.com/bryanwb/hand"
	"github.com/dullgiulio/pingo"
)

type Finger struct{}

// sadly the name parameter is necessary, not sure why
func help() string {
	return "hi, I am easy to use"
}

func (p *Finger) Execute(args []string, resp *hand.Response) error {
	if hand.IsHelpCommand(args) {
		resp.Stdout = help()
		return nil
	}
	resp.Stdout = "Hello, " + args[0]
	resp.Stderr = ""
	return nil
}

func main() {
	plugin := &Finger{}
	pingo.Register(plugin)
	pingo.Run()
}
