package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	maxLen = kingpin.Flag(
		"max-len",
		"Target maximum line length").Short('m').Default("100").Int()
	tabLen      = kingpin.Flag("tab-len", "Length of a tab").Short('t').Default("4").Int()
	writeOutput = kingpin.Flag(
		"write-output",
		"Write output to source instead of stdout").Short('w').Default("false").Bool()
	listFiles = kingpin.Flag(
		"list-files",
		"List files that would be reformatted by this tool").Short('l').Default("false").Bool()
	keepAnnotations = kingpin.Flag(
		"keep-annotations",
		"Keep shortening annotations in final output").Default("false").Bool()
	shortenComments = kingpin.Flag(
		"shorten-comments",
		"Shorten single-line comments").Default("false").Bool()
	ignoreGenerated = kingpin.Flag(
		"ignore-generated",
		"Ignore generated go files").Default("true").Bool()
	debug  = kingpin.Flag("debug", "Show debug output").Short('d').Default("false").Bool()
	dryRun = kingpin.Flag("dry-run", "Show diffs without writing anything").
		Default("false").Bool()
	ignoredDirs = kingpin.Flag("ignored-dirs", "Directories to ignore").
			Default("vendor", "node_modules", ".git").Strings()
	baseFormatter = kingpin.Flag("base-formatter", "Base formatter to use").
			Default("").String()
	paths = kingpin.Arg("paths", "Paths to format").Strings()
)

func main() {
	kingpin.Parse()

	if *debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	log.SetFormatter(&prefixed.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		ForceFormatting: true,
	})

	shortener := NewShortener(
		*maxLen,
		*tabLen,
		*keepAnnotations,
		*shortenComments,
		*ignoreGenerated,
		*baseFormatter,
	)

	if len(*paths) == 0 {
		contents, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}

		result, err := shortener.Shorten(contents)
		if err != nil {
			log.Fatal(err)
		}
		err = handleOutput("", contents, result)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		for _, path := range *paths {
			switch info, err := os.Stat(path); {
			case err != nil:
				log.Fatal(err)
			case info.IsDir():
				err = filepath.Walk(
					path,
					func(subPath string, subInfo os.FileInfo, err error) error {
						if err != nil {
							log.Fatal(err)
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
							contents, result, err := processFile(shortener, subPath)
							if err != nil {
								log.Fatal(err)
							}
							err = handleOutput(subPath, contents, result)
							if err != nil {
								log.Fatal(err)
							}
						}

						return nil
					},
				)
				if err != nil {
					log.Fatal(err)
				}
			default:
				contents, result, err := processFile(shortener, path)
				if err != nil {
					log.Fatal(err)
				}
				err = handleOutput(path, contents, result)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}

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

func handleOutput(path string, contents []byte, result []byte) error {
	if contents == nil {
		return nil
	} else if *dryRun {
		return PrettyDiff(path, contents, result)
	} else if *listFiles {
		if !reflect.DeepEqual(contents, result) {
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

		log.Debugf("Writing output to %s", path)
		return ioutil.WriteFile(path, result, info.Mode())
	} else {
		fmt.Print(string(result))
		return nil
	}
}
