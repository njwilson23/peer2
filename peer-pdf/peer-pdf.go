package main

import (
	//"flag"
	"fmt"
	"github.com/njwilson23/peer2/config"
	"os"
	"path/filepath"
	"strings"
)

// USAGE:
// peer-pdf [query] [-r] [-p] [-o N]

// Walk a root, sending file matches to out channel
func SearchRoot(root string, searchstr string, out *[]string) {
	filepath.Walk(root, func(fnm string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("ERROR:", fnm, err)
			return filepath.SkipDir
		}
		if strings.ToLower(filepath.Ext(fnm)) == ".pdf" {
			if strings.Contains(strings.ToLower(fnm), searchstr) {
				*out = append(*out, fnm)
			}
		}
		return nil
	})
}

func main() {

	var searchstr string
	config := config.ParseConfig(".peer.yaml")
	if len(os.Args) != 1 {
		searchstr = os.Args[1]
	} else {
		fmt.Println("ERROR: must provide at least one search query")
		os.Exit(1)
	}
	if len(config.SearchRoots) == 0 {
		fmt.Println("ERROR: list of search roots is empty")
		os.Exit(1)
	}

	// Search each directory root
	results := make([]string, 0)
	for _, root := range config.SearchRoots {
		SearchRoot(root, searchstr, &results)
	}

	// Print matches
	for i, match := range results {
		fmt.Println(i+1, filepath.Base(match))
	}
}
