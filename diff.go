package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strings"
)

// pyTemplate is a template Python script for pretty-printing diffs
// on the command-line. Unfortunately, there's no equivalent to difflib
// in go, so this is the only way (I think) to provide the same kind of
// user experience.
const pyTemplate = `
import difflib

class bcolors:
    OKBLUE = '\033[94m'
    OKGREEN = '\033[92m'
    FAIL = '\033[91m'
    ENDC = '\033[0m'

def main():
    file_path = '%s'
    orig_lines = %s
    result_lines = %s

    diff = difflib.unified_diff(
        orig_lines,
        result_lines,
        fromfile=file_path,
        tofile=file_path + '.shortened',
    )

    for line in diff:
        line = line.rstrip()
        if line.startswith('+'):
            print(bcolors.OKGREEN + line + bcolors.ENDC)
        elif line.startswith('-'):
            print(bcolors.FAIL + line + bcolors.ENDC)
        elif line.startswith('^'):
            print(bcolors.OKBLUE + line + bcolors.ENDC)
        elif len(line) > 0:
            print(line)

if __name__ == "__main__":
    main()
`

// PrettyDiff prints colored, git-style diffs to the console. It uses an
// embedded Python script (above) to take advantage of Python's difflib package.
func PrettyDiff(path string, contents []byte, results []byte) error {
	if reflect.DeepEqual(contents, results) {
		return nil
	}

	contentLines := strings.Split(string(contents), "\n")
	resultLines := strings.Split(string(results), "\n")

	contentBytes, err := json.Marshal(contentLines)
	if err != nil {
		return err
	}

	resultBytes, err := json.Marshal(resultLines)
	if err != nil {
		return err
	}

	pyScript := fmt.Sprintf(
		pyTemplate,
		path,
		string(contentBytes),
		string(resultBytes),
	)

	// The script should work with either python2 or python3, but prefer the latter
	// if it's available.
	pyPath, err := exec.LookPath("python3")
	if err != nil {
		pyPath, err = exec.LookPath("python")
		if err != nil {
			return errors.New("Could not find python in path")
		}
	}

	cmd := exec.Command(pyPath)
	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err = cmd.Start(); err != nil {
		return err
	}

	_, err = stdinPipe.Write([]byte(pyScript))
	if err != nil {
		return err
	}
	stdinPipe.Close()

	err = cmd.Wait()
	if err != nil {
		return err
	}

	fmt.Print("\n")

	return nil
}
