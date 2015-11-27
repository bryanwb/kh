package kh

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/dullgiulio/pingo"
	"github.com/hashicorp/go-multierror"
	"io/ioutil"
	"os"
	"path"
)

var Logger *log.Logger
var Version = "0.0.1"
var Verbose bool
var HandHome string

type Hand struct {
	Home    string
	Fingers map[string]string
}

func MakeHand(home string) (*Hand, error) {
	h := new(Hand)
	h.Home = home
	Logger.Debugf("Using %s for home", home)
	h.Fingers = make(map[string]string)
	err := h.FindFingers()
	return h, err
}

func (h *Hand) FindFingers() error {
	// find all directories w/ executables of same name
	_, err := os.Stat(h.Home)
	if os.IsNotExist(err) {
		return err
	}
	di, _ := ioutil.ReadDir(h.Home)
	for i := range di {
		fingerPath := di[i].Name()
		name := path.Base(fingerPath)
		fullFingerPath := path.Join(h.Home, fingerPath, name)
		if pathHasFinger(fullFingerPath) {
			Logger.Debugf("Found finger at %s", fullFingerPath)
			h.Fingers[name] = fullFingerPath
		} else {
			Logger.Debugf("No finger found at %s", fullFingerPath)
		}
	}
	return err
}

func (h *Hand) Update(fingers []string) error {
	var result error
	updateList := make([]string, 0)
	if len(fingers) < 1 {
		updateList = h.FingerNames()
	} else {
		updateList = fingers
	}
	for i := range updateList {
		if err := h.UpdateFinger(updateList[i]); err != nil {
			result = multierror.Append(result, err)
		}
	}
	return nil
}

func (h *Hand) UpdateFinger(finger string) error {
	//	p := h.Fingers[finger]
	return nil
}

func (h *Hand) MakeFinger(name string, flags map[string]bool, args []string) (*Finger, error) {
	f := new(Finger)
	f.Path = path.Join(h.Home, name, name)
	f.plugin = pingo.NewPlugin("tcp", f.Path)
	f.args = makeArgs(flags, args)
	return f, nil
}

func (h *Hand) ExecuteFinger(name string, flags map[string]bool, args []string) error {
	f, err := h.MakeFinger(name, flags, args)
	if err != nil {
		return err
	}
	return f.Execute()
}

func makeArgs(flags map[string]bool, args []string) *FingerArgs {
	newArgs := new(FingerArgs)
	fs := new(flagSet)
	fs.Help = flags["help"]
	fs.Verbose = flags["verbose"]
	newArgs.Flags = fs
	newArgs.Args = stripCommonFlags(args)
	return newArgs
}

type flagSet struct {
	Help    bool
	Verbose bool
}

// A finger takes a set of flags plus an arbitrary list of arguments
type FingerArgs struct {
	Flags *flagSet
	Args  []string
}

type Finger struct {
	Path   string
	plugin *pingo.Plugin
	args   *FingerArgs
}

func (f *Finger) Execute() error {
	f.plugin.Start()
	defer f.plugin.Stop()
	resp := new(Response)
	Logger.Debugf("Executing finger %s", f.Path)
	if err := f.plugin.Call("FingerServer.Execute", f.args, &resp); err != nil {
		Logger.Debugf(resp.SprintLog())
		return err
	}
	stdout := resp.SprintStdout()
	stderr := resp.SprintStderr()
	fingerLog := resp.SprintLog()
	if stdout != "" {
		fmt.Println(stdout)
	}
	if stderr != "" {
		fmt.Fprintf(os.Stderr, stderr+"\n")
	}
	if fingerLog != "" {
		fmt.Fprintf(os.Stderr, fingerLog+"\n")
	}
	return nil
}

func (h *Hand) FingerNames() []string {
	names := make([]string, 0)
	for k, _ := range h.Fingers {
		names = append(names, k)
	}
	return names
}
