package wallutils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// has checks if a string slice has the given element
func has(sl []string, e string) bool {
	for _, s := range sl {
		if s == e {
			return true
		}
	}
	return false
}

// unique removes all repeated elements from a slice of strings
func unique(sl []string) []string {
	var nl []string
	for _, s := range sl {
		if !has(nl, s) {
			nl = append(nl, s)
		}
	}
	return nl
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// firstname finds the part of a filename before the extension
func firstname(filename string) string {
	ext := filepath.Ext(filename)
	return filename[:len(filename)-len(ext)]
}

// exists checks if the given path exists
func exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// which tries to find the given executable name in the $PATH
// Returns an empty string if not found.
func which(executable string) string {
	p, err := exec.LookPath(executable)
	if err != nil {
		return ""
	}
	return p
}

// which2 tries to find the given executable name in the $PATH
// This is the same as exec.LookPath.
// Returns an error if not found.
func which2(executable string) (string, error) {
	return exec.LookPath(executable)
}

// run executes the given executable and returns an error if the exit code is
// non-zero. If verbose is true, the command will be printed before running.
func run(executable string, arguments []string, verbose bool) error {
	if verbose {
		fmt.Println(executable + " " + strings.Join(arguments, " "))
	}
	cmd := exec.Command(executable, arguments...)
	if _, err := cmd.CombinedOutput(); err != nil {
		return err
	}
	return nil
}

// output returns the output after running a given executable
// if verbose is true, the command will be printed before running
func output(executable string, arguments []string, verbose bool) string {
	if verbose {
		fmt.Println(executable + " " + strings.Join(arguments, " "))
	}
	cmd := exec.Command(executable, arguments...)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}
	return string(stdoutStderr)
}

// runShell is the same as the "run" function, but runs the commands in a shell.
func runShell(shellCommand string, verbose bool) error {
	if verbose {
		fmt.Println(shellCommand)
	}
	cmd := exec.Command("sh", "-c", shellCommand)
	if _, err := cmd.CombinedOutput(); err != nil {
		return err
	}
	return nil
}

// outputShell is the same as the "output" function,
// but runs the command in a shell
func outputShell(shellCommand string, verbose bool) string {
	if verbose {
		fmt.Println(shellCommand)
	}
	cmd := exec.Command("sh", "-c", shellCommand)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}
	return string(stdoutStderr)
}

// CommonPrefix will find the longest common prefix in a slice of strings
func CommonPrefix(sl []string) string {
	if len(sl) == 0 {
		return ""
	}
	shortestLength := len(sl[0])
	shortestString := sl[0]
	for _, s := range sl {
		if len(s) < shortestLength {
			shortestLength = len(s)
			shortestString = s
		}
	}
	if shortestLength == 0 {
		return ""
	}
	for i := 1; i < shortestLength; i++ {
		for _, s := range sl {
			if !strings.HasPrefix(s, shortestString[:i]) {
				return shortestString[:i-1]
			}
		}
	}
	return shortestString
}

// CommonSuffix will find the longest common suffix in a slice of strings
func CommonSuffix(sl []string) string {
	if len(sl) == 0 {
		return ""
	}
	shortestLength := len(sl[0])
	shortestString := sl[0]
	for _, s := range sl {
		if len(s) < shortestLength {
			shortestLength = len(s)
			shortestString = s
		}
	}
	if shortestLength == 0 {
		return ""
	}
	for i := 1; i < shortestLength; i++ {
		for _, s := range sl {
			if !strings.HasSuffix(s, shortestString[shortestLength-i:]) {
				return shortestString[shortestLength-(i-1):]
			}
		}
	}
	return shortestString
}

// Meat returns the meat of the string: the part that is after the prefix and
// before the suffix. It does not check if the prefix and suffix exists in the
// string. If the given string is too short to contain the prefix and suffix,
// it will be returned as it is.
func Meat(s, prefix, suffix string) string {
	if len(s) < (len(prefix) + len(suffix)) {
		return s
	}
	return s[len(prefix) : len(s)-len(suffix)]
}

// Quit with a nicely formatted error message to stderr
func Quit(err error) {
	msg := err.Error()
	if !strings.HasSuffix(msg, ".") && !strings.HasSuffix(msg, "!") && !strings.Contains(msg, ":") {
		msg += "."
	}
	fmt.Fprintf(os.Stderr, "%s%s\n", strings.ToUpper(string(msg[0])), msg[1:])
	os.Stdout.Sync()
	os.Stderr.Sync()
	os.Exit(1)
}

// expandUser replaces a leading ~ or $HOME with the path
// to the home directory of the current user
func expandUser(path string) string {
	// this is a simpler alternative to using os.UserHomeDir (which requires Go 1.12 or later)
	if strings.HasPrefix(path, "~") {
		path = strings.Replace(path, "~", os.Getenv("HOME"), 1)
	} else if strings.HasPrefix(path, "$HOME") {
		path = strings.Replace(path, "$HOME", os.Getenv("HOME"), 1)
	}
	return path
}
