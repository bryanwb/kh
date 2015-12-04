// Copyright © 2015 Bryan W. Berry <bryan.berry@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package main

import (
	"bytes"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/bryanwb/kh"
	flag "github.com/spf13/pflag"
	"os"
	"os/user"
	"path"
	"strings"
)

var log = logrus.New()
var verboseFlag bool
var helpFlag bool
var flagParsingOut bytes.Buffer

func fingerInvoked(h *kh.Hand, arg string) bool {
	return kh.Contains(h.FingerNames(), arg)
}

func makeHand() (*kh.Hand, error) {
	currentUser, _ := user.Current()
	kh.HandHome = path.Join(currentUser.HomeDir, "/.kh")
	kh.Logger = log
	h, err := kh.MakeHand(kh.HandHome)
	if err != nil {
		log.Debug("Encountered error finding fingers")
		log.Debug("Error was %s", err.Error())
	}
	return h, err
}

func versionCmdInvoked() bool {
	if len(os.Args) < 2 {
		return false
	}
	return os.Args[1] == "version"
}

func showVersion() {
	fmt.Printf("Version %s of The King's Hand\n", kh.Version)
}

func showHelp(h *kh.Hand) {
	helpText := `The King's Hand (kh) is a tool for organizing and executing shellish scripts.
It does your dirty work, so keep it clean.
kh exposes plugins, known as fingers, as subcommands.

Usage:
kh [flags]
kh [finger] [arguments to a finger]

Available subcommands:
install [finger...]   Installs one or more fingers given their import paths
                      uses go get under the covers
version               Print the version number of King's Hand
help                  C'mon, do I have to explain this one?
init                  Creates a ~/.kh directory if one does not exist
                      and installs default fingers
update [finger]       Builds one or more fingers
                      By default, updates all   


Flags:
  -H, --hand-home="/Users/hitman/.kh": Home directory for kh
  -v, --verbose[=false]: verbose output

Use "kh [finger] --help" for more information about a finger.
`
	fmt.Printf(helpText)
	showList(h)
}

func findFingerArgs(args []string) []string {
	if len(args) < 3 {
		return []string{}
	}
	return args[2:]
}

func parseFlagsAndArgs() ([]string, error) {
	// We need to override the default behavior to exit program on parsing error
	// and to not immediately write errors encountered to Stdout
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flag.CommandLine.SetOutput(&flagParsingOut)
	flag.Usage = func() {
		fmt.Fprintf(&flagParsingOut, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.BoolVarP(&verboseFlag, "verbose", "v", false, "Verbose mode")
	flag.BoolVarP(&helpFlag, "help", "h", false, "help")
	err := flag.CommandLine.Parse(os.Args[1:])
	if verboseFlag {
		log.Level = logrus.DebugLevel
	}
	args := flag.Args()
	log.Debugf("args are %v", args)
	if len(args) < 1 {
		return []string{"", ""}, err
	}
	if len(args) < 2 {
		return []string{args[0], ""}, err
	}

	return args, err
}

func executeUpdate(h *kh.Hand, args []string) {
	if err := h.Update(args); err != nil {
		log.Errorf("Encountered error updating fingers")
		log.Errorf("Error message: %v", err)
		os.Exit(1)
	}
}

func executeFinger(h *kh.Hand, fingerName string) {
	remainingArgs := findFingerArgs(os.Args)
	flags := map[string]bool{"help": helpFlag, "verbose": verboseFlag}
	if err := h.ExecuteFinger(fingerName, flags, kh.StripFlags(remainingArgs)); err != nil {
		log.Errorf("Encountered error executing finger %s", fingerName)
		log.Errorf("Error message: %v", err)
		os.Exit(1)
	}

}

func showList(h *kh.Hand) {
	if len(h.Fingers) > 0 {
		fmt.Printf("Available fingers are:\n%s\n",
			strings.Join(h.FingerDescriptions(), "\n"))
	} else {
		fmt.Printf("Currently no fingers available\n")
	}
}

func doInit() {
	if err := kh.Init(); err != nil {
		log.Error("Encountered error executing init")
		log.Errorf("Error message was: %v", err)
		os.Exit(1)
	}
}

func doInstall(h *kh.Hand, args []string) {
	args = kh.StripFlags(args)
	if kh.EmptyArgs(args) {
		log.Error("You must specify one or more fingers to install")
		os.Exit(1)
	}
	if err := h.InstallFingers(args); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

// This doesn't use a cli argument parser because such libraries typically cannot
// handle subcommands that are dynamically loaded
// For this reason cli parsing is done manually
func main() {
	args, parseErr := parseFlagsAndArgs()
	h, _ := makeHand()
	if fingerInvoked(h, args[0]) {
		executeFinger(h, args[0])
		os.Exit(0)
	}
	cmd := args[0]
	// since not executing our finger, throw an error
	// if we see unknown flags
	if parseErr != nil {
		// flush any parsing errors encountered
		// while parsing command line earlier
		fmt.Fprint(os.Stderr, flagParsingOut.String())
		flag.Usage()
		os.Exit(1)
	}
	switch cmd {
	case "version":
		showVersion()
	case "help":
		showHelp(h)
	case "update":
		executeUpdate(h, kh.StripFlags(os.Args[2:]))
	case "init":
		doInit()
	case "install":
		doInstall(h, os.Args[2:])
	case "":
		showHelp(h)
	case "list":
		showList(h)
	default:
		fmt.Printf("Unknown finger or subcommand \"%s\"\nType `kh list` to see available fingers\n", cmd)
	}

}
