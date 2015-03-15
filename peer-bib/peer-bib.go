package main

import (
	"fmt"
	"github.com/njwilson23/peer/bibtex"
	"github.com/njwilson23/peer/config"
	"launchpad.net/gnuflag"
	"sort"
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
	//fnmSearch := gnuflag.String("file", os.Args[1], "Search filenames")
	authorSearch := gnuflag.String("author", "", "Search authors in BibTeX database")
	titleSearch := gnuflag.String("title", "", "Search titles in BibTeX database")
	yearSearch := gnuflag.Int("year", -999, "Restrict matches to YEAR")

	bibfile := gnuflag.String("bibfile", "", "Location of BibTeX database")
	//openMatch := gnuflag.Int("o", -1, "Open best matching file")

	gnuflag.Parse(true)

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

	sort.Sort(bibtex.ByYear(matches))
	for _, entry := range matches {
		fmt.Println(fmt.Sprintf("@%v\n%v (%v)\n\t\"%v\"\n",
			entry.BibTeXkey,
			entry.Author,
			entry.Year,
			entry.Title))
	}

}
