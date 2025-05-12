package main

import (
	"fmt"
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
	assert.Nil(t, err)

	fixturePaths := []string{}

	for _, fileInfo := range info {
		if fileInfo.IsDir() {
			continue
		} else if !strings.HasSuffix(fileInfo.Name(), ".go") {
			continue
		} else if strings.HasSuffix(fileInfo.Name(), "__exp.go") {
			continue
		} else if strings.HasPrefix(fileInfo.Name(), "editorconfig_") {
			continue
		}

		fixturePaths = append(
			fixturePaths,
			filepath.Join(fixturesDir, fileInfo.Name()),
		)
	}

	dotDir, err := os.MkdirTemp("", "dot")
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
		shortenedContents := generateShortenedContents(t, shortener, fixturePath)

		expectedPath := fixturePath[0:len(fixturePath)-3] + "__exp" + ".go"

		if os.Getenv("REGENERATE_TEST_OUTPUTS") == "true" {
			err := os.WriteFile(expectedPath, shortenedContents, 0o644)
			if err != nil {
				t.Fatalf(
					"Unexpected error writing output file %s: %+v",
					expectedPath,
					err,
				)
			}
		}

		expectedContents := readExpectedContents(t, fixturePath)

		assert.Equal(t, string(expectedContents), string(shortenedContents))
	}
}

func generateShortenedContents(t *testing.T, shortener *Shortener, path string) []byte {
	t.Helper()

	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf(
			"Unexpected error reading fixture %s: %+v",
			path,
			err,
		)
	}

	shortenedContents, err := shortener.Shorten(contents)
	assert.Nil(t, err)

	return shortenedContents
}

func readExpectedContents(t *testing.T, path string) []byte {
	t.Helper()

	expectedPath := path[0:len(path)-3] + "__exp" + ".go"

	expectedContents, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatalf(
			"Unexpected error reading expected file %s: %+v",
			expectedPath,
			err,
		)
	}

	return expectedContents
}

func TestShortenerWithEditorConfig(t *testing.T) {
	createEditorconfigFile(t, 150, 120, true)
	defer restoreEditorconfigFile(t)

	dotDir, err := os.MkdirTemp("", "dot")
	if err != nil {
		t.Fatalf("Error creating output directory for dot files: %+v", err)
	}
	defer os.RemoveAll(dotDir)

	shortener := NewShortener(
		ShortenerConfig{
			MaxLen:           0,
			CurrentMaxLen:    0,
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
	// When creating the shortener, the current max length is set to defaultMaxLen
	assert.Equal(t, defaultMaxLen, shortener.config.CurrentMaxLen)

	fixturePath := filepath.Join(fixturesDir, "editorconfig_max_len.go")

	shortener.SetCurrentMaxLen(fixturePath)
	shortenedContents := generateShortenedContents(t, shortener, fixturePath)

	expectedContents := readExpectedContents(t, fixturePath)

	assert.Equal(t, string(expectedContents), string(shortenedContents))
}

func TestCurrentMaxLen(t *testing.T) {
	t.Run("maxlen given by user", func(t *testing.T) {
		userMaxLen := 110

		shortener := NewShortener(ShortenerConfig{
			MaxLen: userMaxLen,
		})

		t.Run("without .editorconfig I should get max length from user", func(t *testing.T) {
			shortener.SetCurrentMaxLen("afile.go")
			assert.Equal(t, userMaxLen, shortener.config.CurrentMaxLen)
		})

		t.Run("with .editorconfig I should get max length from user", func(t *testing.T) {
			createEditorconfigFile(t, 150, 120, true)
			shortener.SetCurrentMaxLen("anotherfile.go")
			assert.Equal(t, userMaxLen, shortener.config.CurrentMaxLen)
			restoreEditorconfigFile(t)
		})
	})

	t.Run("no maxlen given by user", func(t *testing.T) {
		shortener := NewShortener(ShortenerConfig{
			MaxLen:           0,
			TabLen:           4,
			KeepAnnotations:  false,
			ShortenComments:  true,
			ReformatTags:     true,
			IgnoreGenerated:  true,
			BaseFormatterCmd: "gofmt",
			DotFile:          "",
			ChainSplitDots:   true,
		})

		t.Run("without .editorconfig I should get default max length", func(t *testing.T) {
			shortener.SetCurrentMaxLen("/tmp")
			assert.Equal(t, defaultMaxLen, shortener.config.CurrentMaxLen)
		})

		t.Run("with .editorconfig I should get max length from .editorconfig", func(t *testing.T) {
			createEditorconfigFile(t, 150, 0, true)
			shortener.SetCurrentMaxLen("afile.go")
			assert.Equal(t, 150, shortener.config.CurrentMaxLen)
			restoreEditorconfigFile(t)
		})

		t.Run("with Go section in .editorconfig I should get its max length", func(t *testing.T) {
			createEditorconfigFile(t, 150, 120, true)
			shortener.SetCurrentMaxLen("afile.go")
			assert.Equal(t, 120, shortener.config.CurrentMaxLen)
			restoreEditorconfigFile(t)
		})

		t.Run("with invalid .editorconfig I should get default max length", func(t *testing.T) {
			createEditorconfigFile(t, 0, 0, false)
			shortener.SetCurrentMaxLen("afile.go")
			assert.Equal(t, defaultMaxLen, shortener.config.CurrentMaxLen)
			restoreEditorconfigFile(t)
		})
	})
}

func existsEditorconfigFile(t *testing.T, filename string) bool {
	t.Helper()

	pwd, err := os.Getwd()
	require.NoError(t, err)

	_, err = os.Stat(filepath.Join(pwd, filename))
	if os.IsNotExist(err) {
		return false
	}
	require.NoError(t, err)

	return true
}

func backupEditorconfigFile(t *testing.T) {
	t.Helper()

	pwd, err := os.Getwd()
	require.NoError(t, err)

	err = os.Rename(filepath.Join(pwd, editorconfigFilename), filepath.Join(pwd, editorconfigBackupFilename))
	require.NoError(t, err)
}

func createEditorconfigFile(t *testing.T, maxLen int, goSectionMaxLen int, valid bool) {
	t.Helper()

	if existsEditorconfigFile(t, editorconfigFilename) {
		backupEditorconfigFile(t)
	}

	contents := make([]string, 0)

	if valid {
		contents = append(contents, `# EditorConfig is awesome: https://editorconfig.org`)
		contents = append(contents, `# top-most EditorConfig file`, ``)
		contents = append(contents, `root = true`, ``, `[*]`)
		contents = append(contents, fmt.Sprintf(`max_line_length = %d`, maxLen))

		if goSectionMaxLen > 0 {
			contents = append(contents, ``, `[*.go]`)
			contents = append(contents, fmt.Sprintf("max_line_length = %d", goSectionMaxLen))
		}
	} else {
		contents = append(contents, `<html></html>`)
	}

	pwd, err := os.Getwd()
	require.NoError(t, err)

	file, err := os.Create(filepath.Join(pwd, editorconfigFilename))
	require.NoError(t, err)

	_, err = file.WriteString(strings.Join(contents, "\n"))
	require.NoError(t, err)

	err = file.Close()
	require.NoError(t, err)
}

func restoreEditorconfigFile(t *testing.T) {
	t.Helper()

	pwd, err := os.Getwd()
	require.NoError(t, err)

	if existsEditorconfigFile(t, editorconfigBackupFilename) {
		err = os.Rename(filepath.Join(pwd, editorconfigBackupFilename), filepath.Join(pwd, editorconfigFilename))
		require.NoError(t, err)
		return
	}

	err = os.Remove(filepath.Join(pwd, editorconfigFilename))
	require.NoError(t, err)
}

const (
	editorconfigFilename       = ".editorconfig"
	editorconfigBackupFilename = ".editorconfig-backup"
)
