package main

import (
	"flag"
	"fmt"
	"github.com/njwilson23/peer2/bibtex"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// USAGE:
// peer [query] [--author query] [--year query] [--title query] [--file query]

type Config struct {
	Reader      string
	Bibfiles    []string
	SearchRoots []string
}

func main() {

	// Parse configuration
	var config Config
	configData, err := ioutil.ReadFile(".peer.yaml")
	if err != nil {
		fmt.Println("ERROR:", err)
	}
	err = yaml.Unmarshal(configData, &config)

	// Parse command line options
	//fnmSearch := flag.String("file", os.Args[1], "Search filenames")
	authorSearch := flag.String("author", "", "Search authors in BibTeX database")
	titleSearch := flag.String("title", "", "Search titles in BibTeX database")
	yearSearch := flag.Int("year", -999, "Restrict matches to YEAR")

	bibfile := flag.String("bibfile", "", "Location of BibTeX database")
	//openMatch := flag.Int("o", -1, "Open best matching file")

	flag.Parse()

	if *bibfile != "" {
		config.Bibfiles = append(config.Bibfiles, *bibfile)
	}

	// Search BibTeX entries for matches
	//matches := make([]bibtex.Entry, 0)
	//var matches []bibtex.Entry
	for _, bibfile := range config.Bibfiles {

		entries := make(chan bibtex.Entry)
		go bibtex.Read(bibfile, entries)

		for entry := range entries {

			if entry.TestAuthor(*authorSearch) &&
				entry.TestTitle(*titleSearch) &&
				entry.TestYear(*yearSearch) {
				//fmt.Println(entry)
				fmt.Println(fmt.Sprintf("%v (%v), \"%v\"", entry.Author, entry.Year, entry.Title))
				//matches = append(matches, entry)
			}

		}
	}

	// Print the results
	/*for entry := range matches {
		fmt.Println(fmt.Sprintf("%v (%v), \"%v\"", entry.Author, entry.Year, entry.Title))
	}*/

}
