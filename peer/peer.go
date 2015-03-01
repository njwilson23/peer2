package main

import (
	"flag"
	"fmt"
	"github.com/njwilson23/peer2/bibtex"
	"strings"
)

// USAGE:
// peer [query] [--author query] [--year query] [--title query] [--file query]

func main() {

	//fnmSearch := flag.String("file", os.Args[1], "Search filenames")
	authorSearch := flag.String("author", "", "Search authors in BibTeX database")
	titleSearch := flag.String("title", "", "Search titles in BibTeX database")
	yearSearch := flag.Int("year", -999, "Restrict matches to YEAR")

	bibfile := flag.String("bibfile", "../bibtex/test.bib", "Location of BibTeX database")
	//openMatch := flag.Int("o", -1, "Open best matching file")

	flag.Parse()

	entries := make(chan bibtex.Entry)
	go bibtex.Read(*bibfile, entries)

	for entry := range entries {

		if (*titleSearch != "") && strings.Contains(entry.Title, *titleSearch) {
			fmt.Println(entry)
		} else if (*authorSearch != "") && strings.Contains(entry.Author, *authorSearch) {
			fmt.Println(entry)
		} else if entry.Year == *yearSearch {
			fmt.Println(entry)
		}

	}

}
