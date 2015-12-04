package kh

import (
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strings"
)

var commonFlags []string = []string{"-v", "--verbose", "-h", "--help"}
var subcommands []string = []string{"version", "help", "update", "init",
	"list", "install"}

func Map_has_key(m map[string]string, s string) bool {
	for str, _ := range m {
		if s == str {
			return true
		}
	}
	return false
}

func Contains(l []string, s string) bool {
	for _, str := range l {
		if s == str {
			return true
		}
	}
	return false
}

func ContainsAny(list1 []string, list2 []string) bool {
	for _, str := range list1 {
		if Contains(list2, str) {
			return true
		}
	}
	return false
}

func HelpCmdInvoked() bool {
	if len(os.Args) == 1 {
		return true
	}
	helpFlags := []string{"-h", "--help", "help"}
	if Contains(helpFlags, os.Args[1]) {
		return true
	}
	return false
}

// this intended for use by Fingers
func IsHelpCommand(args []string) bool {
	if len(args) == 0 {
		return false
	}
	helpFlags := []string{"-h", "--help", "help"}
	if ContainsAny(helpFlags, args) {
		return true
	}
	return false
}

func stripCommonFlags(args []string) []string {
	strippedArgs := make([]string, 0)
	for i := range args {
		if !Contains(commonFlags, args[i]) {
			strippedArgs = append(strippedArgs, args[i])
		}
	}
	return strippedArgs
}

func StripFlags(args []string) []string {
	newArgs := make([]string, 0)
	for i := range args {
		if !strings.HasPrefix(args[i], "-") {
			newArgs = append(newArgs, args[i])
		}
	}
	return newArgs
}

func SubcommandInvoked(args []string) string {
	args = StripFlags(args)
	if len(args) < 2 {
		return ""
	}
	if Contains(subcommands, args[1]) {
		return args[1]
	}
	return ""
}

func pathHasFinger(fingerPath string) bool {
	if fi, err := os.Stat(fingerPath); err == nil && fi.Mode()&0111 != 0 {
		Logger.Debugf("Found finger at %s", fingerPath)
		return true
	}
	return false
}

func resolvePath(p string) (string, error) {
	fi, err := os.Lstat(p)
	if err != nil {
		Logger.Debugf("Path %s does not exist\n", p)
		return "", err
	}
	if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
		if realPath, err := os.Readlink(p); err == nil {
			return realPath, nil
		}
		return "", err
	}
	return p, nil
}

func linkFingerToHome(source string, name string) error {
	target := path.Join(HandHome, name)
	_, err := os.Lstat(target)
	if os.IsNotExist(err) {
		Logger.Debugf("Linking finger %s with source %s to target %s", name, source, target)
		if symErr := os.Symlink(source, target); symErr != nil {
			return symErr
		}
		return nil
	}
	return err
}

// This is a hack to find the filesystem location of the kh package
// at runtime
type Empty struct{}

//Initializes the HandHome w/ default fingers
func Init() error {
	_, err := os.Stat(HandHome)
	if os.IsNotExist(err) {
		os.Mkdir(HandHome, 0775)
	} else if err != nil {
		return err
	}
	khPath := reflect.TypeOf(Empty{}).PkgPath()
	goPath := os.Getenv("GOPATH")
	fingersPath := path.Join(goPath, "src", khPath, "fingers")
	di, _ := ioutil.ReadDir(fingersPath)
	for i := range di {
		fingerPath := di[i].Name()
		name := path.Base(fingerPath)
		fullFingerPath := path.Join(fingersPath, fingerPath)
		if pathHasFinger(fullFingerPath) {
			Logger.Debugf("Found finger at %s", fullFingerPath)
			if err := linkFingerToHome(fullFingerPath, name); err != nil {
				Logger.Errorf("Couldn not link finger at %s to ./kh/%s",
					fullFingerPath, name)
				Logger.Errorf("Error encountered: %s", err.Error())
			}
		} else {
			Logger.Debugf("No finger found at %s", fullFingerPath)
		}
	}
	return nil
}

func findFingerDescription(p string) string {
	description := ""
	descF := path.Join(p, "DESCRIPTION")
	Logger.Debugf("Looking for Finger description at %s", p)
	_, err := os.Stat(descF)
	if err != nil {
		Logger.Debugf("Finger description not found at %s", p)
		Logger.Debugf("Encountered error %v", err)
		return description
	}
	Logger.Debugf("Found finger description at %s", p)
	if d, err := ioutil.ReadFile(descF); err == nil {
		Logger.Debugf("Found finger description text: %s", string(d))
		return string(d)
	}
	return description
}

// some magic from http://stackoverflow.com/a/16753808 to determine
// if Pipe w/ data has attached to Stdin
func isStdinAttached() bool {
	return !terminal.IsTerminal(0)
}

// wouldn't need to do this shit if Golang supported Monads
func EmptyArgs(args []string) bool {
	if len(args) < 1 || args[0] == "" {
		return true
	}
	return false
}
