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
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/bryanwb/hand"
	"os"
	"os/user"
	"path"
	"strings"
)

var log = logrus.New()

func fingerInvoked(h *hand.Hand, args []string) bool {
	if len(args) < 2 {
		return false
	}
	return hand.Contains(h.FingerNames(), args[1])
}

func verboseSet() bool {
	if hand.ContainsAny(os.Args, []string{"-v", "--verbose"}) {
		return true
	}
	return false
}

func makeHand() *hand.Hand {
	currentUser, _ := user.Current()
	hand.HandHome = path.Join(currentUser.HomeDir, "/.hand")
	hand.Logger = log
	if verboseSet() {
		log.Level = logrus.DebugLevel
	}
	h, err := hand.MakeHand(hand.HandHome)
	if err != nil {
		fmt.Println("Encountered error finding fingers")
		fmt.Printf("Error was %s", err.Error())
		os.Exit(1)
	}
	return h
}

func versionCmdInvoked() bool {
	if len(os.Args) < 2 {
		return false
	}
	return os.Args[1] == "version"
}

func showVersion() {
	fmt.Printf("Version %s of The Hand of the King(hand)\n", hand.Version)
}

func showHelp(h *hand.Hand) {
	helpText := `The Hand of the King (hand) is a tool for organizing and executing shellish scripts.
It does your dirty work, so keep it clean.
Hand exposes plugins, known as fingers, as subcommands.

Usage:
hand [flags]
hand [finger] [arguments to a finger]

Available Commands:
version     Print the version number of Hand
help

Flags:
  -H, --hand-home="/Users/hitman/.hand": Home directory for hand
  -v, --verbose[=false]: verbose output

Use "hand [finger] --help" for more information about a finger.
`
	fmt.Printf(helpText)
	if len(h.Fingers) > 0 {
		fmt.Printf("Available fingers are:\n%s\n",
			strings.Join(h.FingerNames(), "\n"))
	} else {
		fmt.Printf("Currently no fingers available\n")
	}
}

func findFingerArgs(args []string) []string {
	if len(args) < 3 {
		return []string{}
	}
	return args[2:]
}

// This doesn't use a cli argument parser because such libraries typically cannot
// handle subcommands that are dynamically loaded
// For this reason cli parsing is done manually
func main() {
	h := makeHand()
	if fingerInvoked(h, os.Args) {
		fingerName := os.Args[1]
		remainingArgs := findFingerArgs(os.Args)
		if err := h.ExecuteFinger(fingerName, remainingArgs); err != nil {
			fmt.Errorf(err.Error())
		}
		os.Exit(0)
	}
	if versionCmdInvoked() {
		showVersion()
		os.Exit(0)
	}
	if hand.HelpCmdInvoked() {
		showHelp(h)
		os.Exit(0)
	}

}
