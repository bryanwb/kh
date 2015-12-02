package kh

import (
	"github.com/dullgiulio/pingo"
)

// This interface defines what a Finger must implement
// Currently not much, but more than pingo requires
type Finger interface {
	Help(*FingerArgs, *Response) error
	Execute(*FingerArgs, *Response) error
}

// Registers a Finger w/ our plugin system
func Register(f Finger) {
	pingo.Register(f)
}

// Run the finger subsystem
func Run() error {
	return pingo.Run()
}
