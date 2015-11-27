package kh

import (
	"github.com/Sirupsen/logrus"
	"os"
	"path"
	"testing"
)

func setUpTest() string {
	Logger = logrus.New()
	verbose := testing.Verbose()
	if verbose {
		Logger.Level = logrus.DebugLevel
	} else {
		Logger.Level = logrus.PanicLevel
	}
	dir, _ := os.Getwd()
	HandHome := path.Join(dir, "test-fixtures")
	return HandHome
}

func TestFindsFingers(t *testing.T) {
	home := setUpTest()
	h, _ := MakeHand(home)
	for _, k := range []string{"git", "gpg"} {
		if !Map_has_key(h.Fingers, k) {
			t.Errorf("Finger %s missing from set of fingers\n", k)
		}
	}
}
