package kh

import (
	"os"
	"strings"
)

var commonFlags []string = []string{"-v", "--verbose", "-h", "--help"}
var subcommands []string = []string{"version", "help", "update"}

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
