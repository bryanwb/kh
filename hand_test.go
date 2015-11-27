package kh

import (
	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
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

func TestStripCommonFlags(t *testing.T) {
	args := []string{"foo", "bar", "-h", "--help", "-v", "--verbose"}
	expected := []string{"foo", "bar"}
	actual := stripCommonFlags(args)
	assert := assert.New(t)
	assert.Equal(expected, actual,
		"Common flags such as '-h', '--help', '-v', '--verbose' not stripped.")
}

func TestSubcommandInvoked(t *testing.T) {
	var cmd string
	cmd = SubcommandInvoked([]string{"program", "-v", "help", "-v"})
	assert := assert.New(t)
	assert.Equal(cmd, "help", "should have found the help subcommand")
	cmd = SubcommandInvoked([]string{"program", "-v", "--help", "-v"})
	assert.Equal(cmd, "", "No subcommand should be found")
	cmd = SubcommandInvoked([]string{"program", "version", "-v"})
	assert.Equal(cmd, "version", "should have found the version subcommand")
	cmd = SubcommandInvoked([]string{"program", "-h", "-v", "update"})
	assert.Equal(cmd, "update", "should have found the version subcommand")
}
