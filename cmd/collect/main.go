package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"

	mc "github.com/kascote/meta_collector"
)

// CollectorOptions holds the CLI flags
type CollectorOptions struct {
	Twitter  bool
	Facebook bool
}

var version = "unknown"
var flags CollectorOptions
var showHelp bool
var showVersion bool

// Init Initialize CLI flags
func init() {
	flag.BoolVar(&showHelp, "h", false, "show this help text")
	flag.BoolVar(&showVersion, "v", false, "show version information")
}

func main() {

	flag.Parse()
	if showHelp {
		printVersion()
		flag.Usage()
		os.Exit(0)
	}

	if showVersion {
		printVersion()
		os.Exit(0)
	}

	if (flag.NArg() > 1) || (flag.NArg() == 0) {
		// Only support 1 argument
		flag.Usage()
		os.Exit(1)
	}

	var attrs *mc.Attributes
	var err error
	if isPath(flag.Arg(0)) {
		if _, err := os.Stat(flag.Arg(0)); os.IsNotExist(err) {
			fmt.Printf("ERROR: %s\n", err.Error())
			os.Exit(1)
		}
		attrs, err = mc.ParseFile(flag.Arg(0))
		if err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
			os.Exit(1)
		}
	} else {
		attrs, err = mc.ParseHTML(flag.Arg(0))
		if err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
			os.Exit(1)
		}
	}

	js, err := json.MarshalIndent(attrs, "", "\t")
	if err != nil {
		fmt.Printf("ERROR %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("%s\n", js)
	os.Exit(0)

}

func printVersion() {
	fmt.Fprintf(os.Stderr, "collector version: %s (%s)\n", version, runtime.GOOS)
}

// Returns true for things which look like paths (start with "~", "." or "/").
func isPath(path string) bool {
	return strings.HasPrefix(path, "~") ||
		strings.HasPrefix(path, ".") ||
		strings.HasPrefix(path, "/")
}
