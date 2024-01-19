package main

import (
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testFiles = map[string]string{
	"test1.go": `package main

import "fmt"

func main() {
	fmt.Printf("%s %s %s %s %s %s", "argument1", "argument2", "argument3", "argument4", "argument5", "argument6")
}`,
	"test2.go": `package main

func main() {
	myMap := map[string]string{"key1": "value1", "key2": "value2", "key3": "value3", "key4": "value4", "key5", "value5"}
}`,
}

func TestRunDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "go")
	if err != nil {
		t.Fatal("Unexpected error creating temp dir", err)
	}
	defer os.RemoveAll(tmpDir)

	paths = &[]string{tmpDir}
	writeOutput = boolPtr(false)
	listFiles = boolPtr(false)

	writeTestFiles(t, testFiles, false, tmpDir)

	err = run()
	assert.Nil(t, err)

	// Without writeOutput set to true, inputs should be unchanged
	for name, contents := range testFiles {
		path := filepath.Join(tmpDir, name)
		bytes, err := os.ReadFile(path)
		if err != nil {
			t.Fatal("Unexpected error reading test file", err)
		}

		assert.Equal(
			t,
			strings.TrimSpace(contents),
			strings.TrimSpace(string(bytes)),
		)
	}

	writeOutput = boolPtr(true)
	err = run()
	assert.Nil(t, err)

	// Now, files should be modified in place
	for name, contents := range testFiles {
		path := filepath.Join(tmpDir, name)

		bytes, err := os.ReadFile(path)
		if err != nil {
			t.Fatal("Unexpected error reading test file", err)
		}

		assert.NotEqual(
			t,
			strings.TrimSpace(contents),
			strings.TrimSpace(string(bytes)),
		)
	}
}

func TestRunFilePaths(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "go")
	if err != nil {
		t.Fatal("Unexpected error creating temp dir", err)
	}
	defer os.RemoveAll(tmpDir)

	paths = &[]string{}
	writeOutput = boolPtr(true)
	listFiles = boolPtr(false)

	writeTestFiles(t, testFiles, true, tmpDir)

	err = run()
	assert.Nil(t, err)

	// Now, files should be modified in place
	for name, contents := range testFiles {
		path := filepath.Join(tmpDir, name)

		bytes, err := os.ReadFile(path)
		if err != nil {
			t.Fatal("Unexpected error reading test file", err)
		}

		assert.NotEqual(
			t,
			strings.TrimSpace(contents),
			strings.TrimSpace(string(bytes)),
		)
	}
}

func TestRunListFiles(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "go")
	if err != nil {
		t.Fatal("Unexpected error creating temp dir", err)
	}
	defer os.RemoveAll(tmpDir)

	paths = &[]string{}
	listFiles = boolPtr(true)

	updatedTestFiles := map[string]string{
		"test1.go": testFiles["test1.go"],
		"test2.go": testFiles["test2.go"],

		// File that doesn't need to be shortened
		"test3.go": "package main\n",
	}

	writeTestFiles(t, updatedTestFiles, true, tmpDir)

	output, err := captureStdout(t, run)
	assert.Nil(t, err)

	// Only first two files appear in output list
	expectedPaths := []string{
		filepath.Join(tmpDir, "test1.go"),
		filepath.Join(tmpDir, "test2.go"),
	}

	actualPaths := strings.Split(strings.TrimSpace(output), "\n")
	sort.Slice(actualPaths, func(i, j int) bool {
		return actualPaths[i] < actualPaths[j]
	})

	assert.Equal(
		t,
		expectedPaths,
		actualPaths,
	)
}

func boolPtr(b bool) *bool {
	return &b
}

func writeTestFiles(
	t *testing.T,
	fileContents map[string]string,
	addToPaths bool,
	tmpDir string,
) {
	for name, contents := range fileContents {
		path := filepath.Join(tmpDir, name)

		if addToPaths {
			tmpPaths := append(*paths, path)
			paths = &tmpPaths
		}

		err := os.WriteFile(path, []byte(contents), 0644)
		if err != nil {
			t.Fatal("Unexpected error writing test file", err)
		}
	}
}

func captureStdout(t *testing.T, f func() error) (string, error) {
	origStdout := os.Stdout
	defer func() {
		os.Stdout = origStdout
	}()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal("Unexpected error opening pipe", err)
	}
	os.Stdout = w

	resultErr := f()

	w.Close()
	outBytes, err := io.ReadAll(r)
	if err != nil {
		t.Fatal("Unexpected error reading result", err)
	}
	w.Close()

	return string(outBytes), resultErr
}
