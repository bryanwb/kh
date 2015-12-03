package kh

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/dullgiulio/pingo"
	"github.com/hashicorp/go-multierror"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
)

var Logger *log.Logger
var Version = "0.0.1"
var Verbose bool
var HandHome string

type FingerDescriptor struct {
	Path        string
	Description string
}

type Hand struct {
	Home    string
	Fingers map[string]*FingerDescriptor
}

func MakeHand(home string) (*Hand, error) {
	h := new(Hand)
	h.Home = home
	Logger.Debugf("Using %s for home", home)
	h.Fingers = make(map[string]*FingerDescriptor)
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
		name := di[i].Name()
		fullFingerPath := path.Join(h.Home, name, name)
		if pathHasFinger(fullFingerPath) {
			fd := new(FingerDescriptor)
			Logger.Debugf("Found finger at %s", fullFingerPath)
			fd.Path = fullFingerPath
			fd.Description = findFingerDescription(path.Join(h.Home, name))
			h.Fingers[name] = fd
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
		Logger.Debugf("Updating finger %s", updateList[i])
		if err := h.UpdateFinger(updateList[i]); err != nil {
			result = multierror.Append(result, err)
		}
	}
	return nil
}

func (h *Hand) buildFinger(p string) error {
	action := "go build"
	shell := os.Getenv("SHELL")
	goPath := os.Getenv("GOPATH")
	goRoot := os.Getenv("GOROOT")
	Logger.Debugf("Executing command %v in path %s", action, p)
	cmd := exec.Command(shell, "-c", action)
	extendedGoPath := fmt.Sprintf("GOPATH=%s:%s", p, goPath)
	cmd.Env = append(os.Environ(), extendedGoPath, goRoot)
	cmd.Dir = path.Dir(p)
	err := cmd.Run()
	if err != nil {
		Logger.Errorf("Building %s failed with error %s", p, err.Error())
		output, _ := cmd.CombinedOutput()
		Logger.Errorf(string(output))
	} else {
		Logger.Debugf("Building %s completed successfully", p)
	}
	return err
}

func (h *Hand) UpdateFinger(finger string) error {
	p := h.Fingers[finger]
	Logger.Debugf("finger p is %v", p)
	srcP, err := resolvePath(p.Path)
	Logger.Debugf("resolved finger p is %v", p)

	if err != nil {
		return err
	}
	return h.buildFinger(srcP)
}

func (h *Hand) MakeFinger(name string, flags map[string]bool, args []string) (*FingerClient, error) {
	f := new(FingerClient)
	f.Path = path.Join(h.Home, name, name)
	f.finger = pingo.NewPlugin("tcp", f.Path)
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
	Stdin []byte
}

type FingerClient struct {
	Path   string
	finger *pingo.Plugin
	args   *FingerArgs
}

func (f *FingerClient) Execute() error {
	f.finger.Start()
	defer f.finger.Stop()
	resp := new(Response)
	Logger.Debugf("Executing finger %s", f.Path)
	methodName := "Finger.Execute"
	if f.args.Flags.Help {
		methodName = "Finger.Help"
	}
	if err := f.finger.Call(methodName, f.args, &resp); err != nil {
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

func (h *Hand) FingerDescriptions() []string {
	descs := make([]string, 0)
	var desc string
	for k, _ := range h.Fingers {
		desc = fmt.Sprintf("%s\t\t%s", k, h.Fingers[k].Description)
		descs = append(descs, desc)
	}
	return descs
}
