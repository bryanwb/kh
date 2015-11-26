package hand

import (
	"os"
)

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
