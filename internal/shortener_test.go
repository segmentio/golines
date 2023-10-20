package internal

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testdataDir string

func init() {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("Can't get current filename")
	}
	testdataDir = filepath.Join(filepath.Dir(file), "testdata")
}

// TestShortener verifies the core shortening functionality on the files in the testdata
// directory. To update the expected outputs, run tests with the REGENERATE_TEST_OUTPUTS
// environment variable set to "true".
func TestShortener(t *testing.T) {
	info, err := ioutil.ReadDir(testdataDir)
	assert.Nil(t, err)

	fixturePaths := []string{}

	for _, fileInfo := range info {
		if fileInfo.IsDir() {
			continue
		} else if !strings.HasSuffix(fileInfo.Name(), ".go") {
			continue
		} else if strings.HasSuffix(fileInfo.Name(), "__exp.go") {
			continue
		}

		fixturePaths = append(
			fixturePaths,
			filepath.Join(testdataDir, fileInfo.Name()),
		)
	}

	dotDir, err := ioutil.TempDir("", "dot")
	if err != nil {
		t.Fatalf("Error creating output directory for dot files: %+v", err)
	}
	defer os.RemoveAll(dotDir)

	shortener := NewShortener(
		ShortenerConfig{
			MaxLen:           100,
			TabLen:           4,
			KeepAnnotations:  false,
			ShortenComments:  true,
			ReformatTags:     true,
			IgnoreGenerated:  true,
			BaseFormatterCmd: "gofmt",
			DotFile:          filepath.Join(dotDir, "out.dot"),
			ChainSplitDots:   true,
		},
	)

	for _, fixturePath := range fixturePaths {
		contents, err := ioutil.ReadFile(fixturePath)
		if err != nil {
			t.Fatalf(
				"Unexpected error reading fixture %s: %+v",
				fixturePath,
				err,
			)
		}

		shortenedContents, err := shortener.Shorten(contents)
		assert.Nil(t, err)

		expectedPath := fixturePath[0:len(fixturePath)-3] + "__exp" + ".go"

		if os.Getenv("REGENERATE_TEST_OUTPUTS") == "true" {
			err := ioutil.WriteFile(expectedPath, shortenedContents, 0o600)
			if err != nil {
				t.Fatalf(
					"Unexpected error writing output file %s: %+v",
					expectedPath,
					err,
				)
			}
		}

		expectedContents, err := ioutil.ReadFile(expectedPath)
		if err != nil {
			t.Fatalf(
				"Unexpected error reading expected file %s: %+v",
				expectedPath,
				err,
			)
		}

		assert.Equal(t, string(expectedContents), string(shortenedContents))
	}
}
