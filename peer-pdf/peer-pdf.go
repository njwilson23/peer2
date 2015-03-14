package main

import (
	"fmt"
	"github.com/njwilson23/peer/config"
	"launchpad.net/gnuflag"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// USAGE:
// peer-pdf [query] [-r] [-p] [-o N]

// Walk a root, sending file matches to a slice
func SearchRoot(root string, searchstrs *[]string, out *[]string) {
	filepath.Walk(root, func(fnm string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("ERROR:", fnm, err)
			return filepath.SkipDir
		}

		var match bool
		if strings.ToLower(filepath.Ext(fnm)) == ".pdf" {

			match = true
			for _, searchstr := range *searchstrs {
				if !strings.Contains(strings.ToLower(fnm), strings.ToLower(searchstr)) {
					match = false
					break
				}
			}
			if match {
				*out = append(*out, fnm)
			}
		}
		return nil
	})
}

func main() {

	var searchstrs []string
	open := gnuflag.Int("o", 0, "Open option N using the configured reader")
	printpath := gnuflag.Bool("p", false, "Print full paths")
	gnuflag.Parse(true)

	configPath, err := config.FindConfig()
	if err != nil {
		panic(err.Error())
	}
	config := config.ParseConfig(configPath)

	if len(gnuflag.Args()) != 0 {
		searchstrs = gnuflag.Args()
	} else {
		fmt.Println("ERROR: must provide at least one search query")
		os.Exit(1)
	}
	if len(config.SearchRoots) == 0 {
		fmt.Println("WARNING: list of search roots is empty")
	}
	/* curdir, err := os.Getwd()
	if err != nil {
		fmt.Println("ERROR: cannot detect current directory")
		fmt.Println(err)
		os.Exit(1)
	}
	roots := append(config.SearchRoots, curdir) */

	// Search each directory root
	results := make([]string, 0)
	for _, root := range config.SearchRoots {
		SearchRoot(root, &searchstrs, &results)
	}

	if *open != 0 {
		cmd := exec.Command(config.Reader, results[*open-1])
		cmd.Start()
	} else {
		// Print matches
		for i, match := range results {
			if *printpath {
				fmt.Println(i+1, match)
			} else {
				fmt.Println(i+1, filepath.Base(match))
			}
		}
	}
}
