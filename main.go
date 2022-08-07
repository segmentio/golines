package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime/pprof"
	"strings"

	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const (
	versionStr = "0.11.0"
)

var (
	// Flags
	baseFormatterCmd = kingpin.Flag(
		"base-formatter",
		"Base formatter to use").Default("").String()
	chainSplitDots = kingpin.Flag(
		"chain-split-dots",
		"Split chained methods on the dots as opposed to the arguments").
		Default("true").Bool()
	debug = kingpin.Flag(
		"debug",
		"Show debug output").Short('d').Default("false").Bool()
	dotFile = kingpin.Flag(
		"dot-file",
		"Path to dot representation of AST graph").Default("").String()
	dryRun = kingpin.Flag(
		"dry-run",
		"Show diffs without writing anything").Default("false").Bool()
	ignoreGenerated = kingpin.Flag(
		"ignore-generated",
		"Ignore generated go files").Default("true").Bool()
	ignoredDirs = kingpin.Flag(
		"ignored-dirs",
		"Directories to ignore").Default("vendor", "node_modules", ".git").Strings()
	keepAnnotations = kingpin.Flag(
		"keep-annotations",
		"Keep shortening annotations in final output").Default("false").Bool()
	listFiles = kingpin.Flag(
		"list-files",
		"List files that would be reformatted by this tool").Short('l').Default("false").Bool()
	maxLen = kingpin.Flag(
		"max-len",
		"Target maximum line length").Short('m').Default("100").Int()
	profile = kingpin.Flag(
		"profile",
		"Path to profile output").Default("").String()
	reformatTags = kingpin.Flag(
		"reformat-tags",
		"Reformat struct tags").Default("true").Bool()
	shortenComments = kingpin.Flag(
		"shorten-comments",
		"Shorten single-line comments").Default("false").Bool()
	tabLen = kingpin.Flag(
		"tab-len",
		"Length of a tab").Short('t').Default("4").Int()
	version = kingpin.Flag(
		"version",
		"Print out version and exit").Default("false").Bool()
	writeOutput = kingpin.Flag(
		"write-output",
		"Write output to source instead of stdout").Short('w').Default("false").Bool()

	// Args
	paths = kingpin.Arg(
		"paths",
		"Paths to format",
	).Strings()
)

func main() {
	kingpin.Parse()
	if *debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	if *version {
		fmt.Printf("golines v%s\n", versionStr)
		return
	}

	if *profile != "" {
		f, err := os.Create(*profile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	log.SetFormatter(&prefixed.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		ForceFormatting: true,
	})

	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	config := ShortenerConfig{
		MaxLen:           *maxLen,
		TabLen:           *tabLen,
		KeepAnnotations:  *keepAnnotations,
		ShortenComments:  *shortenComments,
		ReformatTags:     *reformatTags,
		IgnoreGenerated:  *ignoreGenerated,
		DotFile:          *dotFile,
		BaseFormatterCmd: *baseFormatterCmd,
		ChainSplitDots:   *chainSplitDots,
	}
	shortener := NewShortener(config)

	if len(*paths) == 0 {
		// Read input from stdin
		contents, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return err
		}

		result, err := shortener.Shorten(contents)
		if err != nil {
			return err
		}
		err = handleOutput("", contents, result)
		if err != nil {
			return err
		}
	} else {
		// Read inputs from paths provided in arguments
		for _, path := range *paths {
			switch info, err := os.Stat(path); {
			case err != nil:
				return err
			case info.IsDir():
				// Path is a directory- walk it
				err = filepath.Walk(
					path,
					func(subPath string, subInfo os.FileInfo, err error) error {
						if err != nil {
							return err
						}

						components := strings.Split(subPath, "/")
						for _, component := range components {
							for _, ignoredDir := range *ignoredDirs {
								if component == ignoredDir {
									return filepath.SkipDir
								}
							}
						}

						if !subInfo.IsDir() && strings.HasSuffix(subPath, ".go") {
							// Shorten file and generate output
							contents, result, err := processFile(shortener, subPath)
							if err != nil {
								return err
							}
							err = handleOutput(subPath, contents, result)
							if err != nil {
								return err
							}
						}

						return nil
					},
				)
				if err != nil {
					return err
				}
			default:
				// Path is a file
				contents, result, err := processFile(shortener, path)
				if err != nil {
					return err
				}
				err = handleOutput(path, contents, result)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// processFile uses the provided Shortener instance to shorten the lines
// in a file. It returns the original contents (useful for debugging), the
// shortened version, and an error.
func processFile(shortener *Shortener, path string) ([]byte, []byte, error) {
	_, fileName := filepath.Split(path)
	if *ignoreGenerated && strings.HasPrefix(fileName, "generated_") {
		return nil, nil, nil
	}

	log.Debugf("Processing file %s", path)

	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}

	result, err := shortener.Shorten(contents)
	return contents, result, err
}

// handleOutput generates output according to the value of the tool's
// flags; depending on the latter, the output might be written over
// the source file, printed to stdout, etc.
func handleOutput(path string, contents []byte, result []byte) error {
	if contents == nil {
		return nil
	} else if *dryRun {
		return PrettyDiff(path, contents, result)
	} else if *listFiles {
		if !bytes.Equal(contents, result) {
			fmt.Println(path)
		}

		return nil
	} else if *writeOutput {
		if path == "" {
			return errors.New("No path to write out to")
		}

		info, err := os.Stat(path)
		if err != nil {
			return err
		}

		if bytes.Equal(contents, result) {
			log.Debugf("Contents unchanged, skipping write")
			return nil
		} else {
			log.Debugf("Contents changed, writing output to %s", path)
			return ioutil.WriteFile(path, result, info.Mode())
		}
	} else {
		fmt.Print(string(result))
		return nil
	}
}
