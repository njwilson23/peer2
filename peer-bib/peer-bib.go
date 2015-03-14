package main

import (
	"flag"
	"fmt"
	"github.com/njwilson23/peer/bibtex"
	"github.com/njwilson23/peer/config"
)

// USAGE:
// peer-bib [query] [--author query] [--year query] [--title query] [--file query]

func main() {

	configPath, err := config.FindConfig()
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}
	config := config.ParseConfig(configPath)

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
	var matches []bibtex.Entry
	for _, bibfile := range config.Bibfiles {

		entries := make(chan bibtex.Entry)
		go bibtex.Read(bibfile, entries)

		for entry := range entries {

			if entry.TestAuthor(*authorSearch) &&
				entry.TestTitle(*titleSearch) &&
				entry.TestYear(*yearSearch) {
				matches = append(matches, entry)
			}

		}
	}

	// Print the results
	for _, entry := range matches {
		fmt.Println(fmt.Sprintf("@%v\n%v (%v)\n\t\"%v\"\n",
			entry.BibTeXkey,
			entry.Author,
			entry.Year,
			entry.Title))
	}

}
