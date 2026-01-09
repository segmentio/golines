package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const fixturesDir = "_fixtures"

// TestShortener verifies the core shortening functionality on the files in the _fixtures
// directory. To update the expected outputs, run tests with the REGENERATE_TEST_OUTPUTS
// environment variable set to "true".
func TestShortener(t *testing.T) {
	info, err := os.ReadDir(fixturesDir)
	require.NoError(t, err)

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
			filepath.Join(fixturesDir, fileInfo.Name()),
		)
	}

	dotDir := t.TempDir()

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
		contents, err := os.ReadFile(fixturePath)
		require.NoErrorf(t, err, "Unexpected error reading fixture %s", fixturePath)

		shortenedContents, err := shortener.Shorten(contents)
		require.NoError(t, err)

		expectedPath := fixturePath[0:len(fixturePath)-3] + "__exp" + ".go"

		if os.Getenv("REGENERATE_TEST_OUTPUTS") == "true" {
			err := os.WriteFile(expectedPath, shortenedContents, 0644)
			require.NoErrorf(t, err, "Unexpected error writing output file %s", expectedPath)
		}

		expectedContents, err := os.ReadFile(expectedPath)
		require.NoErrorf(t, err, "Unexpected error reading expected file %s", expectedPath)

		assert.Equal(t, string(expectedContents), string(shortenedContents))
	}
}
